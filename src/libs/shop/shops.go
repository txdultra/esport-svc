package shop

import (
	"dbs"
	"fmt"
	"libs/dlock"
	"logs"
	"time"
	"utils"

	"github.com/astaxie/beego/orm"
)

func GetProvinces() []*Province {
	ps := []*Province{}
	for _, p := range provinces {
		if p.Enabled {
			ps = append(ps, p)
		}
	}
	return ps
}

func GetProvinceName(provinceId string) string {
	if p, ok := provinces[provinceId]; ok {
		return p.Province
	}
	return ""
}

func GetCities(fatherId string) []*City {
	cs := []*City{}
	for _, c := range citys[fatherId] {
		cs = append(cs, c)
	}
	return cs
}

func GetCityName(cityId string) string {
	for _, cc := range citys {
		for _, c := range cc {
			if c.CityID == cityId {
				return c.City
			}
		}
	}
	return ""
}

func GetAreas(fatherId string) []*Area {
	as := []*Area{}
	for _, c := range areas[fatherId] {
		as = append(as, c)
	}
	return as
}

func GetAreaName(areaId string) string {
	for _, aa := range areas {
		for _, a := range aa {
			if a.AreaID == areaId {
				return a.Area
			}
		}
	}
	return ""
}

func NewShop() *Shops {
	return &Shops{}
}

type Shops struct{}

func (s *Shops) checkItem(item *Item) error {
	if len(item.Name) == 0 {
		return fmt.Errorf("商品名称不能为空")
	}
	if int(item.PriceType) == 0 {
		return fmt.Errorf("未设置币种")
	}
	if int(item.ItemType) == 0 {
		return fmt.Errorf("未设置物品属性")
	}
	if int(item.ItemState) == 0 {
		return fmt.Errorf("未设置物品状态")
	}
	return nil
}

func (s *Shops) itemCacheKey(itemId int64) string {
	return fmt.Sprintf("mobile_shop_item_id:%d", itemId)
}

func (s *Shops) ItemLockKey(itemId int64) string {
	return fmt.Sprintf("item_%d", itemId)
}

func (s *Shops) itemAllCacheKey() string {
	return "mobile_shop_items_all"
}

func (s *Shops) CreateItem(item *Item) error {
	err := s.checkItem(item)
	if err != nil {
		return err
	}

	item.Ts = time.Now().Unix()
	item.ModifyTs = time.Now().Unix()

	//虚拟物品和电子票通过商品码同步库存
	if item.ItemType == ITEM_TYPE_VIRTUAL || item.ItemType == ITEM_TYPE_TICKET {
		item.Stocks = 0
	}
	item.updateAttrs()

	o := dbs.NewOrm(db_aliasname)
	id, err := o.Insert(item)
	if err != nil {
		return err
	}
	item.ItemId = id
	cache := utils.GetCache()
	cache.Add(s.itemCacheKey(item.ItemId), *item, 12*time.Hour)
	cache.Delete(s.itemAllCacheKey())
	return nil
}

func (s *Shops) UpdateItem(item *Item) error {
	err := s.checkItem(item)
	if err != nil {
		return err
	}

	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(s.ItemLockKey(item.ItemId))
	if err != nil {
		logs.Errorf("shop service update item get lock fail:%s", err.Error())
		return fmt.Errorf("系统繁忙,请稍后再试")
	}
	defer locker.Unlock()

	item.ModifyTs = time.Now().Unix()

	item.updateAttrs()
	o := dbs.NewOrm(db_aliasname)
	_, err = o.Update(item, "name", "description", "pricetype", "price", "oprice", "rprice", "img", "imgs", "itemstate",
		"modifyts", "displayorder", "isview", "exattrs", "tagid")
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Set(s.itemCacheKey(item.ItemId), *item, 12*time.Hour)
	return nil
}

