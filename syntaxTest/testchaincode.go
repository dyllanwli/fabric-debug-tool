package main

import "fmt"
import "net/url"

type SimpleChaincode struct {
}

func main() {
	u := "%5B%22a%22%5D"
	var a string
	var err error
	a, err = url.QueryUnescape(u)
	if err == nil {
		fmt.Println(a)
	}
	var args []string
	var A, val string
	args = []string{"a", "100", "b", "200"}
	for key, value := range args {
		if key%2 == 0 {
			A = value
		} else {
			val = value
			fmt.Println(A, val)
		}
	}
}
