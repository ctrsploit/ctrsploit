package vul

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/json"
)

type Vulnerabilities []Vulnerability

func (v Vulnerabilities) Output() (err error) {
	if true {
		var output []byte
		output, err = json.Marshal(v)
		fmt.Println(string(output))
	}
	return
}
