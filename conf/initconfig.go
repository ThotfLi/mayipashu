package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type LogConf struct{
	LogFilePath          string     //log文件位置
	MaxGogroutineNumber  uint32     //log解析的goroutine数量
	LogConsumerChanNumer uint32								   //logConsumerChan的长度
}

var LogConfObj LogConf

func init(){
	//默认初始化配置
	LogConfObj = LogConf{
		LogFilePath:         "../createlog/dig.log",
		MaxGogroutineNumber: 10,
		LogConsumerChanNumer:10,
	}
}

//通过配置文件初始化配置
func InitLogConf() error {
	c,err := ioutil.ReadFile("./tsconfig.json")
	if err != nil {
		fmt.Println("Readall err:",err)
		return err
	}

	json.Unmarshal(c,&LogConfObj)
	return nil
}