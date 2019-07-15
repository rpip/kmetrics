package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Healthcheck represents information about the service health
// can be extended to include system information
type Healthcheck struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

// errorResponse for more user-friendly errors to return to the user or propage
// 404: Not found
// 500: Internal Server Error
type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HandlerFunc extends the http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request, AppContext)

// makeHandler allows us to pass settings/config to our handlers, avoiding global variables
func makeHandler(ctx AppContext, fn func(http.ResponseWriter, *http.Request, AppContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	}
}

// HealthCheckHandler returns info about the app
func HealthCheckHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	check := Healthcheck{
		AppName: appName,
		Version: ctx.Version,
	}
	ctx.Render.JSON(w, http.StatusOK, check)
}

// ListServicesHandler returns a list of pods running in the cluster in namespace default
func ListServicesHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	pods, err := ctx.kube.GetServices("default", "")
	if err != nil {
		response := errorResponse{
			Status:  "500",
			Message: "can't retrieve pods %d",
		}
		log.Println(err)
		ctx.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	ctx.Render.JSON(w, http.StatusOK, pods)
}

// SearchServicesHandler returns a list of pods in the cluster in namespace default
// that are part of the same applicationGroup:
func SearchServicesHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	vars := mux.Vars(req)
	fmt.Println("===> group: ", vars["group"])

	pods, err := ctx.kube.GetServices("default", vars["group"])
	if err != nil {
		response := errorResponse{
			Status:  "500",
			Message: "Pod search failed %d",
		}
		log.Println(err)
		ctx.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	ctx.Render.JSON(w, http.StatusOK, pods)
}
