module github.com/tada-team/tdclient

go 1.14

require (
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/pkg/errors v0.9.1
	github.com/tada-team/tdproto v0.2.1
)

//replace github.com/tada-team/tdproto v0.0.13 => ../tdproto
