package collector

import (
	"fmt"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
)

// PsUtilCollector is a collector for gopsutil metrics.
type PsUtilCollector struct{}

var _ collectWorker = (*PsUtilCollector)(nil)

// NewPsUtilCollector creates new PsUtilCollector.
func NewPsUtilCollector(pollInterval time.Duration, logger *logrus.Logger) *Collector {
	return &Collector{
		collectWorker: &PsUtilCollector{},
		name:          "PsUtilCollector",
		pollInterval:  pollInterval,
		logger:        logger,
	}
}

func (pc *PsUtilCollector) getMetrics() (metric.Metrics, error) {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	resultMetircs := metric.Metrics{
		"TotalMemory": metric.Gauge(memStats.Total),
		"FreeMemory":  metric.Gauge(memStats.Free),
	}

	err = pc.updateWithCPUUtilization(resultMetircs)
	if err != nil {
		return nil, err
	}

	return resultMetircs, nil
}

func (pc *PsUtilCollector) updateWithCPUUtilization(resultMetrics metric.Metrics) error {
	calcUtilization := func(ts cpu.TimesStat) float64 {
		totalTime := ts.User + ts.System + ts.Idle + ts.Nice + ts.Iowait + ts.Irq + ts.Softirq + ts.Steal + ts.Guest + ts.GuestNice
		return (1 - ts.Idle/totalTime) * 100
	}

	cpusTimes, err := cpu.Times(true)
	if err != nil {
		return err
	}

	for i, cpu := range cpusTimes {
		resultMetrics[fmt.Sprintf("CPUutilization%d", i)] = metric.Gauge(calcUtilization(cpu))
	}

	return nil
}

func (pc *PsUtilCollector) collectCleanup() {}
