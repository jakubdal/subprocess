package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jakubdal/subprocess"
)

// This file was only used during my testing
// TODO Provide actually sensible examples
func main() {
	proc, err := subprocess.NewProcess(context.Background(), nil, "sleep", nil, "3s")
	if err != nil {
		log.Fatalf("start process")
	}

	<-time.After(time.Second)
	proc.Signal(os.Interrupt)
	<-time.After(time.Second)
	proc.Signal(os.Interrupt)
	<-time.After(time.Second)
	proc.Stop()
	<-time.After(time.Second)

	<-time.After(time.Second * 5)
}
