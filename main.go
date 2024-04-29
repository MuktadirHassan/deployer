package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
)

var (
	logger = log.Default()
)

func main() {
	logger.Println("Welcome to deployer!")
	var nFlag = flag.Int("n", 1234, "Description of the flag")

	println(*nFlag)
	flag.Parse()

	fmt.Println("word: ", *nFlag)

	cmd := exec.Command("ls", "-lah")

	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		logger.Fatal("ls error", err)
	}

	fmt.Println(out.String())

}
