package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	_ "github.com/antiWalker/golib/db/orm"
	"reflect"
)

/*-------t_risk_engine_sales_order---------*/

//订单表结构
type SalesOrder struct {
	Id                     int64   `json:"id"`                        // 自增ID
	PrepareId              string  `json:"prepare_id"`                // 业务线全局唯一ID
	OrderId                string  `json:"order_id"`                  // 旅游订单ID
	OrderIp                string  `json:"order_ip"`                  // 下单IP（兼容 IP v6）
	OrderTime              string  `json:"order_time"`                // 下单时间
	PayTime                string  `json:"pay_time"`                  // 支付时间
	ConfirmInventoryTime   string  `json:"confirm_inventory_time"`    // 确认库存时间
	IssueTime              string  `json:"issue_time"`                // 出单/出签时间
	TripDate               string  `json:"trip_date"`                 // 出行日期
	FromPoiList            string  `json:"from_poi_list"`             // 出发地列表（json串）
	ToPoiList              string  `json:"to_poi_list"`               // 目的地列表（json串）
	GoodsNum               int16   `json:"goods_num"`                 // 订单购买的商品数量
	TripDays               int8    `json:"trip_days"`                 // 行程天数
	OrderStatus            string  `json:"order_status"`              // 订单状态
	RefundStatus           string  `json:"refund_status"`             // 退款状态
	RefundType             string  `json:"refund_type"`               // 退款原因
	FlagPayment            string  `json:"flag_payment"`              // 支付状态
	CommentStatus          string  `json:"comment_status"`            // 点评状态
	BusinessModel          string  `json:"business_model"`            // 交易模式
	OrderChannel           string  `json:"order_channel"`             // 下单渠道
	PromotionChannel       string  `json:"promotion_channel"`         // 推广渠道
	BusinessType           string  `json:"business_type"`             // 业务线
	MainCategory           string  `json:"main_category"`             // 一级品类
	SubCategory            string  `json:"sub_category"`              // 二级品类
	SalesId                int32   `json:"sales_id"`                  // 产品ID
	SalesName              string  `json:"sales_name"`                // 产品名称
	SkuId                  int32   `json:"sku_id"`                    // 库存ID
	SkuName                string  `json:"sku_name"`                  // 库存名称
	OtaId                  int32   `json:"ota_id"`                    // OTA ID
	OtaName                string  `json:"ota_name"`                  // 店铺名称
	OtaCompany             string  `json:"ota_company"`               // 公司全称
	OrderPrice             float64 `json:"order_price"`               // 订单金额
	PaymentFee             float64 `json:"payment_fee"`               // 支付金额
	CommissionAmount       float64 `json:"commission_amount"`         // 佣金金额
	FlashSaleSupplement    float64 `json:"flash_sale_supplement"`     // 平台蜂抢优惠补贴
	VipPriceSupplement     float64 `json:"vip_price_supplement"`      // 平台会员专享价补贴
	NOffReduce             float64 `json:"n_off_reduce"`              // 商家N人优惠
	EarlyBirdReduce        float64 `json:"early_bird_reduce"`         // 商家早鸟优惠
	EarlyBirdSupplement    float64 `json:"early_bird_supplement"`     // 商家早鸟优惠
	TimeLimitedReduce      float64 `json:"time_limited_reduce"`       // 商家早鸟优惠
	TimeLimitedSupplement  float64 `json:"time_limited_supplement"`   // 商家早鸟优惠
	HoneySupplement        float64 `json:"honey_supplement"`          // 商家早鸟优惠
	NOffSupplement         float64 `json:"n_off_supplement"`          // 马蜂窝N人N折补贴
	RedPacketSupplement    float64 `json:"red_packet_supplement"`     // 红包补贴
	GoldenCardSupplement   float64 `json:"golden_card_supplement"`    // 金卡补贴
	McodeSupplement        float64 `json:"mcode_supplement"`          // M码补贴
	OtaCouponReduce        float64 `json:"ota_coupon_reduce"`         // 商家早鸟优惠
	OtaMfwCouponSupplement float64 `json:"ota_mfw_coupon_supplement"` // 平台商家券补贴
	MfwOtaCouponSupplement float64 `json:"mfw_ota_coupon_supplement"` // 商家早鸟优惠
	MfwCouponSupplement    float64 `json:"mfw_coupon_supplement"`     // 马蜂窝优惠券补贴
	BargainReduce          float64 `json:"bargain_reduce"`            // 商家早鸟优惠
	BargainSupplement      float64 `json:"bargain_supplement"`        // 商家早鸟优惠
	OtaTotalReduce         float64 `json:"ota_total_reduce"`          // 商家优惠总金额
	OtherReduce            float64 `json:"other_reduce"`              // 其他优惠金额
	MfwTotalSupplement     float64 `json:"mfw_total_supplement"`      // 商家早鸟优惠
	MfwCouponSn            string  `json:"mfw_coupon_sn"`             // 马蜂窝优惠券序列号
	MfwCouponName          string  `json:"mfw_coupon_name"`           // 马蜂窝优惠券名称
	Uid                    int64   `json:"uid"`                       // 下单用户ID
	UserName               string  `json:"user_name"`                 // 下单用户名称
	Phone                  string  `json:"phone"`                     // 下单用户手机号
	Email                  string  `json:"email"`                     // 下单用户邮箱
	Wechat                 string  `json:"wechat"`                    // 下单用户微信
	TripUsers              string  `json:"trip_users"`                // 出行人信息
	Uuid                   string  `json:"uuid"`                      // uuid
	OpenUdid               string  `json:"open_udid"`                 // open-udid
	PayAccount             string  `json:"pay_account"`               // 支付账号
	PayChannel             string  `json:"pay_channel"`               // 支付方式
}

//sales-order按照条件出sum
func (salesOrder SalesOrder) SpitCount(job Job) (int64, []string, error) {
	v := func(o orm.Ormer) orm.QuerySeter {
		return o.QueryTable(new(SalesOrder))
	}
	return SplitCountInnerSalesOrder(job, v, salesOrder)
}

//sales-order按照条件出sum
func (salesOrder SalesOrder) SpitSum(job Job) (float64, []string, error) {
	v := func(o orm.Ormer) orm.QuerySeter {
		qs := o.QueryTable(new(SalesOrder))

		return qs
	}
	return SpitSumInnerSalesOrder(job, v, salesOrder)
}

func SplitCountInnerSalesOrder(job Job, countFunc CountFunc, tableRow interface{}) (int64, []string, error) {
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

func SpitSumInnerSalesOrder(job Job, sumFunc SumFunc, tableRow interface{}) (float64, []string, error) {
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

		var ml []SalesOrder
		qs := sumFunc(o)

		qs = qs.Filter("id__gt", id)

		qs = filterMerge(qs, job)
		qs.Limit(100).OrderBy("id")
		qs.All(&ml)

		if len(ml) > 0 {
			for _, atom := range ml {
				id = atom.Id

				value := getCol(atom, name)

				//累加
				result += getNumber(dataType, value)

				ids = append(ids, atom.PrepareId)

			}
		} else {
			break
		}
	}
	return result, ids, nil
}
