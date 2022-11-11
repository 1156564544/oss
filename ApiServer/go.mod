module ApiServer

go 1.18

require github.com/klauspost/reedsolomon v1.11.1

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/klauspost/cpuid/v2 v2.1.1 // indirect
	golang.org/x/sys v0.0.0-20220704084225-05e143d24a9e // indirect
)

require (
	es v1.0.0
	httpTool v1.0.0
	redisTool v1.0.0
	rs v1.0.0
)

replace (
	es => ../util/es
	httpTool => ../util/httpTool
	redisTool => ../util/redisTool
	rs => ../util/rs

)
