package internal_test

import (
	"context"
	"errors"
	"testing"
	"vimpack/internal"
)

var ErrFakeDownloadingFailure = errors.New("faking download error")

type FakeDownloader struct {
	ExpectedDestination string
	T                   *testing.T
	Err                 error
}

func (f *FakeDownloader) Download(ctx context.Context, destination string) error {
	if f.ExpectedDestination != destination {
		f.T.Errorf("wrong destination. got: %s; want: %s", destination, f.ExpectedDestination)
	}
	return f.Err
}

func TestPluginUpdate(t *testing.T) {
	tcs := []struct {
		Name          string
		Destination   string
		MakePlugin    func() *internal.Plugin
		ExpectedError error
	}{
		{
			Name:        "plugin_loaded_onstart",
			Destination: "/home/jdoe/.vim/vimpack",
			MakePlugin: func() *internal.Plugin {
				downloader := &FakeDownloader{
					T:                   t,
					ExpectedDestination: "/home/jdoe/.vim/vimpack/default/start/plugin1",
				}

				return internal.NewPlugin("plugin1", "default", downloader)
			},
		},
		{
			Name:        "plugin_loaded_ondemand",
			Destination: "/home/jdoe/.vim/vimpack",
			MakePlugin: func() *internal.Plugin {
				downloader := &FakeDownloader{
					T:                   t,
					ExpectedDestination: "/home/jdoe/.vim/vimpack/default/opt/plugin1",
				}

				plugin := internal.NewPlugin("plugin1", "default", downloader)
				plugin.Loading = internal.LoadingOnDemand

				return plugin
			},
		},
		{
			Name:          "plugin_download_fail",
			ExpectedError: ErrFakeDownloadingFailure,
			Destination:   "/home/jdoe/.vim/vimpack",
			MakePlugin: func() *internal.Plugin {
				downloader := &FakeDownloader{
					Err:                 ErrFakeDownloadingFailure,
					T:                   t,
					ExpectedDestination: "/home/jdoe/.vim/vimpack/default/opt/plugin1",
				}

				plugin := internal.NewPlugin("plugin1", "default", downloader)
				plugin.Loading = internal.LoadingOnDemand

				return plugin
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			plugin := tc.MakePlugin()

			err := plugin.Update(ctx, tc.Destination)

			if err != tc.ExpectedError {
				t.Errorf("unexpected error. got: %s; want: %s", err, tc.ExpectedError)
			}
		})
	}
}
