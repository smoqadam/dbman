package main

import (
	"strconv"
)

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}
