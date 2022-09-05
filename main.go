package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	snc "github.com/SlothNinja/client"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/sn/v2"
	ucon "github.com/SlothNinja/user-controller/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

const (
	hashKeyLength    = 64
	blockKeyLength   = 32
	sessionName      = "sng-oauth"
	UserProjectIDEnv = "USER_PROJECT_ID"
	UserDSURLEnv     = "USER_DS_URL"
	UserHostURLEnv   = "USER_HOST_URL"
)

func main() {
	ctx := context.Background()

	if sn.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		cl := newClient(ctx)
		defer cl.Close()
		cl.Router.Run()
	} else {
		gin.SetMode(gin.DebugMode)
		cl := newClient(ctx)
		defer cl.Close()
		cl.Router.RunTLS(getPort(), "cert.pem", "key.pem")
	}
}

type client struct {
	*snc.Client
	logClient *log.Client
}

func newClient(ctx context.Context) *client {
	logClient := newLogClient()
	cl := &client{
		logClient: logClient,
		Client: snc.NewClient(ctx, snc.Options{
			ProjectID: getUserProjectID(),
			DSURL:     getUserDSURL(),
			Logger:    logClient.Logger("user-service"),
			Cache:     cache.New(30*time.Minute, 10*time.Minute),
			Router:    gin.New(),
		}),
	}

	store, err := sn.NewCookieClient(cl.Client).NewStore(ctx)
	if err != nil {
		cl.Log.Panicf("unable create cookie store: %v", err)
	}

	cl.Router.Use(
		sessions.Sessions(sessionName, store),
		gin.LoggerWithWriter(cl.Log.StandardLogger(logging.Debug).Writer()),
		gin.RecoveryWithWriter(cl.Log.StandardLogger(logging.Critical).Writer()),
	)

	// User controller
	ucon.NewClient(cl.Client)

	// warmup
	cl.Router.GET("_ah/warmup", func(c *gin.Context) { c.Status(http.StatusOK) })

	return cl.addRoutes()
}

type CloseErrors struct {
	Client    error
	LogClient error
}

func (ce CloseErrors) Error() string {
	return fmt.Sprintf("error closing clients: client: %q logClient: %q", ce.Client, ce.LogClient)
}

func (cl *client) Close() error {
	var ce CloseErrors

	ce.Client = cl.Client.Close()
	ce.LogClient = cl.logClient.Close()

	if ce.Client != nil || ce.LogClient != nil {
		return ce
	}
	return nil
}

// staticHandler for local development since app.yaml is ignored
// static files are handled via app.yaml routes when deployed
func (cl *client) addRoutes() *client {
	if sn.IsProduction() {
		return cl
	}
	cl.Router.StaticFile("/", "dist/index.html")
	cl.Router.StaticFile("/app.js", "dist/app.js")
	cl.Router.StaticFile("/favicon.ico", "dist/favicon.ico")
	cl.Router.Static("/img", "dist/img")
	cl.Router.Static("/js", "dist/js")
	cl.Router.Static("/css", "dist/css")
	return cl
}

func newDSClient(log *log.Logger) *datastore.Client {
	client, err := datastore.NewClient(context.Background(), "")
	if err != nil {
		log.Panicf("unable to create datastore client: %v", err)
	}
	return client
}

func newLogClient() *log.Client {
	client, err := log.NewClient(getUserProjectID())
	if err != nil {
		log.Panicf("unable to create logging client: %v", err)
	}
	return client
}

func getPort() string {
	return ":" + os.Getenv("PORT")
}

func getUserProjectID() string {
	return os.Getenv(UserProjectIDEnv)
}

func getUserDSURL() string {
	return os.Getenv(UserDSURLEnv)
}

func getUserHostURL() string {
	return os.Getenv(UserHostURLEnv)
}
