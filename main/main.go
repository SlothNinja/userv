package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/SlothNinja/sn/v3"
	"github.com/SlothNinja/userv/main/client"
)

func main() {
	ctx := context.Background()
	cl := client.New(ctx)
	defer func() {
		if err := cl.Close(); err != nil {
			slog.Warn(fmt.Sprintf("error when closing client: %w", err))
		}
	}()

	if sn.IsProduction() {
		cl.Router.Run()
	}
	cl.Router.Run(":" + cl.GetPort())
}
