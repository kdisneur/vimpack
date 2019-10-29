package internal_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"vimpack/internal"
)

func TestProgressBarUIFormatter(t *testing.T) {
	tcs := []struct {
		Name    string
		Plugins []*internal.UIPlugin
	}{
		{
			Name:    "first_printing_no_plugins",
			Plugins: []*internal.UIPlugin{},
		},
		{
			Name: "first_printing_without_failed_plugins",
			Plugins: []*internal.UIPlugin{
				{
					Plugin: &internal.Plugin{
						Loading:   internal.LoadingStart,
						Name:      "plugin1",
						Namespace: "default",
					},
					State: internal.UIPluginStateDownloading,
				},
				{
					Plugin: &internal.Plugin{
						Loading:   internal.LoadingStart,
						Name:      "plugin2",
						Namespace: "default",
					},
					State: internal.UIPluginStateDownloaded,
				},
			},
		},
		{
			Name: "first_printing_with_failed_plugins",
			Plugins: []*internal.UIPlugin{
				{
					Plugin: &internal.Plugin{
						Loading:   internal.LoadingStart,
						Name:      "plugin1",
						Namespace: "default",
					},
					State: internal.UIPluginStateDownloading,
				},
				{
					Plugin: &internal.Plugin{
						Loading:   internal.LoadingStart,
						Name:      "plugin2",
						Namespace: "default",
					},
					State: internal.UIPluginStateDownloadFailed,
					Err:   errors.New("can't download"),
				},
				{
					Plugin: &internal.Plugin{
						Loading:   internal.LoadingStart,
						Name:      "plugin3",
						Namespace: "default",
					},
					State: internal.UIPluginStateDownloaded,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			var w bytes.Buffer

			progressBar := internal.NewProgressBarUIFormatter()
			progressBar.Print(&w, tc.Plugins)

			actual := w.Bytes()
			expected := GoldenBytes(t, fmt.Sprintf("progressbaruiformatter-%s-1", tc.Name), actual)

			if !bytes.Equal(actual, expected) {
				t.Errorf("unexpected output.\ngot:\n%s\nwant:\n%s", actual, expected)
			}

			if len(tc.Plugins) > 0 {
				tc.Plugins[0].State = internal.UIPluginStateDownloaded
			}

			progressBar.Print(&w, tc.Plugins)

			actual = w.Bytes()
			expected = GoldenBytes(t, fmt.Sprintf("progressbaruiformatter-%s-2", tc.Name), actual)

			if !bytes.Equal(actual, expected) {
				t.Errorf("unexpected output.\ngot:\n%s\nwant:\n%s", actual, expected)
			}
		})
	}
}
