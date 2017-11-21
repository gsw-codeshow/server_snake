package snake

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nsf/termbox-go"
)

func TestSnake(t *testing.T) {

	context := NewContext()
	context.AddSnake("192.168.1.1", 10, 10)
	fmt.Println(context)
	context.Update()
	fmt.Println(context)
	jsonByte, jsonByteErr := json.Marshal(context)
	if nil != jsonByteErr {
		fmt.Println(jsonByteErr)
	}
	fmt.Println(string(jsonByte))
}

func TestSnakeGrow(t *testing.T) {
	termbox.Init()
	context := NewContext()
	context.AddSnake("192.168.1.1", 10, 10)
	fmt.Println(context)
	context.Update()

	for i := 0; i < 10; i++ {
		context.Update()
		fmt.Println(context)
	}
}
