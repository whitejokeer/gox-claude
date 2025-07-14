package discovery

import (
	"context"
	"fmt"
	"net/url"
)

// ServiceDiscovery defines the interface for service discovery
type ServiceDiscovery interface {
	// Register registers a service instance
	Register(ctx context.Context, service ServiceInstance) error
	
	// Deregister removes a service instance
	Deregister(ctx context.Context, serviceID string) error
	
	// Discover returns healthy instances of a service
	Discover(ctx context.Context, serviceName string) ([]ServiceInstance, error)
	
	// Health performs a health check on a service instance
	Health(ctx context.Context, serviceID string) error
}

// ServiceInstance represents a service instance
type ServiceInstance struct {
	ID          string
	Name        string
	Address     string
	Port        int
	HealthCheck string
	Tags        []string
	Metadata    map[string]string
}

// URL returns the full URL for the service instance
func (s ServiceInstance) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Address, s.Port)
}

// ParsedURL returns a parsed URL for the service instance
func (s ServiceInstance) ParsedURL() (*url.URL, error) {
	return url.Parse(s.URL())
}

// StaticDiscovery implements a simple static service discovery
type StaticDiscovery struct {
	services map[string][]ServiceInstance
}

// NewStaticDiscovery creates a new static discovery instance
func NewStaticDiscovery() *StaticDiscovery {
	return &StaticDiscovery{
		services: make(map[string][]ServiceInstance),
	}
}

// Register registers a service instance
func (s *StaticDiscovery) Register(ctx context.Context, service ServiceInstance) error {
	s.services[service.Name] = append(s.services[service.Name], service)
	return nil
}

// Deregister removes a service instance
func (s *StaticDiscovery) Deregister(ctx context.Context, serviceID string) error {
	for name, instances := range s.services {
		for i, instance := range instances {
			if instance.ID == serviceID {
				s.services[name] = append(instances[:i], instances[i+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("service instance not found: %s", serviceID)
}

// Discover returns healthy instances of a service
func (s *StaticDiscovery) Discover(ctx context.Context, serviceName string) ([]ServiceInstance, error) {
	instances, ok := s.services[serviceName]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", serviceName)
	}
	return instances, nil
}

// Health performs a health check on a service instance
func (s *StaticDiscovery) Health(ctx context.Context, serviceID string) error {
	// In static discovery, we assume services are healthy
	return nil
}
