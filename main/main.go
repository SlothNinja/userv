package main

import (
	"context"
	"log"

	"github.com/SlothNinja/sn/v3"
	"github.com/SlothNinja/userv/main/client"
)

func main() {
	ctx := context.Background()
	cl := client.New(ctx)
	defer func() {
		if err := cl.Close(); err != nil {
			sn.Warnf(ctx, "error when closing client: %v", err)
		}
	}()

	var err error
	if sn.IsProduction() {
		err = cl.Router.Run()
	} else {
		err = cl.Router.Run(":" + cl.GetPort())
	}

	if err != nil {
		log.Panicf("unable to start server: %v", err)
	}
}
