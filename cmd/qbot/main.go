package main

import (
	"context"

	"github.com/Insulince/jlib/pkg/jmain"
	"github.com/pkg/errors"
	_ "modernc.org/sqlite" // SQLite driver

	"github.com/Insulince/qbot/internal/config"
	"github.com/Insulince/qbot/internal/qbot"
	"github.com/Insulince/qbot/internal/store"
)

func main() {
	jmain.MustMain(Main)
}

func Main(ctx context.Context) error {
	cfg := config.MustGetConfig()

	deps := config.MustBuildDependencies(cfg)
	defer deps.MustClose()

	s := store.MustNew(deps.Db)
	defer s.MustClose()

	q := qbot.MustNew(cfg, s)
	defer q.MustClose()

	select {
	case <-ctx.Done():
	case err := <-q.Run(ctx):
		if err != nil {
			return errors.Wrap(err, "running qbot")
		}
	}

	return nil
}
