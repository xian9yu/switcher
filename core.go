package main

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type TimeJudge struct {
	readDone    bool
	timeDone    bool
	twg         *sync.WaitGroup
	firstPacket []byte
}

func listen(rule *rules, wg *sync.WaitGroup) {
	defer wg.Done()
	//监听
	listener, err := net.Listen("tcp", rule.Listen)
	if err != nil {
		log.Printf("[ERROR] [%s] failed to listen at %s", rule.Name, rule.Listen)
		return
	}
	log.Printf("[INFO] [%s] listing at %s", rule.Name, rule.Listen)

	for {
		//处理客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERROR] [%s] failed to accept at %s", rule.Name, rule.Listen)
			break
		}

		go handleRegexp(conn, rule)

	}
}

//timeJudge
func timeJudge(t *TimeJudge) {
	times := 3000
	for i := 0; i < times; i++ {
		time.Sleep(1 * time.Millisecond)
		if t.readDone {
			t.twg.Done()
			return
		}
	}
	t.timeDone = true
	t.twg.Done()
}

func handleRegexp(conn net.Conn, rule *rules) {
	t := &TimeJudge{
		readDone: false,
		timeDone: false,
		twg:      new(sync.WaitGroup),
	}
	t.twg.Add(1)

	//获取第一个数据包
	go waitFirstPacket(conn, t)
	go timeJudge(t)

	t.twg.Wait()

	if t.timeDone {
		target, err := net.Dial("tcp", "127.0.0.1:3306")
		if err != nil {
			log.Println("[ERROR] could not dial mysql:", err)
			return
		}
		_, err = target.Write([]byte{32, 0, 0, 1, 133})
		if err != nil {
			log.Println("[ERROR] conn.Write error :", err)
		}
		go muxBridge(conn, target)
	} else if t.readDone {
		var target net.Conn
		//挨个匹配正则
		for _, v := range rule.Targets {
			if !v.regexp.Match(t.firstPacket) {
				continue
			}

			c, err := net.Dial("tcp", v.Address)
			if err != nil {
				log.Printf("[ERROR] [%s] try to handle connection (%s) failed because target (%s) connected failed, try next match target.",
					rule.Name, conn.RemoteAddr(), v.Address)
				continue
			}
			target = c
			break
		}

		if target == nil {
			log.Printf("[ERROR] [%s] unable to handle connection (%s) because no match target",
				rule.Name, conn.RemoteAddr())
			return
		}
		log.Printf("[INFO] [%s] handle connection (%s) to target (%s)", rule.Name, conn.RemoteAddr(), target.RemoteAddr())

		//把第一个数据包发送给目标
		_, err := target.Write(t.firstPacket)
		if err != nil {
			log.Println("[ERROR] conn.Write error: ", err)
		}

		//io桥
		go tcpBridge(conn, target)
		tcpBridge(target, conn)
	}

}

//获取第一个数据包
func waitFirstPacket(conn net.Conn, t *TimeJudge) {
	length := 5
	buf := make([]byte, length)

	n, err := io.ReadAtLeast(conn, buf, length)
	if err != nil {
		if t.timeDone {
			log.Println("conn read timeout!")
		} else {
			log.Println("conn read error: ", err)
		}
	} else {
		t.readDone = true
	}

	t.firstPacket = buf[:n]
	//fmt.Println(n, " | ", buf[:n])
}

func tcpBridge(a, b net.Conn) {
	defer func() {
		_ = a.Close()
		_ = b.Close()
	}()
	buf := make([]byte, 2048)
	for {
		n, err := a.Read(buf)
		if err != nil {
			log.Println("[ERROR] conn.Read error: ", err)
			return
		}
		_, err = b.Write(buf[:n])
		if err != nil {
			log.Println("[ERROR] conn.Write error: ", err)
		}
	}
}

func muxBridge(conn, mysql net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func(conn, mysql net.Conn) {
		defer wg.Done()
		_, err := io.Copy(conn, mysql)
		if err != nil {
			log.Println("[ERROR] io.Copy error: ", err)
		}
		//conn.Close()
	}(conn, mysql)

	go func(conn, mysql net.Conn) {
		defer wg.Done()
		_, err := io.Copy(mysql, conn)
		if err != nil {
			log.Println("[ERROR] io.Copy error: ", err)
		}

		//mysql.Close()
	}(conn, mysql)

	wg.Wait()
	_ = mysql.Close()
	_ = conn.Close()
}
