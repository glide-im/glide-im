package main

import "testing"

// go test -v -run=TestServerPerf -cpuprofile="cpu2.out"
func TestRunTestServer(t *testing.T) {
	RunTestServer()
}

// go test -v -run=TestRunAnalysisServer -cpuprofile="cpu2.out"
func TestRunAnalysisServer(t *testing.T) {
	RunAnalysisServer()
}

func TestRunClientMsg(t *testing.T) {
	RunClientMsg()
}
