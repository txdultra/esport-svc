package outobjs

import (
	"fmt"
	"libs/shop"
	"libs/vars"
	"strconv"
	"strings"
	"time"
)

var shopp = shop.NewShop()

func GetShopImgUrls(imgs string) []string {
	imgIds := strings.Split(imgs, ",")
	imgUrls := []string{}
	for _, imgId := range imgIds {
		_id, _ := strconv.ParseInt(imgId, 10, 64)
		if _id > 0 {
			imgUrls = append(imgUrls, file.GetFileUrl(_id))
		}
	}
	return imgUrls
}

func GetOutShopItem(item *shop.Item) *OutShopItem {
	out_item := &OutShopItem{
		ItemId:        item.ItemId,
		Name:          item.Name,
		Description:   item.Description,
		PriceType:     item.PriceType,
		Price:         item.Price,
		Jings:         item.Jings,
		OriginalPrice: item.OriginalPrice,
		RmbPrice:      item.RmbPrice,
		Img:           item.Img,
		ImgUrl:        file.GetFileUrl(item.Img),
		ItemType:      item.ItemType,
		ItemState:     item.ItemState,
		Stocks:        item.Stocks,
		Sells:         item.Sells,
	}
	out_item.ShowingImgs = GetShopImgUrls(item.Imgs)
	return out_item
}

func GetOutShopOrder(order *shop.Order) *OutShopOrder {
	out_order := &OutShopOrder{
		OrderNo:     order.OrderNo,
		ItemId:      order.ItemId,
		ItemType:    order.ItemType,
		IssueType:   order.IssueType,
		CreateTime:  time.Unix(order.Ts, 0),
		Uid:         order.Uid,
		OrderStatus: order.OrderStatus,
		PayStatus:   order.PayStatus,
		Nums:        order.Nums,
		Price:       order.Price,
		TotalPrice:  order.TotalPrice,
		PriceType:   order.PriceType,
		SnapId:      order.SnapId,
		Remark:      order.Remark,
		Pay:         "积分",
		PayNo:       order.PayNo,
		Ex1:         order.Ex1,
		Ex2:         order.Ex2,
		Ex3:         order.Ex3,
	}
	item := shopp.GetItem(order.ItemId)
	if item != nil {
		out_order.Item = GetOutShopItem(item)
	}
	snap := shopp.GetItemSnap(order.SnapId)
	if snap != nil {
		out_order.Snap = GetOutShopItemSnap(snap)
	}
	return out_order
}

func GetOutOrderInfo(order *shop.Order) *OutShopOrderInfo {
	out_info := &OutShopOrderInfo{
		OrderNo:     order.OrderNo,
		ItemId:      order.ItemId,
		IssueType:   order.IssueType,
		OrderStatus: order.OrderStatus,
		PayStatus:   order.PayStatus,
		TotalPrice:  order.TotalPrice,
		PriceType:   order.PriceType,
		Remark:      order.Remark,
		Nums:        order.Nums,
		CreateTime:  time.Unix(order.Ts, 0),
		Ex1:         order.Ex1,
		Ex2:         order.Ex2,
		Ex3:         order.Ex3,
	}
	item := shopp.GetItem(order.ItemId)
	snap := shopp.GetItemSnap(order.SnapId)
	transport := shopp.GetOrderTransport(order.OrderNo)
	if snap != nil {
		out_info.Snap = GetOutShopItemSnap(snap)
	}
	if item != nil {
		out_info.Item = GetOutShopItem(item)
	}
	if transport != nil {
		out_info.Transport = &OutShopTransport{
			OrderNo:  order.OrderNo,
			TransNo:  transport.TransNo,
			Company:  getTransportCompanyName(transport.CompanyId),
			TransImg: "",
			Address:  fmt.Sprintf("%s%s%s%s", transport.Province, getCityName(transport.City), transport.Area, transport.Addr1),
			Receiver: transport.Receiver,
			Tel:      transport.Tel1,
		}
	}
	return out_info
}

