package log

import (
	"strconv"
	"sync"
	"testing"
)

func TestLog(t *testing.T) {
	wg := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, flag int) {
			defer wg.Done()
			log := NewLog("zcc" + strconv.Itoa(flag))
			log.InfoMSG("我有重要事件")
			log.WarnMSG("我有警告事件")
			log.ErrorMSG("我有错误事件")
		}(wg, i)

	}
	wg.Wait()
}
