package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/zhufuyi/goctl/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	rootCMD := cmd.NewRootCMD()
	if err := rootCMD.Execute(); err != nil {
		rootCMD.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
