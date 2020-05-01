package distribute

import (
	"fmt"
	"mayipashu/iface"
	"testing"
)

var logConsumerObj iface.ILogConsumer

func testLogConsumer_GetLogToChan(t *testing.T) {
	logConsumerObj = NewLogConsumer()
	defer logConsumerObj.Close()

	go func() {
		for {
			a := <-logConsumerObj.GetLogChan()
			fmt.Printf("%v\n",a)
		}
	}()

	logConsumerObj.GetLogToChan()
}
