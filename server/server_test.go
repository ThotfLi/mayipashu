package server

import (
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	 s := NewServer()
	 s.Start()
	 time.Sleep(10*time.Second)
	 s.Stop()
}
