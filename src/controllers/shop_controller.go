package controllers

import (
	"libs"
	"libs/shop"
	"outobjs"
	"strconv"
	"time"
	"utils"
)

// 商城 API
type ShopController struct {
	BaseController
}

func (c *ShopController) Prepare() {
	c.BaseController.Prepare()
}

func (c *ShopController) URLMapping() {
	c.Mapping("GetProvinces", c.GetProvinces)
	c.Mapping("GetAreas", c.GetAreas)
	c.Mapping("ShowItems", c.ShowItems)
	c.Mapping("GetItem", c.GetItem)
	c.Mapping("GetOrders", c.GetOrders)
	c.Mapping("GetOrder", c.GetOrder)
	c.Mapping("Stocks", c.Stocks)
	c.Mapping("Buy", c.Buy)
	c.Mapping("OrderCancel", c.OrderCancel)
	c.Mapping("Purchaseds", c.Purchaseds)
	c.Mapping("MyTickets", c.MyTickets)
}

// @Title 获取所有省份
// @Description 获取所有省份(数组)
// @Success 200 {object} outobjs.OutShopProvince
// @router /provinces [get]
func (c *ShopController) GetProvinces() {
	cache_key := "front_fast_cache.shop.provinces"
	c_obj := utils.GetLocalFastExpriesTimePartCache(cache_key)
	if c_obj != nil {
		c.Json(c_obj)
		return
	}

	out_ps := []*outobjs.OutShopProvince{}
	ps := shop.GetProvinces()
	for _, p := range ps {
		out_p := &outobjs.OutShopProvince{
			Id:   p.ProvinceID,
			Name: p.Province,
		}
		cs := shop.GetCities(p.ProvinceID)
		out_cs := []*outobjs.OutShopCity{}
		for _, c := range cs {
			out_c := &outobjs.OutShopCity{
				Id:   c.CityID,
				Name: c.City,
			}
			out_cs = append(out_cs, out_c)
		}
		out_p.Cities = out_cs
		out_ps = append(out_ps, out_p)
	}
	utils.SetLocalFastExpriesTimePartCache(2*time.Hour, cache_key, out_ps)
	c.Json(out_ps)
}

// @Title 获取所有地区
// @Description 获取所有地区(数组)
// @Param   cityid   	 path    string  true  "城市id"
// @Success 200 {object} outobjs.OutShopArea
// @router /areas [get]
func (c *ShopController) GetAreas() {
	cityid := c.GetString("cityid")
	if len(cityid) == 0 {
		c.Json([]*outobjs.OutShopArea{})
		return
	}
	areas := shop.GetAreas(cityid)
	out_areas := []*outobjs.OutShopArea{}
	for _, a := range areas {
		out_areas = append(out_areas, &outobjs.OutShopArea{
			Id:   a.AreaID,
			Name: a.Area,
		})
	}
	c.Json(out_areas)
}

// @Title 获取所有可买商品
// @Description 获取所有可买商品(数组)
// @Success 200 {object} outobjs.OutShopItem
// @router /items_show [get]
func (c *ShopController) ShowItems() {
	out_items := []*outobjs.OutShopItem{}
	shopp := shop.NewShop()
	items := shopp.GetItems()
	for _, item := range items {
		if item.IsView {
			out_items = append(out_items, outobjs.GetOutShopItem(item))
		}
	}
	c.Json(out_items)
}

// @Title 获取商品信息
// @Description 获取商品信息
// @Success 200 {object} outobjs.OutShopItem
// @router /item/:id([0-9]+) [get]
func (c *ShopController) GetItem() {
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("shop_get_fail", "SP1001", "id非法", ""))
		return
	}
	shopp := shop.NewShop()
	item := shopp.GetItem(id)
	if item == nil {
		c.Json(libs.NewError("shop_get_fail", "SP1002", "商品不存在", ""))
		return
	}
	c.Json(outobjs.GetOutShopItem(item))
}

// @Title 用户兑换记录
// @Description 用户兑换记录(数组)
// @Param   access_token   path  string  true  "access_token"
// @Param   page   path  int  false  "页"
// @Success 200 {object} outobjs.OutShopOrderPagedList
// @router /orders [get]
func (c *ShopController) GetOrders() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("shop_orders_premission_denied", UNAUTHORIZED_CODE, "必须登录", ""))
		return
	}
	page, _ := c.GetInt("page")
	if page <= 0 {
		page = 1
	}
	shopp := shop.NewShop()
	orders := shopp.GetOrders(uid, page, 20)
	out_orders := []*outobjs.OutShopOrder{}
	for _, order := range orders {
		out_orders = append(out_orders, outobjs.GetOutShopOrder(&order))
	}
	pls := &outobjs.OutShopOrderPagedList{
		CurrentPage: page,
		Orders:      out_orders,
	}
	c.Json(pls)
}

