package client

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (cl *Client) getUser(ctx *gin.Context, uid sn.UID) (*sn.User, error) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	return cl.get(ctx, uid)
}

func (cl *Client) get(ctx *gin.Context, uid sn.UID) (*sn.User, error) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	u, err := cl.mcGet(ctx, uid)
	if err == nil {
		return u, nil
	}

	return cl.dsGet(ctx, uid)
}

func (cl *Client) mcGet(ctx context.Context, uid sn.UID) (*sn.User, error) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	if uid == 0 {
		return nil, ErrMissingUID
	}

	item, found := cl.Cache.Get(newUserKey(uid).Encode())
	if !found {
		return nil, ErrUserNotFound
	}

	u, ok := item.(*sn.User)
	if !ok {
		return nil, ErrInvalidCache
	}
	return u, nil
}

func (cl *Client) mcGetMulti(ctx context.Context, uids []sn.UID) ([]*sn.User, error) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	l := len(uids)
	if l == 0 {
		return nil, ErrMissingUID
	}

	me := make(datastore.MultiError, l)
	us := make([]*sn.User, l)
	isNil := true
	for i, k := range uids {
		us[i], me[i] = cl.mcGet(ctx, k)
		if me[i] != nil {
			isNil = false
		}
	}

	if isNil {
		return us, nil
	}
	return us, me
}

func (cl *Client) dsGet(ctx *gin.Context, uid sn.UID) (*sn.User, error) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	if uid == 0 {
		return nil, ErrMissingUID
	}

	u := new(sn.User)
	err := cl.DS.Get(ctx, newUserKey(uid), u)
	if err != nil {
		sn.Warnf(ctx, "%v", err.Error())
		return nil, err
	}
	u.ID = uid
	cl.cacheUser(ctx, u)
	return u, nil
}

func (cl *Client) dsGetMulti(ctx *gin.Context, uids []sn.UID) ([]*sn.User, error) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	l := len(uids)
	if l == 0 {
		return nil, ErrMissingUID
	}

	us := make([]*sn.User, l)
	ks := pie.Map(uids, func(uid sn.UID) *datastore.Key { return newUserKey(uid) })
	err := cl.DS.GetMulti(ctx, ks, us)
	if err != nil {
		return us, err
	}
	for _, u := range us {
		cl.cacheUser(ctx, u)
	}
	return us, nil
}

func (cl *Client) cacheUser(ctx context.Context, u *sn.User) {
	sn.Debugf(ctx, msgEnter)
	defer sn.Debugf(ctx, msgExit)

	cl.Cache.SetDefault(newUserKey(u.ID).Encode(), u)
}
