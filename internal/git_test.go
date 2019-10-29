package internal_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"vimpack/internal"
)

func TestGitHelperProcess(t *testing.T) {
	if os.Getenv("VIMPACK_GIT_PROCESS_STARTED") != "1" {
		return
	}
	defer os.Exit(0)

	expectedCommand := os.Getenv("VIMPACK_GIT_PROCESS_COMMAND")

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[1:]
	}

	actualCommand := strings.Join(args, " ")

	if actualCommand != expectedCommand {
		fmt.Fprintf(os.Stderr, "unexpected command. got: %s; want: %s", actualCommand, expectedCommand)
		os.Exit(1)
	}
}

func gitCloneHelperProcessCmd(expectedDownloadCommand string) func(ctx context.Context, name string, s ...string) *exec.Cmd {
	return func(ctx context.Context, name string, s ...string) *exec.Cmd {
		cs := []string{"-test.run=TestGitHelperProcess", "--", name}
		cs = append(cs, s...)
		env := []string{
			"VIMPACK_GIT_PROCESS_STARTED=1",
			fmt.Sprintf("VIMPACK_GIT_PROCESS_COMMAND=%s", expectedDownloadCommand),
		}

		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = env
		return cmd
	}
}

func TestGitHubDownload(t *testing.T) {
	tcs := []struct {
		Destination        string
		Err                error
		ExpectedGitCommand string
		Name               string
		RepositoryName     string
		Setup              func() func()
	}{
		{
			Name:               "download_when_repository_never_cloned",
			RepositoryName:     "johndoe/plugin1",
			Destination:        path.Join("test-fixtures", "plugin1"),
			ExpectedGitCommand: "git clone git@github.com:johndoe/plugin1 test-fixtures/plugin1",
			Setup: func() func() {
				return func() { os.Remove(path.Join("test-fixtures", "plugin1")) }
			},
		},
		{
			Name:               "download_when_repository_exists",
			RepositoryName:     "johndoe/plugin1",
			Destination:        path.Join("test-fixtures", "plugin1"),
			ExpectedGitCommand: "git -C test-fixtures/plugin1 pull",
			Setup: func() func() {
				f := path.Join("test-fixtures", "plugin1")
				os.Mkdir(f, 0755)
				return func() { os.Remove(f) }
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			teardown := tc.Setup()
			defer teardown()

			git := internal.NewGitHub(tc.RepositoryName)
			git.CmdExec = gitCloneHelperProcessCmd(tc.ExpectedGitCommand)

			err := git.Download(ctx, tc.Destination)
			if err != nil && tc.Err == nil {
				t.Fatalf("unexpected error. got: %s", err)
			}

			if err == nil && tc.Err != nil {
				t.Fatalf("expected error. got: none; want: %s", tc.Err)
			}

			if tc.Err != nil && err.Error() != tc.Err.Error() {
				t.Fatalf("expected error. got: %s; want: %s", err, tc.Err)
			}
		})
	}
}
