package client

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

var (
	// ErrMissingUID is an error that is returned when a user ID (sn.UID) is required, but missing
	ErrMissingUID = errors.New("missing user ID")
	// ErrUserNotFound is an error that is returned when a user for supplied data cannot be found
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCache is an error that is returned when a value retrieved from cache is invalid
	ErrInvalidCache = errors.New("invalid cache value")
)

const uKind = "User"

func userRootKey() *datastore.Key {
	return datastore.NameKey("Users", "root", nil)
}

func newUserKey(uid sn.UID) *datastore.Key {
	return datastore.IDKey(uKind, int64(uid), userRootKey())
}

func newUser(uid sn.UID) sn.User {
	return sn.User{ID: uid}
}

func (cl *Client) updateUser(ctx *gin.Context, cu, u1, u2 sn.User) (sn.User, bool, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	changed := false
	if cu.Admin {
		if u2.Email != "" && u2.Email != u1.Email {
			hash, err := emailHash(u1.Email)
			if err != nil {
				return u1, changed, err
			}

			u1.Email = u2.Email
			u1.EmailHash = hash
			changed = true
		}

		u1, nameChanged, err := cl.updateUserName(ctx, u1, u2.Name)
		if err != nil {
			return u1, false, err
		}
		changed = changed || nameChanged
	}

	if !cu.Admin && (cu.ID != u1.ID) {
		return u1, changed, nil
	}

	if u1.EmailReminders != u2.EmailReminders {
		u1.EmailReminders = u2.EmailReminders
		changed = true
	}
	if u1.EmailNotifications != u2.EmailNotifications {
		u1.EmailNotifications = u2.EmailNotifications
		changed = true
	}
	if u1.GravType != u2.GravType {
		u1.GravType = u2.GravType
		changed = true
	}
	return u1, changed, nil
}

