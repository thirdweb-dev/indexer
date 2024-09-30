package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Committer Metrics
var (
	SuccessfulCommits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "committer_successful_commits_total",
		Help: "The total number of successful block commits",
	})

	LastCommittedBlock = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "committer_last_committed_block",
		Help: "The last successfully committed block number",
	})

	GapCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "committer_gap_counter",
		Help: "The number of gaps detected during commits",
	})

	MissedBlockNumbers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "committer_first_missed_block_number",
		Help: "The first blocknumber detected in a commit gap",
	})
)

// Worker Metrics
var LastFetchedBlock = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "worker_last_fetched_block_from_rpc",
	Help: "The last block number fetched by the worker from the RPC",
})

// Poller metrics
var (
	PolledBatchSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "polled_batch_size",
		Help: "The number of blocks polled in a single batch",
	})
)

// Failure Recoverer Metrics
var (
	FailureRecovererLastTriggeredBlock = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "failure_recoverer_last_triggered_block",
		Help: "The last block number that the failure recoverer was triggered for",
	})

	FirstBlocknumberInFailureRecovererBatch = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "failure_recoverer_first_block_in_batch",
		Help: "The first block number in the failure recoverer batch",
	})
)