package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

func makeResult(errno errnoType, data interface{}) resultType {
	return resultType{
		Errno:  errno,
		ErrMsg: errnoMsg[errno],
		Data:   data,
	}
}

func makeComplexResult(errno errnoType, errMsg string, data interface{}) resultType {
	return resultType{
		Errno:  errno,
		ErrMsg: errMsgType(errMsg),
		Data:   data,
	}
}

// 发送 JSON 给客户端
func sendJsonToClient(data interface{}) error {

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	if err := enc.Encode(data); err != nil {
		return err
	}else{
		return enc.Encode(data)
	}
}

// 发送结果给客户端
func sendResult(errno errnoType, data interface{}) error {
	result := makeResult(errno, data)

	return sendJsonToClient(result)
}

func sendComplexResult(w http.ResponseWriter, errno errnoType, errMsg string, data interface{}) error {
	result := makeComplexResult(errno, errMsg, data)

	return sendJsonToClient(result)
}

// execute chain of functions
// return the first error if `error` happened
// else return nil if all success
func doChain(fs ...func() error) error {
	for _, f := range fs {
		if err := f(); err != nil {
			return err
		}
	}

	// all is OK
	return nil
}
