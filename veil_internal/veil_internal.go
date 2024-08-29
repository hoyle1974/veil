// This package is for internal use by the veil library and veil generated code
package veil_internal

type ConnectionFactory interface {
	GetConnection() any
}

type ClientInit func(factory ConnectionFactory)
type ServerInit func()

var ClientInits = []ClientInit{}
var ServerInits = []ServerInit{}

// Register code to be run when client library is initialized
func RegisterClientInit(i ClientInit) {
	ClientInits = append(ClientInits, i)
}

// Register code to be run when server library is intiialized
func RegisterServerInit(i ServerInit) {
	ServerInits = append(ServerInits, i)
}

var ServiceRegistry = []serviceRegistry{}

type serviceRegistry interface {
	RPC_Bind_Service(service any) error
}

// Used to register generated RPC services
func RegisterService(service serviceRegistry) {
	ServiceRegistry = append(ServiceRegistry, service)
}

var RemoteImpl = []any{}

// Registers remote implementations for client use
func RegisterRemoteImpl(service any) {
	RemoteImpl = append(RemoteImpl, service)
}
