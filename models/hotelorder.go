package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	_ "github.com/antiWalker/golib/db/orm"
	"reflect"
)

//订单表结构
type HotelOrder struct {
	Id               int64   // 自增id
	PrepareId        string  // 业务线全局唯一ID
	BusiOrderId      int64   // 业务订单ID
	OrderId          string  // 订单ID
	Ota              string  // OTA
	OtaId            int32   // OTA ID
	CheckIn          string  // 入住日期
	CheckOut         string  // 离店日期
	RoomNights       int8    // 间夜数
	Status           int32   // 订单状态
	CommentStatus    int8    // 评论状态
	PoiId            int32   // POI ID
	Promo            float64 // 补贴
	Honey            float64 // 蜂蜜
	CouponDiscount   float64 // 优惠券金额
	OtaPromo         float64 // 供应商补贴
	MfwPromoDiscount float64 // 折扣金额
	MfwPromo         float64 // 折扣
	Profit           float64 // 利润
	ActualRate       float64 // 实收价
	SaleRate         float64 // 销售价
	TotalRate        float64 // 成本价
	Country          string  // 所属目的地国家
	MddName          string  // 所属目的地名称
	RoomNum          int8    // 房间数
	CancelAble       int8    // 是否可取消
	ChildrenNum      int8    // 儿童数
	AdultNum         int8    // 成人数
	MfwRtId          string  // mfw母版房型id
	RoomTypeId       string  // ota房型id
	RoomId           string  // 价格计划id
	Zonetype         int8    // 区域类型
	Uid              int32   // UID
	Ip               string  // IP
	TrueName         string  // 联系人
	Phone            string  // 手机号
	Email            string  // 邮箱
	Uuid             string  // uuid
	OpenUdid         string  // open_udid
	PayType          int8    // 支付方式
	Platform         int8    // 平台
	Mddid            int32   // 目的地ID
	CouponId         string  // 优惠券ID
	HoneyStatus      int8    // 蜂蜜状态
	OrderCtime       string  // 订单创建时间
	Ctime            string  // 创建时间
	Mtime            string  // 更新时间
}

//hotel-order按照条件出sum
func (hotelOrder HotelOrder) SpitCount(job Job) (int64, []string, error) {
	v := func(o orm.Ormer) orm.QuerySeter {
		return o.QueryTable(new(HotelOrder))
	}
	return SplitCountInnerHotelOrder(job, v, hotelOrder)
}

//hotel-order按照条件出sum
func (hotelOrder HotelOrder) SpitSum(job Job) (float64, []string, error) {
	v := func(o orm.Ormer) orm.QuerySeter {
		qs := o.QueryTable(new(HotelOrder))

		return qs
	}
	return SpitSumInnerHotelOrder(job, v, hotelOrder)
}

func SplitCountInnerHotelOrder(job Job, countFunc CountFunc, tableRow interface{}) (int64, []string, error) {
	var id int64 = 0
	/*
		selectStruct := strings.Split(job.Select, "::")
		if len(selectStruct) != 2 {
			return 0, nil, errors.New("select format error")
		}

		var columnName = selectStruct[1]
	*/
	columnName := job.Select[7:]
	//转化为首字母大写
	name := getColName(columnName)

	//提取字类型
	dataType := getColumnDataType(tableRow, name)
	if dataType == reflect.Invalid {
		return 0, nil, errors.New("column " + columnName + " not allow count")
	}

	result := make(map[string]interface{})

	var ids = make([]string, 0)

	//游标遍历条件命中的全量
	for true {
		o := orm.NewOrm()

		var ml []SalesOrder
		qs := countFunc(o)

		qs = qs.Filter("id__gt", id)

		qs = filterMerge(qs, job)
		qs.Limit(100).OrderBy("id")
		qs.All(&ml)

		if len(ml) > 0 {
			for _, atom := range ml {
				id = atom.Id

				value := getCol(atom, name)

				key := getString(value)

				if _, exists := result[key]; exists == false {
					result[key] = value

					ids = append(ids, atom.PrepareId)
				}
			}
		} else {
			break
		}
	}

	return int64(len(result)), ids, nil
}

// sum 算法 TODO 公共抽离出来 抽离的条件是有些字段要保持统一 比如自增id 这样处理起来会方便HEAD
func SpitSumInnerHotelOrder(job Job, sumFunc SumFunc, tableRow interface{}) (float64, []string, error) {
	var id int64 = 0
	/*
		selectStruct := strings.Split(job.Select, "::")
		if len(selectStruct) != 2 {
			return 0, nil, errors.New("select format error")
		}

		var columnName = selectStruct[1]
	*/
	columnName := job.Select[5:]
	//转化为首字母大写
	name := getColName(columnName)

	//提取字类型
	dataType := getColumnDataType(tableRow, name)
	if dataType == reflect.Invalid {
		return 0, nil, errors.New("column " + columnName + " not allow sum")
	}

	var result = 0.0
	var ids = make([]string, 0)

	//游标遍历条件命中的全量
	for true {
		o := orm.NewOrm()

		var ml []HotelOrder
		qs := sumFunc(o)

		qs = qs.Filter("id__gt", id)

		qs = filterMerge(qs, job)
		qs.Limit(100).OrderBy("id")
		qs.All(&ml)

		if len(ml) > 0 {
			for _, atom := range ml {
				id = atom.Id

				value := getCol(atom, name)

				ids = append(ids, atom.PrepareId)

				//累加
				result += getNumber(dataType, value)

			}
		} else {
			break
		}
	}
	return result, ids, nil
}
