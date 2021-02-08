package auth

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestToken(t *testing.T) {
	t1 := Token{
		Iat: time.Now().Unix(),
		Ein: 1, // expires in 1s
	}
	t2 := Token{}
	equals(t, nil, t2.OK(t1.Final()))
	time.Sleep(2 * time.Second)
	equals(t, ErrTokenExpired, t2.OK(t1.Final()))
}

func ExampleToken() {
	t1 := Token{
		Iat: time.Now().Unix(),
		Ein: 1, // expires in 1s
	}
	t2 := Token{}
	e := t2.OKUnless(t1.Final(), func(t *Token) error {
		fmt.Println(1)
		return nil
	})
	fmt.Println(e)
	// Output:
	// 1
	// <nil>
}

func TestTokenMetadata(t *testing.T) {
	t1 := Token{
		Iat: time.Now().Unix(),
		Ein: 1, // expires in 1s
		Meta: Metadata{
			"a": 2,
			"b": "3",
			"t": time.Now(),
		},
		Key: []byte("11101"),
	}
	tk := t1.Final()
	t2 := Token{Key: []byte("11101")}
	e := t2.OK(tk)
	equals(t, nil, e)
}
