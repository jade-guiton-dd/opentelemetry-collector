// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/connector/xconnector"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/xexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/xprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/xreceiver"
	"go.opentelemetry.io/collector/service/internal/builders"
)

// TestMinimumStability verifies that the minimum component stability level
// set in the service settings is enforced when creating components.
func TestMinimumStability(t *testing.T) {
	defaultCfg := struct{}{}

	t.Run("receiver", func(t *testing.T) {
		allType := component.MustNewType("all")
		factories, err := otelcol.MakeFactoryMap([]receiver.Factory{
			xreceiver.NewFactory(
				allType,
				func() component.Config { return &defaultCfg },
				xreceiver.WithTraces(createReceiverTraces, component.StabilityLevelDevelopment),
				xreceiver.WithMetrics(createReceiverMetrics, component.StabilityLevelAlpha),
				xreceiver.WithLogs(createReceiverLogs, component.StabilityLevelDeprecated),
			),
		}...)
		require.NoError(t, err)
		allID := component.NewID(allType)
		cfgs := map[component.ID]component.Config{
			allID: defaultCfg,
		}

		b := builders.NewReceiver(cfgs, factories, component.StabilityLevelAlpha)
		_, err = b.CreateTraces(context.Background(), settings(allID), consumertest.NewNop())
		require.ErrorContains(t, err, `component is too unstable for use: stability level is Development versus required Alpha`)

		me, err := b.CreateMetrics(context.Background(), settings(allID), consumertest.NewNop())
		require.NoError(t, err)
		assert.Equal(t, nopReceiverInstance, me)

		// Deprecated components are allowed regardless of the minimum stability level.
		le, err := b.CreateLogs(context.Background(), settings(allID), consumertest.NewNop())
		require.NoError(t, err)
		assert.Equal(t, nopReceiverInstance, le)
	})

	t.Run("processor", func(t *testing.T) {
		allType := component.MustNewType("all")
		factories, err := otelcol.MakeFactoryMap([]processor.Factory{
			xprocessor.NewFactory(
				allType,
				func() component.Config { return &defaultCfg },
				xprocessor.WithTraces(createProcessorTraces, component.StabilityLevelDevelopment),
				xprocessor.WithMetrics(createProcessorMetrics, component.StabilityLevelAlpha),
			),
		}...)
		require.NoError(t, err)
		allID := component.NewID(allType)
		cfgs := map[component.ID]component.Config{allID: defaultCfg}

		b := builders.NewProcessor(cfgs, factories, component.StabilityLevelAlpha)
		_, err = b.CreateTraces(context.Background(), createProcessorSettings(allID), consumertest.NewNop())
		require.ErrorContains(t, err, `component is too unstable for use: stability level is Development versus required Alpha`)

		me, err := b.CreateMetrics(context.Background(), createProcessorSettings(allID), consumertest.NewNop())
		require.NoError(t, err)
		assert.Equal(t, nopProcessorInstance, me)
	})

	t.Run("exporter", func(t *testing.T) {
		allType := component.MustNewType("all")
		factories, err := otelcol.MakeFactoryMap([]exporter.Factory{
			xexporter.NewFactory(
				allType,
				func() component.Config { return &defaultCfg },
				xexporter.WithTraces(createExporterTraces, component.StabilityLevelDevelopment),
				xexporter.WithMetrics(createExporterMetrics, component.StabilityLevelAlpha),
			),
		}...)
		require.NoError(t, err)
		allID := component.NewID(allType)
		cfgs := map[component.ID]component.Config{allID: defaultCfg}

		b := builders.NewExporter(cfgs, factories, component.StabilityLevelAlpha)
		_, err = b.CreateTraces(context.Background(), createExporterSettings(allID))
		require.ErrorContains(t, err, `component is too unstable for use: stability level is Development versus required Alpha`)

		me, err := b.CreateMetrics(context.Background(), createExporterSettings(allID))
		require.NoError(t, err)
		assert.Equal(t, nopExporterInstance, me)
	})

	t.Run("connector", func(t *testing.T) {
		allType := component.MustNewType("all")
		factories, err := otelcol.MakeFactoryMap([]connector.Factory{
			xconnector.NewFactory(
				allType,
				func() component.Config { return &defaultCfg },
				xconnector.WithTracesToTraces(createConnectorTracesToTraces, component.StabilityLevelDevelopment),
				xconnector.WithMetricsToTraces(createConnectorMetricsToTraces, component.StabilityLevelAlpha),
			),
		}...)
		require.NoError(t, err)
		allID := component.NewID(allType)
		cfgs := map[component.ID]component.Config{allID: defaultCfg}

		b := builders.NewConnector(cfgs, factories, component.StabilityLevelAlpha)
		_, err = b.CreateTracesToTraces(context.Background(), createConnectorSettings(allID), consumertest.NewNop())
		require.ErrorContains(t, err, `component is too unstable for use: stability level is Development versus required Alpha`)

		ce, err := b.CreateMetricsToTraces(context.Background(), createConnectorSettings(allID), consumertest.NewNop())
		require.NoError(t, err)
		assert.Equal(t, nopConnectorInstance, ce)
	})

	t.Run("extension", func(t *testing.T) {
		testType := component.MustNewType("test")
		factories, err := otelcol.MakeFactoryMap([]extension.Factory{
			extension.NewFactory(
				testType,
				func() component.Config { return &defaultCfg },
				func(_ context.Context, set extension.Settings, _ component.Config) (extension.Extension, error) {
					return nopExtension{Settings: set}, nil
				},
				component.StabilityLevelDevelopment,
			),
		}...)
		require.NoError(t, err)
		testID := component.NewID(testType)
		cfgs := map[component.ID]component.Config{testID: defaultCfg}

		b := builders.NewExtension(cfgs, factories, component.StabilityLevelBeta)
		_, err = b.Create(context.Background(), createExtensionSettings(testID))
		require.ErrorContains(t, err, `component is too unstable for use: stability level is Development versus required Beta`)

		b = builders.NewExtension(cfgs, factories, component.StabilityLevelDevelopment)
		e, err := b.Create(context.Background(), createExtensionSettings(testID))
		require.NoError(t, err)
		assert.NotNil(t, e)
	})
}
