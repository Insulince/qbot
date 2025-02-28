package version

import (
	"os"
	"strings"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/pkg/errors"
)

const Filename = "VERSION.md"

func Get() (string, error) {
	raw, err := os.ReadFile(Filename)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read %q", Filename)
	}
	contents := string(raw)

	version := strings.TrimSpace(contents)

	return version, nil
}

func MustGet() string {
	return jmust.Must[string](Get)[0]
}
