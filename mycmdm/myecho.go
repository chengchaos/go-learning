package main

import (
	"fmt"
	"github.com/chengchaos/go-learning/helper"
	"io"
	"log"
	"os/exec"
)


// "curl", "https://www.baidu.com")
func start1(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	stdout , err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln("E002 ", err)
	}
	cmd.Stderr = cmd.Stdout
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		log.Fatal("E001 ",err)
	}

	buf := make([]byte, 4)
	content := make([]byte, 0)
	var count int = 0

	for {
		n, err := stdout.Read(buf)
		if err != nil {
			break
		}
		count += n
		content = append(content, buf[0:n]...)
	}

	fmt.Print(helper.GBK(string(content[0:count])))

	if err := cmd.Wait(); err != nil {
		log.Fatalln("E003", err)
	}
}

func start2(name string, arg ...string) {

	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	cmd.Stderr = cmd.Stdout


	defer stdout.Close()

	if err = cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	buff := make([]byte, 1024)

	for {
		n, err := stdout.Read(buff)
		if err != nil && err != io.EOF {
			log.Printf("err in stdout.Read =>%s\n", err)
			break
		}
		fmt.Print(string(buff[0:n]))
	}

	fmt.Println("Unreachable *** ... ***")
	if err := cmd.Wait(); err != nil {
		log.Fatalln(err)
	}
}


func main() {
	start2("curl" , "https://www.chaos.luxe")
}