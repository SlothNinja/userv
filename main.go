package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	"github.com/SlothNinja/cookie"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/sn"
	ucon "github.com/SlothNinja/user-controller"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

const (
	hashKeyLength      = 64
	blockKeyLength     = 32
	sessionName        = "sng-oauth"
	googleCloudProject = "GOOGLE_CLOUD_PROJECT"
)

func main() {
	setGinMode()

	logClient := newLogClient()
	defer logClient.Close()

	logger := logClient.Logger("user-service")

	client := sn.NewClient(newDSClient(logger), logger, cache.New(30*time.Minute, 10*time.Minute), gin.New())

	store, err := cookie.NewClient(client.Log, client.Cache).NewStore()
	if err != nil {
		logger.Panicf("unable create cookie store: %v", err)
	}

	client.Router.Use(
		sessions.Sessions(sessionName, store),
		gin.LoggerWithWriter(logger.StandardLogger(logging.Debug).Writer()),
		gin.RecoveryWithWriter(logger.StandardLogger(logging.Critical).Writer()),
	)

	// user controller
	ucon.NewClient(client.DS, logger, client.Cache, client.Router)

	// warmup
	client.Router.GET("_ah/warmup", func(c *gin.Context) { c.Status(http.StatusOK) })

	client.Router = staticRoutes(client.Router)
	client.Router.Run()
}

// staticHandler for local development since app.yaml is ignored
// static files are handled via app.yaml routes when deployed
func staticRoutes(r *gin.Engine) *gin.Engine {
	if sn.IsProduction() {
		return r
	}
	r.StaticFile("/", "dist/index.html")
	r.StaticFile("/app.js", "dist/app.js")
	r.StaticFile("/favicon.ico", "dist/favicon.ico")
	r.Static("/img", "dist/img")
	r.Static("/js", "dist/js")
	r.Static("/css", "dist/css")
	return r
}

func setGinMode() {
	if sn.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		return
	}
	gin.SetMode(gin.DebugMode)
}

func getProjectID() string {
	return os.Getenv(googleCloudProject)
}

func newDSClient(log *log.Logger) *datastore.Client {
	client, err := datastore.NewClient(context.Background(), "")
	if err != nil {
		log.Panicf("unable to create datastore client: %v", err)
	}
	return client
}

func newLogClient() *log.Client {
	client, err := log.NewClient(getProjectID())
	if err != nil {
		log.Panicf("unable to create logging client: %v", err)
	}
	return client
}
