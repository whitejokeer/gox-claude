package routing

import (
	"net/http"
)

// RegisterRoutes registers all app routes
func RegisterRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Home page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(homePageHTML))
	})

	// Static files
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
}

var homePageHTML = `<!DOCTYPE html>
<html>
<head>
	<title>GOX Framework</title>
	<style>
		body { font-family: system-ui; max-width: 800px; margin: 50px auto; padding: 20px; }
		h1 { color: #333; }
		.info { background: #f0f0f0; padding: 20px; border-radius: 8px; }
	</style>
</head>
<body>
	<h1>Welcome to GOX Framework</h1>
	<div class="info">
		<p>Your app is running successfully!</p>
		<p>Start building your application:</p>
		<ul>
			<li>Create pages in <code>app/pages/</code></li>
			<li>Add components in <code>app/components/</code></li>
			<li>Generate services with <code>gox generate service [name]</code></li>
		</ul>
	</div>
</body>
</html>`
