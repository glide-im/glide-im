package main

import "testing"

// go test -v -run=TestServerPerf -cpuprofile="cpu.out"
func TestRunTestServer(t *testing.T) {
	RunTestServer()
}

func TestRunAnalysisServer(t *testing.T) {
	RunAnalysisServer()
}

func TestRunClientMsg(t *testing.T) {
	RunClientMsg()
}
