package rpc

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"strconv"
	"testing"
	"time"
	"xiuianserver/systemfunc"
)

//	func TestName(t *testing.T) {
//		conn, err := jsonrpc.Dial("tcp", ":7001")
//		if err != nil {
//			log.Fatalln("dailing error: ", err)
//		}
//		req := ArithRequest{9, 2}
//		var res ArithResponse
//		err = conn.Call("Arith.Multiply", req, &res)
//		if err != nil {
//			log.Fatalln("arith error: ", err)
//		}
//		fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Mul)
//		err = conn.Call("Arith.Divide", req, &res)
//		if err != nil {
//			log.Fatalln("arith error: ", err)
//		}
//		fmt.Printf("%d / %d= %d\n", req.A, req.B, res.Div)
//	}
func TestSystemNotify(t *testing.T) {
	conn, err := jsonrpc.Dial("tcp", ":7001")
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}
	tick := time.NewTicker(5 * time.Second)
	min := 65
	defer tick.Stop()
	for {
		select {
		case <-time.After(1 * time.Minute):
			return
		case <-tick.C:
			min -= 5
			if min <= 0 {
				return
			}
			req := systemfunc.SystemMessageRequest{"系统将在" + strconv.Itoa(min) + "秒后更新。。。"}
			var res systemfunc.SystemMessageResponse
			err = conn.Call("SystemFunc.Broadcast", req, &res)
			if err != nil {
				log.Fatalln("SystemFunc error: ", err)
			}
			fmt.Println(res.Status)
		}
	}

}
