// Code generated by mdatagen. DO NOT EDIT.

package metadatatest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/service/internal/metadata"
)

func TestSetupTelemetry(t *testing.T) {
	testTel := componenttest.NewTelemetry()
	tb, err := metadata.NewTelemetryBuilder(testTel.NewTelemetrySettings())
	require.NoError(t, err)
	defer tb.Shutdown()
	require.NoError(t, tb.RegisterProcessCPUSecondsCallback(func(_ context.Context, observer metric.Float64Observer) error {
		observer.Observe(1)
		return nil
	}))
	require.NoError(t, tb.RegisterProcessMemoryRssCallback(func(_ context.Context, observer metric.Int64Observer) error {
		observer.Observe(1)
		return nil
	}))
	require.NoError(t, tb.RegisterProcessRuntimeHeapAllocBytesCallback(func(_ context.Context, observer metric.Int64Observer) error {
		observer.Observe(1)
		return nil
	}))
	require.NoError(t, tb.RegisterProcessRuntimeTotalAllocBytesCallback(func(_ context.Context, observer metric.Int64Observer) error {
		observer.Observe(1)
		return nil
	}))
	require.NoError(t, tb.RegisterProcessRuntimeTotalSysMemoryBytesCallback(func(_ context.Context, observer metric.Int64Observer) error {
		observer.Observe(1)
		return nil
	}))
	require.NoError(t, tb.RegisterProcessUptimeCallback(func(_ context.Context, observer metric.Float64Observer) error {
		observer.Observe(1)
		return nil
	}))
	tb.ConnectorConsumedItems.Add(context.Background(), 1)
	tb.ConnectorConsumedSize.Add(context.Background(), 1)
	tb.ConnectorProducedItems.Add(context.Background(), 1)
	tb.ConnectorProducedSize.Add(context.Background(), 1)
	tb.ExporterConsumedItems.Add(context.Background(), 1)
	tb.ExporterConsumedSize.Add(context.Background(), 1)
	tb.ProcessorConsumedItems.Add(context.Background(), 1)
	tb.ProcessorConsumedSize.Add(context.Background(), 1)
	tb.ProcessorProducedItems.Add(context.Background(), 1)
	tb.ProcessorProducedSize.Add(context.Background(), 1)
	tb.ReceiverProducedItems.Add(context.Background(), 1)
	tb.ReceiverProducedSize.Add(context.Background(), 1)
	AssertEqualConnectorConsumedItems(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualConnectorConsumedSize(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualConnectorProducedItems(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualConnectorProducedSize(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualExporterConsumedItems(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualExporterConsumedSize(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessCPUSeconds(t, testTel,
		[]metricdata.DataPoint[float64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessMemoryRss(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessRuntimeHeapAllocBytes(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessRuntimeTotalAllocBytes(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessRuntimeTotalSysMemoryBytes(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessUptime(t, testTel,
		[]metricdata.DataPoint[float64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessorConsumedItems(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessorConsumedSize(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessorProducedItems(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualProcessorProducedSize(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualReceiverProducedItems(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())
	AssertEqualReceiverProducedSize(t, testTel,
		[]metricdata.DataPoint[int64]{{Value: 1}},
		metricdatatest.IgnoreTimestamp())

	require.NoError(t, testTel.Shutdown(context.Background()))
}
