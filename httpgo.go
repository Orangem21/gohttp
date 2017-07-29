package httpgo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"io/ioutil"
	"github.com/astaxie/beego/logs/alils"
	debug2 "runtime/debug"
)
//定义常量
const (
	contentType     = "Content-Type"
	jsonContentType =  "application/json"
	formContentType = "application/x-www-form-urlencoded"
)

//定义默认超时时间

const DefaultTimeout = 3 * time.Second

//定义开发模式

const debugEnv = "GOHTTP_DEBUG"

type basicAuth struct {
		username string
		password string
}

type fileForm struct {
	fieldName string
	filename string
	file *os.File
}

type GoRsponse struct {
	*http.Response
}


func (resp *GoRsponse) Asstring() (string,error) {
	data,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}
	return string(data),nil
}


func (resp *GoRsponse) AsBytes() ([]bytes,error){
	return ioutil.ReadAll(resp.Body)
}

func (resp *GoRsponse) AsJson(v interface{}) error{
	data,err := resp.AsBytes()
	if err != nil{
		return err
	}
	return json.Unmarshal(data,v)
}



type Client struct {
	c *http.Client
	query map[string]string
	queryStructs []interface{}
	headers map[string]string
	url string
	path string
	body io.Reader
	auth basicAuth
	cookies []*http.Cookie
	files []*fileForm
	proxy string
	timeout time.Duration
	tlsHandshakeTimeouut time.Duration
	retries int
	transport *http.Transport
	debug bool
	logger *log.Logger
}

var DefaultClient New()

func New() *Client {
	t := &http.Transport{}
	debug := os.Getenv(debugEnv) == "1"
	logger := log.New(os.Stderr,"[gohttp]",log.Ldate|log.Ltime|log.Lshortfile)

	return &Client{
		query:make(map[string]string),
		queryStructs:make([]interface{},0),
		headers:make(map[string]string),
		auth:basicAuth{},
		cookies:make([]*http.Cookie,0),
		files:make([]*fileForm, 0),
		timeout:DefaultTimeout,
		transport:t,
		debug:debug,
		logger:logger,
	}
}

//拷贝map
func copyMap(src map[string]string) map[string]string{
	newmap := make(map[string]string)
	for k,v := range src{
		newmap[k] = v
	}
	return newmap
}


func (c *Client) New() *Client {
	newClient := &Client{}
	//简单拷贝
	newClient.url = c.url
	newClient.path = c.path
	newClient.auth = c.auth
	newClient.proxy = c.proxy
	newClient.timeout = c.timeout
	newClient.tlsHandshakeTimeouut = c.tlsHandshakeTimeouut
	newClient.retries = c.retries
	newClient.debug = c.debug

	//mapcopy
	newClient.query = copyMap(c.query)
	newClient.headers = copyMap(c.headers)

	//use the same transport
	newClient.transport = c.transport

	return newClient
}

//log writes a formatted log message record to stderr

func (c *Client) logf(format string,v interface{}){
	if c.debug {
		c.logger.Printf(format,v)
	}
}


//setupClient handles the connection details from http client to TCP connections
//timeout ,proxy,TLS config ....these are very important but rarely used directly
//by httpclient users

func (c *Client) setupClient() error{
	//creat the transport and client instance first
	if c.proxy != ""{
		//used proxy
		proxy,err := url.Parse(c.proxy)
		if err != nil {
			return err
		}
		c.transport.Proxy = http.ProxyURL(proxy)
	}
	if c.tlsHandshakeTimeouut != time.Duration(0){
		c.transport.TLSHandshakeTimeout = c.tlsHandshakeTimeouut
	}

	c.c = &http.Client{Transport:c.transport}

	if c.timeout != time.Duration(0){
		c.c.Timeout = c.timeout
	}
	return nil //没有err
}


func (c *Client) preparefiles() error{
	if len(c.files) >0 {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		//for each file, write field name,filename,and file content to the body

		for _,file := range c.files {
			part, err := writer.CreateFormFile(file.fieldName, file.filename)
			if err != nil {
				return err
			}
			_ , err := io.copy(part, file.file)
			if err != nil {
				return err
			}
		}
		err := writer.Close()
		if err !=nil {
			return err
		}
		//finally ,get the real content type,and set the header
		multipartContentType := writer.FromDataContentType()
		c.Header(contentType,multipartContentType)
		c.body = body

	}
	return nil
}

//prepareRequests does all the preparation jobs for "gohttp"












