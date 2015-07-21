package shop

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

func (s ShopTicketQRCodeService) Process(qrCode string) {

}
