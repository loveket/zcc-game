package utils

// type Message struct {
// }
// type pool struct {
// 	TaskQueue []chan Message
// }

// func (p *pool) StartWorkerPool() {
// 	for i := 0; i < 100; i++ {
// 		p.TaskQueue[i] = make(chan Message, 4096)
// 		//启动当前worker,阻塞等待消息从channel传入
// 		//go p.Worker(i, p.TaskQueue[i])
// 	}
// }
