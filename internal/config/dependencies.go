package config

import (
	"database/sql"
	"io"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/pkg/errors"
)

type Dependencies struct {
	Db *sql.DB
}

var _ io.Closer = Dependencies{}

func BuildDependencies(cfg Config) (Dependencies, error) {
	var deps Dependencies

	db, err := sql.Open(cfg.DbDriver, cfg.DbFile)
	if err != nil {
		return Dependencies{}, errors.Wrap(err, "error opening database")
	}
	deps.Db = db

	return deps, nil
}

func MustBuildDependencies(cfg Config) Dependencies {
	return jmust.Must[Dependencies](BuildDependencies, cfg)[0]
}

func (deps Dependencies) Close() error {
	if err := deps.Db.Close(); err != nil {
		return errors.Wrap(err, "db close")
	}

	return nil
}

func (deps Dependencies) MustClose() {
	jmust.MustClose(deps)
}