func GetOutShopItemSnap(snap *shop.OrderItemSnap) *OutShopItemSnap {
	if snap == nil {
		return nil
	}
	tag := shopp.GetItemTag(snap.TagId)
	return &OutShopItemSnap{
		Name:        snap.Name,
		Description: snap.Description,
		PriceType:   snap.PriceType,
		Price:       snap.Price,
		Jings:       snap.Jings,
		ImgUrl:      file.GetFileUrl(snap.Img),
		ShowingImgs: GetShopImgUrls(snap.Imgs),
		TagId:       snap.TagId,
		Tag:         GetOutShopItemTag(tag),
		Attrs:       snap.GetAttrsMap(),
		RmbPrice:    snap.RmbPrice,
	}
}

func getCityName(city string) string {
	if city == "市辖区" || city == "省直辖县级行政单位" || city == "县" || city == "市" {
		return ""
	}
	return city
}

func getTransportCompanyName(id int) string {
	switch id {
	case 1:
		return "顺丰"
	case 2:
		return "中通"
	default:
		return "未知"
	}
}

func GetOutShopItemTag(tag *shop.ItemTag) *OutShopItemTag {
	if tag == nil {
		return nil
	}
	return &OutShopItemTag{
		Id:          tag.Id,
		Title:       tag.Title,
		Description: tag.Description,
		Img1:        tag.Img1,
		Img1Url:     file.GetFileUrl(tag.Img1),
		Img2:        tag.Img2,
		Img2Url:     file.GetFileUrl(tag.Img2),
		Img3:        tag.Img3,
		Img3Url:     file.GetFileUrl(tag.Img3),
	}
}

func GetOutShopTicket(ticket *shop.ItemTicket) *OutShopTicket {
	if ticket == nil {
		return nil
	}
	tag := shopp.GetItemTag(ticket.TagId)
	status := ticket.Status
	if status == shop.ITEM_TICKET_STATUS_NOUSE && ticket.EndTime < time.Now().Unix() {
		status = shop.ITEM_TICKET_STATUS_EXPIRED
	}

	return &OutShopTicket{
		Id:        ticket.Id,
		ItemId:    ticket.ItemId,
		Code:      ticket.Code,
		Img1:      ticket.Img1,
		Img1Url:   file.GetFileUrl(ticket.Img1),
		Img2:      ticket.Img2,
		Img2Url:   file.GetFileUrl(ticket.Img2),
		Img3:      ticket.Img3,
		Img3Url:   file.GetFileUrl(ticket.Img3),
		StartTime: time.Unix(ticket.StartTime, 0),
		EndTime:   time.Unix(ticket.EndTime, 0),
		TagId:     int(ticket.TagId),
		Tag:       GetOutShopItemTag(tag),
		Status:    status,
		TType:     ticket.TType,
		OrderNo:   ticket.OrderNo,
		BuyTime:   time.Unix(ticket.BuyTime, 0),
	}
}

func GetOutShopTransport(transport *shop.OrderTransport) *OutShopTransport {
	return &OutShopTransport{
		OrderNo:  transport.OrderNo,
		TransNo:  transport.TransNo,
		TransImg: "",
		Company:  getTransportCompanyName(transport.CompanyId),
		Address:  transport.Province + " " + transport.City + " " + transport.Area + " " + transport.Addr1 + " " + transport.Addr2,
		Receiver: transport.Receiver,
		Tel:      transport.Tel1 + " " + transport.Tel2,
	}
}

