package internal

import (
	"context"
	"path"
)

type PluginDownloader interface {
	Download(ctx context.Context, destination string) error
}

type Plugin struct {
	Loading   Loading
	Name      string
	Namespace Namespace
	Source    PluginDownloader
}

func NewPlugin(name string, namespace Namespace, source PluginDownloader) *Plugin {
	return &Plugin{
		Name:      name,
		Namespace: namespace,
		Source:    source,
	}
}

func (p *Plugin) Update(ctx context.Context, baseFolder string) error {
	return p.Source.Download(ctx, p.downloadPath(baseFolder))
}

func (p *Plugin) downloadPath(baseFolder string) string {
	return path.Join(baseFolder, p.Namespace.String(), p.Loading.String(), p.Name)
}
