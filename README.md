# indexer

The Indexer is a blockchain data processing tool designed to fetch, process, and store on-chain data. It provides a solution for indexing blockchain data, facilitating efficient querying of transactions, logs through a simple API.

The Indexer's architecture consists of five main components:
<img width="728" alt="architecture-diagram" src="https://github.com/user-attachments/assets/8ad19477-c18f-442e-a899-0f502a00bae6">

1. **Poller**: The Poller is responsible for continuously fetching new blocks from the configured RPC and processing them. It uses multiple worker goroutines to concurrently retrieve block data, handles successful and failed results, and stores the processed block data and any failures.
   
2. **Worker**: The Worker is responsible for processing batches of block numbers, fetching block data, logs, and traces (if supported) from the configured RPC. It divides the work into chunks, processes them concurrently, and returns the results as a collection of WorkerResult structures, which contain the block data, transactions, logs, and traces for each processed block.

3. **Committer**: The Committer is responsible for periodically moving data from the staging storage to the main storage. It ensures that blocks are committed sequentially, handling any gaps in the data, and updates various metrics while performing concurrent inserts of blocks, logs, transactions, and traces into the main storage.

4. **Failure Recoverer**: The FailureRecoverer is responsible for recovering from block processing failures. It periodically checks for failed blocks, attempts to reprocess them using a worker, and either removes successfully processed blocks from the failure list or updates the failure count for blocks that continue to fail.

5. **Orchestrator**: The Orchestrator is responsible for coordinating and managing the poller, failure recoverer, and committer. It initializes these components based on configuration settings and starts them concurrently, ensuring they run independently while waiting for all of them to complete their tasks.

The Indexer's modular architecture and configuration options allow for adaptation to various evm chains and use cases.


## Getting started

### Pre-requites
1. Install Golang v1.23
2. Run the DB scripts (change the db name)
3. Create a clickhouse instance

