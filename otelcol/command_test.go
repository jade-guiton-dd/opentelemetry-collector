// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otelcol

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/featuregate"
)

func TestNewCommandVersion(t *testing.T) {
	cmd := NewCommand(CollectorSettings{BuildInfo: component.BuildInfo{Version: "test_version"}})
	assert.Equal(t, "test_version", cmd.Version)
}

func TestNewCommandNoConfigURI(t *testing.T) {
	cmd := NewCommand(CollectorSettings{Factories: nopFactories})
	require.Error(t, cmd.Execute())
}

// This test emulates usage of Collector in Jaeger all-in-one, which
// allows running the binary with no explicit configuration.
func TestNewCommandProgrammaticallyPassedConfig(t *testing.T) {
	cmd := NewCommand(CollectorSettings{Factories: nopFactories, ConfigProviderSettings: ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			ProviderFactories: []confmap.ProviderFactory{confmap.NewProviderFactory(newFailureProvider)},
			DefaultScheme:     "file",
		},
	}})
	otelRunE := cmd.RunE
	cmd.RunE = func(c *cobra.Command, args []string) error {
		configFlag := c.Flag("config")
		cfg := `
service:
  extensions: [invalid_component_name]
receivers:
  invalid_component_name:
`
		require.NoError(t, configFlag.Value.Set("yaml:"+cfg))
		return otelRunE(cmd, args)
	}
	// verify that cmd.Execute was run with the implicitly provided config.
	require.ErrorContains(t, cmd.Execute(), "invalid_component_name")
}

func TestAddFlagToSettings(t *testing.T) {
	filePath := filepath.Join("testdata", "otelcol-invalid.yaml")
	fileProvider := newFakeProvider("file", func(_ context.Context, _ string, _ confmap.WatcherFunc) (*confmap.Retrieved, error) {
		return confmap.NewRetrieved(newConfFromFile(t, filePath))
	})
	set := CollectorSettings{
		ConfigProviderSettings: ConfigProviderSettings{
			ResolverSettings: confmap.ResolverSettings{
				URIs:              []string{filePath},
				ProviderFactories: []confmap.ProviderFactory{fileProvider},
			},
		},
	}
	flgs := flags(featuregate.NewRegistry())
	err := flgs.Parse([]string{"--config=a.yaml", "--config=b.yaml"})
	require.NoError(t, err)

	err = updateSettingsUsingFlags(&set, flgs)
	require.NoError(t, err)
	require.Equal(t, []string{"a.yaml", "b.yaml"}, set.ConfigProviderSettings.ResolverSettings.URIs)
}

func TestNewCommandMultipleConfigs(t *testing.T) {
	fileProvider := newFakeProvider("file", func(_ context.Context, uri string, _ confmap.WatcherFunc) (*confmap.Retrieved, error) {
		path := strings.TrimPrefix(uri, "file:")
		return confmap.NewRetrieved(newConfFromFile(t, path))
	})
	set := CollectorSettings{
		Factories: nopFactories,
		ConfigProviderSettings: ConfigProviderSettings{
			ResolverSettings: confmap.ResolverSettings{
				ProviderFactories: []confmap.ProviderFactory{fileProvider},
				DefaultScheme:     "file",
			},
		},
	}
	cmd := NewCommand(set)
	otelRunE := cmd.RunE
	cmd.RunE = func(c *cobra.Command, args []string) error {
		require.NoError(t, c.Flag("config").Value.Set(filepath.Join("testdata", "otelcol-nop.yaml")))
		require.NoError(t, c.Flag("config").Value.Set(filepath.Join("testdata", "otelcol-invalid.yaml")))

		return otelRunE(cmd, args)
	}
	require.ErrorContains(t, cmd.Execute(), "references processor \"invalid\" which is not configured")
}

func TestInvalidCollectorSettings(t *testing.T) {
	set := CollectorSettings{
		ConfigProviderSettings: ConfigProviderSettings{
			ResolverSettings: confmap.ResolverSettings{
				URIs: []string{"--config=otelcol-nop.yaml"},
			},
		},
	}

	cmd := NewCommand(set)
	require.Error(t, cmd.Execute())
}

func TestNewCommandInvalidComponent(t *testing.T) {
	filePath := filepath.Join("testdata", "otelcol-invalid.yaml")
	fileProvider := newFakeProvider("file", func(_ context.Context, _ string, _ confmap.WatcherFunc) (*confmap.Retrieved, error) {
		return confmap.NewRetrieved(newConfFromFile(t, filePath))
	})
	set := ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs:              []string{filePath},
			ProviderFactories: []confmap.ProviderFactory{fileProvider},
		},
	}

	cmd := NewCommand(CollectorSettings{Factories: nopFactories, ConfigProviderSettings: set})
	require.Error(t, cmd.Execute())
}

func TestNoProvidersReturnsError(t *testing.T) {
	set := CollectorSettings{
		ConfigProviderSettings: ConfigProviderSettings{
			ResolverSettings: confmap.ResolverSettings{
				URIs: []string{filepath.Join("testdata", "otelcol-invalid.yaml")},
			},
		},
	}
	flgs := flags(featuregate.NewRegistry())
	err := flgs.Parse([]string{"--config=otelcol-nop.yaml"})
	require.NoError(t, err)

	err = updateSettingsUsingFlags(&set, flgs)
	require.ErrorContains(t, err, "at least one Provider must be supplied")
}

