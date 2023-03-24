package log

import "testing"

func TestLogger(t *testing.T) {
	GetLogger()

	for i := 0; i < 10; i++ {
		LoggerSingle.Info("我挺重要的")
		LoggerSingle.Warning("我好像有问题的")
		LoggerSingle.Error("我好像有错误的")
	}
	LoggerSingle.StopLogger()
}
