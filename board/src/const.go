package main

const (
	TIMESTAMP_INF = 4102415999000

	COUNT_COMM_CACHE = 200

	LEVEL_COMM_FIRST = 1
	LEVEL_COMM_CHILD = 2

	TTL_ZSET_CRITICAL = 1                // 2 * (redis read/write timeout 500ms)
	TTL_ZSET_KEY      = 3600 * 9         //
	TTL_HASH_KEY      = TTL_ZSET_KEY * 3 //
)