func (s *Shops) UpdateItemStocks(itemId int64, addStockNums int, addSellNums int, locking bool) error {
	//分布式锁
	if !locking {
		lock := dlock.NewDistributedLock()
		locker, err := lock.Lock(s.ItemLockKey(itemId))
		if err != nil {
			logs.Errorf("shop service update item stocks get lock fail:%s", err.Error())
			return fmt.Errorf("系统繁忙,请稍后再试")
		}
		defer locker.Unlock()
	}

	o := dbs.NewOrm(db_aliasname)
	_, err := o.QueryTable(&Item{}).Filter("itemid", itemId).Update(orm.Params{
		"stocks": orm.ColValue(orm.Col_Add, addStockNums),
		"sells":  orm.ColValue(orm.Col_Add, addSellNums),
	})
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Delete(s.itemCacheKey(itemId))
	return nil
}

func (s *Shops) DeleteItem(itemId int64) error {
	item := s.GetItem(itemId)
	if item == nil {
		return fmt.Errorf("商品不存在")
	}
	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(s.ItemLockKey(item.ItemId))
	if err != nil {
		logs.Errorf("shop service delete item get lock fail:%s", err.Error())
		return fmt.Errorf("系统繁忙,请稍后再试")
	}
	defer locker.Unlock()

	item.Enabled = false
	o := dbs.NewOrm(db_aliasname)
	num, err := o.Update(item, "enabled")
	if err != nil {
		return err
	}
	if num <= 0 {
		return fmt.Errorf("商品不存在")
	}
	cache := utils.GetCache()
	cache.Delete(s.itemCacheKey(item.ItemId))
	cache.Delete(s.itemAllCacheKey())
	return nil
}

func (s *Shops) GetItem(itemId int64) *Item {
	cache := utils.GetCache()
	item := Item{}
	err := cache.Get(s.itemCacheKey(itemId), &item)
	if err == nil {
		return &item
	}
	o := dbs.NewOrm(db_aliasname)
	item.ItemId = itemId
	err = o.Read(&item)
	if err != nil {
		return nil
	}
	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(s.ItemLockKey(item.ItemId))
	if err != nil {
		logs.Errorf("shop service get item in db get lock fail:%s", err.Error())
		return nil
	}
	defer locker.Unlock()

	//防止重复设置cache
	err = cache.Get(s.itemCacheKey(itemId), &item)
	if err == nil {
		return &item
	}

	cache.Set(s.itemCacheKey(itemId), item, 12*time.Hour)
	return &item
}

func (s *Shops) GetItems() []*Item {
	cache := utils.GetCache()
	res := []int64{}
	err := cache.Get(s.itemAllCacheKey(), &res)
	if err != nil {
		o := dbs.NewOrm(db_aliasname)
		var items []*Item
		o.QueryTable(&Item{}).Filter("enabled", true).OrderBy("-displayorder", "-price").All(&items, "ItemId")
		for _, item := range items {
			res = append(res, item.ItemId)
		}
		cache.Set(s.itemAllCacheKey(), res, 1*time.Hour)
	}
	items := []*Item{}
	for _, id := range res {
		item := s.GetItem(id)
		if item != nil {
			items = append(items, item)
		}
	}
	return items
}

//商品快照
func (s *Shops) itemSnapCacheKey(snapId int64) string {
	return fmt.Sprintf("mobile_shop_itemsnap_id:%d", snapId)
}

func (s *Shops) ItemSnap(item *Item) (int64, error) {
	snap := &OrderItemSnap{
		Ts:          time.Now().Unix(),
		Name:        item.Name,
		Description: item.Description,
		PriceType:   item.PriceType,
		Price:       item.Price,
		Img:         item.Img,
		Imgs:        item.Imgs,
		TagId:       item.TagId,
		ItemId:      item.ItemId,
		ExAttrs:     item.ExAttrs,
		RmbPrice:    item.RmbPrice,
	}
	o := dbs.NewOrm(db_aliasname)
	id, err := o.Insert(snap)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Shops) GetItemSnap(snapId int64) *OrderItemSnap {
	cache := utils.GetCache()
	snap := OrderItemSnap{}
	err := cache.Get(s.itemSnapCacheKey(snapId), &snap)
	if err == nil {
		return &snap
	}
	o := dbs.NewOrm(db_aliasname)
	snap.SnapId = snapId
	err = o.Read(&snap)
	if err != nil {
		return nil
	}
	cache.Set(s.itemSnapCacheKey(snapId), snap, 2*time.Hour)
	return &snap
}

