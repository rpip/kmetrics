package main

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/unrolled/render"
)

const (
	localEnv = "local"
	appName  = "KMetrics"
)

// default path to kube config
var kubeConfigPath = filepath.Join(homeDir(), ".kube", "config")

// AppContext holds service configuration data
type AppContext struct {
	Render        *render.Render
	Version       string
	Env           string
	Port          string
	isDevelopment bool
	kube          *kubeClient
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
		kubeconfig  = getEnv("kubeconfig", kubeConfigPath)
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
