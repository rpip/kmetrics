package main

import (
	"net/http"
)

// HandlerFunc extends the http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request, AppContext)

// makeHandler allows us to pass settings/config to our handlers, avoiding global variables
func makeHandler(ctx AppContext, fn func(http.ResponseWriter, *http.Request, AppContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	}
}

// HealthHandler returns info about the app
func HealthHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	check := Healthcheck{
		AppName: appName,
		Version: ctx.Version,
	}
	ctx.Render.JSON(w, http.StatusOK, check)
}