// @Title 用户兑换记录详情
// @Description 用户兑换记录详情
// @Param   access_token   path  string  true  "access_token"
// @Param   orderno   path  string  true  "订单号"
// @Success 200 {object} outobjs.OutShopOrderInfo
// @router /order [get]
func (c *ShopController) GetOrder() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("shop_order_premission_denied", UNAUTHORIZED_CODE, "必须登录", ""))
		return
	}
	orderno := utils.StripSQLInjection(c.GetString("orderno"))
	if len(orderno) == 0 {
		c.Json(libs.NewError("shop_order_orderno_fail", "SP1010", "订单号错误", ""))
		return
	}
	shopp := shop.NewShop()
	order := shopp.GetOrder(orderno)
	if order == nil {
		c.Json(libs.NewError("shop_order_notexist", "SP1011", "订单不存在", ""))
		return
	}
	out_info := outobjs.GetOutOrderInfo(order)
	c.Json(out_info)
}

// @Title 获取库存数
// @Description 获取库存数
// @Success 200 {object} libs.Error
// @router /stocks/:id([0-9]+) [get]
func (c *ShopController) Stocks() {
	itemId, _ := c.GetInt64(":id")
	shopp := shop.NewShop()
	item := shopp.GetItem(itemId)
	if item == nil {
		c.Json(libs.NewError("shop_order_stocks_noexist", "SP1020", "0", ""))
		return
	}
	stocks := strconv.Itoa(item.Stocks)
	c.Json(libs.NewError("shop_order_stocks", RESPONSE_SUCCESS, stocks, ""))
}

// @Title 购买商品
// @Description 购买商品
// @Param   access_token   path  string  true  "access_token"
// @Param   item_id   path  int  true  "商品id"
// @Param 	price_type path int  false  "货币类型(默认积分购买)"
// @Param   nums   path  int  true  "商品数量"
// @Param   remark   path  string  false  "描述"
// @Param   province   path  string  false  "省"
// @Param   city   path  string  false  "市"
// @Param   area   path  string  false  "区"
// @Param   address   path  string  false  "地址"
// @Param   receiver   path  string  false  "收货人"
// @Param   tel   path  string  false  "收货人电话"
// @Success 200 {object} libs.Error
// @router /buy [post]
func (c *ShopController) Buy() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("shop_buy_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能购买", ""))
		return
	}
	itemId, _ := c.GetInt64("item_id")
	nums, _ := c.GetInt("nums")
	remark, _ := utils.UrlDecode(c.GetString("remark"))
	province, _ := utils.UrlDecode(c.GetString("province"))
	city, _ := utils.UrlDecode(c.GetString("city"))
	area, _ := utils.UrlDecode(c.GetString("area"))
	address, _ := utils.UrlDecode(c.GetString("address"))
	receiver, _ := utils.UrlDecode(c.GetString("receiver"))
	tel, _ := utils.UrlDecode(c.GetString("tel"))
	priceType, _ := c.GetInt("price_type")

	if len(tel) > 15 {
		c.Json(libs.NewError("shop_buy_tel_fail", "SP1033", "电话号码太长", ""))
		return
	}
	switch priceType {
	case int(libs.PRICE_TYPE_CREDIT):
		break
	case int(libs.PRICE_TYPE_JING):
		break
	case int(libs.PRICE_TYPE_RMB):
		break
	default:
		priceType = int(libs.PRICE_TYPE_CREDIT) //默认积分购买
	}
	var pt = libs.PRICE_TYPE(priceType)

	shopp := shop.NewShop()
	item := shopp.GetItem(itemId)
	if item == nil {
		c.Json(libs.NewError("shop_buy_item_notexist", "SP1030", "商品不存在", ""))
		return
	}
	if nums <= 0 {
		c.Json(libs.NewError("shop_buy_nums_fail", "SP1031", "购买数量不能少于等于0", ""))
		return
	}

	_province := shop.GetProvinceName(province)
	_city := shop.GetCityName(city)
	_area := shop.GetAreaName(area)

	if item.ItemType == shop.ITEM_TYPE_ENTITY {
		if len(_province) == 0 || len(_city) == 0 || len(_area) == 0 || len(address) == 0 || len(receiver) == 0 || len(tel) == 0 {
			c.Json(libs.NewError("shop_buy_address_fail", "SP1032", "实物商品收货地址信息必须填写完整", ""))
			return
		}
	}

	addrInfo := &shop.ConsumerInfo{
		ProvinceName: _province,
		CityName:     _city,
		AreaName:     _area,
		Address:      address,
		Receiver:     receiver,
		Tel:          tel,
	}
	purchaser := shop.NewShopPurchaser()

	buyResult := purchaser.Buy(itemId, pt, uid, nums, remark, addrInfo)
	if buyResult.Error == nil {
		if buyResult.ItemType == shop.ITEM_TYPE_VIRTUAL {
			c.Json(libs.NewError("shop_buy_success", RESPONSE_SUCCESS, buyResult.Code, ""))
			return
		}
		if buyResult.ItemType == shop.ITEM_TYPE_ENTITY {
			c.Json(libs.NewError("shop_buy_success", RESPONSE_SUCCESS, buyResult.OrderNo, ""))
			return
		}
		if buyResult.ItemType == shop.ITEM_TYPE_TICKET {
			c.Json(libs.NewError("shop_buy_success", RESPONSE_SUCCESS, buyResult.OrderNo, ""))
			return
		}
	}
	errMsg := ""
	if buyResult.Error != nil {
		errMsg = buyResult.Error.Error()
	}
	c.Json(libs.NewError("shop_buy_fail", REPSONSE_FAIL, errMsg, ""))
}

