package internal_test

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"testing"
	"vimpack/internal"
)

func testStringPlugins(plugins []*internal.Plugin) string {
	var teststring strings.Builder

	for _, plugin := range plugins {
		source := plugin.Source.(*internal.Git)
		teststring.WriteString(fmt.Sprintf("loading=%s plugin=%s namespace=%s source=%s\n", plugin.Loading, plugin.Name, plugin.Namespace, source.URL))
	}

	return teststring.String()
}

func TestParser(t *testing.T) {
	tcs := []struct {
		Name     string
		Filepath string
		Err      error
		Plugins  []*internal.Plugin
	}{
		{
			Name: "missing_namespace",
			Err:  errors.New("test-fixtures/parser-missing_namespace.vimpack:1:1: expected onstart to be in a namespace"),
		},
		{
			Name: "namespace_not_string",
			Err:  errors.New("test-fixtures/parser-namespace_not_string.vimpack:1:11: expected string got: anamespace"),
		},
		{
			Name: "unsupported_loading_method",
			Err:  errors.New("test-fixtures/parser-unsupported_loading_method.vimpack:5:1: unexpected token: onanything"),
		},
		{
			Name: "plugin_not_string",
			Err:  errors.New("test-fixtures/parser-plugin_not_string.vimpack:3:9: expected string got: johndoe"),
		},
		{
			Name: "plugin_wrong_format",
			Err:  errors.New("test-fixtures/parser-plugin_wrong_format.vimpack:3:18: expected repositoy name to be '<repository>/<name>', got: plugin1"),
		},
		{
			Name: "valid_file",
			Err:  nil,
			Plugins: []*internal.Plugin{
				&internal.Plugin{Name: "plugin1", Namespace: "n1", Loading: internal.LoadingStart, Source: internal.NewGitHub("johndoe/plugin1")},
				&internal.Plugin{Name: "plugin2", Namespace: "n1", Loading: internal.LoadingOnDemand, Source: internal.NewGitHub("janedoe/plugin2")},
				&internal.Plugin{Name: "otherplugin", Namespace: "n2", Loading: internal.LoadingStart, Source: internal.NewGitHub("johndoe/otherplugin")},
				&internal.Plugin{Name: "otherplugin", Namespace: "n2", Loading: internal.LoadingOnDemand, Source: internal.NewGitHub("janedoe/otherplugin")},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			parser := internal.NewParser()
			plugins, err := parser.ParseFile(path.Join("test-fixtures", fmt.Sprintf("parser-%s.vimpack", tc.Name)))

			if err == nil && tc.Err != nil || err != nil && tc.Err == nil {
				t.Fatalf("unexpected error. got: %#+v; want: %#+v", err, tc.Err)
			}

			if tc.Err != nil && err.Error() != tc.Err.Error() {
				t.Fatalf("wrong error. got: %s; want: %s", err, tc.Err)
			}

			if len(plugins) != len(tc.Plugins) {
				t.Errorf("wrong number of plugins. got: %d; want: %d", len(plugins), len(tc.Plugins))
			}

			actualPlugins := testStringPlugins(plugins)
			expectedPlugins := testStringPlugins(tc.Plugins)

			if actualPlugins != expectedPlugins {
				t.Errorf("wrong plugins.\ngot:\n%s\nwant:\n%s", actualPlugins, expectedPlugins)
			}
		})
	}
}
