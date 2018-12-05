package pgtune

import (
	"fmt"
	"math"
)

// Keys in the conf file that are tuned related to parallelism
const (
	MaxWorkerProcessesKey       = "max_worker_processes"
	MaxParallelWorkersGatherKey = "max_parallel_workers_per_gather"
	MaxParallelWorkers          = "max_parallel_workers" // pg10+

	errOneCPU = "cannot make recommendations with just 1 CPU"
)

// ParallelLabel is the label used to refer to the parallelism settings group
const ParallelLabel = "parallelism"

// ParallelKeys is an array of keys that are tunable for parallelism
var ParallelKeys = []string{
	MaxWorkerProcessesKey,
	MaxParallelWorkersGatherKey,
	MaxParallelWorkers,
}

// ParallelRecommender gives recommendations for ParallelKeys based on system resources
type ParallelRecommender struct {
	cpus int
}

// NewParallelRecommender returns a ParallelRecommender that recommends based on the given
// number of cpus.
func NewParallelRecommender(cpus int) *ParallelRecommender {
	return &ParallelRecommender{cpus}
}

// IsAvailable returns whether this Recommender is usable given the system resources. True when number of CPUS > 1.
func (r *ParallelRecommender) IsAvailable() bool {
	return r.cpus > 1
}

// Recommend returns the recommended PostgreSQL formatted value for the conf
// file for a given key.
func (r *ParallelRecommender) Recommend(key string) string {
	var val string
	if r.cpus <= 1 {
		panic(errOneCPU)
	}
	if key == MaxWorkerProcessesKey || key == MaxParallelWorkers {
		val = fmt.Sprintf("%d", r.cpus)
	} else if key == MaxParallelWorkersGatherKey {
		val = fmt.Sprintf("%d", int(math.Round(float64(r.cpus)/2.0)))
	} else {
		panic(fmt.Sprintf("unknown key: %s", key))
	}
	return val
}

// ParallelSettingsGroup is the SettingsGroup to represent parallelism settings.
type ParallelSettingsGroup struct {
	pgVersion string
	cpus      int
}

// Label should always return the value ParallelLabel.
func (sg *ParallelSettingsGroup) Label() string { return ParallelLabel }

// Keys should always return the ParallelKeys slice.
func (sg *ParallelSettingsGroup) Keys() []string {
	if sg.pgVersion == "9.6" {
		return ParallelKeys[:2]
	}
	return ParallelKeys
}

// GetRecommender should return a new ParallelRecommender.
func (sg *ParallelSettingsGroup) GetRecommender() Recommender {
	return NewParallelRecommender(sg.cpus)
}
