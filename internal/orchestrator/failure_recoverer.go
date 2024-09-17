package orchestrator

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/thirdweb-dev/indexer/internal/common"
	"github.com/thirdweb-dev/indexer/internal/worker"
)

const DEFAULT_FAILURES_PER_POLL = 10
const DEFAULT_FAILURE_TRIGGER_INTERVAL = 1000

type FailureRecoverer struct {
	failuresPerPoll     int
	triggerIntervalMs   int
	orchestratorStorage OrchestratorStorage
	rpc                 common.RPC
}

func NewFailureRecoverer(rpc common.RPC, orchestratorStorage OrchestratorStorage) *FailureRecoverer {
	failuresPerPoll, err := strconv.Atoi(os.Getenv("FAILURES_PER_POLL"))
	if err != nil || failuresPerPoll == 0 {
		failuresPerPoll = DEFAULT_FAILURES_PER_POLL
	}
	triggerInterval, err := strconv.Atoi(os.Getenv("FAILURE_TRIGGER_INTERVAL"))
	if err != nil || triggerInterval == 0 {
		triggerInterval = DEFAULT_FAILURE_TRIGGER_INTERVAL
	}
	return &FailureRecoverer{
		triggerIntervalMs:   triggerInterval,
		failuresPerPoll:     failuresPerPoll,
		orchestratorStorage: orchestratorStorage,
		rpc:                 rpc,
	}
}

func (fr *FailureRecoverer) Start() error {
	interval := time.Duration(fr.triggerIntervalMs) * time.Millisecond
	ticker := time.NewTicker(interval)

	go func() error {
		for t := range ticker.C {
			fmt.Println("Failure Recovery running at", t)

			blockFailures, err := fr.orchestratorStorage.GetBlockFailures()

			if err != nil {
				log.Printf("Failed to get block failures: %s", err)
				continue
			}

			log.Printf("Triggering workers for %d block failures", len(blockFailures))

			var wg sync.WaitGroup
			for _, failure := range blockFailures {
				wg.Add(1)
				go func(failure BlockFailure) {
					defer wg.Done()
					fr.triggerWorker(failure.BlockNumber)
				}(failure)
			}
			wg.Wait()
		}
		return nil
	}()

	// Keep the program running (otherwise it will exit)
	select {}
}

func (fr *FailureRecoverer) triggerWorker(blockNumber uint64) {
	worker := worker.NewWorker(fr.rpc, blockNumber)
	err := worker.FetchData()
	if err != nil {
		log.Printf("Error retrying block %d: %v", blockNumber, err)
	} else {
		log.Printf("Successfully retried block %d", blockNumber)
	}
}
