package otelcol

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/featuregate"
)

type DefaultStabilitySettings struct {
	Components   component.StabilityLevel
	FeatureGates featuregate.Stage
}
