package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	gc "xiuianserver/config"
	"xiuianserver/systemfunc"
)

var RpcServer *Rpc

func init() {
	RpcServer = NewRpc()
	//RpcServer.RpcMap["Arith"] = &Arith{}
	RpcServer.RpcMap["SystemFunc"] = &systemfunc.SystemFunc{}
}

type Rpc struct {
	IpAddr string
	RpcMap map[interface{}]interface{}
}

func NewRpc() *Rpc {
	return &Rpc{
		IpAddr: gc.GlobalConfig.RpcConfig.IpAddr,
		RpcMap: make(map[interface{}]interface{}),
	}
}

func (r *Rpc) registers() bool {
	for _, structVal := range r.RpcMap {
		if err := rpc.Register(structVal); err != nil {
			return false
		}
	}
	return true
}
func (r *Rpc) Start() {
	if !r.registers() {
		fmt.Println("rpc register failed")
		return
	}
	lis, err := net.Listen("tcp", r.IpAddr)
	if err != nil {
		fmt.Println("start rpc err", err)
		return
	}
	fmt.Println("[rpc already start]...")
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go func(conn net.Conn) {
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}