func (cl *Client) updateUserName(ctx *gin.Context, u sn.User, n string) (sn.User, bool, error) {
	matcher := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9._%+\-]+$`)

	switch {
	case n == u.Name:
		return u, false, nil
	case len(n) > 15:
		return u, false, fmt.Errorf("%q is too long: %w", n, sn.ErrValidation)
	case !matcher.MatchString(n):
		return u, false, fmt.Errorf("%q is not a valid user name: %w", n, sn.ErrValidation)
	default:
		uniq, err := cl.nameIsUnique(ctx, n)
		if err != nil {
			return u, false, err
		}
		if !uniq {
			return u, false, fmt.Errorf("%q is not a unique user name: %w", n, sn.ErrValidation)
		}
		u.Name = n
		u.LCName = strings.ToLower(n)
		return u, true, nil
	}
}

func (cl *Client) nameIsUnique(ctx *gin.Context, name string) (bool, error) {
	LCName := strings.ToLower(name)

	q := datastore.NewQuery("User").Filter("LCName=", LCName)

	cnt, err := cl.DS.Count(ctx, q)
	if err != nil {
		return false, err
	}
	return cnt == 0, nil
}

func getUID(ctx *gin.Context, param string) (sn.UID, error) {
	id, err := strconv.ParseInt(ctx.Param(param), 10, 64)
	return sn.UID(id), err
}

func (cl *Client) userJSONHandler(uidParam string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cl.Log.Debugf(msgEnter)
		defer cl.Log.Debugf(msgExit)

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
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.RequireLogin(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	u, err := cl.Session(ctx).GetNewUser()
	if err != nil {
		cl.Log.Errorf(err.Error())
		sn.JErr(ctx, err)
		return
	}

	u.EmailReminders = true
	u.EmailNotifications = true
	u.GravType = "monsterid"
	hash, err := emailHash(u.Email)
	if err != nil {
		cl.Log.Warningf("email hash error: %v", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
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
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, u, err := cl.createUser(ctx)
	if err != nil {
		cl.Log.Warningf("cannot create user: %w", err)
		ctx.Redirect(http.StatusSeeOther, cl.GetHome())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"CU":      cu,
		"User":    u,
		"Message": "account created for " + u.Name,
	})
}

// returns current user, created user, and error
func (cl *Client) createUser(ctx *gin.Context) (sn.User, sn.User, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.RequireLogin(ctx)
	if err == nil && cu.ID != 0 {
		cl.Log.Warningf("%s(%d) already has an account", cu.Name, cu.ID)
		return sn.User{}, sn.User{}, err
	}

	token, ok := cl.Session(ctx).GetUserToken()
	if !ok {
		return sn.User{}, sn.User{}, errors.New("missing token")
	}

	if token.ID != 0 {
		// ctx.Redirect(http.StatusSeeOther, homePath)
		return sn.User{}, sn.User{}, errors.New("user present, no need for new one")
	}

	u := newUser(0)
	err = ctx.ShouldBind(&u)
	if err != nil {
		return sn.User{}, sn.User{}, err
	}

	u, _, err = cl.updateUser(ctx, u, u, u)
	if err != nil {
		return sn.User{}, sn.User{}, err
	}

	ks, err := cl.DS.AllocateIDs(ctx, []*datastore.Key{newUserKey(u.ID)})
	if err != nil {
		return sn.User{}, sn.User{}, err
	}

	u.ID = sn.UID(ks[0].ID)
	u.LCName = strings.ToLower(u.Name)

	t := time.Now()
	oaid := genOAuthID(token.Sub)
	oa := newOAuth(oaid)
	oa.ID = u.ID

	oa.UpdatedAt = t
	oa.CreatedAt = t
	u.UpdatedAt = t
	u.CreatedAt = t
	u.Joined = t

	_, err = cl.DS.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		ks := []*datastore.Key{oa.Key, newUserKey(u.ID)}
		es := []interface{}{&oa, &u}
		_, err := tx.PutMulti(ks, es)
		return err

	})

	if err != nil {
		return sn.User{}, sn.User{}, err
	}

	token.Data = u.Data
	if token.ID == 0 {
		token.ID = u.ID
	}
	err = cl.Session(ctx).SaveToken(token)
	if err != nil {
		return sn.User{}, sn.User{}, err
	}
	return u, u, nil

}

func (cl *Client) updateUserHandler(uidParam string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cl.Log.Debugf(msgEnter)
		defer cl.Log.Debugf(msgExit)

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

		obj := newUser(0)
		err = ctx.ShouldBind(&obj)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		cl.Log.Debugf("obj: %#v", obj)
		u, changed, err := cl.updateUser(ctx, cu, u, obj)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		if !changed {
			ctx.JSON(http.StatusOK, gin.H{"Message": "no change to user"})
			return
		}

		u.UpdatedAt = time.Now()
		_, err = cl.DS.Put(ctx, newUserKey(u.ID), &u)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}

		token, _ := cl.Session(ctx).GetUserToken()
		token.ID = u.ID
		token.Data = u.Data

		err = cl.Session(ctx).SaveToken(token)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}
		cl.Cache.SetDefault(newUserKey(u.ID).Encode(), u)

		if cu.ID == u.ID {
			ctx.JSON(http.StatusOK, gin.H{"CU": u, "Message": "user updated"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"Message": "user updated"})
	}
}

func (cl *Client) getUser(ctx *gin.Context, uid sn.UID) (sn.User, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	return cl.get(ctx, uid)
}

func (cl *Client) get(ctx *gin.Context, uid sn.UID) (sn.User, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	u, err := cl.mcGet(uid)
	if err == nil {
		return u, nil
	}

	return cl.dsGet(ctx, uid)
}

func (cl *Client) mcGet(uid sn.UID) (sn.User, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	if uid == 0 {
		return sn.User{}, ErrMissingUID
	}

	item, found := cl.Cache.Get(newUserKey(uid).Encode())
	if !found {
		return sn.User{}, ErrUserNotFound
	}

	u, ok := item.(sn.User)
	if !ok {
		return sn.User{}, ErrInvalidCache
	}
	return u, nil
}

func (cl *Client) mcGetMulti(uids []sn.UID) ([]sn.User, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	l := len(uids)
	if l == 0 {
		return nil, ErrMissingUID
	}

	me := make(datastore.MultiError, l)
	us := make([]sn.User, l)
	isNil := true
	for i, k := range uids {
		us[i], me[i] = cl.mcGet(k)
		if me[i] != nil {
			isNil = false
		}
	}

	if isNil {
		return us, nil
	}
	return us, me
}

func (cl *Client) dsGet(ctx *gin.Context, uid sn.UID) (sn.User, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	if uid == 0 {
		return sn.User{}, ErrMissingUID
	}

	var u sn.User
	err := cl.DS.Get(ctx, newUserKey(uid), &u)
	if err != nil {
		return sn.User{}, err
	}
	u.ID = uid
	cl.cacheUser(u)
	return u, nil
}

func (cl *Client) dsGetMulti(ctx *gin.Context, uids []sn.UID) ([]sn.User, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	l := len(uids)
	if l == 0 {
		return nil, ErrMissingUID
	}

	us := make([]sn.User, l)
	ks := pie.Map(uids, func(uid sn.UID) *datastore.Key { return newUserKey(uid) })
	err := cl.DS.GetMulti(ctx, ks, us)
	if err != nil {
		return us, err
	}
	for _, u := range us {
		cl.cacheUser(u)
	}
	return us, nil
}

func (cl *Client) cacheUser(u sn.User) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cl.Cache.SetDefault(newUserKey(u.ID).Encode(), u)
}

func emailHash(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	hash := md5.New()
	_, err := hash.Write([]byte(email))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
