module github.com/SlothNinja/userv

go 1.14

require (
	cloud.google.com/go/datastore v1.2.0
	github.com/SlothNinja/log v0.0.2
	github.com/SlothNinja/restful v1.0.0
	github.com/SlothNinja/sn v1.0.0
	github.com/SlothNinja/sn/v2 v2.0.0-alpha.9 // indirect
	github.com/SlothNinja/user-controller/v2 v2.0.0-alpha.4
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/gorilla/securecookie v1.1.1
	golang.org/x/sys v0.0.0-20200727154430-2d971f7391a4 // indirect
	golang.org/x/tools v0.0.0-20200727215259-7b4c4ad3dc39 // indirect
	google.golang.org/genproto v0.0.0-20200726014623-da3ae01ef02d // indirect
)

replace github.com/SlothNinja/user-controller/v2 => ./user-controller

replace github.com/SlothNinja/user/v2 => ./user
