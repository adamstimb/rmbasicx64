package main

import (
	"fmt"
	"regexp"
)

// GetStringValues finds string values in a line of code
func GetStringValues(code string) {
	// use a regex to find a string between double quotes
	r, _ := regexp.Compile("\"(.*?)\"")
	b := []byte(code)
	fmt.Println(r.FindAll(b, -1))      // get string values themselves including "
	fmt.Println(r.FindAllIndex(b, -1)) // get start and end indexes of each match
}
