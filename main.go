package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/sn"
	ucon "github.com/SlothNinja/user-controller"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/patrickmn/go-cache"
)

const (
	hashKeyLength        = 64
	blockKeyLength       = 32
	sessionName          = "sng-oauth"
	GOOGLE_CLOUD_PROJECT = "GOOGLE_CLOUD_PROJECT"
)

func main() {
	setGinMode()

	logClient := newLogClient()
	defer logClient.Close()

	logger := logClient.Logger("user-service")

	client := sn.NewClient(newDSClient(), logger, cache.New(30*time.Minute, 10*time.Minute), gin.New())

	s, err := getSecrets()
	if err != nil {
		panic(err.Error())
	}

	client.Router.Use(
		sessions.Sessions(sessionName, createCookieStore(s)),
		gin.LoggerWithWriter(logger.StandardLogger(logging.Debug).Writer()),
		gin.RecoveryWithWriter(logger.StandardLogger(logging.Critical).Writer()),
	)

	// user controller
	ucon.NewClient(client.DS, logger, client.Cache, client.Router)

	// cookie route
	client.Router.GET("cookie", cookieHandler(s))

	// warmup
	client.Router.GET("_ah/warmup", func(c *gin.Context) { c.Status(http.StatusOK) })

	client.Router = staticRoutes(client.Router)
	client.Router.Run()
}

// type client struct {
// 	dsClient  *datastore.Client
// 	logClient *log.Client
// 	mcache    *cache.Cache
// 	router    *gin.Engine
// }
//
// func newClient() *client {
// 	return &client{
// 		dsClient:  newDSClient(),
// 		logClient: newLogClient(),
// 		mcache:    cache.New(30*time.Minute, 10*time.Minute),
// 		router:    gin.Default(),
// 	}
// }

type secrets struct {
	HashKey   []byte         `json:"hashKey"`
	BlockKey  []byte         `json:"blockKey"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Key       *datastore.Key `datastore:"__key__" json:"-"`
}

func getSecrets() (secrets, error) {
	s := secrets{
		Key: secretsKey(),
	}

	c := context.Background()
	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return s, err
	}

	err = dsClient.Get(c, s.Key, &s)
	if err == nil {
		return s, nil
	}

	if err != datastore.ErrNoSuchEntity {
		return s, err
	}

	log.Warningf("generated new secrets")
	s, err = genSecrets()
	if err != nil {
		return s, err
	}

	_, err = dsClient.Put(c, s.Key, &s)
	return s, err
}

func secretsKey() *datastore.Key {
	return datastore.NameKey("Secrets", "root", nil)
}

func genSecrets() (secrets, error) {
	s := secrets{
		HashKey:  securecookie.GenerateRandomKey(hashKeyLength),
		BlockKey: securecookie.GenerateRandomKey(blockKeyLength),
		Key:      secretsKey(),
	}

	if s.HashKey == nil {
		return s, fmt.Errorf("generated hashKey was nil")
	}

	if s.BlockKey == nil {
		return s, fmt.Errorf("generated blockKey was nil")
	}

	return s, nil
}

func (s *secrets) Load(ps []datastore.Property) error {
	return datastore.LoadStruct(s, ps)
}

func (s *secrets) Save() ([]datastore.Property, error) {
	s.UpdatedAt = time.Now()
	return datastore.SaveStruct(s)
}

func (s *secrets) LoadKey(k *datastore.Key) error {
	s.Key = k
	return nil
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

func createCookieStore(s secrets) cookie.Store {
	if !sn.IsProduction() {
		log.Debugf("hashKey: %s\nblockKey: %s",
			base64.StdEncoding.EncodeToString(s.HashKey),
			base64.StdEncoding.EncodeToString(s.BlockKey),
		)
		store := cookie.NewStore(s.HashKey, s.BlockKey)
		opts := sessions.Options{
			Domain: "fake-slothninja.com",
			Path:   "/",
		}
		store.Options(opts)
		return store
	}
	store := cookie.NewStore(s.HashKey, s.BlockKey)
	opts := sessions.Options{
		Domain: "slothninja.com",
		Path:   "/",
		Secure: true,
	}
	store.Options(opts)
	return store
}

func cookieHandler(s secrets) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, s)
	}
}

func setGinMode() {
	if sn.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		return
	}
	gin.SetMode(gin.DebugMode)
}

func getProjectID() string {
	return os.Getenv(GOOGLE_CLOUD_PROJECT)
}

func newDSClient() *datastore.Client {
	client, err := datastore.NewClient(context.Background(), "")
	if err != nil {
		log.Panicf("unable to create datastore client: %w", err)
	}
	return client
}

func newLogClient() *log.Client {
	client, err := log.NewClient(getProjectID())
	if err != nil {
		log.Panicf("unable to create logging client: %w", err)
	}
	return client
}
