package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	firebase "firebase.google.com/go"
	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofrs/uuid"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// Environment variable for specifying host url of service
	HOST           = "HOST"
	UserHostURLEnv = "USER_HOST_URL"
	tokenLength    = 32
	oauthsKind     = "OAuths"
	oauthKind      = "OAuth"
	root           = "root"
	stateKey       = "state"
	redirectKey    = "redirect"
)

func getRedirectionPath(ctx *gin.Context) (string, bool) {
	return ctx.GetQuery("redirect")
}

func (cl client) login(path string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state := randToken(tokenLength)
		cl.Session(ctx).Set(stateKey, state)

		redirect, found := getRedirectionPath(ctx)
		if !found {
			redirect = base64.StdEncoding.EncodeToString([]byte(ctx.Request.Header.Get("Referer")))
		}
		cl.Session(ctx).Set(redirectKey, redirect)
		err := cl.Session(ctx).Save()
		if err != nil {
			cl.Log.Warningf("unable to save session: %v", err)
		}
		ctx.Redirect(http.StatusSeeOther, cl.getLoginURL(path, state))
	}
}

func (cl client) logout(ctx *gin.Context) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// get redirection path from session before clearing
	path, pathFound := getRedirectionPath(ctx)

	cl.Session(ctx).Clear()
	err := cl.Session(ctx).Save()
	if err != nil {
		cl.Log.Warningf("unable to save session: %v", err)
	}

	if pathFound {
		ctx.Redirect(http.StatusSeeOther, path)
		return
	}

	ctx.Redirect(http.StatusSeeOther, homePath)
}

func randToken(length int) string {
	key := securecookie.GenerateRandomKey(length)
	return base64.StdEncoding.EncodeToString(key)
}

func (cl client) getLoginURL(path, state string) string {
	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	return cl.oauth2Config(path, scopes()...).AuthCodeURL(state)
}

func (cl client) oauth2Config(path string, scopes ...string) *oauth2.Config {
	redirectURL := fmt.Sprintf("%s:%s/%s", cl.GetBackEndURL(), cl.GetBackEndPort(), strings.TrimPrefix(path, "/"))
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes:       scopes,
		RedirectURL:  redirectURL,
	}
}

func scopes() []string {
	return []string{"email", "profile", "openid"}
}

func getHost() string {
	return os.Getenv(HOST)
}

func getUserHostURL() string {
	s := os.Getenv(UserHostURLEnv)
	if s != "" {
		return s
	}
	return getHost()
}

type oaInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	LoggedIn      bool
	Admin         bool
}

const fqdn = "www.slothninja.com"

var namespaceUUID = uuid.NewV5(uuid.NamespaceDNS, fqdn)

// Generates ID for User from ID obtained from OAuth OpenID Connect
func genOAuthID(s string) string {
	return uuid.NewV5(namespaceUUID, s).String()
}

type oauth struct {
	Key       *datastore.Key `datastore:"__key__"`
	ID        sn.UID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (o *oauth) Load(ps []datastore.Property) error {
	return datastore.LoadStruct(o, ps)
}

func (o *oauth) Save() ([]datastore.Property, error) {
	t := time.Now()
	if o.CreatedAt.IsZero() {
		o.CreatedAt = t
	}
	o.UpdatedAt = t
	return datastore.SaveStruct(o)
}

func (o *oauth) LoadKey(k *datastore.Key) error {
	o.Key = k
	return nil
}

func pk() *datastore.Key {
	return datastore.NameKey(oauthsKind, root, nil)
}

func newOAuthKey(id string) *datastore.Key {
	return datastore.NameKey(oauthKind, id, pk())
}

func newOAuth(id string) oauth {
	return oauth{Key: newOAuthKey(id)}
}

func (cl client) redirectPathFrom(ctx *gin.Context) string {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	retrievedPath, ok := cl.Session(ctx).Get(redirectKey).(string)
	if !ok {
		return ""
	}

	bs, err := base64.StdEncoding.DecodeString(retrievedPath)
	if err != nil {
		return ""
	}
	return string(bs)
}

// returns whether user present in database and any error resulting from trying to create session
func (cl client) loginSessionByOAuthSub(ctx *gin.Context, sub string) (bool, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	oaid := genOAuthID(sub)
	oa, err := cl.getOAuth(ctx, oaid)
	if err != nil {
		return false, err
	}

	// Succesfully pulled uid from datastore using OAuth Sub
	u, err := cl.getUser(ctx, oa.ID)
	if err != nil {
		return false, err
	}

	// created new token and save to session
	token := cl.Session(ctx).NewToken(u.ID, sub, u.Data)
	return true, cl.Session(ctx).SaveToken(token)
}

// returns whether user present in datastore and any error resulting for trying to create session
func (cl client) loginSessionByEmailAndSub(ctx *gin.Context, email, sub string) (bool, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	u, err := cl.getByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	oa := newOAuth(genOAuthID(sub))
	oa.ID = u.ID

	_, err = cl.DS.Put(ctx, oa.Key, &oa)
	if err != nil {
		return true, err
	}

	token := cl.Session(ctx).NewToken(u.ID, sub, u.Data)
	// st, err := newSessionToken(ctx, u.ID, sub)
	// if err != nil {
	// 	return true, err
	// }

	return true, cl.Session(ctx).SaveToken(token)
}

// returns error resulting for trying to create session
func (cl client) loginSessionNewUser(ctx *gin.Context, email, sub string) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	u := newUser(0)
	u.Name = strings.Split(email, "@")[0]
	u.Email = email
	cl.Session(ctx).SetUserName(u.Name)
	cl.Session(ctx).SetUserEmail(u.Email)
	token := cl.Session(ctx).NewToken(u.ID, sub, u.Data)
	return cl.Session(ctx).SaveToken(token)
}

