module Client

go 1.18

replace utills => ../utills

require redisTool v1.0.0

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/garyburd/redigo v1.6.4 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
)

replace redisTool => ../util/redisTool
