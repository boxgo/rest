package main

import (
	"github.com/boxgo/box"
	"github.com/boxgo/logger"
)

func main() {
	server := newRestServer()
	config := newConfig()

	box := box.
		NewBox(
			box.WithConfig(config),
			box.WithBoxes(server),
		).
		Mount(
			logger.Default,
		)

	box.Serve()
}
