package client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

func getUID(ctx *gin.Context, param string) (sn.UID, error) {
	id, err := strconv.ParseInt(ctx.Param(param), 10, 64)
	return sn.UID(id), err
}

func (cl *Client) userJSONHandler(uidParam string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sn.Debugf(ctx, msgEnter)
		defer sn.Debugf(ctx, msgExit)

		cu, err := cl.RequireLogin(ctx)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		uid, err := getUID(ctx, uidParam)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		if cu.ID == uid {
			ctx.JSON(http.StatusOK, gin.H{"User": cu})
			return
		}

		u, err := cl.getUser(ctx, uid)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"User": u})
	}
}

func (cl *Client) newUserHandler(ctx *gin.Context) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	cu, err := cl.RequireLogin(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	u, err := cl.getNewUser(ctx)
	if err != nil {
		sn.Errorf(ctx, "%v", err.Error())
		sn.JErr(ctx, err)
		return
	}

	u.EmailReminders = true
	u.EmailNotifications = true
	u.GravType = "monsterid"
	hash, err := emailHash(u.Email)
	if err != nil {
		sn.Warnf(ctx, "email hash error: %v", err)
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	u.EmailHash = hash

	if !cu.Admin {
		cu = u
	}

	ctx.JSON(http.StatusOK, gin.H{
		"CU":      cu,
		"User":    u,
		"Message": fmt.Sprintf("user created for %s", u.Name),
	})
}

func (cl *Client) createUserHandler(ctx *gin.Context) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	cu, u, err := cl.createUser(ctx)
	if err != nil {
		sn.Errorf(ctx, "%v", err.Error())
		sn.JErr(ctx, fmt.Errorf("cannot create user: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"CU":      cu,
		"User":    u,
		"Message": "account created for " + u.Name,
	})
}

func (cl *Client) updateUserHandler(uidParam string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sn.Debugf(ctx, msgEnter)
		defer sn.Debugf(ctx, msgExit)

		cu, err := cl.RequireLogin(ctx)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		uid, err := getUID(ctx, uidParam)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		u := cu
		if cu.ID != uid {
			if _, err := cl.RequireAdmin(ctx); err != nil {
				sn.JErr(ctx, err)
				return
			}

			u, err = cl.getUser(ctx, uid)
			if err != nil {
				sn.JErr(ctx, err)
				return
			}
		}

		obj := new(struct {
			User *sn.User
		})
		err = ctx.ShouldBind(obj)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		u, changed, err := cl.updateUser(ctx, cu, u, obj.User)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		if !changed {
			ctx.JSON(http.StatusOK, gin.H{"Message": "no change to user"})
			return
		}

		u.UpdatedAt = time.Now()
		_, err = cl.DS.Put(ctx, newUserKey(u.ID), u)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		token := cl.GetSessionToken(ctx)
		cl.SetSessionToken(ctx, u, token.Sub)

		if err := cl.SaveSession(ctx); err != nil {
			sn.JErr(ctx, err)
			return
		}
		cl.Cache.SetDefault(newUserKey(u.ID).Encode(), u)

		if cu.ID == u.ID {
			ctx.JSON(http.StatusOK, gin.H{"CU": u})
			return
		}
		ctx.JSON(http.StatusOK, nil)
	}
}
