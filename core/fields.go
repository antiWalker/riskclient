package core

// 字段类型
type FieldType string

// definition or all supported field
const (
	// DO NOT forget add to `fieldSupportOperationMap` list
	FieldOtaName FieldType = "ota_name"
	FieldOtaId   FieldType = "ota_id"

	FieldUid        FieldType = "uid"
	FieldUserName   FieldType = "user_name"
	FieldUserMobile FieldType = "user_mobile"
	FieldUserEMail  FieldType = "user_email"

	FieldBusinessLine   FieldType = "business_line"   // 业务线
	FieldFirstCategory  FieldType = "first_category"  // 一级品类
	FieldSecondCategory FieldType = "second_category" // 二级品类

	FieldSalesId   FieldType = "sales_id"
	FieldSalesName FieldType = "sales_name"
	// todo add all field definition to here
)

// field -> support operator map
// global value
var fieldSupportOperationMap = map[FieldType][]AOT{
	FieldOtaId:   {compareEqual, compareNotEqual, compareInList},
	FieldOtaName: {stringContain, stringInWordList},

	FieldUid:        {compareEqual, compareNotEqual, compareInList},
	FieldUserName:   {stringContain, stringInWordList},
	FieldUserMobile: {stringContain, stringInWordList},
	FieldUserEMail:  {stringContain, stringInWordList},

	FieldBusinessLine:   {compareEqual, compareNotEqual, compareInList},
	FieldFirstCategory:  {compareEqual, compareNotEqual, compareInList},
	FieldSecondCategory: {compareEqual, compareNotEqual, compareInList},

	FieldSalesId:   {compareEqual, compareNotEqual, compareInList},
	FieldSalesName: {stringContain, stringInWordList},
}

// 所有支持的字段 使用 fieldSupportOperationMap 生成
var allFields []FieldType

func init() {
	allFields = make([]FieldType, len(fieldSupportOperationMap))

	var i = 0
	for key := range fieldSupportOperationMap {
		allFields[i] = key
		i += 1
	}
}

// 所有支持的字段
func AllSupportedField() []FieldType {
	return allFields
}

// 检测字段是否有效
func FieldCheckValid(fieldName FieldType) bool {
	for _, field := range allFields {
		if field == fieldName {
			return true
		}
	}
	return false
}
