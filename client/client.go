package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"server_snake/snake"
	"time"

	"github.com/nsf/termbox-go"
)

func Update(tcpConReader *bufio.Reader) {
	var context snake.Context
	time.Sleep(snake.Speed)
	msg, _ := tcpConReader.ReadString('\n')
	json.Unmarshal([]byte(msg), &context)
	context.Draw()
}

func main() {
	terminalErr := termbox.Init()
	if nil != terminalErr {
		fmt.Println(terminalErr)
		return
	}
	tcpAddr, tcpAddrErr := net.ResolveTCPAddr("tcp4", "127.0.0.1:9090")
	if nil != tcpAddrErr {
		fmt.Println(tcpAddrErr)
		return
	}
	tcpCon, tcpConErr := net.DialTCP("tcp4", nil, tcpAddr)
	tcpConReader := bufio.NewReader(tcpCon)
	if nil != tcpConErr {
		fmt.Println(tcpConErr)
	}

	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case e := <-events:
			switch e.Type {
			case termbox.EventKey:
				switch e.Key {
				case termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight, termbox.KeyArrowUp:
					eventByte, eventByteErr := json.Marshal(e)
					fmt.Println(string(eventByte))
					if nil != eventByteErr {
						continue
					}
					tcpCon.Write(eventByte)
					tcpCon.Write([]byte("\n"))
				}
			}
		default:
			Update(tcpConReader)
		}
	}
}
