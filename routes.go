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
	Route{"ListServices", "GET", "/services", ListServicesHandler},
	Route{"SearchServices", "GET", "/services/{group}", SearchServicesHandler},

	// service health endpoint
	Route{"Health", "GET", "/health", HealthHandler},
}
