package helper

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"runtime"
	"strings"
)

func ChineseGBK(format string, a ... interface{}) string {


	if strings.ToLower(runtime.GOOS) == "windows" {
		f1, _ := simplifiedchinese.GBK.NewEncoder().String(format)
		argsLen := len(a)
		others := make([]interface{}, argsLen)

		for i := 0; i < argsLen ; i++ {
			ai := a[i].(string)
			others[i] , _ = simplifiedchinese.GBK.NewEncoder().String(ai)
		}

		return fmt.Sprintf(f1, others...)
	}

	return fmt.Sprintf(format, a ...)
}


func GBK(input string) string {


	if strings.ToLower(runtime.GOOS) == "windows" {
		f1, _ := simplifiedchinese.GBK.NewEncoder().String(input)
		return f1
	}
	return input

}
