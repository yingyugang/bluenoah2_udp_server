package main

import (
	"encoding/json"
	"fmt"
	"net"
	_ "sync"
	"time"
)

type BaseMessage struct {
	API int
	Uuid string
	Message string
}

type UserMessage struct {
	  X float32
	  Y float32
	  Z float32
	  Hp int
	  MaxHp int
	  Damage float64
	  UserName string
}

type UserConn struct {
	uuid string
	tick int
	address *net.UDPAddr
	userMessage UserMessage
}

type AllUserMessage struct {
	 UserMessages [4]string
}

var userMap map[string]UserConn

var listen *net.UDPConn

func Send() {
	for {
		time.Sleep(1000 * time.Millisecond)
		var users AllUserMessage
		var i int
		for  k := range userMap{
			if len(userMap[k].uuid) != 0{
				if userMap[k].tick > 0{
					result,err := json.Marshal(userMap[k].userMessage)
					if err != nil {
						fmt.Printf("userMap[v].userMessage, err:%v\n", err)
					}else {
						users.UserMessages[i] = string(result)
					}
					var user = userMap[k]
					user.tick = userMap[k].tick - 1
					userMap[k] = user
				}else{
					delete(userMap,k)
				}
			}
			i++
		}
		if i > 0{
			result,err := json.Marshal(users)
			fmt.Println(string(result))
			if err!=nil{
				fmt.Printf("users, err:%v\n", err)
			}
			for k := range userMap {
				if len(userMap[k].uuid) != 0{
					_, err := listen.WriteToUDP(result, userMap[k].address)
					if err != nil {
						fmt.Printf("write udp failed, err:%v\n", err)
						continue
					}
				}
			}
		}
	}
}

func main() {

	userMap = make(map[string]UserConn)
	//建立一个UDP的监听，这里使用的是ListenUDP，并且地址是一个结构体
	lis, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8090,
	})
	listen = lis
	if err != nil {
		fmt.Printf("listen failed, err:%v\n", err)
		return
	}
	go Send()
	for {
		var data [4096]byte
		//读取UDP数据
		count, addr, err := listen.ReadFromUDP(data[:])

		if err != nil {
			fmt.Printf("read udp failed, err:%v\n", err)
			continue
		}
		var baseMessage BaseMessage
		json.Unmarshal(data[0:count],&baseMessage)
		if baseMessage.API == 10{
			var userMessage UserMessage
			json.Unmarshal([]byte(baseMessage.Message),&userMessage)
			if val,ok := userMap[baseMessage.Uuid]; ok{
				val.address = addr
				val.userMessage = userMessage
				val.tick = 3
				userMap[baseMessage.Uuid] = val
			}else{
				if len(userMap) <= 4{
					userConn := UserConn{uuid:baseMessage.Uuid,address:addr,userMessage:userMessage,tick: 3}
					userMap[baseMessage.Uuid] = userConn
					fmt.Println(userConn.address)
				}
			}
		}
		fmt.Printf("data:%s addr:%v count:%d\n", string(data[0:count]), addr, count)
	}
}