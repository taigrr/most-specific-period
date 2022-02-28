package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		log.Println("line", s.Text())
	}
}
