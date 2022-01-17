package taskrunner

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	d := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i
			log.Printf("Dispatcher sent: %v", i)
		}

		return nil
	}

	e := func(dc dataChan) error {
		// 如果这么写的话，for循坏 break 下面的代码就可以继续执行，从而退出程序（死循坏）
	forshendu:
		for {
			select {
			case d := <-dc:
				log.Printf("Executor received: %v", d)
			default:
				break forshendu
			}
		}
		fmt.Println("执行函数执行")
		// 执行结束时要返回错误终止程序
		return errors.New("executor")
		//return nil

	}
	// 任务通过这个函数，初始化并交给后台调度
	runner := NewRunner(30, false, d, e)
	go runner.StartAll()
	time.Sleep(3 * time.Second)
}
