package shop

import (
	"encoding/json"
	"fmt"
	"libs"
	"time"
)

type ITEM_TYPE int

const (
	ITEM_TYPE_ENTITY  ITEM_TYPE = 1
	ITEM_TYPE_VIRTUAL ITEM_TYPE = 2
	ITEM_TYPE_TICKET  ITEM_TYPE = 3
)

type ITEM_STATE int

const (
	ITEM_STATE_NORMAL ITEM_STATE = 1
)

type ISSUE_TYPE int

const (
	ISSUE_TYPE_EXPRESS ISSUE_TYPE = 1
	ISSUE_TYPE_DIRECT  ISSUE_TYPE = 2
)

type ORDER_STATUS string

const (
	ORDER_STATUS_AUDITING   ORDER_STATUS = "auditing"
	ORDER_STATUS_AUDITSUCC  ORDER_STATUS = "auditsucc"
	ORDER_STATUS_AUDITFAIL  ORDER_STATUS = "auditfail"
	ORDER_STATUS_SENDED     ORDER_STATUS = "sended"
	ORDER_STATUS_USERCANCEL ORDER_STATUS = "usercancel"
	ORDER_STATUS_COMPLETED  ORDER_STATUS = "completed"
)

type PAY_STATUS string

const (
	PAY_STATUS_UNPAID PAY_STATUS = "unpaid"
	PAY_STATUS_PAIED  PAY_STATUS = "paid"
	PAY_STATUS_FAIL   PAY_STATUS = "fail"
)

var PAYID_CREDIT int = 1

type Item struct {
	ItemId        int64                  `orm:"column(itemid);pk"`
	Name          string                 `orm:"column(name)"`
	Description   string                 `orm:"column(description)"`
	PriceType     int                    `orm:"column(pricetype)"`
	Price         float64                `orm:"column(price)"`
	Jings         int64                  `orm:"column(jings)"`
	OriginalPrice float64                `orm:"column(oprice)"`
	RmbPrice      float64                `orm:"column(rprice)"`
	Img           int64                  `orm:"column(img)"`
	Imgs          string                 `orm:"column(imgs)"`
	ItemType      ITEM_TYPE              `orm:"column(itemtype)"`
	ItemState     ITEM_STATE             `orm:"column(itemstate)"`
	Ts            int64                  `orm:"column(ts)"`
	ModifyTs      int64                  `orm:"column(modifyts)"`
	DisplayOrder  int                    `orm:"column(displayorder)"`
	Stocks        int                    `orm:"column(stocks)"`
	Sells         int                    `orm:"column(sells)"`
	Enabled       bool                   `orm:"column(enabled)"`
	IsView        bool                   `orm:"column(isview)"`
	ExAttrs       string                 `orm:"column(exattrs)"`
	TagId         int                    `orm:"column(tagid)"`
	attrs         map[string]interface{} `orm:"-"`
}

func (self *Item) TableName() string {
	return "items"
}

func (self *Item) TableEngine() string {
	return "INNODB"
}

func (self *Item) updateAttrs() {
	data, _ := json.Marshal(self.attrs)
	self.ExAttrs = string(data)
}

func (self *Item) SetAttr(name string, val interface{}) {
	if self.attrs == nil {
		self.attrs = make(map[string]interface{})
	}
	self.attrs[name] = val
	self.updateAttrs()
}

func (self *Item) GetAttr(name string) interface{} {
	if self.attrs == nil {
		self.attrs = make(map[string]interface{})
		json.Unmarshal([]byte(self.ExAttrs), &self.attrs)
	}
	if obj, ok := self.attrs[name]; ok {
		return obj
	}
	return nil
}

func (self *Item) GetAttrsMap() map[string]interface{} {
	if self.attrs == nil {
		self.attrs = make(map[string]interface{})
		json.Unmarshal([]byte(self.ExAttrs), &self.attrs)
	}
	return self.attrs
}

