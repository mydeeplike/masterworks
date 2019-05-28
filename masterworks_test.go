package masterworks

import (
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {

	arr := []int {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	obj := New(2, 4)
	obj.ProcessFunc = func() {
		fmt.Printf("num: %v\r", obj.ProcessNumber)
		time.Sleep(time.Second)
	}
	obj.WorkFunc = func() {
		for v := range obj.DataChan {
			i := v.(int)
			fmt.Printf("i: %v\n", i)
		}
	}
	obj.MasterFunc = func() bool {
		for _, v := range arr {
			obj.SendData(v)
		}
		return false
	}

	obj.Run()

}
