package qrcode

import (
	"fmt"
	"strings"
	"sync"
)

type IQRCodeService interface {
	Flag() string
	EncodeCode(code string) string
	DecodeCode(fromUid int64, code string) (*QRCodeResult, error)
}

type QRCodeResult struct {
	Mod    string
	Action string
	Result string
	Msg    string
	Args   []interface{}
}

var qrcodeServices map[string]IQRCodeService = make(map[string]IQRCodeService)
var lock *sync.RWMutex = new(sync.RWMutex)

func RegisterQRProcessor(flag string, service IQRCodeService) {
	if service == nil {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	if _, ok := qrcodeServices[flag]; !ok {
		qrcodeServices[flag] = service
	} else {
		panic(flag + "二维码服务已注册")
	}
}

func processor(qrCode string) (IQRCodeService, error) {
	lock.RLock()
	defer lock.RUnlock()
	args := strings.Split(qrCode, ":")
	if len(args) < 2 {
		return nil, fmt.Errorf("没有对应的处理器")
	}
	flag := strings.ToLower(args[0])
	if service, ok := qrcodeServices[flag]; ok {
		return service, nil
	}
	return nil, fmt.Errorf("没有对应的处理器")
}

func DecodeCode(fromUid int64, qrCode string) (*QRCodeResult, error) {
	service, err := processor(qrCode)
	if err != nil {
		return nil, err
	}
	return service.DecodeCode(fromUid, qrCode)
}

func GetClientCode(flag string, srcCode string) (string, error) {
	var service IQRCodeService
	if ser, ok := qrcodeServices[flag]; !ok {
		return "", fmt.Errorf("没有对应的处理器")
	} else {
		service = ser
	}
	newcode := fmt.Sprintf("%s:%s", service.Flag(), service.EncodeCode(srcCode))
	return newcode, nil
}
