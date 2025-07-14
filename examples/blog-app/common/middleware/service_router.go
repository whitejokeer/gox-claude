package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// ServiceRouter handles routing requests to backend services
type ServiceRouter struct {
	services map[string]*ServiceConfig
	mu       sync.RWMutex
}

// ServiceConfig holds configuration for a backend service
type ServiceConfig struct {
	URL          *url.URL
	HealthCheck  string
	CircuitBreaker *CircuitBreaker
}

// CircuitBreaker implements basic circuit breaker pattern
type CircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration
	failures         int
	lastFailureTime  time.Time
	state            string // "closed", "open", "half-open"
	mu               sync.Mutex
}

// NewServiceRouter creates a new service router
func NewServiceRouter() *ServiceRouter {
	sr := &ServiceRouter{
		services: make(map[string]*ServiceConfig),
	}
	sr.loadStaticServices()
	return sr
}

// loadStaticServices loads service URLs from environment variables
func (sr *ServiceRouter) loadStaticServices() {
	// Load services from environment variables
	// Format: SERVICE_<NAME>=http://host:port
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := parts[0]
		value := parts[1]
		
		if strings.HasPrefix(key, "SERVICE_") && !strings.HasSuffix(key, "_REGISTRY") {
			serviceName := strings.ToLower(strings.TrimPrefix(key, "SERVICE_"))
			if u, err := url.Parse(value); err == nil {
				sr.RegisterService(serviceName, u.String())
			}
		}
	}
}

// RegisterService registers a new service
func (sr *ServiceRouter) RegisterService(name string, serviceURL string) error {
	u, err := url.Parse(serviceURL)
	if err != nil {
		return fmt.Errorf("invalid service URL: %%w", err)
	}

	sr.mu.Lock()
	defer sr.mu.Unlock()

	sr.services[name] = &ServiceConfig{
		URL:         u,
		HealthCheck: u.String() + "/health",
		CircuitBreaker: &CircuitBreaker{
			failureThreshold: 5,
			resetTimeout:     60 * time.Second,
			state:            "closed",
		},
	}

	return nil
}

// Middleware returns the HTTP middleware
func (sr *ServiceRouter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is an API route
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		// Extract service name from path
		// Format: /api/service-name/...
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
		if len(parts) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		serviceName := parts[0]
		
		// Get service config
		sr.mu.RLock()
		service, exists := sr.services[serviceName]
		sr.mu.RUnlock()

		if !exists {
			// No service found, continue to next handler
			next.ServeHTTP(w, r)
			return
		}

		// Check circuit breaker
		if !service.CircuitBreaker.CanRequest() {
			http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(service.URL)
		
		// Customize the director to rewrite the path
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			// Remove /api/service-name prefix
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api/"+serviceName)
			if !strings.HasPrefix(req.URL.Path, "/") {
				req.URL.Path = "/" + req.URL.Path
			}
			// Add X-Forwarded headers
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.Header.Set("X-Forwarded-Proto", "http")
			if r.TLS != nil {
				req.Header.Set("X-Forwarded-Proto", "https")
			}
		}

		// Handle errors
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			service.CircuitBreaker.RecordFailure()
			http.Error(w, "Service unavailable", http.StatusBadGateway)
		}

		// Add response modifier to record success
		proxy.ModifyResponse = func(resp *http.Response) error {
			if resp.StatusCode < 500 {
				service.CircuitBreaker.RecordSuccess()
			} else {
				service.CircuitBreaker.RecordFailure()
			}
			return nil
		}

		proxy.ServeHTTP(w, r)
	})
}

// Circuit Breaker methods
func (cb *CircuitBreaker) CanRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case "open":
		// Check if we should transition to half-open
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = "half-open"
			cb.failures = 0
			return true
		}
		return false
	case "half-open":
		// Allow one request to test
		return true
	default: // closed
		return true
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == "half-open" {
		cb.state = "closed"
	}
	cb.failures = 0
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailureTime = time.Now()

	if cb.failures >= cb.failureThreshold {
		cb.state = "open"
	}
}
