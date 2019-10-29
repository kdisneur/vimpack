package internal

import (
	"context"
	"io"
	"sync"
)

type UpdateFn func()

type UIFormatter interface {
	Print(w io.Writer, plugins []*UIPlugin)
}

type UIPluginState string

type UIPlugin struct {
	*Plugin
	State UIPluginState
	Err   error
}

const (
	UIPluginStateWaiting        UIPluginState = "awaiting"
	UIPluginStateDownloading    UIPluginState = "downloading"
	UIPluginStateDownloaded     UIPluginState = "downloaded"
	UIPluginStateDownloadFailed UIPluginState = "failed"
)

type UI struct {
	Formatter    UIFormatter
	Out          io.Writer
	ctx          context.Context
	hasStopped   bool
	mutex        sync.Mutex
	plugins      []*UIPlugin
	shouldStop   chan interface{}
	shouldUpdate chan interface{}
	stopped      chan interface{}
}

func NewUI(ctx context.Context, out io.Writer) *UI {
	return &UI{
		Formatter:    NewProgressBarUIFormatter(),
		Out:          out,
		ctx:          ctx,
		plugins:      []*UIPlugin{},
		shouldStop:   make(chan interface{}),
		shouldUpdate: make(chan interface{}),
		stopped:      make(chan interface{}),
	}
}

func (ui *UI) AddPlugin(plugin *Plugin) (*UIPlugin, UpdateFn) {
	ui.mutex.Lock()

	defer ui.TriggerUpdate()
	defer ui.mutex.Unlock()

	uiPlugin := &UIPlugin{Plugin: plugin, State: UIPluginStateWaiting}

	ui.plugins = append(ui.plugins, uiPlugin)

	return uiPlugin, ui.TriggerUpdate
}

func (ui *UI) StartEventLoop() {
	ui.Print()

	for {
		select {
		case <-ui.ctx.Done():
			return
		case <-ui.shouldUpdate:
			ui.Print()
		case <-ui.shouldStop:
			ui.hasStopped = true
			ui.Print()
			ui.stopped <- true
			close(ui.stopped)
		}
	}
}

func (ui *UI) StopEventLoop() <-chan interface{} {
	if !ui.hasStopped {
		ui.shouldStop <- true
	}

	return ui.stopped
}

func (ui *UI) TriggerUpdate() {
	if ui.hasStopped {
		return
	}

	ui.shouldUpdate <- true
}

func (ui *UI) Print() {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()

	ui.Formatter.Print(ui.Out, ui.plugins)
}
