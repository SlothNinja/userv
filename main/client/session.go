package client

import (
	"errors"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

const (
	emailKey    = "email"
	nameKey     = "name"
	stateKey    = "state"
	redirectKey = "redirect"
	sessionName = "sng-oauth"
	sessionKey  = "session"
)

func (cl *Client) setSessionState(ctx *gin.Context, state string) {
	cl.Session(ctx).Set(stateKey, state)
}

func (cl *Client) getSessionState(ctx *gin.Context) string {
	state, ok := cl.Session(ctx).Get(stateKey).(string)
	if !ok {
		return ""
	}
	return state
}

func (cl *Client) getNewUser(ctx *gin.Context) (*sn.User, error) {
	token := cl.GetSessionToken(ctx)
	if token == nil {
		return nil, sn.ErrNotLoggedIn
	}

	if token.ID != 0 {
		return nil, errors.New("user present, no need for new one")
	}

	u := token.ToUser()
	u.Name = cl.getSessionUserName(ctx)
	if u.Name == "" {
		return nil, errors.New("session missing user name")
	}

	u.Email = cl.getSessionUserEmail(ctx)
	if u.Email == "" {
		return nil, errors.New("session missing user email")
	}

	return u, nil
}

func (cl *Client) setSessionUserEmail(ctx *gin.Context, email string) {
	cl.Session(ctx).Set(emailKey, email)
}

func (cl *Client) getSessionUserEmail(ctx *gin.Context) string {
	email, ok := cl.Session(ctx).Get(emailKey).(string)
	if !ok {
		return ""
	}
	return email
}

func (cl *Client) setSessionUserName(ctx *gin.Context, name string) {
	cl.Session(ctx).Set(nameKey, name)
}

func (cl *Client) getSessionUserName(ctx *gin.Context) string {
	name, ok := cl.Session(ctx).Get(nameKey).(string)
	if !ok {
		return ""
	}
	return name
}

func (cl *Client) setSessionRedirect(ctx *gin.Context, redirect string) {
	cl.Session(ctx).Set(redirectKey, redirect)
}

func (cl *Client) getSessionRedirect(ctx *gin.Context) string {
	redirect, ok := cl.Session(ctx).Get(redirectKey).(string)
	if !ok {
		return ""
	}
	return redirect
}
