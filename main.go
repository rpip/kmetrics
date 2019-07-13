package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/unrolled/render"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"
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
	ctx := loadCtx()

	// start application
	StartServer(ctx, routes)
}

func loadCtx() AppContext {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		// environment variables
		env         = os.Getenv("ENV")
		port        = os.Getenv("PORT")
		versionFile = os.Getenv("VERSION")
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
	}

}

// StartServer wraps the routers with the app context and add-in middlewares
func StartServer(ctx AppContext, routes []Route) {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = makeHandler(ctx, route.HandlerFunc)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	// secure the server
	secureMiddleware := secure.New(secure.Options{
		IsDevelopment:      ctx.isDevelopment, // if in development, causes the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored
		AllowedHosts:       []string{},        // list of domain names that are allowed (CORS)
		ContentTypeNosniff: true,              // adds the X-Content-Type-Options header with the value `nosniff`
		BrowserXssFilter:   true,              // adds the X-XSS-Protection header with the value `1; mode=block`
	})

	// init http server and start
	n := negroni.Classic()
	n.Use(negroni.NewLogger())
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.UseHandler(router)
	log.Printf("===> Starting app (v%s) on port %s, env: %s", ctx.Version, ctx.Port, ctx.Env)
	if ctx.Env == localEnv {
		n.Run("localhost:" + ctx.Port)
	} else {
		n.Run(":" + ctx.Port)
	}
}

// ParseVersionFile returns the version as a string, parsing and validating a file given the path
func parseVersionFile(versionPath string) (string, error) {
	dat, err := ioutil.ReadFile(versionPath)
	if err != nil {
		return "", errors.Wrap(err, "error reading version file")
	}
	version := string(dat)
	version = strings.Trim(strings.Trim(version, "\n"), " ")
	// taken from https://github.com/sindresorhus/semver-regex
	semverRegex := `^v?(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)(?:-[\da-z\-]+(?:\.[\da-z\-]+)*)?(?:\+[\da-z\-]+(?:\.[\da-z\-]+)*)?$`
	match, err := regexp.MatchString(semverRegex, version)
	if err != nil {
		return "", errors.Wrap(err, "error executing regex match")
	}
	if !match {
		return "", errors.New("invalid version number")
	}
	return version, nil
}
