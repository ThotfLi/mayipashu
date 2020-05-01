package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type LogConf struct {
	LogFilePath          string //log文件位置
	MaxGogroutineNumber  uint32 //log解析的goroutine数量
	LogConsumerChanNumer uint32 //logConsumerChan的长度
	TimeInterval	     time.Duration  //定时任务时间

	MinTaskPool          uint32
	MaxTaskPool          uint32

}

var LogConfObj LogConf

func init() {
	//默认初始化配置
	LogConfObj = LogConf{
		LogFilePath:          "../createlog/dig.log",
		MaxGogroutineNumber:  10,
		LogConsumerChanNumer: 10,
		TimeInterval:5,
		MinTaskPool:2,
		MaxTaskPool:10,
	}
}

//通过配置文件初始化配置
func InitLogConf() error {
	c, err := ioutil.ReadFile("./tsconfig.json")
	if err != nil {
		fmt.Println("Readall err:", err)
		return err
	}

	json.Unmarshal(c, &LogConfObj)
	return nil
}
