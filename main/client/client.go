// Package client provides a user service client
package client

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
)

const msgEnter = "Entering"
const msgExit = "Exiting"

// Client represents user service client
type Client struct {
	*sn.Client
	DS *datastore.Client
}

// New returns a user service client
func New(ctx context.Context, opts ...sn.Option) *Client {
	cl := &Client{Client: sn.NewClient(ctx, opts...)}
	return cl.initUserDatastore(ctx).addRoutes()
}

func (cl *Client) initUserDatastore(ctx context.Context) *Client {
	var err error
	cl.DS, err = datastore.NewClient(ctx, cl.GetProjectID())
	if err != nil {
		log.Panicf("unable to connect to user database: %v", err)
	}
	return cl
}

// AddRoutes addes routing for game.
func (cl *Client) addRoutes() *Client {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// New
	cl.Router.GET(cl.GetPrefix()+"/user/new", cl.newUserHandler)

	// Create
	cl.Router.PUT(cl.GetPrefix()+"/user/new", cl.createUserHandler)

	// Update
	cl.Router.PUT(cl.GetPrefix()+"/user/:uid/update", cl.updateUserHandler("uid"))
	cl.Router.PUT(cl.GetPrefix()+"/user/:uid/as", cl.as)

	// Get
	cl.Router.GET(cl.GetPrefix()+"/user/:uid/json", cl.userJSONHandler("uid"))

	cl.Router.GET(cl.GetPrefix()+"/user/login", cl.login(cl.GetPrefix()+"/user/auth"))

	cl.Router.GET(cl.GetPrefix()+"/user/logout", cl.logout)

	cl.Router.GET(cl.GetPrefix()+"/user/auth", cl.auth(cl.GetPrefix()+"/user/auth"))

	return cl
}

// Close closes the user service client
func (cl *Client) Close() error {
	err1 := cl.Client.Close()
	err2 := cl.DS.Close()
	return errors.Join(err1, err2)
}
