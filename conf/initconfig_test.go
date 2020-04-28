package conf

import (
	"fmt"
	"testing"
)

func TestInitLogConf(t *testing.T) {
	err := InitLogConf()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v",LogConfObj)
}