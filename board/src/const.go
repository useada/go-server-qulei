package main

const (
	COUNT_COMM_CACHE = 200

	TTL_ZSET_CRITICAL = 1            // 2 * (redis read/write timeout 500ms)
	TTL_ZSET_KEY      = 3600 * 9     //
	TTL_HASH_KEY      = 3600 * 9 * 3 //
)
