package model

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"strings"
)

func IsZeroOrInvalid(v interface{}, fields ...string) error {
	ref := reflect.Indirect(reflect.ValueOf(v))
	for _, field := range fields {
		if f := ref.FieldByName(field); f.IsValid() {
			if f.IsZero() {
				return fmt.Errorf("request body field '%s' has no value", strings.ToLower(field))
			}
		} else {
			return fmt.Errorf("request body has no field '%s'", strings.ToLower(field))
		}
	}
	return nil
}

// func exportalize(s string) string {
// 	b := strings.Builder{}
// 	b.Grow(len(s))
// 	c := s[0]
// 	if 'a' <= c && c <= 'z' {
// 		c -= 'a' - 'A'
// 	}
// 	b.WriteByte(c)
// 	b.WriteString(s[1:])
// 	return b.String()
// }

func RandStr(n int) string {
	b := make([]byte, n/2)
	rand.Read(b) // nolint:errcheck
	return fmt.Sprintf("%x", b)
}
