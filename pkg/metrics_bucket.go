package pkg

import (
	"math"
)

type metricsAggregate struct {
	min   float32
	max   float32
	sum   float32
	count int
}

func newMetricsAggregate() *metricsAggregate {
	return &metricsAggregate{
		min:   math.MaxFloat32,
		max:   0,
		sum:   0,
		count: 0,
	}
}

type MetricsBucket struct {
	stations map[string]*metricsAggregate
}

func NewMetricsBucket() *MetricsBucket {
	return &MetricsBucket{
		stations: map[string]*metricsAggregate{},
	}
}

func (m *MetricsBucket) GetMetrics() []*StationMetrics {

	metrics := make([]*StationMetrics, len(m.stations))
	index := 0

	for station, aggregate := range m.stations {
		metrics[index] = &StationMetrics{
			Station: station,
			Min:     aggregate.min,
			Mean:    aggregate.sum / float32(aggregate.count),
			Max:     aggregate.max,
		}

		index++
	}

	return metrics
}

func (m *MetricsBucket) Aggregate(batch map[string]*metricsAggregate) {
	for station, aggregate := range batch {
		value, ok := m.stations[station]

		if !ok {
			m.stations[station] = aggregate
			continue
		}

		value.min = min(value.min, aggregate.min)
		value.max = max(value.max, aggregate.max)
		value.sum += aggregate.sum
		value.count += aggregate.count
	}
}

func aggregateFromMsg(msg []*valueMessage) map[string]*metricsAggregate {
	localMetricsAggregate := map[string]*metricsAggregate{}

	for _, value := range msg {
		aggregate, ok := localMetricsAggregate[value.station]
		if !ok {
			aggregate = newMetricsAggregate()
			localMetricsAggregate[value.station] = aggregate
		}

		aggregate.min = min(aggregate.min, value.value)
		aggregate.max = max(aggregate.max, value.value)
		aggregate.sum += value.value
		aggregate.count++
	}

	return localMetricsAggregate
}
