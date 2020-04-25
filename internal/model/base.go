package model

import (
	"net/url"
	"strings"
)

// BlogGin hello BlogGin.
type GiNana struct {
	Hello string
}

func Rawurlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}
