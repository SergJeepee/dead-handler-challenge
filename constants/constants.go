package constants

import "net/http"

const (
	MillisMultilayer    = 1e6
	IterationConfName   = "iterations"
	PoolSizeConfName    = "poolsize"
	UrlConfName         = "url"
	MethodConfName      = "method"
	PayloadConfName     = "payload"
	ContentTypeConfName = "contentType"
)

var (
	AllowedHttpMethods = []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut}
)
