package example

import (
	"context"
	"net"
	"time"
)

type Server interface {
	SetAcceptor(Acceptor)               // 处理握手等相关工作
	SetMessageListener(MessageListener) // 设置消息监听器
	SetStateListener(StateListener)     // 设置状态监听器
	SetReadWait(time.Duration)          // 设置连接超时,用于控制心跳逻辑
	SetChannelMap(ChannelMap)           // 设置连接器

	Start() error
	Push(string, []byte) error
	Shutdown(context.Context) error
}

type Acceptor interface {
	Accept()
}

type Conn interface {
	net.Conn
	ReadFrame() (Frame, error)
	WriteFrame(OpCode, []byte) error
	Flush() error
}

type ChannelMap interface {
	Add(channel Channel)
	Remove(id string)
	Get(id string) (Channel, bool)
	All() []Channel
}

type Channel interface {
	Conn
	Agent
	Close() error
	Readloop(lst MessageListener) error
	SetWriteWait(time.Duration)
	SetReadWait(time.Duration)
}

type OpCode byte

const (
	OpContinuation OpCode = 0x0
	OpText         OpCode = 0x1
	OpBinary       OpCode = 0x2
	OpClose        OpCode = 0x8
	OpPing         OpCode = 0x9
	OpPong         OpCode = 0xa
)

//这是一个数据帧,WebSocket这一块的

type Frame interface {
	SetOpCode(OpCode)
	GetOpCode() OpCode
	SetPayload([]byte)
	GetPayload() []byte
}

//接收消息接口

type MessageListener interface {
	Receive(Agent, []byte)
}

//状态监听接口

type StateListener interface {
	Disconnect(string) error
}

//发送方

type Agent interface {
	ID() string        //返回连接的channelID
	Push([]byte) error //用于上层业务返回信息
}

//拨号连接接口

type Dialer interface {
	DialAndHandshake(DialerContext) (net.Conn, error)
}

//拨号连接结构体

type DialerContext struct {
	Id      string
	Name    string
	Address string
	Timeout time.Duration
}

type Client interface {
	ID() string
	Name() string
	Connect(string) error
	SetDialer(Dialer)
	Send([]byte) error
	Read() (Frame, error)
	Close()
}