//运单
func (s *Shops) orderTransportCacheKey(orderNo string) string {
	return fmt.Sprintf("mobile_shop_transport_orderno:%s", orderNo)
}

func (s *Shops) CreateOrderTransport(trans *OrderTransport) error {
	if len(trans.OrderNo) == 0 {
		return fmt.Errorf("订单号或运单号不能为空")
	}
	trans.Ts = time.Now().Unix()
	o := dbs.NewOrm(db_aliasname)
	_, err := o.Insert(trans)
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Add(s.orderTransportCacheKey(trans.OrderNo), *trans, 6*time.Hour)
	return nil
}

func (s *Shops) UpdateOrderTransport(trans *OrderTransport) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.Update(trans)
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Set(s.orderTransportCacheKey(trans.OrderNo), *trans, 6*time.Hour)
	return nil
}

func (s *Shops) GetOrderTransport(orderNo string) *OrderTransport {
	cache := utils.GetCache()
	trans := OrderTransport{}
	err := cache.Get(s.orderTransportCacheKey(orderNo), &trans)
	if err == nil {
		return &trans
	}
	o := dbs.NewOrm(db_aliasname)
	trans.OrderNo = orderNo
	err = o.Read(&trans)
	if err != nil {
		return nil
	}
	cache.Set(s.orderTransportCacheKey(orderNo), trans, 2*time.Hour)
	return &trans
}

//虚拟商品码
func (s *Shops) CreateItemCodes(itemCodes []*ItemCode) (int, error) {
	var itemId int64 = 0
	if len(itemCodes) == 0 {
		return 0, fmt.Errorf("参数不能为空")
	}
	itemId = itemCodes[0].ItemId
	for _, item := range itemCodes {
		if item.ItemId != itemId {
			return 0, fmt.Errorf("批量添加的商品码必须是同一商品")
		}
	}

	o := dbs.NewOrm(db_aliasname)
	o.Begin()
	qs := o.QueryTable(&ItemCode{})
	i, _ := qs.PrepareInsert()
	err_count := 0
	var err error = nil
	for _, ic := range itemCodes {
		_, err = i.Insert(ic)
		if err != nil {
			err_count++
		}
	}
	i.Close() // 别忘记关闭 statement
	err = o.Commit()
	if err != nil {
		return 0, err
	}

	//更新库存
	s.UpdateItemStocks(itemId, len(itemCodes)-err_count, 0, false)

	return err_count, err
}

func (s *Shops) UseItemCode(itemId int64, orderNo string) (string, error) {
	item := s.GetItem(itemId)
	if item == nil {
		return "", fmt.Errorf("商品不存在")
	}
	if item.ItemType != ITEM_TYPE_VIRTUAL {
		return "", fmt.Errorf("商品必须是虚拟物品")
	}
	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(fmt.Sprintf("get_item_%d_code", itemId))
	if err != nil {
		logs.Errorf("shop service use item code get lock fail:%s", err.Error())
		return "", fmt.Errorf("系统繁忙,请稍后再试")
	}
	defer locker.Unlock()

	o := dbs.NewOrm(db_aliasname)
	var ic ItemCode
	err = o.QueryTable(&ItemCode{}).Filter("used", false).OrderBy("ts").Limit(1).One(&ic)
	if err != nil {
		return "", fmt.Errorf("商品已售罄")
	}
	ic.Used = true
	ic.UsedTs = time.Now().Unix()
	ic.OrderNo = orderNo
	num, _ := o.Update(&ic)
	if num <= 0 {
		return "", fmt.Errorf("领取失败")
	}
	return ic.Code, nil
}

