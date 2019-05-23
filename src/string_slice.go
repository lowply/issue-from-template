package main

import "strings"

type StringSlice string

func (s StringSlice) splitAndTrimSpace() []string {
	arr := strings.Split(string(s), ",")
	for i, v := range arr {
		arr[i] = strings.TrimSpace(v)
	}
	return arr
}
