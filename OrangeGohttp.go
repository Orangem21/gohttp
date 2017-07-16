package main

import (
	"sync"
	"math/rand"
	"time"
	"fmt"
)

var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().UnixNano())
}


func main(){
	//创建无缓冲的通道
	court := make(chan int )
	//等待2个goroutine
	wg.Add(2)

	go player("tom",court)
	go player("jam",court)

	//发球
	court <- 1

	//等待游戏结束
	wg.Wait()
}

//player 模拟打球

func player(name string,court chan int ){
	defer wg.Done()

	for {
		ball,ok := <-court
		if !ok {
			//如果通道关闭则我们胜利
			fmt.Printf("Player %s win \n",name)
			return
		}
		n := rand.Intn(100)
		if n%13 == 0{
			fmt.Printf("Play %s loss \n",name)
			close(court)
			return
		}
		fmt.Printf("Player %s Hit %d \n",name,ball)
		ball++
		court <- ball
	}
}

