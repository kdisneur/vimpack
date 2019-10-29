package internal

import (
	"fmt"
	"io"
)

type ProgressBarUIFormatter struct {
	dirty              bool
	registeredFailures map[string]bool
}

func NewProgressBarUIFormatter() *ProgressBarUIFormatter {
	return &ProgressBarUIFormatter{
		registeredFailures: make(map[string]bool),
	}
}

func (p *ProgressBarUIFormatter) Print(w io.Writer, plugins []*UIPlugin) {
	var downloaded, failed int

	if p.dirty {
		fmt.Fprintf(w, "\033[A")
	}

	for _, plugin := range plugins {
		switch plugin.State {
		case UIPluginStateDownloaded:
			downloaded++
		case UIPluginStateDownloadFailed:
			name := fmt.Sprintf("%s;%s", plugin.Namespace, plugin.Name)
			if !p.registeredFailures[name] {
				p.registeredFailures[name] = true
				fmt.Fprintf(w, "%s[namespace: %s]: %s\n", plugin.Name, plugin.Namespace, plugin.Err.Error())
			}
			failed++
		default:
		}
	}

	fmt.Fprintf(w, "%d/%d downloaded. %d failed\n", downloaded, len(plugins), failed)

	p.dirty = true
}