func (cl client) auth(authPath string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cl.Log.Debugf(msgEnter)
		defer cl.Log.Debugf(msgExit)

		uInfo, err := cl.getUInfo(ctx, authPath)
		if err != nil {
			cl.Log.Errorf(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if userExists, err := cl.loginSessionByOAuthSub(ctx, uInfo.Sub); userExists && err != nil {
			cl.Log.Errorf(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		} else if err == nil {
			ctx.Redirect(http.StatusSeeOther, cl.redirectPathFrom(ctx))
			return
		} else {
			cl.Log.Debugf(err.Error())
		}

		// OAuth sub not associated with UID in datastore
		// Check to see if other entities exist for same email address.
		// If so, use old entities for user
		if userExists, err := cl.loginSessionByEmailAndSub(ctx, uInfo.Email, uInfo.Sub); userExists && err != nil {
			cl.Log.Errorf(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		} else if err == nil {
			ctx.Redirect(http.StatusSeeOther, cl.redirectPathFrom(ctx))
			return
		} else {
			cl.Log.Debugf(err.Error())
		}

		// Create New User
		if err := cl.loginSessionNewUser(ctx, uInfo.Email, uInfo.Sub); err != nil {
			cl.Log.Errorf(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		cl.Log.Debugf("cl: %#v", cl)
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("https://%s/user/new", cl.GetUserFrontURL()))
	}
}

func getFBToken(ctx *gin.Context, uid sn.UID) (string, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgEnter)

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("error initializing app: %w", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting Auth client: %w", err)
	}

	token, err := client.CustomToken(ctx, fmt.Sprintf("%d", uid))
	if err != nil {
		return "", fmt.Errorf("error minting custom token: %w", err)
	}

	return token, err
}

func (cl client) as(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	admin, err := cl.GetAdmin(ctx)
	if err != nil {
		cl.Log.Errorf(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !admin {
		cl.Log.Errorf("must be admin")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("uid"), 10, 64)
	if err != nil {
		cl.Log.Errorf(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	u, err := cl.getUser(ctx, sn.UID(id))
	if err != nil {
		cl.Log.Errorf(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token := cl.Session(ctx).NewToken(u.ID, "", u.Data)
	err = cl.Session(ctx).SaveToken(token)
	if err != nil {
		cl.Log.Errorf(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	return
}

func (cl client) getUInfo(ctx *gin.Context, path string) (oaInfo, error) {
	// Handle the exchange code to initiate a transport.
	retrievedState := cl.Session(ctx).Get("state")
	if retrievedState != ctx.Query("state") {
		return oaInfo{}, fmt.Errorf("Invalid session state: %s", retrievedState)
	}

	conf := cl.oauth2Config(path, scopes()...)
	tok, err := conf.Exchange(ctx, ctx.Query("code"))
	if err != nil {
		return oaInfo{}, fmt.Errorf("tok error: %s", err.Error())
	}

	client := conf.Client(ctx, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return oaInfo{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return oaInfo{}, err
	}

	uInfo := oaInfo{}
	var b binding.BindingBody = binding.JSON
	err = b.BindBody(body, &uInfo)
	if err != nil {
		return oaInfo{}, err
	}
	return uInfo, nil
}

func (cl client) getOAuth(ctx *gin.Context, id string) (oauth, error) {
	return cl.getOAuthByKey(ctx, newOAuthKey(id))
}

func (cl client) getOAuthByKey(ctx *gin.Context, k *datastore.Key) (oauth, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	oauth, found := cl.getCachedOAuth(k)
	if found {
		return oauth, nil
	}

	oauth = newOAuth(k.Name)
	err := cl.DS.Get(ctx, k, &oauth)
	if err != nil {
		return oauth, err
	}
	cl.cacheOAuth(oauth)
	return oauth, nil
}

func (cl client) getCachedOAuth(k *datastore.Key) (oauth, bool) {
	auth := newOAuth(k.Name)
	if k == nil {
		return auth, false
	}

	data, found := cl.Cache.Get(k.Encode())
	if !found {
		return auth, false
	}

	auth, ok := data.(oauth)
	if !ok {
		return auth, false
	}
	return auth, true
}

func (cl client) cacheOAuth(auth oauth) {
	if auth.Key == nil {
		return
	}
	cl.Cache.SetDefault(auth.Key.Encode(), auth)
}

func (cl client) getByEmail(ctx *gin.Context, email string) (sn.User, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	email = strings.ToLower(strings.TrimSpace(email))
	q := datastore.NewQuery(uKind).
		Ancestor(userRootKey()).
		Filter("Email=", email).
		KeysOnly()

	ks, err := cl.DS.GetAll(ctx, q, nil)
	if err != nil {
		return sn.User{}, err
	}

	for i := range ks {
		if ks[i].ID != 0 {
			return cl.getUser(ctx, sn.UID(ks[i].ID))
		}
	}
	return sn.User{}, errors.New("unable to find user")
}
