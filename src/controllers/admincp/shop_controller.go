package admincp

import (
	"controllers"
	"fmt"
	"libs"
	"libs/qrcode"
	"libs/shop"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 商城管理 API
type ShopCPController struct {
	AdminController
	storage libs.IFileStorage
}

func (c *ShopCPController) Prepare() {
	c.AdminController.Prepare()
	c.storage = libs.NewFileStorage()
}

func (c *ShopCPController) getItem(item *shop.Item) *outobjs.OutShopItemForAdmin {
	imgIdstrs := strings.Split(item.Imgs, ",")
	imgIds := []int64{}
	for _, imgId := range imgIdstrs {
		i, err := strconv.ParseInt(imgId, 10, 64)
		if err == nil {
			imgIds = append(imgIds, i)
		}
	}
	return &outobjs.OutShopItemForAdmin{
		ItemId:        item.ItemId,
		Name:          item.Name,
		Description:   item.Description,
		PriceType:     item.PriceType,
		Price:         item.Price,
		OriginalPrice: item.OriginalPrice,
		RmbPrice:      item.RmbPrice,
		Img:           item.Img,
		ImgUrl:        c.storage.GetFileUrl(item.Img),
		Imgs:          imgIds,
		ShowingImgs:   outobjs.GetShopImgUrls(item.Imgs),
		ItemType:      item.ItemType,
		ItemState:     item.ItemState,
		Stocks:        item.Stocks,
		Sells:         item.Sells,
		Ts:            item.Ts,
		ModifyTs:      item.ModifyTs,
		DisplayOrder:  item.DisplayOrder,
		Enabled:       item.Enabled,
		IsView:        item.IsView,
		Attrs:         item.GetAttrsMap(),
		TagId:         item.TagId,
	}
}

// @Title 获取商品信息
// @Description 获取商品信息
// @Param   item_id   path	int true  "商品id"
// @Success 200 {object} outobjs.OutShopItemForAdmin
// @router /item [get]
func (c *ShopCPController) GetItem() {
	itemid, _ := c.GetInt64("item_id")
	if itemid <= 0 {
		c.Json(libs.NewError("admincp_shop_getitem_fail", "GM030_081", "未提供商品id", ""))
		return
	}
	shopp := shop.NewShop()
	item := shopp.GetItem(itemid)
	if item == nil {
		c.Json(libs.NewError("admincp_shop_getitem_fail", "GM030_082", "商品不存在", ""))
		return
	}
	c.Json(c.getItem(item))
}

// @Title 获取所有商品
// @Description 获取所有商品
// @Success 200 {object} outobjs.OutShopItemForAdmin
// @router /items [get]
func (c *ShopCPController) GetItems() {
	shopp := shop.NewShop()
	items := shopp.GetItems()
	out_items := []*outobjs.OutShopItemForAdmin{}
	for _, item := range items {
		out_items = append(out_items, c.getItem(item))
	}
	c.Json(out_items)
}

// @Title 添加新商品
// @Description 添加新商品
// @Param   name   path	string true  "名称"
// @Param   description   path	string true  "描述"
// @Param   price_type   path	int true  "价格类型"
// @Param   price   path	float true  "价格"
// @Param   original_price   path	float true  "原价"
// @Param   rmb_price   path	float true  "人民币价格"
// @Param   img   path	int true  "图片"
// @Param   imgs   path	string true  "图片集(,隔开)"
// @Param   item_type   path	int true  "商品类型"
// @Param   item_state   path	int true  "商品状态"
// @Param   stocks   path	int true  "库存"
// @Param   display_order   path	int true  "排序"
// @Param   is_view   path	bool true  "是否显示"
// @Param   tag_id   path	int false  "标签编号"
// @Param   attrs   path	string false  "其他参数(格式:name1|value2,name1|value3)"
// @Success 200 {object} libs.Error
// @router /item/add [post]
func (c *ShopCPController) AddItem() {
	name, _ := utils.UrlDecode(c.GetString("name"))
	description, _ := utils.UrlDecode(c.GetString("description"))
	price_type, _ := c.GetInt("price_type")
	price, _ := c.GetFloat("price")
	ori_price, _ := c.GetFloat("original_price")
	rmb_price, _ := c.GetFloat("rmb_price")
	img, _ := c.GetInt64("img")
	imgs := c.GetString("imgs")
	item_type, _ := c.GetInt("item_type")
	item_state, _ := c.GetInt("item_state")
	stocks, _ := c.GetInt("stocks")
	display_order, _ := c.GetInt("display_order")
	is_view, _ := c.GetBool("is_view")
	tag_id, _ := c.GetInt("tag_id")
	attrstr, _ := utils.UrlDecode(c.GetString("attrs"))

	if len(name) == 0 {
		c.Json(libs.NewError("admincp_shop_add_fail", "GM030_001", "参数name错误", ""))
		return
	}
	if price_type <= 0 {
		c.Json(libs.NewError("admincp_shop_add_fail", "GM030_002", "参数price_type错误", ""))
		return
	}
	if price <= 0 {
		c.Json(libs.NewError("admincp_shop_add_fail", "GM030_003", "参数price错误", ""))
		return
	}
	if img <= 0 {
		c.Json(libs.NewError("admincp_shop_add_fail", "GM030_004", "参数img错误", ""))
		return
	}
	if item_type <= 0 {
		c.Json(libs.NewError("admincp_shop_add_fail", "GM030_005", "参数item_type错误", ""))
		return
	}
	if item_state <= 0 {
		c.Json(libs.NewError("admincp_shop_add_fail", "GM030_006", "参数item_state错误", ""))
		return
	}

	shopp := shop.NewShop()
	item := &shop.Item{
		Name:          name,
		Description:   description,
		PriceType:     shop.PRICE_TYPE(price_type),
		Price:         price,
		OriginalPrice: ori_price,
		RmbPrice:      rmb_price,
		Img:           img,
		Imgs:          imgs,
		ItemType:      shop.ITEM_TYPE(item_type),
		ItemState:     shop.ITEM_STATE(item_state),
		Ts:            time.Now().Unix(),
		ModifyTs:      time.Now().Unix(),
		DisplayOrder:  display_order,
		Stocks:        stocks,
		Enabled:       true,
		IsView:        is_view,
		TagId:         tag_id,
	}
	attrkvs := strings.Split(attrstr, ",")
	for _, kv := range attrkvs {
		_kv := strings.Split(kv, "|")
		if len(_kv) == 2 {
			item.SetAttr(_kv[0], _kv[1])
		}
	}
	err := shopp.CreateItem(item)
	if err == nil {
		c.Json(libs.NewError("admincp_shop_add_succ", controllers.RESPONSE_SUCCESS, "新商品添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_shop_add_fail", "GM030_007", "新商品添加失败:"+err.Error(), ""))
}

// @Title 更新商品
// @Description 更新商品
// @Param   item_id   path	int true  "商品id"
// @Param   name   path	string true  "名称"
// @Param   description   path	string true  "描述"
// @Param   price_type   path	int true  "价格类型"
// @Param   price   path	float true  "价格"
// @Param   original_price   path	float true  "原价"
// @Param   rmb_price   path	float true  "人民币价格"
// @Param   img   path	int true  "图片"
// @Param   imgs   path	string true  "图片集(,隔开)"
// @Param   item_state   path	int true  "商品状态"
// @Param   display_order   path	int true  "排序"
// @Param   is_view   path	bool true  "是否显示"
// @Param   tag_id   path	int false  "标签编号"
// @Param   attrs   path	string false  "其他参数(格式:name1|value2,name1|value3)"
// @Success 200 {object} libs.Error
// @router /item/update [post]
func (c *ShopCPController) UpdateItem() {
	item_id, _ := c.GetInt64("item_id")
	name, _ := utils.UrlDecode(c.GetString("name"))
	//&nbsp; bug
	description, _ := utils.UrlDecode(c.GetString("description"))
	description, _ = utils.UrlDecode(description)

	price_type, _ := c.GetInt("price_type")
	price, _ := c.GetFloat("price")
	ori_price, _ := c.GetFloat("original_price")
	rmb_price, _ := c.GetFloat("rmb_price")
	img, _ := c.GetInt64("img")
	imgs := c.GetString("imgs")
	//item_type, _ := c.GetInt("item_type")
	item_state, _ := c.GetInt("item_state")
	display_order, _ := c.GetInt("display_order")
	is_view, _ := c.GetBool("is_view")
	tag_id, _ := c.GetInt("tag_id")
	attrstr, _ := utils.UrlDecode(c.GetString("attrs"))

	if item_id <= 0 {
		c.Json(libs.NewError("admincp_shop_update_fail", "GM030_020", "参数item_id错误", ""))
		return
	}
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_shop_update_fail", "GM030_021", "参数name错误", ""))
		return
	}
	if price_type <= 0 {
		c.Json(libs.NewError("admincp_shop_update_fail", "GM030_022", "参数price_type错误", ""))
		return
	}
	if price <= 0 {
		c.Json(libs.NewError("admincp_shop_update_fail", "GM030_023", "参数price错误", ""))
		return
	}
	if img <= 0 {
		c.Json(libs.NewError("admincp_shop_update_fail", "GM030_024", "参数img错误", ""))
		return
	}
	//	if item_type <= 0 {
	//		c.Json(libs.NewError("admincp_shop_update_fail", "GM030_025", "参数item_type错误", ""))
	//		return
	//	}
	if item_state <= 0 {
		c.Json(libs.NewError("admincp_shop_update_fail", "GM030_026", "参数item_state错误", ""))
		return
	}
	shopp := shop.NewShop()
	item := shopp.GetItem(item_id)
	if item == nil {
		c.Json(libs.NewError("admincp_shop_update_fail", "GM030_027", "商品不存在", ""))
		return
	}
	item.Name = name
	item.Description = description
	item.PriceType = shop.PRICE_TYPE(price_type)
	item.Price = price
	item.OriginalPrice = ori_price
	item.RmbPrice = rmb_price
	item.Img = img
	item.Imgs = imgs
	//item.ItemType = shop.ITEM_TYPE(item_type)
	item.ItemState = shop.ITEM_STATE(item_state)
	item.ModifyTs = time.Now().Unix()
	item.DisplayOrder = display_order
	item.IsView = is_view
	item.TagId = tag_id

	attrkvs := strings.Split(attrstr, ",")
	for _, kv := range attrkvs {
		_kv := strings.Split(kv, "|")
		if len(_kv) == 2 {
			item.SetAttr(_kv[0], _kv[1])
		}
	}

	err := shopp.UpdateItem(item)
	if err == nil {
		c.Json(libs.NewError("admincp_shop_update_succ", controllers.RESPONSE_SUCCESS, "商品更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_shop_update_fail", "GM030_028", "商品更新失败:"+err.Error(), ""))
}

// @Title 删除商品
// @Description 删除商品
// @Param   item_id   path	int true  "商品id"
// @Success 200 {object} libs.Error
// @router /item/del [delete]
func (c *ShopCPController) DeleteItem() {
	item_id, _ := c.GetInt64("item_id")
	if item_id <= 0 {
		c.Json(libs.NewError("admincp_shop_delete_fail", "GM030_040", "参数item_id错误", ""))
		return
	}
	shopp := shop.NewShop()
	err := shopp.DeleteItem(item_id)
	if err == nil {
		c.Json(libs.NewError("admincp_shop_delete_succ", controllers.RESPONSE_SUCCESS, "商品删除成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_shop_delete_fail", "GM030_041", "商品删除失败:"+err.Error(), ""))
}

// @Title 增加商品库存
// @Description 增加商品库存
// @Param   item_id   path	int true  "商品id"
// @Param   oper   path	string true  "操作符(incr,decr)"
// @Param   nums   path	int true  "数量"
// @Success 200 {object} libs.Error
// @router /item/add_stock [post]
func (c *ShopCPController) AddItemStock() {
	itemid, _ := c.GetInt64("item_id")
	oper := c.GetString("oper")
	nums, _ := c.GetInt("nums")
	if itemid <= 0 {
		c.Json(libs.NewError("admincp_shop_addstock_fail", "GM030_060", "参数item_id错误", ""))
		return
	}
	if len(oper) == 0 {
		c.Json(libs.NewError("admincp_shop_addstock_fail", "GM030_061", "操作符不能为空", ""))
		return
	}
	if nums <= 0 {
		c.Json(libs.NewError("admincp_shop_addstock_fail", "GM030_062", "增加数量不能小于等于0", ""))
		return
	}
	shopp := shop.NewShop()
	item := shopp.GetItem(itemid)
	if item == nil {
		c.Json(libs.NewError("admincp_shop_addstock_fail", "GM030_063", "商品不存在", ""))
		return
	}
	if item.ItemType == shop.ITEM_TYPE_VIRTUAL {
		c.Json(libs.NewError("admincp_shop_addstock_fail", "GM030_064", "虚拟物品不能手动添加库存", ""))
		return
	}
	if oper == "decr" && nums > item.Stocks {
		c.Json(libs.NewError("admincp_shop_addstock_fail", "GM030_065", "减去的数量不能大于库存数", ""))
		return
	}
	if oper == "decr" {
		nums = -nums
	}
	err := shopp.UpdateItemStocks(itemid, nums, 0, true)
	if err == nil {
		c.Json(libs.NewError("admincp_shop_addstock_succ", controllers.RESPONSE_SUCCESS, "更改库存成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_shop_addstock_fail", "GM030_066", "更改库存失败:"+err.Error(), ""))
}

// @Title 载入商品码
// @Description 载入商品码
// @Param   item_id   path	int true  "商品id"
// @Param   codes   path	string true  "虚拟物品码(,分隔)"
// @Success 200 {object} libs.Error
// @router /item/add_codes [post]
func (c *ShopCPController) AddItemCodes() {
	item_id, _ := c.GetInt64("item_id")
	if item_id <= 0 {
		c.Json(libs.NewError("admincp_shop_addcodes_fail", "GM030_050", "参数item_id错误", ""))
		return
	}
	codestr := c.GetString("codes")
	codes := strings.Split(codestr, ",")
	if len(codes) == 0 {
		c.Json(libs.NewError("admincp_shop_addcodes_fail", "GM030_051", "商品码不能为空", ""))
		return
	}
	itemCodes := []*shop.ItemCode{}
	for _, code := range codes {
		_c := strings.Trim(code, " ")
		if len(_c) == 0 {
			continue
		}
		itemCodes = append(itemCodes, &shop.ItemCode{
			ItemId: item_id,
			Code:   _c,
			Ts:     time.Now().Unix(),
		})
	}
	shopp := shop.NewShop()
	i, err := shopp.CreateItemCodes(itemCodes)
	errinfo := fmt.Sprintf("导入失败个数为:%d,最后错误信息:", i)
	if err != nil {
		errinfo += err.Error()
	}

	c.Json(libs.NewError("admincp_shop_addcodes_succ", controllers.RESPONSE_SUCCESS, errinfo, ""))
}

// @Title 获取所有订单
// @Description 获取所有订单
// @Param   page   path	int false  "页"
// @Param   uid   path	int false  "用户uid"
// @Param   order_status   path	string false  "订单状态"
// @Success 200 {object} outobjs.OutShopOrderPagedListForAdmin
// @router /orders [get]
func (c *ShopCPController) GetOrders() {
	page, _ := c.GetInt("page")
	status := c.GetString("order_status")
	uid, _ := c.GetInt64("uid")
	shopp := shop.NewShop()
	total, orders := shopp.GetOrdersByStatus(uid, page, 20, status)
	out_orders := []*outobjs.OutShopOrderForAdmin{}
	for _, order := range orders {
		snap := shopp.GetItemSnap(order.SnapId)
		out_snap := outobjs.GetOutShopItemSnap(snap)
		out_orders = append(out_orders, &outobjs.OutShopOrderForAdmin{
			OrderNo:     order.OrderNo,
			ItemId:      order.ItemId,
			ItemType:    order.ItemType,
			IssueType:   order.IssueType,
			Ts:          order.Ts,
			Uid:         order.Uid,
			Member:      outobjs.GetOutSimpleMember(order.Uid),
			OrderStatus: order.OrderStatus,
			PayStatus:   order.PayStatus,
			Nums:        order.Nums,
			Price:       order.Price,
			TotalPrice:  order.TotalPrice,
			PriceType:   order.PriceType,
			SnapId:      order.SnapId,
			Snap:        out_snap,
			Remark:      order.Remark,
			PayId:       order.PayId,
			PayNo:       order.PayNo,
			Ex1:         order.Ex1,
			Ex2:         order.Ex2,
			Ex3:         order.Ex3,
		})
	}
	c.Json(&outobjs.OutShopOrderPagedListForAdmin{
		CurrentPage: page,
		Total:       total,
		Size:        20,
		Orders:      out_orders,
	})
}

// @Title 获取订单信息
// @Description 获取订单信息
// @Param   order_no   path	int false  "订单号"
// @Success 200 {object} outobjs.OutShopOrderForAdmin
// @router /order [get]
func (c *ShopCPController) GetOrder() {
	orderNo := c.GetString("order_no")
	if len(orderNo) == 0 {
		c.Json(libs.NewError("admincp_shop_get_order_fail", "GM030_055", "订单号不能为空", ""))
		return
	}
	shopp := shop.NewShop()
	order := shopp.GetOrder(orderNo)
	if order == nil {
		c.Json(libs.NewError("admincp_shop_get_order_fail", "GM030_056", "订单不存在", ""))
		return
	}
	snap := shopp.GetItemSnap(order.SnapId)
	out_snap := outobjs.GetOutShopItemSnap(snap)
	c.Json(&outobjs.OutShopOrderForAdmin{
		OrderNo:     order.OrderNo,
		ItemId:      order.ItemId,
		ItemType:    order.ItemType,
		IssueType:   order.IssueType,
		Ts:          order.Ts,
		Uid:         order.Uid,
		Member:      outobjs.GetOutSimpleMember(order.Uid),
		OrderStatus: order.OrderStatus,
		PayStatus:   order.PayStatus,
		Nums:        order.Nums,
		Price:       order.Price,
		TotalPrice:  order.TotalPrice,
		PriceType:   order.PriceType,
		SnapId:      order.SnapId,
		Snap:        out_snap,
		Remark:      order.Remark,
		PayId:       order.PayId,
		PayNo:       order.PayNo,
		Ex1:         order.Ex1,
		Ex2:         order.Ex2,
		Ex3:         order.Ex3,
	})
}

// @Title 设置订单状态
// @Description 设置订单状态
// @Param   order_no   path	string true  "订单号"
// @Param   order_status   path	string true  "订单状态"
// @Param   return_stock   path	bool false  "是否退回到库存"
// @Param   transid   path	int false  "货运公司编号"
// @Param   transno   path	string false  "运单编号"
// @Success 200 {object} libs.Error
// @router /order_status [post]
func (c *ShopCPController) SetOrderStatus() {
	orderNo := c.GetString("order_no")
	status := c.GetString("order_status")
	returnStock, _ := c.GetBool("return_stock")
	transid, _ := c.GetInt("transid")
	transno := c.GetString("transno")
	if len(orderNo) == 0 {
		c.Json(libs.NewError("admincp_shop_set_orderstatus_fail", "GM030_062", "订单号不能为空", ""))
		return
	}
	if len(status) == 0 {
		c.Json(libs.NewError("admincp_shop_set_orderstatus_fail", "GM030_063", "状态码不能为空", ""))
		return
	}

	setOrderStatus := shop.ORDER_STATUS(status)
	if setOrderStatus == shop.ORDER_STATUS_USERCANCEL || setOrderStatus == shop.ORDER_STATUS_AUDITING {
		c.Json(libs.NewError("admincp_shop_set_orderstatus_fail", "GM030_060", "不支持的状态", ""))
		return
	}
	var transInfo *shop.TransportInfo = nil
	if setOrderStatus == shop.ORDER_STATUS_SENDED {
		transInfo = &shop.TransportInfo{
			TransId: transid,
			TransNo: transno,
		}
	}
	purchaser := shop.NewShopPurchaser()
	err := purchaser.ChangeOrderStatus(orderNo, setOrderStatus, returnStock, transInfo)
	if err == nil {
		c.Json(libs.NewError("admincp_shop_set_orderstatus_succ", controllers.RESPONSE_SUCCESS, "设置订单状态成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_shop_set_orderstatus_fail", "GM030_061", "设置订单状态失败:"+err.Error(), ""))
}

// @Title 商品快照
// @Description 商品快照
// @Param   snap_id   path	int true  "快照id"
// @Success 200 {object} outobjs.OutShopItemSnap
// @router /snap [get]
func (c *ShopCPController) Snap() {
	snapid, _ := c.GetInt64("snap_id")
	if snapid <= 0 {
		c.Json(libs.NewError("admincp_shop_snap_fail", "GM030_070", "id非法", ""))
		return
	}
	shopp := shop.NewShop()
	snap := shopp.GetItemSnap(snapid)
	if snap == nil {
		c.Json(libs.NewError("admincp_shop_snap_fail", "GM030_071", "快照不存在", ""))
		return
	}
	out_snap := outobjs.GetOutShopItemSnap(snap)
	c.Json(out_snap)
}

// @Title 生成电子票
// @Description 生成电子票
// @Param   item_id   path	int true  "商品id"
// @Param   nums   path	int true  "电子票数量"
// @Param   img1   path	int false  "电子票图片1"
// @Param   img2   path	int false  "电子票图片2"
// @Param   img3   path	int false  "电子票图片3"
// @Param   start_time_attrname   path	string true  "开始时间名称"
// @Param   end_time_attrname   path	string true  "结束时间名称"
// @Success 200 {object} libs.Error
// @router /item/build_tickets [post]
func (c *ShopCPController) BuildTickets() {
	itemId, _ := c.GetInt64("item_id", 0)
	nums, _ := c.GetInt("nums", 0)
	img1, _ := c.GetInt64("img1", 0)
	img2, _ := c.GetInt64("img2", 0)
	img3, _ := c.GetInt64("img3", 0)
	stimeName := c.GetString("start_time_attrname")
	etimeName := c.GetString("end_time_attrname")

	if itemId <= 0 || nums <= 0 {
		c.Json(libs.NewError("admincp_shop_build_ticket_fail", "GM030_080", "参数错误", ""))
		return
	}
	if len(stimeName) == 0 || len(etimeName) == 0 {
		c.Json(libs.NewError("admincp_shop_build_ticket_fail", "GM030_081", "时间名称错误", ""))
		return
	}
	shopp := shop.NewShop()
	item := shopp.GetItem(itemId)
	if item == nil || item.ItemType != shop.ITEM_TYPE_TICKET {
		c.Json(libs.NewError("admincp_shop_build_ticket_fail", "GM030_082", "商品不存在或不是电子票", ""))
		return
	}
	loc, _ := time.LoadLocation("Local")
	_val := item.GetAttr(stimeName).(string)
	stime, terr := time.ParseInLocation("2006-01-02 15:04", _val, loc) //time.Parse("2006-01-02 15:04", _val)
	_val = item.GetAttr(etimeName).(string)
	etime, terr := time.ParseInLocation("2006-01-02 15:04", _val, loc) //time.Parse("2006-01-02 15:04", _val)
	if terr != nil {
		c.Json(libs.NewError("admincp_shop_build_ticket_fail", "GM030_083", "日期转换失败", ""))
		return
	}

	tickets := []*shop.ItemTicket{}
	for i := 0; i < nums; i++ {
		_code, _ := qrcode.GetClientCode(shop.QRCODE_TICKET_FLAG, utils.RandomStrings(10))
		tickets = append(tickets, &shop.ItemTicket{
			ItemId:    item.ItemId,
			Code:      _code,
			Img1:      img1,
			Img2:      img2,
			Img3:      img3,
			StartTime: stime.Unix(),
			EndTime:   etime.Unix(),
			TagId:     item.TagId,
			Status:    shop.ITEM_TICKET_STATUS_NOUSE,
			TType:     shop.ITEM_TICKET_TYPE_VIRUAL,
		})
	}
	err := shopp.MultiCreateItemTickets(itemId, tickets)
	if err == nil {
		c.Json(libs.NewError("admincp_shop_build_ticket_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_shop_build_ticket_fail", "GM030_084", "添加失败:"+err.Error(), ""))
}

// @Title 添加商品标签
// @Description 添加商品标签
// @Param   title   path  string  true  "标题"
// @Param   description   path	string true  "描述"
// @Param   img1   path	int false  "图片1"
// @Param   img2   path	int false  "图片2"
// @Param   img3   path	int false  "图片3"
// @Success 200 {object} libs.Error
// @router /tag/add [post]
func (c *ShopCPController) AddTag() {
	title, _ := utils.UrlDecode(c.GetString("title"))
	description, _ := utils.UrlDecode(c.GetString("description"))
	img1, _ := c.GetInt64("img1")
	img2, _ := c.GetInt64("img2")
	img3, _ := c.GetInt64("img3")

	if len(title) == 0 || len(description) == 0 {
		c.Json(libs.NewError("admincp_shop_addtag_fail", "GM030_090", "参数错误", ""))
		return
	}
	tag := &shop.ItemTag{
		Title:       title,
		Description: description,
		Img1:        img1,
		Img2:        img2,
		Img3:        img3,
	}
	shopp := shop.NewShop()
	err := shopp.CreateItemTag(tag)
	if err == nil {
		c.Json(libs.NewError("admincp_shop_addtag_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_shop_addtag_fail", "GM030_091", "添加失败:"+err.Error(), ""))
}

// @Title 更新商品标签
// @Description 更新商品标签
// @Param   id   path  int  true  "标签id"
// @Param   title   path  string  true  "标题"
// @Param   description   path	string true  "描述"
// @Param   img1   path	int false  "图片1"
// @Param   img2   path	int false  "图片2"
// @Param   img3   path	int false  "图片3"
// @Success 200 {object} libs.Error
// @router /tag/update [post]
func (c *ShopCPController) UpdateTag() {
	id, _ := c.GetInt("id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	description, _ := utils.UrlDecode(c.GetString("description"))
	img1, _ := c.GetInt64("img1")
	img2, _ := c.GetInt64("img2")
	img3, _ := c.GetInt64("img3")

	if id <= 0 || len(title) == 0 || len(description) == 0 {
		c.Json(libs.NewError("admincp_shop_updatetag_fail", "GM030_100", "参数错误", ""))
		return
	}
	shopp := shop.NewShop()
	tag := shopp.GetItemTag(id)
	if tag == nil {
		c.Json(libs.NewError("admincp_shop_updatetag_fail", "GM030_101", "参数错误", ""))
		return
	}
	tag.Title = title
	tag.Description = description
	tag.Img1 = img1
	tag.Img2 = img2
	tag.Img3 = img3
	err := shopp.UpdateItemTag(tag)
	if err == nil {
		c.Json(libs.NewError("admincp_shop_updatetag_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_shop_updatetag_fail", "GM030_102", "更新失败:"+err.Error(), ""))
}

func (c *ShopCPController) getOutItemTag(tag *shop.ItemTag) *outobjs.OutShopItemTagForAdmin {
	if tag == nil {
		return nil
	}
	return &outobjs.OutShopItemTagForAdmin{
		Id:          tag.Id,
		Title:       tag.Title,
		Description: tag.Description,
		Img1:        tag.Img1,
		Img1Url:     c.storage.GetFileUrl(tag.Img1),
		Img2:        tag.Img2,
		Img2Url:     c.storage.GetFileUrl(tag.Img2),
		Img3:        tag.Img3,
		Img3Url:     c.storage.GetFileUrl(tag.Img3),
		PostTime:    time.Unix(tag.PostTime, 0),
	}
}

// @Title 获取商品标签
// @Description 获取商品标签
// @Param   id   path  int  true  "标签id"
// @Success 200 {object} outobjs.OutShopItemTagForAdmin
// @router /tag/:id([0-9]+) [get]
func (c *ShopCPController) GetTag() {
	id, _ := c.GetInt(":id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_shop_get_fail", "GM030_110", "参数错误", ""))
		return
	}
	shopp := shop.NewShop()
	tag := shopp.GetItemTag(id)
	if tag == nil {
		c.Json(libs.NewError("admincp_shop_get_fail", "GM030_111", "标签不存在", ""))
		return
	}
	c.Json(c.getOutItemTag(tag))
}

// @Title 获取商品标签集合
// @Description 获取商品标签集合
// @Success 200 {object} outobjs.OutShopItemTagForAdmin
// @router /tag/all [get]
func (c *ShopCPController) GetTags() {
	shopp := shop.NewShop()
	tags := shopp.GetItemTags()
	out_t := make([]*outobjs.OutShopItemTagForAdmin, len(tags), len(tags))
	for i, tag := range tags {
		out_t[i] = c.getOutItemTag(tag)
	}
	c.Json(out_t)
}
