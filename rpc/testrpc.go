package rpc

import "errors"

type Arith struct {
}
type ArithRequest struct {
	A int
	B int
}

type ArithResponse struct {
	Mul int
	Div int
}

func (a *Arith) Multiply(req ArithRequest, res *ArithResponse) error {
	res.Mul = req.A * req.B
	return nil
}

// 除法运算方法
func (a *Arith) Divide(req ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("divide by zero")
	}
	res.Div = req.A / req.B
	return nil
}
