package main

import (
	"context"
	"os"

	"github.com/SlothNinja/sn/v3"
)

func main() {
	cl := sn.NewUserServiceClient(
		context.Background(),
		sn.WithLoggerID("user-service"),
		sn.WithCORSAllow("https://user.fake-slothninja.com:8088/*"),
	)
	defer cl.Close()

	if sn.IsProduction() {
		cl.Router.Run()
	} else {
		cl.Router.RunTLS(getPort(), "cert.pem", "key.pem")
	}
}

func getPort() string {
	return ":" + os.Getenv("PORT")
}
