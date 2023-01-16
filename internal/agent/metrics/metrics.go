package metrics

import (
	"github.com/devldavydov/promytheus/internal/common/types"
)

type Metrics struct {
	Alloc         types.Gauge
	BuckHashSys   types.Gauge
	Frees         types.Gauge
	GCCPUFraction types.Gauge
	GCSys         types.Gauge
	HeapAlloc     types.Gauge
	HeapIdle      types.Gauge
	HeapInuse     types.Gauge
	HeapObjects   types.Gauge
	HeapReleased  types.Gauge
	HeapSys       types.Gauge
	LastGC        types.Gauge
	Lookups       types.Gauge
	MCacheInuse   types.Gauge
	MCacheSys     types.Gauge
	MSpanInuse    types.Gauge
	MSpanSys      types.Gauge
	Mallocs       types.Gauge
	NextGC        types.Gauge
	NumForcedGC   types.Gauge
	NumGC         types.Gauge
	OtherSys      types.Gauge
	PauseTotalNs  types.Gauge
	StackInuse    types.Gauge
	StackSys      types.Gauge
	Sys           types.Gauge
	TotalAlloc    types.Gauge
	PollCount     types.Counter
	RandomValue   types.Gauge
}

type MetricsItem struct {
	metricName string
	typeName   string
	value      string
}

func (metrics Metrics) ToItemsList() []MetricsItem {
	return []MetricsItem{
		{"Alloc", metrics.Alloc.TypeName(), metrics.Alloc.StringValue()},
		{"BuckHashSys", metrics.BuckHashSys.TypeName(), metrics.BuckHashSys.StringValue()},
		{"Frees", metrics.Frees.TypeName(), metrics.Frees.StringValue()},
		{"GCCPUFraction", metrics.GCCPUFraction.TypeName(), metrics.GCCPUFraction.StringValue()},
		{"GCSys", metrics.GCSys.TypeName(), metrics.GCSys.StringValue()},
		{"HeapAlloc", metrics.HeapAlloc.TypeName(), metrics.HeapAlloc.StringValue()},
		{"HeapIdle", metrics.HeapIdle.TypeName(), metrics.HeapIdle.StringValue()},
		{"HeapInuse", metrics.HeapInuse.TypeName(), metrics.HeapInuse.StringValue()},
		{"HeapObjects", metrics.HeapObjects.TypeName(), metrics.HeapObjects.StringValue()},
		{"HeapReleased", metrics.HeapReleased.TypeName(), metrics.HeapReleased.StringValue()},
		{"HeapSys", metrics.HeapSys.TypeName(), metrics.HeapSys.StringValue()},
		{"LastGC", metrics.LastGC.TypeName(), metrics.LastGC.StringValue()},
		{"Lookups", metrics.Lookups.TypeName(), metrics.Lookups.StringValue()},
		{"MCacheInuse", metrics.MCacheInuse.TypeName(), metrics.MCacheInuse.StringValue()},
		{"MCacheSys", metrics.MCacheSys.TypeName(), metrics.MCacheSys.StringValue()},
		{"MSpanInuse", metrics.MSpanInuse.TypeName(), metrics.MSpanInuse.StringValue()},
		{"MSpanSys", metrics.MSpanSys.TypeName(), metrics.MSpanSys.StringValue()},
		{"Mallocs", metrics.Mallocs.TypeName(), metrics.Mallocs.StringValue()},
		{"NextGC", metrics.NextGC.TypeName(), metrics.NextGC.StringValue()},
		{"NumForcedGC", metrics.NumForcedGC.TypeName(), metrics.NumForcedGC.StringValue()},
		{"NumGC", metrics.NumGC.TypeName(), metrics.NumGC.StringValue()},
		{"OtherSys", metrics.OtherSys.TypeName(), metrics.OtherSys.StringValue()},
		{"PauseTotalNs", metrics.PauseTotalNs.TypeName(), metrics.PauseTotalNs.StringValue()},
		{"StackInuse", metrics.StackInuse.TypeName(), metrics.StackInuse.StringValue()},
		{"StackSys", metrics.StackSys.TypeName(), metrics.StackSys.StringValue()},
		{"Sys", metrics.Sys.TypeName(), metrics.Sys.StringValue()},
		{"TotalAlloc", metrics.TotalAlloc.TypeName(), metrics.TotalAlloc.StringValue()},
		{"PollCount", metrics.PollCount.TypeName(), metrics.PollCount.StringValue()},
		{"RandomValue", metrics.RandomValue.TypeName(), metrics.RandomValue.StringValue()},
	}
}
