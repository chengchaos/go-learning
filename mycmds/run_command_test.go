package mycmds

import (
	"log"
	"os"
	"os/exec"
	"testing"
)

func Test_RunCommand001(t *testing.T) {


	cmd := exec.Command("java", "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("result =>", string(output))

}

func Test_RunCommand002(t *testing.T) {
	err := RunCommand("curl", "https://www.chaos.luxe")
	if err != nil {
		log.Fatalln(err)
	}
}

func Test_RunCommandRedirect(t *testing.T) {

	//var stdout io.ReadCloser
	//var err error

	stdout, err :=os.OpenFile("stdout.log", os.O_CREATE|os.O_WRONLY, os.FileMode.Perm(0600))
	if err != nil {
		log.Fatalln(err)
	}
	defer stdout.Close()


	// 重定向标准输出到文件
	cmd := exec.Command("java.exe", "-h")
	cmd.Stdout = stdout

	// 执行命令
	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}
}