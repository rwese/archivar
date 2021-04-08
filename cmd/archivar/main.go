package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rwese/archivar/archivar"
)

func main() {
	s := archivar.New()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	s.Run(ctx)
}
