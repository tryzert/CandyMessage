package service

import (
	"net"
	"fmt"
	"time"
	"io"
	"sync"
	"encoding/json"
	"can/msg"
)


//每个客户端对象
type Client struct {
	conn net.Conn
	deadline time.Time
}


//数据流，由各个客户端发送过来，并交给处理器统一处理
type DataStream struct {
	Data []byte
	Conn net.Conn
}


//工作池
type WorkPool struct {
	Lock sync.RWMutex

	ClientQueue map[int]*Client
	MaxClientQueueSize int
	ClientEnterKey chan struct{}

	MessageQueue chan *DataStream

	// Handler func()
}


func NewWorkPool(poolsize, mqsize int) *WorkPool {
	return &WorkPool{
		ClientQueue: make(map[int]*Client),
		MaxClientQueueSize: poolsize,
		ClientEnterKey: make(chan struct{}, poolsize),
		MessageQueue: make(chan *DataStream, mqsize),
		// Handler: handler,
	}
}


//只是单纯增加，不包括改
func (this *WorkPool) addClient(cid int, conn net.Conn) {
	this.Lock.Lock()
	this.ClientQueue[cid] = &Client{
		conn: conn,
		deadline: time.Now(),
	}
	this.Lock.Unlock()
	<- this.ClientEnterKey
}


//删除客户端
func (this *WorkPool) deleteClient(cid int) {
	if ct, ok := this.findClient(cid); ok {
		this.Lock.Lock()
		ct.conn.Close()
		delete(this.ClientQueue, cid)	
		this.Lock.Unlock()
		this.ClientEnterKey <- struct{}{}
	}
}

//重置客户端
func (this *WorkPool) resetClient(cid int, conn net.Conn) {
	this.Lock.Lock()
	this.ClientQueue[cid].conn = conn
	this.ClientQueue[cid].deadline = time.Now()
	this.Lock.Unlock()
}


//读取内容不需要加锁
func (this *WorkPool) findClient(cid int) (*Client, bool) {
	v, ok := this.ClientQueue[cid]
	return v, ok
}


//刷新客户端活跃时间
func (this *WorkPool) updateClientDeadline(cid int) {
	if ct, ok := this.findClient(cid); ok {
		this.Lock.Lock()
		ct.deadline = time.Now()
		this.Lock.Unlock()
	}
}

//用于清除一段时间未响应的客户端
func (this *WorkPool) kickoutDeadClient(alive_period, check_interval time.Duration) {
	for {
		for cid, v := range this.ClientQueue {
			if time.Since(v.deadline) > alive_period {
				this.deleteClient(cid)
			}			
		}
		// fmt.Println("[kickDeadClient]接入的客户端数量：", len(this.ClientQueue))
		time.Sleep(check_interval)
	}
}





//启动服务
func (this *WorkPool) Serve(port string) {
	server, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("启动服务失败...", err)
		return
	}
	fmt.Println("服务正在运行...")
	defer server.Close()

	go this.kickoutDeadClient(time.Minute * 30, time.Minute * 3)

	for i := 0; i < this.MaxClientQueueSize; i++ {
		this.ClientEnterKey <- struct{}{}
	}
	clientId := 0
	for {
		client, err := server.Accept()
		if err != nil || len(this.ClientEnterKey) == 0 {
			client.Write([]byte("服务器连接过载..."))
			client.Close()
			continue
		}
		<-this.ClientEnterKey
		go HandleClientConnection(this, client, clientId)
		clientId++
	}
}


//消息分类调度器
func (this *WorkPool) MessageDispatcher() {
	mb := &msg.BaseMessage{}
	for {
		data_stream := <- this.MessageQueue
		err := json.Unmarshal(data_stream.Data, mb)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch mb.Type {
			case msg.REGISTER:
				TODO()
			case msg.LOGIN:
				TODO()
			case msg.HEART_BEAT:
				// message.SendData(data, conn, false)
			case msg.NORMAL:
				TODO()
			case msg.GROUP:
				TODO()
			case msg.TEMPORARY:
				TODO()
			default:
				fmt.Println("收到消息格式出错...")
		}
		fmt.Println()
	}
}


//客户端处理
func HandleClientConnection(this *WorkPool, c net.Conn, id int) {
	this.ClientQueue[id] = &Client{
		conn: c,
		deadline: time.Now(),
	}
	defer func() {
		c.Close()
		delete(this.ClientQueue, id)
		this.ClientEnterKey <- struct{}{}
	}()

	buff := make([]byte, 1024)
	for {		
		size, err := c.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil{
			continue
		}
		if size > 0 {
			this.Lock.Lock()
			this.ClientQueue[id].deadline = time.Now()
			this.Lock.Unlock()
			fmt.Println("[来自客户端的消息]", id, string(buff[:size]))
		}		
		time.Sleep(time.Second * 2)
	}
}


func TODO() {
}