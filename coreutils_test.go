package task_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-task/task/v3"
)

func TestCoreutils(t *testing.T) {
	tests := []struct {
		name         string
		taskfile     string
		aliasEnv     string
		expectedOut  string
		wantSetupErr bool
	}{
		{
			name:        "taskfile turtle",
			taskfile:    "coreutils: turtle\n",
			expectedOut: "charm.land/turtle",
		},
		{
			name:        "env alias",
			aliasEnv:    "turtle",
			expectedOut: "charm.land/turtle",
		},
		{
			name:        "taskfile wins over env",
			taskfile:    "coreutils: turtle\n",
			aliasEnv:    "off",
			expectedOut: "charm.land/turtle",
		},
		{
			name:         "taskfile invalid",
			taskfile:     "coreutils: bogus\n",
			wantSetupErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			taskfile := "version: '3'\n" + tt.taskfile + `tasks:
  default:
    cmds:
      - basename --version
`
			require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.yml"), []byte(taskfile), 0o644))
			t.Setenv("TASK_COREUTILS", tt.aliasEnv)

			var buff bytes.Buffer
			e := task.NewExecutor(
				task.WithDir(dir),
				task.WithStdout(&buff),
				task.WithStderr(&buff),
				task.WithSilent(true),
			)
			err := e.Setup()
			if tt.wantSetupErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid coreutils")
				return
			}
			require.NoError(t, err)
			require.NoError(t, e.Run(t.Context(), &task.Call{Task: "default"}))
			assert.Contains(t, buff.String(), tt.expectedOut)
		})
	}
}
