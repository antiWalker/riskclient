package core

import (
	"encoding/json"
	"testing"
	"time"
)

func Test_DemoHello(t *testing.T) {
	if d, err := time.ParseDuration("-24h"); err == nil {
		now := time.Now()
		a := now.Add(d)
		t.Log(a)
	} else {
		t.Error(err.Error())
	}
}

func Test_TypeChange(t *testing.T) {
	a := struct {
		Type string `json:"type"`
	}{Type: "hello"}

	data, _ := json.Marshal(a)

	f := func(bs []byte, v interface{}) {
		c := v

		_ = json.Unmarshal(bs, &c)

		t.Log(c)
	}

	t.Log(&a)
	f(data, a)
}
