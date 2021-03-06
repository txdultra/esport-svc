package shop

import (
	"fmt"
	credit_client "libs/credits/client"
	credit_proxy "libs/credits/proxy"
	"libs/dlock"
	"libs/message"
	"libs/vars"
	"logs"
	"time"
)

type ConsumerInfo struct {
	ProvinceName string
	CityName     string
	AreaName     string
	Address      string
	Receiver     string
	Tel          string
}

type TransportInfo struct {
	TransId int
	TransNo string
}

type BuyResult struct {
	OrderNo  string
	Code     string
	ItemType ITEM_TYPE
	Error    error
}

func NewShopPurchaser() *ShopPurchaser {
	return &ShopPurchaser{}
}

type ShopPurchaser struct{}

func (sp *ShopPurchaser) issueType(itemType ITEM_TYPE) ISSUE_TYPE {
	if itemType == ITEM_TYPE_ENTITY {
		return ISSUE_TYPE_EXPRESS
	}
	if itemType == ITEM_TYPE_VIRTUAL {
		return ISSUE_TYPE_DIRECT
	}
	if itemType == ITEM_TYPE_TICKET {
		return ISSUE_TYPE_DIRECT
	}
	panic("商品类型没有对应的发货方式")
}

func (sp *ShopPurchaser) ChangeOrderStatus(orderNo string, status ORDER_STATUS, returnStock bool, tpInfo *TransportInfo) error {
	shopp := NewShop()
	order := shopp.GetOrder(orderNo)
	if order == nil {
		return fmt.Errorf("订单号不存在")
	}
	original_status := order.OrderStatus
	switch status {
	case ORDER_STATUS_USERCANCEL: //用户取消订单
		if order.OrderStatus != ORDER_STATUS_AUDITING {
			return fmt.Errorf("订单已锁定,不能撤销")
		}
		err := shopp.UpdateOrderStatus(orderNo, ORDER_STATUS_USERCANCEL)
		if err != nil {
			return err
		}
		if order.PayStatus == PAY_STATUS_PAIED {
			err = sp.rollbackCredit(order.PayNo, order.PriceType) //回滚订单
			if err != nil {
				shopp.UpdateOrderStatus(orderNo, original_status)
				logs.Errorf("订单:%s货币返还失败", orderNo)
			}
		}
		//退回库存和销量
		if returnStock {
			shopp.UpdateItemStocks(order.ItemId, order.Nums, -order.Nums, false)
		}
		return err
	case ORDER_STATUS_AUDITSUCC: //审核成功
		return shopp.UpdateOrderStatus(orderNo, ORDER_STATUS_AUDITSUCC)
	case ORDER_STATUS_AUDITFAIL: //审核失败
		err := shopp.UpdateOrderStatus(orderNo, ORDER_STATUS_AUDITFAIL)
		if err != nil {
			return err
		}
		if order.PayStatus == PAY_STATUS_PAIED {
			err = sp.rollbackCredit(order.PayNo, order.PriceType) //回滚订单
			if err != nil {
				shopp.UpdateOrderStatus(orderNo, original_status)
				logs.Errorf("订单:%s积分返还失败", orderNo)
			}
		}
		//退回库存和销量
		if returnStock {
			shopp.UpdateItemStocks(order.ItemId, order.Nums, -order.Nums, false)
		}
		return err
	case ORDER_STATUS_SENDED: //发货
		err := shopp.UpdateOrderStatus(orderNo, ORDER_STATUS_SENDED)
		if err != nil {
			return err
		}
		transport := shopp.GetOrderTransport(orderNo)
		if transport != nil && tpInfo != nil {
			transport.CompanyId = tpInfo.TransId
			transport.TransNo = tpInfo.TransNo
			shopp.UpdateOrderTransport(transport)
		}
		if order.PayStatus == PAY_STATUS_PAIED {
			err = sp.enterCredit(order.PayNo, order.PriceType) //确认扣除积分
			if err != nil {
				shopp.UpdateOrderStatus(orderNo, original_status)
				logs.Errorf("订单:%s积分返还失败", orderNo)
			}
		}
		go func() {
			message.SendMsgV2(0, order.Uid, vars.MSG_TYPE_SYS, "恭喜您！您在商城提交的订单已发货。", order.OrderNo, nil)
		}()
		return err
	default:
		return nil
	}
}

