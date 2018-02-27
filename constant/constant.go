package constant

const DEBUG bool = true

//const name  =

type ContentType = int
type HttpMethod = int
type HttpParamType = int
type HttpParamFrom = int

const (
	BODY_XFORM_TYPE = 1
	BODY_JSON_TYPE  = 2
	BODY_FORM_TYPE  = 3
)

const (
	GET  = 1
	POST = 2
)
const (
	BODY   = 1
	HEADER = 2
)

const (
	DB    = 1
	VALUE = 2
)

const CMSP_IP = "172.22.16.138"
const CMSP_PORT = 1216