type Order struct {
	OrderNo     string          `orm:"column(orderno);pk"`
	ItemId      int64           `orm:"column(itemid)"`
	ItemType    ITEM_TYPE       `orm:"column(itemtype)"` //新增
	IssueType   ISSUE_TYPE      `orm:"column(issuetype)"`
	Ts          int64           `orm:"column(ts)"`
	Uid         int64           `orm:"column(uid)"`
	OrderStatus ORDER_STATUS    `orm:"column(orderstatus)"`
	PayStatus   PAY_STATUS      `orm:"column(paystatus)"`
	Nums        int             `orm:"column(nums)"`
	Price       float64         `orm:"column(price)"`
	TotalPrice  float64         `orm:"column(totalprice)"`
	PriceType   libs.PRICE_TYPE `orm:"column(pricetype)"`
	SnapId      int64           `orm:"column(snap)"`
	Remark      string          `orm:"column(remark)"`
	PayId       int             `orm:"column(payid)"`
	PayNo       string          `orm:"column(payno)"`
	Ex1         string          `orm:"column(ex1)"`
	Ex2         string          `orm:"column(ex2)"`
	Ex3         string          `orm:"column(ex3)"`
}

func (self *Order) TableName() string {
	return "orders"
}

func (self *Order) TableEngine() string {
	return "INNODB"
}

type OrderIncerment struct {
	Id int64 `orm:"column(id);pk"`
	Ts int64 `orm:"column(ts)"`
}

func (self *OrderIncerment) TableName() string {
	return "order_incerment"
}

func (self *OrderIncerment) TableEngine() string {
	return "INNODB"
}

type OrderTransport struct {
	OrderNo   string `orm:"column(orderno);pk"`
	TransNo   string `orm:"column(transno)"`
	CompanyId int    `orm:"column(transid)"`
	Country   string `orm:"column(country)"`
	Province  string `orm:"column(province)"`
	City      string `orm:"column(city)"`
	Area      string `orm:"column(area)"`
	Addr1     string `orm:"column(addr1)"`
	Addr2     string `orm:"column(addr2)"`
	Receiver  string `orm:"column(receiver)"`
	Tel1      string `orm:"column(tel1)"`
	Tel2      string `orm:"column(tel2)"`
	Ts        int64  `orm:"column(ts)"`
}

func (self *OrderTransport) TableName() string {
	return "order_trans"
}

func (self *OrderTransport) TableEngine() string {
	return "INNODB"
}

type ItemCode struct {
	CodeId  int64  `orm:"column(codeid);pk"`
	ItemId  int64  `orm:"column(itemid)"`
	Code    string `orm:"column(code)"`
	Ts      int64  `orm:"column(ts)"`
	UsedTs  int64  `orm:"column(usedts)"`
	Used    bool   `orm:"column(used)"`
	OrderNo string `orm:"column(orderno)"`
}

func (self *ItemCode) TableName() string {
	return "item_codes"
}

func (self *ItemCode) TableEngine() string {
	return "INNODB"
}

type OrderItemSnap struct {
	SnapId      int64                  `orm:"column(snapid);pk"`
	Ts          int64                  `orm:"column(ts)"`
	Name        string                 `orm:"column(name)"`
	Description string                 `orm:"column(description)"`
	PriceType   int                    `orm:"column(pricetype)"`
	Price       float64                `orm:"column(price)"`
	Jings       int64                  `orm:"column(jings)"`
	Img         int64                  `orm:"column(img)"`
	Imgs        string                 `orm:"column(imgs)"`
	TagId       int                    `orm:"column(tagid)"`
	ItemId      int64                  `orm:"column(itemid)"`
	ExAttrs     string                 `orm:"column(exattrs)"`
	attrs       map[string]interface{} `orm:"-"`
	RmbPrice    float64                `orm:"column(rmbprice)"`
}

func (self *OrderItemSnap) TableName() string {
	return "order_itemsnap"
}

func (self *OrderItemSnap) TableEngine() string {
	return "INNODB"
}

func (self *OrderItemSnap) updateAttrs() {
	data, _ := json.Marshal(self.attrs)
	self.ExAttrs = string(data)
}

func (self *OrderItemSnap) GetAttrsMap() map[string]interface{} {
	if self.attrs == nil {
		self.attrs = make(map[string]interface{})
		json.Unmarshal([]byte(self.ExAttrs), &self.attrs)
	}
	return self.attrs
}

func (self *OrderItemSnap) SetAttr(name string, val interface{}) {
	if self.attrs == nil {
		self.attrs = make(map[string]interface{})
	}
	self.attrs[name] = val
	self.updateAttrs()
}

