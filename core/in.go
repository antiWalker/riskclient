package core

import (
	"bigrisk/models"
	"strings"
)

var BaseTime string
type TraceContext struct {
	CostTime int64 `json:"costTime,omitempty"`
	TraceId string `json:"trace_id"`
	ConditionRule *complexNode `json:"condition_rule,omitempty"`
	MatchRule	*complexNode `json:"match_rule,omitempty"`
	RunStack    *Stack `json:"run_stack,omitempty"`
	Where	[]models.Where	`json:"where,omitempty"`
	ExecuteNode	*simpleNode	`json:"execute_node,omitempty"`
}
func changeField(data string)string{
	var w,fieldName string
	for i:= range data{
		s := data[i]
		if 'A' <= s && s <= 'Z' {
			if i==0{
				s = s - 'A' + 'a'
				w = string(s)
			}else{
				s = s - 'A' + 'a'
				w = "_"+string(s)
			}
		}else{
			w = string(s)
		}
		fieldName +=w
	}

	fieldName = strings.Replace(fieldName, "_i_p", "_ip", -1 )
	fieldName = strings.Replace(fieldName, "_i_d", "_id", -1 )

	return fieldName
}
