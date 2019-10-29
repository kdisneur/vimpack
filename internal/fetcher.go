package internal

import (
	"context"
	"fmt"
	"sync"
)

type FetcherDisplay interface {
	AddPlugin(*Plugin) (*UIPlugin, UpdateFn)
	StartEventLoop()
	StopEventLoop() <-chan interface{}
}

type Fetcher struct {
	UI FetcherDisplay
}

func NewFetcher(ui FetcherDisplay) *Fetcher {
	return &Fetcher{UI: ui}
}

func (f *Fetcher) All(ctx context.Context, plugins []*Plugin, destination string) {
	downloadFinished := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(len(plugins))

	go f.UI.StartEventLoop()

	for _, plugin := range plugins {
		uiPlugin, updateFn := f.UI.AddPlugin(plugin)
		go f.updatePlugin(ctx, &wg, updateFn, uiPlugin, destination)
	}

	go f.waitDownloads(&wg, downloadFinished)

	select {
	case <-ctx.Done():
	case <-downloadFinished:
		<-f.UI.StopEventLoop()
	}
}

func (f *Fetcher) updatePlugin(ctx context.Context, wg *sync.WaitGroup, updateFn UpdateFn, plugin *UIPlugin, destination string) {
	defer updateFn()
	defer wg.Done()

	plugin.State = UIPluginStateDownloading

	updateFn()

	if err := plugin.Plugin.Update(ctx, destination); err != nil {
		plugin.Err = fmt.Errorf("can't download %s: %s", plugin.Plugin.Name, err)
		plugin.State = UIPluginStateDownloadFailed
		return
	}

	plugin.State = UIPluginStateDownloaded
}

func (f *Fetcher) waitDownloads(wg *sync.WaitGroup, downloadFinished chan<- interface{}) {
	wg.Wait()
	downloadFinished <- true
}