func Test_UseUnifiedEnvVarExpansionRules(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "default scheme set",
			input:    "file",
			expected: "file",
		},
		{
			name:     "default scheme not set",
			input:    "",
			expected: "env",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileProvider := newFakeProvider("file", func(_ context.Context, _ string, _ confmap.WatcherFunc) (*confmap.Retrieved, error) {
				return &confmap.Retrieved{}, nil
			})
			set := CollectorSettings{
				ConfigProviderSettings: ConfigProviderSettings{
					ResolverSettings: confmap.ResolverSettings{
						ProviderFactories: []confmap.ProviderFactory{fileProvider},
						DefaultScheme:     tt.input,
					},
				},
			}
			flgs := flags(featuregate.NewRegistry())
			err := flgs.Parse([]string{"--config=otelcol-nop.yaml"})
			require.NoError(t, err)

			err = updateSettingsUsingFlags(&set, flgs)
			require.NoError(t, err)
			require.Equal(t, tt.expected, set.ConfigProviderSettings.ResolverSettings.DefaultScheme)
		})
	}
}

func TestNewFeatureGateCommand(t *testing.T) {
	t.Run("list all featuregates", func(t *testing.T) {
		cmd := newFeatureGateCommand()
		require.NotNil(t, cmd)

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := cmd.RunE(cmd, []string{})
		require.NoError(t, err)

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = oldStdout

		output := string(out)
		assert.Contains(t, output, "ID")
		assert.Contains(t, output, "Enabled")
		assert.Contains(t, output, "Stage")
		assert.Contains(t, output, "Description")
	})
	t.Run("specific featuregate details", func(t *testing.T) {
		cmd := newFeatureGateCommand()

		// Register a test feature gate in the global registry
		featuregate.GlobalRegistry().MustRegister("test.feature", featuregate.StageBeta,
			featuregate.WithRegisterDescription("Test feature description"))

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := cmd.RunE(cmd, []string{"test.feature"})
		require.NoError(t, err)

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = oldStdout

		output := string(out)
		assert.Contains(t, output, "Feature: test.feature")
		assert.Contains(t, output, "Description: Test feature description")
		assert.Contains(t, output, "Stage: Beta")
	})

	t.Run("non-existent featuregate", func(t *testing.T) {
		cmd := newFeatureGateCommand()
		err := cmd.RunE(cmd, []string{"non.existent.feature"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "feature \"non.existent.feature\" not found")
	})

	t.Run("rejects multiple arguments", func(t *testing.T) {
		cmd := newFeatureGateCommand()
		cmd.SetArgs([]string{"gate1", "gate2"})
		err := cmd.Execute()
		require.Error(t, err)
	})
}

func TestFeatureGateStability(t *testing.T) {
	reg := featuregate.GlobalRegistry()
	alphaGate := reg.MustRegister("test.channel.alpha", featuregate.StageAlpha)
	betaGate := reg.MustRegister("test.channel.beta", featuregate.StageBeta)
	stableGate := reg.MustRegister("test.channel.stable", featuregate.StageStable,
		featuregate.WithRegisterToVersion("v1.0.0"))
	t.Cleanup(func() {
		reg.SetStableOnly(false)
	})

	newSet := func(stability featuregate.Stage) CollectorSettings {
		return CollectorSettings{
			DefaultStabilitySettings: DefaultStabilitySettings{FeatureGates: stability},
			ConfigProviderSettings: ConfigProviderSettings{
				ResolverSettings: confmap.ResolverSettings{
					ProviderFactories: []confmap.ProviderFactory{confmap.NewProviderFactory(newFailureProvider)},
				},
			},
		}
	}

	tests := []struct {
		name                string
		settingsStability   featuregate.Stage
		channel             string
		featureGates        string
		expectedStability   featuregate.Stage
		expectedBetaEnabled bool
		expectedErr         string
	}{
		{
			name:                "stable channel keeps only stable gates enabled",
			channel:             "stable",
			expectedStability:   featuregate.StageStable,
			expectedBetaEnabled: false,
		},
		{
			name:                "beta channel enables beta and stable gates",
			channel:             "beta",
			expectedStability:   featuregate.StageBeta,
			expectedBetaEnabled: true,
		},
		{
			name:         "stable channel rejects customized feature gates",
			channel:      "stable",
			featureGates: "+test.channel.beta",
			expectedErr:  "cannot customize feature gates in stable release channel; consider switching to --channel beta",
		},
		{
			name:                "feature gate stability in settings keeps only stable gates enabled",
			settingsStability:   featuregate.StageStable,
			expectedStability:   featuregate.StageStable,
			expectedBetaEnabled: false,
		},
		{
			name:                "channel flag takes precedence over feature gate stability in settings",
			settingsStability:   featuregate.StageStable,
			channel:             "beta",
			expectedStability:   featuregate.StageBeta,
			expectedBetaEnabled: true,
		},
		{
			name:              "feature gate stability in settings rejects customized feature gates",
			settingsStability: featuregate.StageStable,
			featureGates:      "+test.channel.beta",
			expectedErr:       "cannot customize feature gates in stable release channel; consider switching to --channel beta",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg.SetStableOnly(false)

			set := newSet(tt.settingsStability)
			args := []string{"--config=a.yaml"}
			if tt.channel != "" {
				args = append(args, "--channel="+tt.channel)
			}
			if tt.featureGates != "" {
				args = append(args, "--feature-gates="+tt.featureGates)
			}
			flgs := flags(reg)
			require.NoError(t, flgs.Parse(args))

			err := updateSettingsUsingFlags(&set, flgs)
			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tt.expectedStability, set.DefaultStabilitySettings.FeatureGates)
			assert.False(t, alphaGate.IsEnabled())
			assert.Equal(t, tt.expectedBetaEnabled, betaGate.IsEnabled())
			assert.True(t, stableGate.IsEnabled())
		})
	}
}
