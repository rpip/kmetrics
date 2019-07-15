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
	Route{"Health", "GET", "/health", HealthHandler},
	// Kubernetes
	Route{"ListServices", "GET", "/services", ListServicesHandler},
	//Route{"GetUser", "GET", "/users/{uid:[0-9]+}", GetUserHandler},
}