type OutShopItem struct {
	ItemId        int64           `json:"item_id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	PriceType     int             `json:"price_type"`
	Price         float64         `json:"price"`
	Jings         int64           `json:"jings"`
	OriginalPrice float64         `json:"original_price"`
	RmbPrice      float64         `json:"rmb_price"`
	Img           int64           `json:"img_id"`
	ImgUrl        string          `json:"img_url"`
	ShowingImgs   []string        `json:"showing_imgs_url"`
	ItemType      shop.ITEM_TYPE  `json:"item_type"`
	ItemState     shop.ITEM_STATE `json:"item_state"`
	Stocks        int             `json:"stocks"`
	Sells         int             `json:"sells"`
}

type OutShopTransport struct {
	OrderNo  string `json:"order_no"`
	TransNo  string `json:"transport_no"`
	TransImg string `json:"transport_imgurl"`
	Company  string `json:"transport_company"`
	Address  string `json:"address"`
	Receiver string `json:"receiver"`
	Tel      string `json:"tel"`
}

type OutShopItemSnap struct {
	Id          int64                  `json:"id"`
	ItemId      int64                  `json:"item_id"`
	Name        string                 `json:"item_name"`
	Description string                 `json:"item_description"`
	PriceType   int                    `json:"item_price_type"`
	Price       float64                `json:"item_price"`
	Jings       int64                  `json:"item_jings"`
	ImgUrl      string                 `json:"img_url"`
	ShowingImgs []string               `json:"showing_imgs_url"`
	TagId       int                    `json:"tag_id"`
	Tag         *OutShopItemTag        `json:"tag"`
	Attrs       map[string]interface{} `json:"attrs"`
	RmbPrice    float64                `json:"rmb_price"`
}

type OutShopOrder struct {
	OrderNo     string             `json:"order_no"`
	ItemId      int64              `json:"item_id"`
	ItemType    shop.ITEM_TYPE     `json:"item_type"`
	IssueType   shop.ISSUE_TYPE    `json:"issue_type"`
	CreateTime  time.Time          `json:"create_time"`
	Uid         int64              `json:"uid"`
	OrderStatus shop.ORDER_STATUS  `json:"order_status"`
	PayStatus   shop.PAY_STATUS    `json:"pay_status"`
	Nums        int                `json:"nums"`
	Price       float64            `json:"price"`
	TotalPrice  float64            `json:"total_price"`
	PriceType   vars.CURRENCY_TYPE `json:"price_type"`
	SnapId      int64              `json:"snap_id"`
	Snap        *OutShopItemSnap   `json:"item_snap"`
	Remark      string             `json:"remark"`
	Pay         string             `json:"pay"`
	PayNo       string             `json:"pay_no"`
	Ex1         string             `json:"ex1"`
	Ex2         string             `json:"ex2"`
	Ex3         string             `json:"ex3"`
	Item        *OutShopItem       `json:"item"`
}

type OutShopOrderPagedList struct {
	CurrentPage int             `json:"current_page"`
	Orders      []*OutShopOrder `json:"orders"`
}

type OutShopOrderInfo struct {
	OrderNo     string             `json:"order_no"`
	ItemId      int64              `json:"item_id"`
	IssueType   shop.ISSUE_TYPE    `json:"issue_type"`
	OrderStatus shop.ORDER_STATUS  `json:"order_status"`
	PayStatus   shop.PAY_STATUS    `json:"pay_status"`
	TotalPrice  float64            `json:"total_price"`
	PriceType   vars.CURRENCY_TYPE `json:"price_type"`
	Remark      string             `json:"remark"`
	Nums        int                `json:"nums"`
	CreateTime  time.Time          `json:"create_time"`
	Item        *OutShopItem       `json:"item"`
	Snap        *OutShopItemSnap   `json:"item_snap"`
	Transport   *OutShopTransport  `json:"transport"`
	Ex1         string             `json:"ex1"`
	Ex2         string             `json:"ex2"`
	Ex3         string             `json:"ex3"`
}

type OutShopProvince struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	Cities []*OutShopCity `json:"cities"`
}

type OutShopCity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OutShopArea struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OutShopTicketPagedList struct {
	CurrentPage int              `json:"current_page"`
	Tickets     []*OutShopTicket `json:"tickets"`
}

type OutShopTicket struct {
	Id        int64                   `json:"id"`
	ItemId    int64                   `json:"item_id"`
	Code      string                  `json:"code"`
	Img1      int64                   `json:"img1_id"`
	Img1Url   string                  `json:"img1_url"`
	Img2      int64                   `json:"img2_id"`
	Img2Url   string                  `json:"img2_url"`
	Img3      int64                   `json:"img3_id"`
	Img3Url   string                  `json:"img3_url"`
	StartTime time.Time               `json:"start_time"`
	EndTime   time.Time               `json:"end_time"`
	TagId     int                     `json:"tag_id"`
	Tag       *OutShopItemTag         `json:"tag"`
	Status    shop.ITEM_TICKET_STATUS `json:"status"`
	TType     shop.ITEM_TICKET_TYPE   `json:"ttype"`
	OrderNo   string                  `json:"order_no"`
	BuyTime   time.Time               `json:"buy_time"`
}

type OutShopItemTag struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Img1        int64  `json:"img1_id"`
	Img1Url     string `json:"img1_url"`
	Img2        int64  `json:"img2_id"`
	Img2Url     string `json:"img2_url"`
	Img3        int64  `json:"img3_id"`
	Img3Url     string `json:"img3_url"`
}

type OutShopItemForAdmin struct {
	ItemId        int64                  `json:"item_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	PriceType     int                    `json:"price_type"`
	Price         float64                `json:"price"`
	Jings         int64                  `json:"jings"`
	OriginalPrice float64                `json:"original_price"`
	RmbPrice      float64                `json:"rmb_price"`
	Img           int64                  `json:"img_id"`
	ImgUrl        string                 `json:"img_url"`
	Imgs          []int64                `json:"showing_imgs_id"`
	ShowingImgs   []string               `json:"showing_imgs_url"`
	ItemType      shop.ITEM_TYPE         `json:"item_type"`
	ItemState     shop.ITEM_STATE        `json:"item_state"`
	Stocks        int                    `json:"stocks"`
	Sells         int                    `json:"sells"`
	Ts            int64                  `json:"ts"`
	ModifyTs      int64                  `json:"modifyts"`
	DisplayOrder  int                    `json:"displayorder"`
	Enabled       bool                   `json:"enabled"`
	IsView        bool                   `json:"is_view"`
	Attrs         map[string]interface{} `json:"attrs"`
	TagId         int                    `json:"tag_id"`
}

