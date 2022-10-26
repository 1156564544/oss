module ApiServer

go 1.16

require (
	redisTool v1.0.0
)

replace (
	redisTool => ../DataServer/redisTool
)