func (sp *ShopPurchaser) verifyItemPriceType(item *Item, priceType vars.CURRENCY_TYPE) error {
	if priceType == vars.CURRENCY_TYPE_RMB {
		return fmt.Errorf("暂不支持人民币购买")
	}
	if item.PriceType&int(priceType) != int(priceType) {
		return fmt.Errorf("商品不支持此货币购买")
	}
	return nil
}

func (sp *ShopPurchaser) itemPrice(item *Item, priceType vars.CURRENCY_TYPE) float64 {
	switch priceType {
	case vars.CURRENCY_TYPE_CREDIT, vars.CURRENCY_TYPE_RMB:
		return item.Price
	case vars.CURRENCY_TYPE_JING:
		return float64(item.Jings)
	}
	return -1
}

//现阶段只支持积分购买
func (sp *ShopPurchaser) Buy(itemId int64, priceType vars.CURRENCY_TYPE, uid int64, nums int, remark string, cInfo *ConsumerInfo) *BuyResult {
	shopp := NewShop()
	item := shopp.GetItem(itemId)
	if item == nil {
		return &BuyResult{
			OrderNo: "",
			Code:    "",
			Error:   fmt.Errorf("商品不存在"),
		}
	}
	//检查购买货币
	err := sp.verifyItemPriceType(item, priceType)
	if err != nil {
		return &BuyResult{
			OrderNo: "",
			Code:    "",
			Error:   err,
		}
	}

	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(shopp.ItemLockKey(itemId))
	if err != nil {
		logs.Errorf("shop buy item get lock fail:%s", err.Error())
		return &BuyResult{
			OrderNo:  "",
			Code:     "",
			ItemType: item.ItemType,
			Error:    fmt.Errorf("系统繁忙"),
		}
	}
	defer locker.Unlock()

	if !item.Enabled {
		return &BuyResult{
			OrderNo:  "",
			Code:     "",
			ItemType: item.ItemType,
			Error:    fmt.Errorf("商品已被删除"),
		}
	}

	//检测库存
	if item.Stocks < nums {
		return &BuyResult{
			OrderNo:  "",
			Code:     "",
			ItemType: item.ItemType,
			Error:    fmt.Errorf("商品已售完"),
		}
	}

	//计算总价
	itemPrice := sp.itemPrice(item, priceType)
	totalPrice := itemPrice * float64(nums)
	if itemPrice < 0 {
		return &BuyResult{
			OrderNo: "",
			Code:    "",
			Error:   fmt.Errorf("商品价格错误"),
		}
	}

	m_credits := sp.getCredit(uid, priceType)
	if m_credits < int64(totalPrice) {
		return &BuyResult{
			OrderNo:  "",
			Code:     "",
			ItemType: item.ItemType,
			Error:    fmt.Errorf("您使用的货币不足购买"),
		}
	}

	snapId, err := shopp.ItemSnap(item)
	if err != nil {
		return &BuyResult{
			OrderNo:  "",
			Code:     "",
			ItemType: item.ItemType,
			Error:    fmt.Errorf("创建快照失败"),
		}
	}

	order := &Order{
		ItemId:      item.ItemId,
		ItemType:    item.ItemType,
		IssueType:   sp.issueType(item.ItemType),
		Ts:          time.Now().Unix(),
		Uid:         uid,
		OrderStatus: ORDER_STATUS_AUDITING,
		PayStatus:   PAY_STATUS_UNPAID,
		Nums:        nums,
		Price:       itemPrice,
		TotalPrice:  totalPrice,
		PriceType:   vars.CURRENCY_TYPE(priceType),
		SnapId:      snapId,
		Remark:      remark,
		PayId:       PAYID_CREDIT,
	}

	orderNo, err := shopp.CreateOrder(order)
	if err != nil {
		return &BuyResult{
			OrderNo:  "",
			Code:     "",
			ItemType: item.ItemType,
			Error:    fmt.Errorf("创建订单失败"),
		}
	}

	//扣款
	no, err := sp.lockCredit(int64(totalPrice), priceType, uid, credit_proxy.OPERATION_ACTOIN_LOCKDECR, item.Name)
	if err != nil {
		return &BuyResult{
			OrderNo:  "",
			Code:     "",
			ItemType: item.ItemType,
			Error:    fmt.Errorf("扣除货币失败,无法购买商品"),
		}
	}

	if item.ItemType == ITEM_TYPE_VIRTUAL {
		code, err := shopp.UseItemCode(itemId, orderNo)
		if err != nil {
			sp.rollbackCredit(no, priceType)
			shopp.DeleteOrder(orderNo)
			//需加修正库存机制
			return &BuyResult{
				OrderNo:  "",
				Code:     "",
				ItemType: item.ItemType,
				Error:    fmt.Errorf("商品码已无货"),
			}
		}
		//确定扣除积分
		sp.enterCredit(no, priceType)
		//更新订单状态
		order.PayId = PAYID_CREDIT
		order.PayNo = no
		order.PayStatus = PAY_STATUS_PAIED
		order.OrderStatus = ORDER_STATUS_COMPLETED
		order.Ex1 = code
		shopp.UpdateOrder(order)

		//扣除库存和销量
		shopp.UpdateItemStocks(itemId, -nums, nums, true)

		//添加用户计数
		go shopp.UpdateMemberCount(uid, 1, MEMBER_COUNT_UPDATE_ITEM_PURCHASEDS)

		return &BuyResult{
			OrderNo:  orderNo,
			Code:     code,
			ItemType: item.ItemType,
			Error:    nil,
		}
	}

	if item.ItemType == ITEM_TYPE_ENTITY {
		shopp.CreateOrderTransport(&OrderTransport{
			OrderNo:  orderNo,
			Country:  "zh",
			Province: cInfo.ProvinceName,
			City:     cInfo.CityName,
			Area:     cInfo.AreaName,
			Addr1:    cInfo.Address,
			Receiver: cInfo.Receiver,
			Tel1:     cInfo.Tel,
			Ts:       time.Now().Unix(),
		})
		//更新订单状态
		order.PayId = PAYID_CREDIT
		order.PayNo = no
		order.PayStatus = PAY_STATUS_PAIED
		order.OrderStatus = ORDER_STATUS_AUDITING
		shopp.UpdateOrder(order)

		//扣除库存和销量
		shopp.UpdateItemStocks(itemId, -nums, nums, true)

		return &BuyResult{
			OrderNo:  orderNo,
			Code:     "",
			ItemType: item.ItemType,
			Error:    nil,
		}
	}
	if item.ItemType == ITEM_TYPE_TICKET {
		it, err := shopp.GetNoUseItemTicket(itemId, uid, orderNo)
		if err != nil {
			sp.rollbackCredit(no, priceType)
			shopp.DeleteOrder(orderNo)
			//需加修正库存机制
			return &BuyResult{
				OrderNo:  "",
				Code:     "",
				ItemType: item.ItemType,
				Error:    fmt.Errorf("商品码已无货"),
			}
		}
		//确定扣除积分
		sp.enterCredit(no, priceType)
		//更新订单状态
		order.PayId = PAYID_CREDIT
		order.PayNo = no
		order.PayStatus = PAY_STATUS_PAIED
		order.OrderStatus = ORDER_STATUS_COMPLETED
		ttype := "real"
		if it.TType != ITEM_TICKET_TYPE_REAL {
			ttype = "virual"
		}
		order.Ex1 = it.Code
		order.Ex2 = ttype
		shopp.UpdateOrder(order)

		//扣除库存和销量
		shopp.UpdateItemStocks(itemId, -nums, nums, true)

		//添加用户计数
		go shopp.UpdateMemberCount(uid, 1, MEMBER_COUNT_UPDATE_ITEM_PURCHASEDS)

		return &BuyResult{
			OrderNo:  orderNo,
			Code:     it.Code,
			ItemType: item.ItemType,
			Error:    nil,
		}
	}
	//商品类型没有可用流程
	sp.rollbackCredit(no, priceType)
	shopp.DeleteOrder(orderNo)
	return &BuyResult{
		OrderNo:  "",
		Code:     "",
		ItemType: item.ItemType,
		Error:    fmt.Errorf("商品类型不存在购买方式"),
	}
}

