package conf

import (
	"regexp"
	"net/url"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"log"
	"errors"
	"strings"
)

type MatchExp struct{
	Exp *regexp.Regexp
	Match int
	Folder interface{} //可选Url,title,none,正则表达式
}

type Config struct {
	Root		*url.URL
	ImageRegex 		[]*MatchExp
	PageRegex		[]*regexp.Regexp
	imgPageRegex	[]*regexp.Regexp
	HerfRegex		[]*MatchExp
}

//{
//"root":"sexy.faceks.com",
//	"regex":{
//		"image":[
//				{
//					"exp":"bigimgsrc=\"([^\"?]+)",
//					"match":1,
//					"folder":"none" ##存放图片的文件夹，可选值url,title,none,正则表达式
//				}
//			],
//		"page":[],
//		"imgInPage":["\S+/post/\S+"],
//		"href":[
//				{
//					"exp":"\s+href=\"([a-zA-Z0-9_\-/:\.%?=]+)\"",
//					"match":1
//				}
//			]
//		}
//}

func (c *Config) Load(file string) error{
	content,err:= ioutil.ReadFile(file)
	if err!=nil{
		return err
	}
	content = bytes.Replace(content, []byte("\\"), []byte("\\\\"), -1 )
	content = bytes.Replace(content, []byte("\\\\\""),[]byte("\\\""),-1) //什么意思

	comRegex := regexp.MustCompile(`\s*##.*`)
	content = comRegex.ReplaceAll(content,[]byte{})

	nRegex := regexp.MustCompile(`\n|\t|\r`)
	content = nRegex.ReplaceAll(content, []byte{})

	json0bj := make(map[string]interface{})

	err = json.Unmarshal(content,&json0bj)

	if err != nil{
		log.Println(string(content))
		return errors.New("[1]配置文件格式有误:"+err.Error())
	}
	root := json0bj["root"].(string)
	temp := strings.ToLower(root)

	if !strings.HasPrefix(temp, "http://") && !strings.HasPrefix(temp, "https://") {
		root = "http://"+root
	}
	c.Root, err = url.Parse(root)

	if err!=nil{
		return err
	}
	reg,ok := json0bj["regex"].(map[string]interface{}){
		if !ok {
			return errors.New("解析正则表达式错误")
		}
	}

	imgRegs,ok := reg["images"].([]interface{})
	if ok{
		c.ImageRegex = make([]*MatchExp, len(imgRegs))
		for i ,val := range imgRegs {
			obj, ok := val.(map[string]interface{})
			if !ok {
				return errors.New("解析图片配置失败")
			}

			exp := obj["exp"].(string)

			c.ImageRegex[i] = &MatchExp{}
			c.ImageRegex[i].Match = int(obj["match"].(float64))

			folder := strings.ToLower(obj["folder"].(string))

			if folder != "none" && folder != "url" && folder != "title" {
				c.ImageRegex[i].Folder, err = regexp.Compile(folder)
				if err != nil {
					return errors.New("[4]解析正则表达式" + folder + "时出错")
				}
			} else {
				c.ImageRegex[i].Folder = folder
			}
			c.ImageRegex[i].Exp, err = regexp.Compile(exp)

			if err != nil {
				return errors.New("[5]解析正则表达式" + exp + "时出错")
			}
		}
	}else {
		return errors.New("[6]解析regex.image时出错")
	}

}




