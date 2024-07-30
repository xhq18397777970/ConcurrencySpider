package Concurrency_Spider

import (
	"fmt"
	"strconv"
	"sync"
)

var maxWorker int = 4

func Start() {

	ch := make(chan string, 10)
	for i := 0; i < 10; i++ {
		ch <- "xhq" + strconv.Itoa(i)
	}
	close(ch)
	var wg sync.WaitGroup
	for i := 0; i < maxWorker; i++ {
		wg.Add(1)
		go worker(&wg, ch)
	}

	wg.Wait()
}

// ch <-chan string 定义了一个只能发送数据的带缓冲区的管道
func worker(wg *sync.WaitGroup, ch chan string) {
	for task := range ch {
		fmt.Printf("任务 %s 已完成\n", task)
	}
	wg.Done()
}
