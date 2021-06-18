package libs

import (
	"crypto/md5"
	"fmt"
)
// Token
func NewToken(str string) string {
	//
	data := []byte(str)
	s := fmt.Sprintf("%x", md5.Sum(data))
	return s
}