func (sp *ShopPurchaser) getCreditHost(priceType vars.CURRENCY_TYPE) string {
	switch priceType {
	case vars.CURRENCY_TYPE_CREDIT:
		return credit_service_host
	case vars.CURRENCY_TYPE_JING:
		return jing_service_host
	default:
		return ""
	}
}

func (sp *ShopPurchaser) getCredit(uid int64, priceType vars.CURRENCY_TYPE) int64 {
	host := sp.getCreditHost(priceType)
	if len(host) == 0 {
		return 0
	}
	client, transport, err := credit_client.NewClient(host)
	if err != nil {
		return 0
	}
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	_credits, err := client.GetCredit(uid)
	return _credits
}

func (sp *ShopPurchaser) rollbackCredit(creditNo string, priceType vars.CURRENCY_TYPE) error {
	host := sp.getCreditHost(priceType)
	if len(host) == 0 {
		return fmt.Errorf("货币系统不存在")
	}
	client, transport, err := credit_client.NewClient(host)
	if err != nil {
		return err
	}
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	operParameter := &credit_proxy.OperationCreditParameter{
		No:     creditNo,
		Action: credit_proxy.OPERATION_ACTOIN_ROLLBACKLOCK,
	}
	result, err := client.Do(operParameter)
	if err != nil || result.State != credit_proxy.OPERATION_STATE_SUCCESS {
		return fmt.Errorf("回滚积分失败")
	}
	return nil
}

