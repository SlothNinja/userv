package main

import (
	"context"
	"os"

	"github.com/SlothNinja/sn/v3"
)

const homePath = "/"

func main() {
	cl := newClient(
		context.Background(),
		sn.WithLoggerID("user-service"),
		sn.WithCORSAllow("https://user.fake-slothninja.com:8088/*"),
	)
	defer func() {
		if err := cl.close(); err != nil {
			sn.Warningf("error when closing client: %w", err)
		}
	}()

	if sn.IsProduction() {
		cl.Router.Run()
	} else {
		cl.Router.RunTLS(getPort(), "cert.pem", "key.pem")
	}
}

func getPort() string {
	return ":" + os.Getenv("PORT")
}
