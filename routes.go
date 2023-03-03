package main

import (
	"github.com/SlothNinja/sn/v2"
	"github.com/gin-gonic/gin"
)

type Client struct {
	*sn.Client
	User *sn.UserClient
}

func NewClient(snClient *sn.Client) *Client {
	snClient.Log.Debugf(msgEnter)
	defer snClient.Log.Debugf(msgExit)

	cl := &Client{
		Client: snClient,
		User:   sn.NewUserClient(snClient),
	}
	return cl.addRoutes()
}

func userFrom(c *gin.Context) (*sn.User, error) {
	return sn.From(c), nil
}

func (cl *Client) addRoutes() *Client {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	// New
	cl.Router.GET("new", cl.NewAction)

	// Create
	cl.Router.PUT("new", cl.Create)

	// Update
	cl.Router.PUT("update/:uid", cl.Update("uid"))

	// Get
	cl.Router.GET("json/:uid", cl.JSON("uid"))

	// // User Games
	// cl.Router.POST("show/:uid/games/json",
	// 	cl.Game.GetFiltered(gtype.All),
	// 	cl.Game.JSONIndexAction,
	// )

	// cl.Router.POST("edit/:uid/games/json",
	// 	// user.RequireLogin(),
	// 	cl.Game.GetFiltered(gtype.All),
	// 	cl.Game.JSONIndexAction,
	// )

	cl.Router.GET("as/:uid", cl.User.As)

	cl.Router.GET("current", cl.Current)

	cl.Router.GET("login", sn.Login("auth"))

	cl.Router.GET("logout", sn.Logout)

	cl.Router.GET("auth", cl.User.Auth("auth"))

	// Index
	cl.Router.GET("index", cl.Index)

	return cl
}
