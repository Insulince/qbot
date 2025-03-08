package main

import (
	"context"
	"fmt"
	"github.com/Insulince/jlib/pkg/jmain"
	"github.com/Insulince/qbot/internal/qbot"
	"github.com/pkg/errors"
	"os"

	_ "modernc.org/sqlite" // SQLite driver
)

const (
	EnvVarDiscordBotToken = "DISCORD_BOT_TOKEN"
)

func main() {
	jmain.MustMain(Main)
}

func Main(ctx context.Context) error {
	token, found := os.LookupEnv(EnvVarDiscordBotToken)
	if !found {
		return fmt.Errorf("environment variable %q not set", EnvVarDiscordBotToken)
	}

	// Initialize and start the bot
	q, err := qbot.New(token)
	if err != nil {
		return errors.Wrap(err, "new qbot")
	}

	// Block main so the bot can run.
	<-ctx.Done()

	// Graceful shutdown
	fmt.Println("Q is shutting down...")
	q.Close()

	return nil
}
