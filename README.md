基于golang
一个简单的信息转发服务端

基本思路：
信息传输格式暂定为json,因为格式比较好解析:
	比如：{
			Type:1,
			Data:{
					Sender:100101,
					Receiver:888888,
					Content:"Hello world",
					Date:"2020-6-2",
					Time:"13:22:35",
					Received:"False",
				}
		  }
一个草案。