//订单
func (s *Shops) checkOrder(order *Order) error {
	if order.ItemId == 0 {
		return fmt.Errorf("购买商品未设置")
	}
	if int(order.IssueType) == 0 {
		return fmt.Errorf("发放类型未设置")
	}
	if order.Uid <= 0 {
		return fmt.Errorf("未设置购买用户")
	}
	if order.SnapId <= 0 {
		return fmt.Errorf("未设置物品快照")
	}
	if order.Nums <= 0 {
		return fmt.Errorf("数量必须大于0")
	}
	return nil
}

func (s *Shops) orderCacheKey(orderNo string) string {
	return fmt.Sprintf("mobile_shop_order:%s", orderNo)
}

func (s *Shops) CreateOrder(order *Order) (string, error) {
	err := s.checkOrder(order)
	if err != nil {
		return "", err
	}
	o := dbs.NewOrm(db_aliasname)
	id, _ := o.Insert(&OrderIncerment{
		Ts: time.Now().Unix(),
	})
	if id <= 0 {
		return "", fmt.Errorf("订单创建失败")
	}

	orderNo := buildOrderNo(id)
	order.OrderNo = orderNo
	_, err = o.Insert(order)
	if err != nil {
		return "", err
	}
	cache := utils.GetCache()
	cache.Add(s.orderCacheKey(orderNo), *order, 12*time.Hour)
	cache.Delete(s.ordersByUidFirstPage(order.Uid))
	return orderNo, nil
}

func (s *Shops) UpdateOrder(order *Order) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.Update(order)
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Delete(s.orderCacheKey(order.OrderNo))
	cache.Delete(s.ordersByUidFirstPage(order.Uid))
	return nil
}

func (s *Shops) DeleteOrder(orderNo string) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.QueryTable(&Order{}).Filter("orderno", orderNo).Delete()
	return err
}

func (s *Shops) UpdateOrderStatus(orderNo string, status ORDER_STATUS) error {
	order := s.GetOrder(orderNo)
	if order == nil {
		return fmt.Errorf("订单不存在")
	}
	o := dbs.NewOrm(db_aliasname)
	_, err := o.QueryTable(&Order{}).Filter("orderno", orderNo).Update(orm.Params{
		"orderstatus": string(status),
	})
	cache := utils.GetCache()
	cache.Delete(s.orderCacheKey(orderNo))
	cache.Delete(s.ordersByUidFirstPage(order.Uid))

	if status == ORDER_STATUS_COMPLETED || status == ORDER_STATUS_SENDED {
		go s.UpdateMemberCount(order.Uid, 1, MEMBER_COUNT_UPDATE_ITEM_PURCHASEDS)
	}

	return err
}

func (s *Shops) GetOrder(orderNo string) *Order {
	cache := utils.GetCache()
	order := Order{}
	err := cache.Get(s.orderCacheKey(orderNo), &order)
	if err == nil {
		return &order
	}
	o := dbs.NewOrm(db_aliasname)
	order.OrderNo = orderNo
	err = o.Read(&order)
	if err != nil {
		return nil
	}
	cache.Set(s.orderCacheKey(orderNo), order, 12*time.Hour)
	return &order
}

func (s *Shops) ordersByUidFirstPage(uid int64) string {
	return fmt.Sprintf("mobile_shop_orders_uid_%d_p_1", uid)
}

func (s *Shops) GetOrders(uid int64, page int, size int) []Order {
	cache := utils.GetCache()
	if page <= 0 {
		page = 1
	}
	orders := []Order{}
	if page == 1 {
		key := s.ordersByUidFirstPage(uid)
		err := cache.Get(key, &orders)
		if err == nil {
			return orders
		}
	}

	o := dbs.NewOrm(db_aliasname)
	o.QueryTable(&Order{}).Filter("uid", uid).OrderBy("-ts").Limit(size, (page-1)*size).All(&orders)
	if page == 1 { //只缓存第一页
		cache.Set(s.ordersByUidFirstPage(uid), orders, 10*time.Hour)
	}
	return orders
}

