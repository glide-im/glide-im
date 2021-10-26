package main

import (
	"go_im/cmd/test"
	"testing"
)

// go test -v -run=TestServerPerf -cpuprofile="cpu.out"
func TestServerPerf(t *testing.T) {
	test.RunAnalysisServer()
}

func TestRunClient(t *testing.T) {
	test.RunClientMsg()
}
