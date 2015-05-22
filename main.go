package main

import (
	"flag"
	"log"

	"github.com/StreamMeBots/meep/pkg/config"
	"github.com/StreamMeBots/meep/pkg/db"
	"github.com/StreamMeBots/meep/routes"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	// setup loggger
	log.SetPrefix("[BOT] ")

	// parse command line flags
	flag.Parse()

	// load from config path?
	config.CheckConfigPath()

	// Open bolt DB
	db.Open()

	// routes and server get attached to the gin engine
	r := gin.Default()

	// setup routes
	routes.Init(r)

	// no cache please
	r.Use(func(ctx *gin.Context) {
		ctx.Request.Header.Del("If-Modified-Since")
		ctx.Writer.Header().Add("Cache-Control", "no-cache")
	})

	// All undefined routes will get served from the client directory.
	// If a file is not found the client/index.html gets served
	r.Use(static.Serve("/", Assets()))

	// start server
	if err := r.Run(":8888"); err != nil {
		log.Fatal(err)
	}
}
