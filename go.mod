module github.com/tada-team/tdclient

go 1.14

require (
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/manifoldco/promptui v0.7.0
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/tada-team/kozma v1.1.0
	github.com/tada-team/tdproto v1.2.8
)

//replace github.com/tada-team/tdproto v0.3.0 => ../tdproto
