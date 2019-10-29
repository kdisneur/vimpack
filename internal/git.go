package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type Git struct {
	CmdExec func(ctx context.Context, name string, cs ...string) *exec.Cmd
	URL     string
}

func NewGitHub(repositoryName string) *Git {
	return &Git{
		CmdExec: exec.CommandContext,
		URL:     fmt.Sprintf("git@github.com:%s", repositoryName),
	}
}

func (g *Git) Download(ctx context.Context, destination string) error {
	stat, err := os.Stat(destination)

	if os.IsNotExist(err) {
		return g.clone(ctx, destination)
	}

	if !stat.IsDir() {
		return fmt.Errorf("can't clone in '%s', it's not a directory", destination)
	}

	return g.pull(ctx, destination)
}

func (g *Git) clone(ctx context.Context, destination string) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("can't clone: %s", err)
	}

	if err := g.CmdExec(ctx, "git", "clone", g.URL, destination).Run(); err != nil {
		return fmt.Errorf("can't clone: %s", err)
	}

	return nil
}

func (g *Git) pull(ctx context.Context, destination string) error {
	if err := g.CmdExec(ctx, "git", "-C", destination, "pull").Run(); err != nil {
		return fmt.Errorf("can't pull: %s", err)
	}

	return nil
}
