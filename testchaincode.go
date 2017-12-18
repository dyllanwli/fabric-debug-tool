package main

import "fmt"
import "net/url"

type SimpleChaincode struct {
}

func main() {
	u := "%5B%22a%22%5D"
	var a string
	a, _ = url.QueryUnescape(u)
	fmt.Println(a)
}