func (sp *ShopPurchaser) lockCredit(credits int64, priceType vars.CURRENCY_TYPE, uid int64, oper credit_proxy.OPERATION_ACTOIN, product string) (string, error) {
	host := sp.getCreditHost(priceType)
	if len(host) == 0 {
		return "", fmt.Errorf("货币系统不存在")
	}
	client, transport, err := credit_client.NewClient(host)
	if err != nil {
		return "连接积分系统失败", err
	}
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	operParameter := &credit_proxy.OperationCreditParameter{
		Uid:    uid,
		Points: credits,
		Desc:   "购买商品" + product,
		Action: oper,
	}
	result, err := client.Do(operParameter)
	if err != nil {
		return "扣除货币失败", err
	}
	return result.No, nil
}

func (sp *ShopPurchaser) enterCredit(creditNo string, priceType vars.CURRENCY_TYPE) error {
	host := sp.getCreditHost(priceType)
	if len(host) == 0 {
		return fmt.Errorf("货币系统不存在")
	}
	client, transport, err := credit_client.NewClient(host)
	if err != nil {
		return err
	}
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	operParameter := &credit_proxy.OperationCreditParameter{
		No:     creditNo,
		Action: credit_proxy.OPERATION_ACTOIN_UNLOCK,
	}
	result, err := client.Do(operParameter)
	if err != nil || result.State != credit_proxy.OPERATION_STATE_SUCCESS {
		return fmt.Errorf("确定扣除积分失败")
	}
	return nil
}
