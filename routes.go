package main

// Route represents an HTTP router with info and the actual request handler
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerFunc
}

// Routes holds the http routers defined in this service
type Routes []Route

var routes = Routes{
	// main endpoints
	Route{"SearchServices", "GET", "/services/{group}", SearchServicesHandler},
	Route{"ListServices", "GET", "/services", ListServicesHandler},

	// service health check endpoint
	Route{"Health", "GET", "/health", HealthCheckHandler},
}