// @Title 取消已购买的订单
// @Description 取消已购买的订单
// @Param   access_token   path  string  true  "access_token"
// @Param   order_no   path  string  true  "订单号"
// @Success 200 {object} libs.Error
// @router /order/cancel [post]
func (c *ShopController) OrderCancel() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("shop_order_cancel_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能取消订单", ""))
		return
	}
	orderNo := c.GetString("order_no")
	if len(orderNo) == 0 {
		c.Json(libs.NewError("shop_order_cancel_fail", "SP1040", "订单号不能为空", ""))
		return
	}
	shopp := shop.NewShop()
	order := shopp.GetOrder(orderNo)
	if order == nil {
		c.Json(libs.NewError("shop_order_cancel_fail", "SP1041", "订单不存在", ""))
		return
	}
	if order.Uid != uid {
		c.Json(libs.NewError("shop_order_cancel_fail", "SP1042", "被取消的订单买家必须和当前用户一致", ""))
		return
	}
	purchaser := shop.NewShopPurchaser()
	err := purchaser.ChangeOrderStatus(orderNo, shop.ORDER_STATUS_USERCANCEL, true, nil)
	if err != nil {
		c.Json(libs.NewError("shop_order_cancel_fail", "SP1045", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("shop_order_cancel_success", RESPONSE_SUCCESS, "订单取消成功", ""))
}

// @Title 获取用户购买数
// @Description 获取用户购买数
// @Param   access_token   path  string  true  "access_token"
// @Success 200 {object} libs.Error
// @router /purchaseds [get]
func (c *ShopController) Purchaseds() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("shop_purchaseds_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能查询", ""))
		return
	}
	shopp := shop.NewShop()
	counts := shopp.GetMemberCount(uid, shop.MEMBER_COUNT_UPDATE_ITEM_PURCHASEDS)
	c.Json(libs.NewError("shop_purchaseds", RESPONSE_SUCCESS, strconv.Itoa(counts), ""))
}

// @Title 已购买的电子票
// @Description 已购买的电子票
// @Param   access_token   path  string  true  "access_token"
// @Param   page   path  int  false  "page"
// @Success 200 {object} outobjs.OutShopTicketPagedList
// @router /my_tickets [get]
func (c *ShopController) MyTickets() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("shop_mytickets_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能查询", ""))
		return
	}
	page, _ := c.GetInt("page", 1)
	shopp := shop.NewShop()
	tickets := shopp.GetItemTickets(uid, page, 20)
	out_t := make([]*outobjs.OutShopTicket, len(tickets), len(tickets))
	for i, t := range tickets {
		out_t[i] = outobjs.GetOutShopTicket(t)
	}
	c.Json(&outobjs.OutShopTicketPagedList{
		CurrentPage: page,
		Tickets:     out_t,
	})
}
