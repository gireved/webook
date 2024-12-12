package channel

import (
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	// 声明了一个放int类型的channel,并没有初始化
	/*var ch chan int
	ch <- 123
	val := <-ch
	println(val)*/
	// 放空结构体，一般用来做信号
	//var chv1 chan struct{}

	/*// 不带容量的要小心操作
	ch1 := make(chan int)
	// 带容量的*/
	ch2 := make(chan int, 2)
	ch2 <- 123
	// 关闭ch2
	// 不能再发送，但是可以读取

	close(ch2)
	val, ok := <-ch2
	println(val, ok)

}

func TestChannelClose(t *testing.T) {
	ch := make(chan int, 2)
	ch <- 123
	ch <- 234
	val, ok := <-ch
	t.Log(val, ok)
	close(ch)
	//ch <- 1234
	val1, ok1 := <-ch
	t.Log(val1, ok1)
}

func TestLoopChannel(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
			time.Sleep(time.Millisecond * 10)
		}
	}()
	for val := range ch {
		t.Log(val)
	}
}
