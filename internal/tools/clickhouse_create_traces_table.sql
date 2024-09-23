CREATE TABLE base.traces (
    `id` UUID,
    `chain_id` UInt256,
    `block_number` UInt256,
    `block_hash` FixedString(66),
    `block_timestamp` UInt64 CODEC(Delta, ZSTD),
    `transaction_hash` FixedString(66),
    `transaction_index` UInt64,
    `call_type` String,
    `error` Nullable(String),
    `from_address` FixedString(42),
    `to_address` FixedString(42),
    `gas` UInt128,
    `gas_used` UInt128,
    `input` String,
    `output` Nullable(String),
    `subtraces` UInt64,
    `trace_address` String,
    `trace_type` String,
    `value` UInt256,
    `is_deleted` UInt8 DEFAULT 0,
    `insert_timestamp` DateTime DEFAULT now(),
    INDEX hash_idx transaction_hash TYPE bloom_filter GRANULARITY 1,
    INDEX to_address_idx to_address TYPE bloom_filter GRANULARITY 1,
    INDEX from_address_idx from_address TYPE bloom_filter GRANULARITY 1,
) ENGINE = SharedReplacingMergeTree(
    '/clickhouse/tables/{uuid}/{shard}',
    '{replica}',
    insert_timestamp,
    is_deleted
)
ORDER BY (block_number) SETTINGS index_granularity = 8192
SETTINGS allow_experimental_replacing_merge_with_cleanup = 1;