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
		{"Alloc", metrics.Alloc.TypeName(), metrics.Alloc.String()},
		{"BuckHashSys", metrics.BuckHashSys.TypeName(), metrics.BuckHashSys.String()},
		{"Frees", metrics.Frees.TypeName(), metrics.Frees.String()},
		{"GCCPUFraction", metrics.GCCPUFraction.TypeName(), metrics.GCCPUFraction.String()},
		{"GCSys", metrics.GCSys.TypeName(), metrics.GCSys.String()},
		{"HeapAlloc", metrics.HeapAlloc.TypeName(), metrics.HeapAlloc.String()},
		{"HeapIdle", metrics.HeapIdle.TypeName(), metrics.HeapIdle.String()},
		{"HeapInuse", metrics.HeapInuse.TypeName(), metrics.HeapInuse.String()},
		{"HeapObjects", metrics.HeapObjects.TypeName(), metrics.HeapObjects.String()},
		{"HeapReleased", metrics.HeapReleased.TypeName(), metrics.HeapReleased.String()},
		{"HeapSys", metrics.HeapSys.TypeName(), metrics.HeapSys.String()},
		{"LastGC", metrics.LastGC.TypeName(), metrics.LastGC.String()},
		{"Lookups", metrics.Lookups.TypeName(), metrics.Lookups.String()},
		{"MCacheInuse", metrics.MCacheInuse.TypeName(), metrics.MCacheInuse.String()},
		{"MCacheSys", metrics.MCacheSys.TypeName(), metrics.MCacheSys.String()},
		{"MSpanInuse", metrics.MSpanInuse.TypeName(), metrics.MSpanInuse.String()},
		{"MSpanSys", metrics.MSpanSys.TypeName(), metrics.MSpanSys.String()},
		{"Mallocs", metrics.Mallocs.TypeName(), metrics.Mallocs.String()},
		{"NextGC", metrics.NextGC.TypeName(), metrics.NextGC.String()},
		{"NumForcedGC", metrics.NumForcedGC.TypeName(), metrics.NumForcedGC.String()},
		{"NumGC", metrics.NumGC.TypeName(), metrics.NumGC.String()},
		{"OtherSys", metrics.OtherSys.TypeName(), metrics.OtherSys.String()},
		{"PauseTotalNs", metrics.PauseTotalNs.TypeName(), metrics.PauseTotalNs.String()},
		{"StackInuse", metrics.StackInuse.TypeName(), metrics.StackInuse.String()},
		{"StackSys", metrics.StackSys.TypeName(), metrics.StackSys.String()},
		{"Sys", metrics.Sys.TypeName(), metrics.Sys.String()},
		{"TotalAlloc", metrics.TotalAlloc.TypeName(), metrics.TotalAlloc.String()},
		{"PollCount", metrics.PollCount.TypeName(), metrics.PollCount.String()},
		{"RandomValue", metrics.RandomValue.TypeName(), metrics.RandomValue.String()},
	}
}