### Usage
To run the indexer and the associated API, follow these steps:
1. Clone the repo
```
git clone https://github.com/thirdweb-dev/indexer.git
```
2. Run the migration scripts here
3. Create `config.yml` from `config.example.yml` and set the values by following the [config guide](#supported-configurations)
4. Create `secrets.yml` from `secrects.example.yml` and set the needed credentials
5. Build an instance
```
go build -o main
```
5. Run the indexer
```
./main orchestrator
```
6. Run the Data API
```
./main api
```
7. API is available at `http://localhost:3000`

## Metrics

The indexer node exposes Prometheus metrics at `http://localhost:2112/metrics`. Here the exposed metrics [metrics.go](https://github.com/thirdweb-dev/indexer/blob/main/internal/metrics/metrics.go)

## Configuration

You can configure the application in 3 ways.
The order of priority is
1. Command line arguments
2. Environment variables
3. Configuration files

### Configuration using command line arguments
You can configure the application using command line arguments.
For example to configure the `rpc.url` configuration, you can use the `--rpc-url` command line argument.
Only select arguments are implemented. More on this below.

### Configuration using environment variables
You can also configure the application using environment variables. You can configure any configuration in the `config.yml` file using environment variables by replacing the `.` in the configuration with `_` and making the variable name uppercase.  
For example to configure the `rpc.url` configuration to `https://my-rpc.com`, you can set the `RPC_URL` environment variable to `https://my-rpc.com`.  

### Configuration using configuration files
The default configuration should live in `configs/config.yml`. Copy `configs/config.example.yml` to get started.  
Or you can use the `--config` flag to specify a different configuration file.  
If you want to add secrets to the configuration file, you can copy `configs/secrets.example.yml` to `configs/secrets.yml` and add the secrets. They won't be committed to the repository or the built image.

### Supported configurations:

#### RPC URL
URL to use as the RPC client.

cmd: `--rpc-url`
env: `RPC_URL`
yaml:
```yaml
rpc:
  url: https://rpc.com
```

#### RPC Blocks Per Request
How many blocks at a time to fetch from the RPC. Default is 1000.

cmd: `--rpc-blocks-blocksPerRequest`
env: `RPC_BLOCKS_BLOCKSPERREQUEST`
yaml:
```yaml
rpc:
  blocks:
    blocksPerRequest: 1000
```

#### RPC Blocks Batch Delay
Milliseconds to wait between batches of blocks when fetching from the RPC. Default is 0.

cmd: `--rpc-blocks-batchDelay`
env: `RPC_BLOCKS_BATCHDELAY`
yaml:
```yaml
rpc:
  blocks:
    batchDelay: 100
```

#### RPC Logs Blocks Per Request
How many blocks at a time to query logs for from the RPC. Default is 100.
Has no effect if it's larger than RPC blocks per request.

cmd: `--rpc-logs-blocksPerRequest`
env: `RPC_LOGS_BLOCKSPERREQUEST`
yaml:
```yaml
rpc:
  logs:
    blocksPerRequest: 100
```

#### RPC Logs Batch Delay
Milliseconds to wait between batches of logs when fetching from the RPC. Default is 0.

cmd: `--rpc-logs-batchDelay`
env: `RPC_LOGS_BATCHDELAY`
yaml:
```yaml
rpc:
  logs:
    batchDelay: 100
```

#### RPC Traces Enabled
Whether to enable fetching traces from the RPC. Default is `true`, but it will try to detect if the RPC supports traces automatically.

cmd: `--rpc-traces-enabled`
env: `RPC_TRACES_ENABLED`
yaml:
```yaml
rpc:
  traces:
    enabled: true
```

#### RPC Traces Blocks Per Request
How many blocks at a time to fetch traces for from the RPC. Default is 100.
Has no effect if it's larger than RPC blocks per request.

cmd: `--rpc-traces-blocksPerRequest`
env: `RPC_TRACES_BLOCKSPERREQUEST`
yaml:
```yaml
rpc:
  traces:
    blocksPerRequest: 100
```

#### RPC Traces Batch Delay
Milliseconds to wait between batches of traces when fetching from the RPC. Default is 0.

cmd: `--rpc-traces-batchDelay`
env: `RPC_TRACES_BATCHDELAY`
yaml:
```yaml
rpc:
  traces:
    batchDelay: 100
```

#### Log Level
Log level for the logger. Default is `warn`.

cmd: `--log-level`
env: `LOG_LEVEL`
yaml:
```yaml
log:
  level: debug
```

#### Prettify logs
Whether to print logs in a prettified format. Affects performance. Default is `false`.

cmd: `--log-prettify`
env: `LOG_PRETTIFY`
yaml:
```yaml
log:
  prettify: true
```

#### Poller
Whether to enable the poller. Default is `true`.

cmd: `--poller-enabled`
env: `POLLER_ENABLED`
yaml:
```yaml
poller:
  enabled: true
```

#### Poller Interval
Poller trigger interval in milliseconds. Default is `1000`.

cmd: `--poller-interval`
env: `POLLER_INTERVAL`
yaml:
```yaml
poller:
  interval: 3000
```

#### Poller Blocks Per Poll
How many blocks to poll each interval. Default is `10`.

cmd: `--poller-blocks-per-poll`
env: `POLLER_BLOCKSPERPOLL`
yaml:
```yaml
poller:
  blocksPerPoll: 3
```

#### Poller From Block
From which block to start polling. Default is `0`.

cmd: `--poller-from-block`
env: `POLLER_FROMBLOCK`
yaml:
```yaml
poller:
  fromBlock: 20000000
```

#### Poller Force Start Block
From which block to start polling. Default is `false`.

cmd: `--poller-force-from-block`
env: `POLLER_FORCEFROMBLOCK`
yaml:
```yaml
poller:
  forceFromBlock: false
```

#### Poller Until Block
Until which block to poll. If not set, it will poll until the latest block.

cmd: `--poller-until-block`
env: `POLLER_UNTILBLOCK`
yaml:
```yaml
poller:
  untilBlock: 20000010
```

#### Committer
Whether to enable the committer. Default is `true`.

cmd: `--committer-enabled`
env: `COMMITTER_ENABLED`
yaml:
```yaml
committer:
  enabled: true
```

#### Committer Interval
Committer trigger interval in milliseconds. Default is `250`.

cmd: `--committer-interval`
env: `COMMITTER_INTERVAL`
yaml:
```yaml
committer:
  interval: 3000
```

#### Committer Blocks Per Commit
How many blocks to commit each interval. Default is `10`.

cmd: `--committer-blocks-per-commit`
env: `COMMITTER_BLOCKSPERCOMMIT`
yaml:
```yaml
committer:
  blocksPerCommit: 1000
```

#### Failure Recoverer
Whether to enable the failure recoverer. Default is `true`.

cmd: `--failure-recoverer-enabled`
env: `FAILURERECOVERER_ENABLED`
yaml:
```yaml
failureRecoverer:
  enabled: true
```

#### Failure Recoverer Interval
Failure recoverer trigger interval in milliseconds. Default is `1000`.

cmd: `--failure-recoverer-interval`
env: `FAILURERECOVERER_INTERVAL`
yaml:
```yaml
failureRecoverer:
  interval: 3000
```

#### Failure Recoverer Blocks Per Run
How many blocks to recover each interval. Default is `10`.

cmd: `--failure-recoverer-blocks-per-run`
env: `FAILURERECOVERER_BLOCKSPERRUN`
yaml:
```yaml
failureRecoverer:
  blocksPerRun: 100
```

#### Storage
This application has 3 kinds of strorage. Main, Staging and Orchestrator.
Each of them takes similar configuration depending on the driver you want to use.

There are no defaults, so this needs to be configured.

```yaml
storage:
  main:
    driver: "clickhouse"
    clickhouse:
      host: "localhost"
      port: 3000
      user: "admin"
      password: "admin"
      database: "base"
  staging:
    driver: "memory"
    memory:
      maxItems: 1000000
  orchestrator:
    driver: "memory"
```