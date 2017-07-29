package main

import (
	"net/url"
	"gohttp/conf"
	"time"
	"regexp"
)

type image struct {
	imageURL string
	fileName string
	retry 	 int
	folder   string
}
type page struct {
	url    string
	body   string
	retry   int
	parse   bool //解析规则
}

type contxt struct {
	page       map[string]int //记录处理状态，key是地址，value是状态
	image      map[string]int //记录图片的处理状态
	pageChan   chan *page
	imgChan    chan *image
	parseChan  chan *page
	imageCount chan int
	savePath   string
	rootURL   *url.URL
	config   *conf.Config
}

const  (
	bufferSize		=	64*1024
	numpuller		=	5		//并发数
	numDownloader	=	10
	maxRetry		=	5
	statusInterval	=	15	*	time.Second
	chanBuffersize	=	80

//图片处理的状态
	readt = iota	//待处理
	done
	fail
)

var (
	titleExp  = regexp.MustCompile(`<title>([^<>]+)</title>`)          //regexp.MustCompile(`<img\s+src="([^"'<>]*)"/?>`)
	invalidCharExp = regexp.MustCompile(`[\\/*?:><|]`)
)
//主程序
func main(){
	configFile := "config.json"
	cf := &conf.Config{}
	if err := cf.load
}



