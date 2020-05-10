package main

import (
	"flag"
	"os"
	"time"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "path to kubeconfig")
	flag.Parse()

	if *kubeconfig == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	time.Sleep(time.Second * 60)
}
