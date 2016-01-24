package main

import (
	"fmt"
	"regexp"
)

func isEmail(email string) bool {
	if email != "" {
		if isOk, _ := regexp.MatchString("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", email); isOk {
			return true
		}
		return false
	}
	return false
}

func main() {
	eMaileStr := "wxw113147088163.com"
	ret := isEmail(eMaileStr)
	if ret {
		fmt.Println("is Email!!!")
	} else {
		fmt.Println("is not Email!!!!")
	}
}
