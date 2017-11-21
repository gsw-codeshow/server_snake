package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	_ "math/rand"
	"net"
	"server_snake/snake"
	"sync"

	"github.com/nsf/termbox-go"
)

type Server struct {
	ServerTcpAddr   *net.TCPAddr
	ServerListener  *net.TCPListener
	ClientConnMutex *sync.Mutex
	contextMutex    *sync.Mutex
	ClientConn      map[string]*net.TCPConn
	ConText         *snake.Context

	Err error
}

func (server *Server) InitServer() {
	server.ServerTcpAddr, server.Err = net.ResolveTCPAddr("tcp4", "127.0.0.1:9090")
	server.ClientConnMutex = new(sync.Mutex)
	server.contextMutex = new(sync.Mutex)
	server.ClientConn = make(map[string]*net.TCPConn)
	server.ConText = snake.NewContext()
	if nil != server.Err {
		return
	}
	server.ServerListener, server.Err = net.ListenTCP("tcp4", server.ServerTcpAddr)
}

func (server *Server) BroadCast() {
	for {
		server.contextMutex.Lock()
		server.ConText.Update()
		context, contextErr := json.Marshal(&server.ConText)
		if nil != contextErr {
			continue
		}
		server.contextMutex.Unlock()
		server.ClientConnMutex.Lock()
		for connName, conn := range server.ClientConn {
			fmt.Println(string(context))
			_, connErr := conn.Write(context)
			if nil != connErr {
				delete(server.ClientConn, connName)
				delete(server.ConText.Snake, connName)
				conn.Close()
				fmt.Println(connName + "off-line")
			}
			_, connErr = conn.Write([]byte("\n"))
			if nil != connErr {
				delete(server.ClientConn, connName)
				delete(server.ConText.Snake, connName)
				conn.Close()
				fmt.Println(connName + "off-line")
			}
		}
		server.ClientConnMutex.Unlock()
	}
}

func (server *Server) JoinClient() {
	for {
		clientConn, clientConnErr := server.ServerListener.AcceptTCP()
		if nil != clientConnErr {
			continue
		}
		server.contextMutex.Lock()
		server.ConText.AddSnake(clientConn.RemoteAddr().String(), 10, 10)
		server.contextMutex.Unlock()
		fmt.Println(clientConn.RemoteAddr().String() + "join the game")
		server.ClientConnMutex.Lock()
		server.ClientConn[clientConn.RemoteAddr().String()] = clientConn
		go server.AcceptKeyBorad(clientConn)
		server.ClientConnMutex.Unlock()
	}
	return
}

func NewServer() *Server {
	server := new(Server)
	server.InitServer()
	return server
}

func (server *Server) AcceptKeyBorad(clientConn *net.TCPConn) {
	for {
		var event termbox.Event
		clientRead := bufio.NewReader(clientConn)
		msg, msgErr := clientRead.ReadString('\n')
		if nil != msgErr {
			fmt.Println(msg)
			return
		}
		json.Unmarshal([]byte(msg), &event)
		server.contextMutex.Lock()
		_, snakeOk := server.ConText.Snake[clientConn.RemoteAddr().String()]
		if snakeOk {
			server.ConText.HandleKey(clientConn.RemoteAddr().String(), event.Key)
		}
		server.contextMutex.Unlock()

	}
}

func main() {
	server := NewServer()
	if nil != server.Err {
		fmt.Println(server.Err)
	}
	go server.BroadCast()
	server.JoinClient()
}
