package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)


func main() {


	if strings.ToLower(runtime.GOOS) == "windows" {
		t1, _ := simplifiedchinese.GBK.NewEncoder().String("我靠")
		fmt.Println(t1)
		log.Printf("%s\n", t1)
	}
}


func nothing() {
	log.Println("log")
	cmd0 := exec.Command("C:\\works\\local\\bin\\curl.exe", " http://www.chaos.luxe")

	//if err := cmd0.Start(); err != nil {
	//	fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
	//	return
	//}
	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("error: Couldn't obtain the stdout pipe for command No.0:%s\n", err)
		return
	}
	defer stdout0.Close()

	output0 := make([]byte, 30)
	n, err := stdout0.Read(output0)

	if err != nil {
		fmt.Printf("Error: Could't read data from the pipe: %s\n", err)
		return
	}

	fmt.Printf("output => %s\n", output0[:n])
}


func sayHello() {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Println("Please input your namne:")

		// delim 是分隔符， 遇到分隔符读取就完成。
		input, err := inputReader.ReadString('\n')

		if err != nil {
			fmt.Printf("Found an error %s\n", err)
		} else {
			// 对 input 进行切片，去掉内容中最后一个（分隔符）
			input = input[:len(input)-1]
			fmt.Printf("Hello, %s\n", input)
		}

}