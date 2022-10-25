package main

type Conf struct {
	LoopCount   int      // loop count
	ValueSize   int      // values size
	ClientCount int      // client count
	Command     []string // redis command
	Addrs       []string // redis Addrs "ip:port"
}

var conf = Conf{
	LoopCount:   10000,
	ValueSize:   250000,
	ClientCount: 200,
	Command:     []string{"SET", "GET"},
	//Addrs:       []string{"127.0.0.1:7011"}, // cluster 구성이 되어 있으면, 접속 주소를 하나만 넣어도 샤딩됨.
	Addrs: []string{"127.0.0.1:7011", "127.0.0.1:7013", "127.0.0.1:7015"}, // 하지만, flushall을 하기 위해 다 넣어줌.
}
