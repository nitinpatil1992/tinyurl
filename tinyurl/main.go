package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tinyurl/app"
	"tinyurl/config"
)

var (
	appConfig *config.Config
)

func init() {
	env := os.Getenv("env")
	appConfig = config.Init(env)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	app.RegisterHandlers(appConfig)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appConfig.Port), nil))
}