func (s *Shops) GetOrdersByStatus(uid int64, page int, size int, orderStatus string) (int, []*Order) {
	var orders []*Order
	o := dbs.NewOrm(db_aliasname)
	query := o.QueryTable(&Order{})
	if uid > 0 {
		query = query.Filter("uid", uid)
	}
	if len(orderStatus) > 0 {
		query = query.Filter("orderstatus", orderStatus)
	}
	if page <= 0 {
		page = 1
	}
	total, _ := query.Count()
	query.OrderBy("-ts").Limit(size, (page-1)*size).All(&orders)
	return int(total), orders
}

const (
	cache_mobile_shop_counts_uid = "mobile_shop_counts_uid_%d"
)

func (s *Shops) UpdateMemberCount(uid int64, n int, item MEMBER_COUNT_UPDATE_ITEM) {
	o := dbs.NewOrm(db_aliasname)
	mc := &ShopMemberCount{}
	existed := o.QueryTable(mc).Filter("uid", uid).Exist()
	colname := ""
	switch item {
	case MEMBER_COUNT_UPDATE_ITEM_PURCHASEDS:
		colname = "count1"
		break
	default:
		break
	}
	if len(colname) == 0 {
		panic("shop count update item not existed")
	}
	if existed {
		o.QueryTable(mc).Filter("uid", uid).Update(orm.Params{
			colname: orm.ColValue(orm.Col_Add, n),
		})
	} else {
		tblname := mc.TableName()
		insertsql := fmt.Sprintf("insert into %s(uid,%s) values(%d,%d)", tblname, colname, uid, n)
		o.Raw(insertsql).Exec()
	}
	ckey := fmt.Sprintf(cache_mobile_shop_counts_uid, uid)
	cache := utils.GetCache()
	cache.Delete(ckey)
}

func (s *Shops) GetMemberCount(uid int64, item MEMBER_COUNT_UPDATE_ITEM) int {
	ckey := fmt.Sprintf(cache_mobile_shop_counts_uid, uid)
	cc := &ShopMemberCount{}
	cache := utils.GetCache()
	mc := func(_c *ShopMemberCount, _item MEMBER_COUNT_UPDATE_ITEM) int {
		if _c == nil {
			return 0
		}
		switch _item {
		case MEMBER_COUNT_UPDATE_ITEM_PURCHASEDS:
			return _c.Count1
		default:
			return 0
		}
	}
	err := cache.Get(ckey, cc)
	if err == nil {
		return mc(cc, item)
	} else {
		o := dbs.NewOrm(db_aliasname)
		_mc := &ShopMemberCount{Uid: uid}
		err := o.Read(_mc)
		cache.Add(ckey, _mc, 24*time.Hour)
		if err == nil {
			return mc(_mc, item)
		}
	}
	return 0
}

func (s *Shops) itemTagCacheKey(id int) string {
	return fmt.Sprintf("mobile_shop_itemtag_%d", id)
}

func (s *Shops) checkItemTag(tt *ItemTag) error {
	if len(tt.Title) == 0 {
		return fmt.Errorf("标题不能为空")
	}
	return nil
}

func (s *Shops) CreateItemTag(tt *ItemTag) error {
	err := s.checkItemTag(tt)
	if err != nil {
		return err
	}
	tt.PostTime = time.Now().Unix()
	o := dbs.NewOrm(db_aliasname)
	id, err := o.Insert(tt)
	if err == nil {
		tt.Id = int(id)
		cache := utils.GetCache()
		cache.Add(s.itemTagCacheKey(tt.Id), *tt, 48*time.Hour)
		return nil
	}
	return err
}

