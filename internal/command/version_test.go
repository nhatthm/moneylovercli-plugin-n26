package command_test

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nhatthm/moneylovercli-plugin-n26/internal/command"
)

func TestNewVersion(t *testing.T) {
	t.Parallel()

	stdout := new(bytes.Buffer)
	cmd := command.New(command.WithStdout(stdout))

	cmd.SetArgs([]string{"version"})

	err := cmd.Execute()

	require.NoError(t, err)

	expected := fmt.Sprintf("dev (rev: ; %s; %s/%s)\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	assert.Equal(t, expected, stdout.String())
}
