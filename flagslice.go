package main

import "strings"

type StringSlice struct {
	Slice *[]string
}

func (v StringSlice) String() string {
	if v.Slice != nil {
		strings.Join(*v.Slice, ",")
	}
	return ""
}

func (v StringSlice) Set(s string) error {
	*v.Slice = strings.Split(s, ",")
	return nil
}