func (s *Shops) UpdateItemTag(tt *ItemTag) error {
	err := s.checkItemTag(tt)
	if err != nil {
		return err
	}
	o := dbs.NewOrm(db_aliasname)
	_, err = o.Update(tt)
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Set(s.itemTagCacheKey(tt.Id), *tt, 48*time.Hour)
	return nil
}

func (s *Shops) GetItemTag(id int) *ItemTag {
	if id <= 0 {
		return nil
	}
	ckey := s.itemTagCacheKey(id)
	cache := utils.GetCache()
	tt := &ItemTag{}
	err := cache.Get(ckey, tt)
	if err == nil {
		return tt
	}
	o := dbs.NewOrm(db_aliasname)
	tt.Id = id
	err = o.Read(tt)
	if err == nil {
		cache.Add(ckey, *tt, 48*time.Hour)
		return tt
	}
	return nil
}

func (s *Shops) GetItemTags() []*ItemTag {
	o := dbs.NewOrm(db_aliasname)
	var tags []*ItemTag
	o.QueryTable(&ItemTag{}).OrderBy("-posttime").All(&tags)
	return tags
}

func (s *Shops) CreateItemTicket(it *ItemTicket) error {
	if it.ItemId == 0 {
		return fmt.Errorf("商品id未指定")
	}
	o := dbs.NewOrm(db_aliasname)
	id, err := o.Insert(it)
	if err != nil {
		return err
	}
	it.Id = id
	return nil
}

func (s *Shops) MultiCreateItemTickets(itemId int64, its []*ItemTicket) error {
	if len(its) == 0 {
		return nil
	}
	o := dbs.NewOrm(db_aliasname)
	n, err := o.InsertMulti(len(its), its)

	//更新库存
	s.UpdateItemStocks(itemId, int(n), 0, false)

	return err
}

func (s *Shops) GetNoUseItemTicket(itemId int64, uid int64, orderNo string) (*ItemTicket, error) {
	item := s.GetItem(itemId)
	if item == nil {
		return nil, fmt.Errorf("商品不存在")
	}
	if item.ItemType != ITEM_TYPE_TICKET {
		return nil, fmt.Errorf("商品必须是电子票")
	}
	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(fmt.Sprintf("get_itemticket_%d", itemId))
	if err != nil {
		logs.Errorf("shop service use item code get lock fail:%s", err.Error())
		return nil, fmt.Errorf("系统繁忙,请稍后再试")
	}
	defer locker.Unlock()

	o := dbs.NewOrm(db_aliasname)
	var it ItemTicket
	err = o.QueryTable(&ItemTicket{}).Filter("uid", 0).OrderBy("-id").Limit(1).One(&it)
	if err != nil {
		return nil, fmt.Errorf("商品已售罄")
	}
	it.Uid = uid
	it.BuyTime = time.Now().Unix()
	it.OrderNo = orderNo
	num, _ := o.Update(&it)
	if num <= 0 {
		return nil, fmt.Errorf("购买失败")
	}
	return &it, nil
}

func (s *Shops) UsedItemTicket(it *ItemTicket) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.Update(it, "f_uid", "f_time", "status")
	return err
}

func (s *Shops) GetItemTicket(code string) *ItemTicket {
	it := ItemTicket{}
	o := dbs.NewOrm(db_aliasname)
	err := o.QueryTable(&it).Filter("code", code).One(&it)
	if err == nil {
		return &it
	}
	return nil
}

func (s *Shops) DeleteItemTicket(id int64) error {
	o := dbs.NewOrm(db_aliasname)
	o.QueryTable(&ItemTicket{}).Filter("id", id).Delete()
	return nil
}

func (s *Shops) GetItemTickets(uid int64, page int, size int) []*ItemTicket {
	o := dbs.NewOrm(db_aliasname)
	offset := (page - 1) * size
	var its []*ItemTicket
	o.QueryTable(&ItemTicket{}).Filter("uid", uid).OrderBy("status", "end_time").Limit(size, offset).All(&its)
	return its
}
