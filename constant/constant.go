package constant

const DEBUG = false

//const name  =

type ContentType = int
type HttpMethod = int
type HttpParamType = int
type HttpParamFrom = int

const (
	BodyXFormType = 1
	BodyJsonType  = 2
	BodyFormType  = 3
)

const (
	Get  = 1
	Post = 2
)
const (
	BODY   = 1
	HEADER = 2
)

const (
	DB    = 1
	Value = 2
)

const CMSPIP = "172.22.16.138"
const CMSPPort = 1216

const MongoIp = "172.22.16.144:27017"
