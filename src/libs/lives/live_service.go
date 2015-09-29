package lives

import (
	"libs/lives/service"
	"strconv"
)

var live_spms = &LiveSubPrograms{}

type LiveServiceImpl struct{}

func (l LiveServiceImpl) SetRefId(pid int64, refId string, refType string) (r *service.ActionResult_, err error) {
	if refType == "BET" {
		betId, err := strconv.ParseInt(refId, 10, 64)
		if err != nil {
			return &service.ActionResult_{
				Success:   false,
				Exception: "id非法",
			}, nil
		}
		err = live_spms.UpdateBetId(pid, betId)
		if err != nil {
			return &service.ActionResult_{
				Success:   false,
				Exception: err.Error(),
			}, nil
		}
		return &service.ActionResult_{
			Success:   true,
			Exception: "",
		}, nil
	}
	return &service.ActionResult_{
		Success:   false,
		Exception: refType + "类型不支持",
	}, nil
}