func (self *OrderItemSnap) GetAttr(name string) interface{} {
	if self.attrs == nil {
		self.attrs = make(map[string]interface{})
		json.Unmarshal([]byte(self.ExAttrs), &self.attrs)
	}
	if obj, ok := self.attrs[name]; ok {
		return obj
	}
	return nil
}

type Province struct {
	Id         int64  `orm:"column(id);pk"`
	ProvinceID string `orm:"column(provinceID)"`
	Province   string `orm:"column(province)"`
	Enabled    bool   `orm:"column(enabled)"`
}

func (self *Province) TableName() string {
	return "province"
}

func (self *Province) TableEngine() string {
	return "INNODB"
}

type City struct {
	Id     int64  `orm:"column(id);pk"`
	CityID string `orm:"column(cityID)"`
	City   string `orm:"column(city)"`
	Father string `orm:"column(father)"`
}

func (self *City) TableName() string {
	return "city"
}

func (self *City) TableEngine() string {
	return "INNODB"
}

type Area struct {
	Id     int64  `orm:"column(id);pk"`
	AreaID string `orm:"column(areaID)"`
	Area   string `orm:"column(area)"`
	Father string `orm:"column(father)"`
}

func (self *Area) TableName() string {
	return "area"
}

func (self *Area) TableEngine() string {
	return "INNODB"
}

func buildOrderNo(orderId int64) string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d001%09d", now.Year(), now.Month(), now.Day(), orderId)
}

type MEMBER_COUNT_UPDATE_ITEM int

const (
	MEMBER_COUNT_UPDATE_ITEM_PURCHASEDS MEMBER_COUNT_UPDATE_ITEM = 1
)

type ShopMemberCount struct {
	Uid    int64 `orm:"column(uid);pk"`
	Count1 int   `orm:"column(count1)"`
	Count2 int   `orm:"column(count2)"`
	Count3 int   `orm:"column(count3)"`
	Count4 int   `orm:"column(count4)"`
	Count5 int   `orm:"column(count5)"`
}

func (self *ShopMemberCount) TableName() string {
	return "shop_member_count"
}

func (self *ShopMemberCount) TableEngine() string {
	return "INNODB"
}

type ItemTag struct {
	Id          int    `orm:"column(id);pk"`
	Title       string `orm:"column(title)"`
	Description string `orm:"column(description)"`
	Img1        int64  `orm:"column(img1)"`
	Img2        int64  `orm:"column(img2)"`
	Img3        int64  `orm:"column(img3)"`
	PostTime    int64  `orm:"column(posttime)"`
}

func (self *ItemTag) TableName() string {
	return "item_tag"
}

func (self *ItemTag) TableEngine() string {
	return "INNODB"
}

type ITEM_TICKET_STATUS int

const (
	ITEM_TICKET_STATUS_NOUSE   ITEM_TICKET_STATUS = 1
	ITEM_TICKET_STATUS_USED    ITEM_TICKET_STATUS = 2
	ITEM_TICKET_STATUS_EXPIRED ITEM_TICKET_STATUS = 3
)

type ITEM_TICKET_TYPE int

const (
	ITEM_TICKET_TYPE_REAL   ITEM_TICKET_TYPE = 1
	ITEM_TICKET_TYPE_VIRUAL ITEM_TICKET_TYPE = 2
)

type ItemTicket struct {
	Id        int64              `orm:"column(id);pk"`
	ItemId    int64              `orm:"column(itemid)"`
	Code      string             `orm:"column(code)"`
	Pwd       string             `orm:"column(pwd)"`
	Uid       int64              `orm:"column(uid)"`
	FUid      int64              `orm:"column(f_uid)"`
	FTime     int64              `orm:"column(f_time)"`
	Img1      int64              `orm:"column(img_1)"`
	Img2      int64              `orm:"column(img_2)"`
	Img3      int64              `orm:"column(img_3)"`
	StartTime int64              `orm:"column(start_time)"`
	EndTime   int64              `orm:"column(end_time)"`
	TagId     int                `orm:"column(tag_id)"`
	Status    ITEM_TICKET_STATUS `orm:"column(status)"`
	TType     ITEM_TICKET_TYPE   `orm:"column(ttype)"`
	OrderNo   string             `orm:"column(orderno)"`
	BuyTime   int64              `orm:"column(buytime)"`
}

func (self *ItemTicket) TableName() string {
	return "item_tickets"
}

func (self *ItemTicket) TableEngine() string {
	return "INNODB"
}