type OutShopOrderPagedListForAdmin struct {
	CurrentPage int                     `json:"current_page"`
	Total       int                     `json:"total"`
	Size        int                     `json:"size"`
	Orders      []*OutShopOrderForAdmin `json:"orders"`
}

type OutShopOrderForAdmin struct {
	OrderNo     string             `json:"order_no"`
	ItemId      int64              `json:"item_id"`
	ItemType    shop.ITEM_TYPE     `json:"item_type"`
	IssueType   shop.ISSUE_TYPE    `json:"issue_type"`
	Ts          int64              `json:"ts"`
	Uid         int64              `json:"uid"`
	Member      *OutSimpleMember   `json:"member"`
	OrderStatus shop.ORDER_STATUS  `json:"order_status"`
	PayStatus   shop.PAY_STATUS    `json:"pay_status"`
	Nums        int                `json:"nums"`
	Price       float64            `json:"price"`
	TotalPrice  float64            `json:"total_price"`
	PriceType   vars.CURRENCY_TYPE `json:"price_type"`
	SnapId      int64              `json:"snap_id"`
	Snap        *OutShopItemSnap   `json:"snap"`
	Remark      string             `json:"remark"`
	PayId       int                `json:"payid"`
	PayNo       string             `json:"payno"`
	Ex1         string             `json:"ex1"`
	Ex2         string             `json:"ex2"`
	Ex3         string             `json:"ex3"`
}

type OutShopItemTagForAdmin struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Img1        int64     `json:"img1_id"`
	Img1Url     string    `json:"img1_url"`
	Img2        int64     `json:"img2_id"`
	Img2Url     string    `json:"img2_url"`
	Img3        int64     `json:"img3_id"`
	Img3Url     string    `json:"img3_url"`
	PostTime    time.Time `json:"post_time"`
}
