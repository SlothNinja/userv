package main

import (
	"context"

	"github.com/SlothNinja/sn/v3"
	"github.com/SlothNinja/userv/main/client"
)

func main() {
	if sn.IsProduction() {
		cl := client.New(
			context.Background(),
			sn.WithLoggerID("user-service"),
		)
		defer func() {
			if err := cl.Close(); err != nil {
				sn.Warningf("error when closing client: %w", err)
			}
		}()

		cl.Router.Run()
		return
	}

	cl := client.New(
		context.Background(),
		sn.WithLoggerID("user-service"),
		sn.WithCORSAllow("https://user.fake-slothninja.com:8088/*"),
	)
	defer func() {
		if err := cl.Close(); err != nil {
			sn.Warningf("error when closing client: %w", err)
		}
	}()

	cl.Router.RunTLS(":"+cl.GetPort(), "cert.pem", "key.pem")
}
