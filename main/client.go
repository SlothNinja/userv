// Main provides a user service
package main

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const msgEnter = "Entering"
const msgExit = "Exiting"

type client struct {
	sn.Client
	DS *datastore.Client
}

func newClient(ctx context.Context, opts ...sn.Option) client {
	cl := client{Client: sn.NewClient(ctx, opts...)}
	return cl.initUserDatastore(ctx).addRoutes()
}

func (cl client) initUserDatastore(ctx context.Context) client {
	if sn.IsProduction() {
		return cl.getUserDatastoreForProduction(ctx)
	}
	return cl.getUserDataStoreForDevelopment(ctx)
}

func (cl client) getUserDatastoreForProduction(ctx context.Context) client {
	dsClient, err := datastore.NewClient(ctx, cl.GetUserProjectID())
	if err != nil {
		panic(fmt.Errorf("unable to connect to user database: %w", err))
	}
	cl.DS = dsClient
	return cl
}

func (cl client) getUserDataStoreForDevelopment(ctx context.Context) client {
	dsClient, err := datastore.NewClient(
		ctx,
		cl.GetUserProjectID(),
		option.WithEndpoint(cl.GetUserDSURL()),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		option.WithGRPCConnectionPool(50),
	)
	if err != nil {
		panic(fmt.Errorf("unable to connect to user database: %w", err))
	}
	cl.DS = dsClient
	return cl
}

// AddRoutes addes routing for game.
func (cl client) addRoutes() client {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	// New
	cl.Router.GET(cl.GetPrefix()+"/user/new", cl.newUserHandler)

	// Create
	cl.Router.PUT(cl.GetPrefix()+"/user/new", cl.createUserHandler)

	// Update
	cl.Router.PUT(cl.GetPrefix()+"/user/:uid/update", cl.updateUserHandler("uid"))

	// Get
	cl.Router.GET(cl.GetPrefix()+"/user/:uid/json", cl.userJSONHandler("uid"))

	cl.Router.GET(cl.GetPrefix()+"/user/:uid/as", cl.as)

	cl.Router.GET(cl.GetPrefix()+"/user/login", cl.login(cl.GetPrefix()+"/user/auth"))

	cl.Router.GET(cl.GetPrefix()+"/user/logout", cl.logout)

	cl.Router.GET(cl.GetPrefix()+"/user/auth", cl.auth(cl.GetPrefix()+"/user/auth"))

	return cl
}

func (cl client) close() error {
	err1 := cl.Client.Close()
	err2 := cl.DS.Close()
	return errors.Join(err1, err2)
}
