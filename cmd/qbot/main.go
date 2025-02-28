package main

import (
	"context"
	"fmt"
	"github.com/Insulince/jlib/pkg/jmain"
	"github.com/Insulince/qbot/internal/qbot"
	"log"
)

func main() {
	jmain.MustMain(Main)
}

func Main(ctx context.Context) error {
	// Initialize and start the bot
	b, err := qbot.NewQBot()
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}

	<-ctx.Done()

	// Graceful shutdown
	fmt.Println("Shutting down bot...")
	b.Close()

	return nil
}
