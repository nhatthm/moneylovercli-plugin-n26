package version

import (
	"runtime"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	t.Parallel()

	expected := Information{
		Version:      "dev",
		Revision:     "",
		Branch:       "",
		BuildUser:    "",
		BuildDate:    "",
		GoVersion:    runtime.Version(),
		GoOS:         runtime.GOOS,
		GoArch:       runtime.GOARCH,
		Dependencies: nil,
	}

	assert.Equal(t, expected, Info())
}

func TestWriteInformation(t *testing.T) {
	t.Parallel()

	info := Information{
		Version:   "dev",
		Revision:  "d586cddb296a43e2ba51748a5c2e55412626ccf3",
		BuildUser: "user",
		BuildDate: "2021-02-03T04:05:06Z",
		GoVersion: "1.16",
		GoOS:      "darwin",
		GoArch:    "amd64",
		Dependencies: []*debug.Module{
			{
				Path:    "github.com/nhatthm/n26api",
				Version: "v0.3.3",
			},
		},
	}

	var sb strings.Builder

	expected := "dev (rev: d586cddb296a43e2ba51748a5c2e55412626ccf3; 1.16; darwin/amd64)\n"

	WriteInformation(&sb, info, false)

	assert.Equal(t, expected, sb.String())
}

func TestWriteInformation_Full(t *testing.T) {
	t.Parallel()

	info := Information{
		Version:   "dev",
		Revision:  "d586cddb296a43e2ba51748a5c2e55412626ccf3",
		BuildUser: "user",
		BuildDate: "2021-02-03T04:05:06Z",
		GoVersion: "1.16",
		GoOS:      "darwin",
		GoArch:    "amd64",
		Dependencies: []*debug.Module{
			{
				Path:    "github.com/nhatthm/n26api",
				Version: "v0.3.3",
			},
		},
	}

	var sb strings.Builder

	expected := `dev (rev: d586cddb296a43e2ba51748a5c2e55412626ccf3; 1.16; darwin/amd64)

build user: user
build date: 2021-02-03T04:05:06Z

dependencies:
  github.com/nhatthm/n26api: v0.3.3
`

	WriteInformation(&sb, info, true)

	assert.Equal(t, expected, sb.String())
}
