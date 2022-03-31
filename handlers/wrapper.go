package handlers

import (
	"context"
	"encoding/json"
	"github.com/antiWalker/golib/common"
	"net/http"
	"reflect"
	"time"
)

// args is a point to `Type`
type wrapperHandler = func(ctx context.Context, response http.ResponseWriter, args interface{}) error

type TimeContext struct {
	CostTime int64         `json:"costTime,omitempty"`
	TraceId  string        `json:"trace_id"`
	Params   *DetectFormV2 `json:"params,omitempty"`
}

// WARNING:--------------------------------------------------
// args is a Plain struct data [not pointer]
// but args in `wrapperHandler` is a pointer point to `args`
// args is json encode-able data
// ----------------------------------------------------------
// and must be have the same field with http request form data
// handler is the real handler of request
// args field must be a struct and field type are must be string
// [if you need convert to other, manually do it]
func Wrapper(args interface{}, handler wrapperHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UnixNano()

		// we do need deep copy  [防止并发的时候出现问题]
		// data is a pointer to args's `type`
		// let's say: args type's real type is: ArgType
		// the data have type: *ArgType

		data := reflect.New(reflect.TypeOf(args)).Interface()

		// parse function  Convert r.Forms  -> args [args is a pointer]
		p := func() error {
			m := map[string]string{}
			// do not allowed for multi data
			for key, value := range r.Form {
				m[key] = value[0]
			}
			for key, value := range r.PostForm {
				m[key] = value[0]
			}

			var bytes []byte
			var err error

			return doChain(
				func() error {
					bytes, err = json.Marshal(m)
					return err
				},

				func() error {
					err = json.Unmarshal(bytes, data)
					return err
				})
		}

		// do love lambda if we have
		if err := doChain(r.ParseForm, p); err != nil {
			_ = sendResult(errnoParseParamArg, nil)
			return
		}
		TraceId := r.Header.Get("TRACE_ID")
		ctx := context.WithValue(context.Background(), "TraceId", TraceId)
		// do process request
		// may process error or record time
		if err := handler(ctx, w, data); err != nil {
			common.ErrorLogger.Error(err.Error())
		}

		//request_id

		common.InfoLogger.Info("Wrapper Cost Time: ", &TimeContext{
			TraceId:  TraceId,
			CostTime: (time.Now().UnixNano() - start) / 1000,
		})
	}
}
