//有缓冲的通道

package main
import(
	"sync"
	"time"
	"fmt"
	"math/rand"
)

const (
	numberGoroutine = 4 //goroutine 数量
	taskLoad        = 10 //任务数量
)

var wg sync.WaitGroup

func init(){
 rand.Seed(time.Now().Unix())
}

func main() {
	//创建有缓冲的通道
	tasks := make(chan string, taskLoad)
	//启动goroutine来启动工作
	wg.Add(numberGoroutine)

	for gr := 1; gr <= numberGoroutine; gr ++ {
		go worker(tasks, gr)
	}
	for post := 2; post <= taskLoad; post ++ {
		tasks <- fmt.Sprintf("Tasks:%d", post)
	}
	//当所有工作处理完毕的时候关闭通道
	close(tasks)
	wg.Wait()
}
//worker 作为goroutine 启动处理
func worker(tasks chan string,worker int){
	defer wg.Done()
	for {
		//等待分配工作
		task,ok := <- tasks
		if !ok {
			fmt.Printf("Worker:%d : Shutting Down\n",worker)
			return
		}
		fmt.Printf("Worker:%d : Started %s \n",worker,task)
		sleep:=rand.Int63n(100)
		time.Sleep(time.Duration(sleep)*time.Millisecond)

		fmt.Printf("Worker:%d Complated %s \n",worker,task)
	}
}
