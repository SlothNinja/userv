package client

import (
	"crypto/md5"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
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

func newUser(uid sn.UID) *sn.User {
	return &sn.User{ID: uid}
}

func (cl *Client) updateUser(ctx *gin.Context, cu, u1, u2 *sn.User) (*sn.User, bool, error) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	changed := false
	// If admin or newly created user
	if cu.Admin || (cu.ID == 0 && u1.ID == 0) {
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

func (cl *Client) updateUserName(ctx *gin.Context, u *sn.User, n string) (*sn.User, bool, error) {
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

	q := datastore.NewQuery("User").FilterField("LCName", "=", LCName)

	cnt, err := cl.DS.Count(ctx, q)
	if err != nil {
		return false, err
	}
	return cnt == 0, nil
}

// returns current user, created user, and error
func (cl *Client) createUser(ctx *gin.Context) (*sn.User, *sn.User, error) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	cu, err := cl.RequireLogin(ctx)
	if err == nil && cu.ID != 0 {
		sn.Warnf(ctx, "%s(%d) already has an account", cu.Name, cu.ID)
		return nil, nil, err
	}

	token := cl.GetSessionToken(ctx)
	if token == nil {
		return nil, nil, sn.ErrNotLoggedIn
	}

	if token.ID != 0 {
		return nil, nil, errors.New("user present, no need for new one")
	}

	obj := new(struct {
		User *sn.User
	})
	err = ctx.ShouldBind(obj)
	if err != nil {
		return nil, nil, err
	}

	sn.Debugf(ctx, "obj.User: %#v", obj.User)

	u := newUser(0)
	u, _, err = cl.updateUser(ctx, u, u, obj.User)
	if err != nil {
		return nil, nil, err
	}

	ks, err := cl.DS.AllocateIDs(ctx, []*datastore.Key{newUserKey(u.ID)})
	if err != nil {
		return nil, nil, err
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
		es := []any{oa, u}
		_, err := tx.PutMulti(ks, es)
		return err
	})
	if err != nil {
		return nil, nil, err
	}

	cl.SetSessionToken(ctx, u, token.Sub)
	if err := cl.SaveSession(ctx); err != nil {
		return nil, nil, err
	}
	return u, u, nil
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
