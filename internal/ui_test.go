package internal_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"io"
	"strings"
	"testing"
	"vimpack/internal"
	"vimpack/internal/mock_internal"
)

func TestUIAddPlugin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var out strings.Builder
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	formatter := mock_internal.NewMockUIFormatter(ctrl)

	formatter.
		EXPECT().
		Print(gomock.Any(), gomock.Any()).
		Do(func(w io.Writer, plugins []*internal.UIPlugin) {
			if w != &out {
				t.Errorf("wrong writer. got: %#+v, want: %#+v", w, &out)
			}
		}).
		Times(3) // onstart, onaddplugin, onstop

	ui := internal.NewUI(ctx, &out)
	ui.Formatter = formatter
	go ui.StartEventLoop()

	plugin := internal.NewPlugin("plugin1", "n1", internal.NewGitHub("johndoe/plugin1"))

	actual, _ := ui.AddPlugin(plugin)
	if actual.Plugin != plugin {
		t.Errorf("wrong plugin. got: %#+v; want: %#+v", actual.Plugin, plugin)
	}

	<-ui.StopEventLoop()
}

func TestUIAddPluginUpdateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var out strings.Builder
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	formatter := mock_internal.NewMockUIFormatter(ctrl)

	formatter.
		EXPECT().
		Print(gomock.Any(), gomock.Any()).
		Do(func(w io.Writer, plugins []*internal.UIPlugin) {
			if w != &out {
				t.Errorf("wrong writer. got: %#+v, want: %#+v", w, &out)
			}
		}).
		Times(4) // onstart, onaddplugin, onupdate, onstop

	ui := internal.NewUI(ctx, &out)
	ui.Formatter = formatter
	go ui.StartEventLoop()

	plugin := internal.NewPlugin("plugin1", "n1", internal.NewGitHub("johndoe/plugin1"))

	actual, updateFn := ui.AddPlugin(plugin)
	if actual.Plugin != plugin {
		t.Errorf("wrong plugin. got: %#+v; want: %#+v", actual.Plugin, plugin)
	}

	updateFn()

	<-ui.StopEventLoop()
}
