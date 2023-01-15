package metrics

import (
	"strconv"
)

type gauge float64

func (g gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 6, 64)
}

func (g gauge) TypeName() string {
	return "gauge"
}

func (g gauge) StringValue() string {
	return g.String()
}

type counter int64

func (c counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (c counter) TypeName() string {
	return "counter"
}

func (c counter) StringValue() string {
	return c.String()
}

type Metrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
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
		{"HeapIdle", metrics.HeapIdle.TypeName(), metrics.HeapAlloc.StringValue()},
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
