package shop

import (
	"fmt"
	"libs/qrcode"
	"time"
	"utils"
)

const (
	QRCODE_TICKET_FLAG = "shop_ticket"
)

type ShopTicketQRCodeService struct{}

func (s ShopTicketQRCodeService) Flag() string {
	return QRCODE_TICKET_FLAG
}

func (s ShopTicketQRCodeService) EncodeCode(code string) string {
	return code
}

func (s ShopTicketQRCodeService) DecodeCode(fromUid int64, code string) (*qrcode.QRCodeResult, error) {
	_code := utils.StripSQLInjection(code)
	if _code != code {
		return nil, fmt.Errorf("二维码错误")
	}
	shopp := NewShop()
	it := shopp.GetItemTicket(code)
	if it == nil {
		return &qrcode.QRCodeResult{
			Result: "ticket_noexist",
			Msg:    "电子票不存在",
		}, nil
	}

	if it.Status == ITEM_TICKET_STATUS_USED {
		return &qrcode.QRCodeResult{
			Result: "ticket_used",
			Msg:    "电子票已使用",
		}, nil
	}
	if it.Status == ITEM_TICKET_STATUS_EXPIRED ||
		(it.Status == ITEM_TICKET_STATUS_NOUSE && it.EndTime < time.Now().Unix()) {
		return &qrcode.QRCodeResult{
			Result: "ticket_expired",
			Msg:    "电子票已过期",
		}, nil
	}
	//是否是当天有效票
	stime := time.Unix(it.StartTime, 0)
	y, m, d := stime.Date()
	ny, nm, nd := time.Now().Date()

	if y != ny || m != nm || d != nd {
		return &qrcode.QRCodeResult{
			Result: "ticket_unenforced",
			Msg:    "电子票未到有效期",
		}, nil
	}

	it.FTime = time.Now().Unix()
	it.FUid = fromUid
	it.Status = ITEM_TICKET_STATUS_USED
	err := shopp.UsedItemTicket(it)
	if err == nil {
		return &qrcode.QRCodeResult{
			Result: "success",
			Msg:    "电子票扫描成功",
		}, nil
	}
	return nil, err
}
