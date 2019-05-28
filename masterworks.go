package masterworks

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

type setting struct {
	ProcessFunc func()
	WorkFunc func()
	MasterFunc func() bool
	works int
	chanCap int // 通道容量
	ProcessNumber uint64 // 进度
	DataChan chan interface{}
}

func New(works int, chanCap int) *setting {
	r := &setting{works: works, chanCap: chanCap}
	r.DataChan = make(chan interface{}, chanCap)
	return r
}

// 发送数据
func (this *setting) SendData(v interface{}) {
	this.DataChan <- v
	atomic.AddUint64(&this.ProcessNumber, 1)
}

func (this *setting) Run() {
	var abortProcess bool

	go func() {
		for !abortProcess {
			this.ProcessFunc()
		}
	}()

	// 消费者
	wg := sync.WaitGroup{}
	wg.Add(this.works)
	for i := 0; i < this.works; i++ {
		go func() {
			defer wg.Done()
			this.WorkFunc()
			fmt.Println("Work exit!")
		}()
	}

	signChan := make(chan os.Signal, 2)

	// 生产者
	wgMaster := sync.WaitGroup{}
	wgMaster.Add(1)
	go func() {
		for !abortProcess {
			if this.MasterFunc() == false {
				break
			}
		}
		abortProcess = true
		fmt.Println("Master exit!")
		signChan <- os.Kill
		wgMaster.Done()
	}()

	// 收到信号以后，让生产者停止生产！
	signal.Notify(signChan, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	<-signChan
	// signal.Stop(signChan)

	// 终止生产者，并且等待返回
	abortProcess = true
	wgMaster.Wait()

	// 终止工作协程，等待返回
	close(this.DataChan)
	wg.Wait()

	fmt.Println("done")
}
