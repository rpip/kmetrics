package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"
)

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

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
