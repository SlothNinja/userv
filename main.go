package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	"github.com/SlothNinja/sn/v2"
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
		cl.Router.TrustedPlatform = gin.PlatformGoogleAppEngine
		cl.Router.Run()
	} else {
		gin.SetMode(gin.DebugMode)
		cl := newClient(ctx)
		defer cl.Close()
		cl.Router.SetTrustedProxies(nil)
		cl.Router.RunTLS(getPort(), "cert.pem", "key.pem")
	}
}

type client struct {
	*sn.Client
}

func newClient(ctx context.Context) *client {
	logClient := newLogClient()
	cl := &client{
		Client: sn.NewClient(ctx, sn.Options{
			ProjectID: getUserProjectID(),
			DSURL:     getUserDSURL(),
			Logger:    logClient.Logger("user-service"),
			Cache:     cache.New(30*time.Minute, 10*time.Minute),
			Router:    gin.Default(),
		}),
	}

	store, err := sn.NewCookieClient(cl.Client).NewStore(ctx)
	if sn.IsProduction() {
		opts := sessions.Options{
			Domain: "slothninja.com",
			Path:   "/",
			MaxAge: 60 * 60 * 24 * 7, // 7 days in seconds
			Secure: true,
		}
		store.Options(opts)
	}
	if err != nil {
		cl.Log.Panicf("unable create cookie store: %v", err)
	}

	cl.Router.Use(
		sessions.Sessions(sessionName, store),
		gin.LoggerWithWriter(cl.Log.StandardLogger(logging.Debug).Writer()),
		gin.RecoveryWithWriter(cl.Log.StandardLogger(logging.Critical).Writer()),
	)

	// User controller
	NewClient(cl.Client)

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

func newDSClient(log *sn.Logger) *datastore.Client {
	client, err := datastore.NewClient(context.Background(), "")
	if err != nil {
		sn.Panicf("unable to create datastore client: %v", err)
	}
	return client
}

func newLogClient() *sn.LogClient {
	client, err := sn.NewLogClient(getUserProjectID())
	if err != nil {
		sn.Panicf("unable to create logging client: %v", err)
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
