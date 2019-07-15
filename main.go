package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/unrolled/render"
)

const (
	localEnv = "local"
	appName  = "KMetrics"
)

// AppContext holds service configuration data
type AppContext struct {
	Render        *render.Render
	Version       string
	Env           string
	Port          string
	isDevelopment bool
	kube          *kubeClient
}

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

func main() {

	// load application configs
	ctx := loadConfig()

	// start application
	StartServer(ctx, routes)
}

func loadConfig() AppContext {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		// environment variables
		env         = getEnv("ENV", localEnv)
		port        = getEnv("PORT", "3001")
		versionFile = getEnv("VERSION", "VERSION")
		kubeconfig  = getEnv("kubeconfig", filepath.Join(homeDir(), ".kube", "config"))
	)

	// read version from file
	version, err := parseVersionFile(versionFile)

	if err != nil {
		log.Fatal(err)
	}

	// initialize application context
	return AppContext{
		Render:        render.New(),
		Version:       version,
		Env:           env,
		Port:          port,
		isDevelopment: env == localEnv,
		kube:          NewKubeClient(kubeconfig),
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
