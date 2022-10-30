module ApiServer

go 1.16

require (
	httpTool v1.0.0
	redisTool v1.0.0
	es v1.0.0
)

replace (
	httpTool => ../util/httpTool
	redisTool => ../util/redisTool
	es => ../util/es
)
