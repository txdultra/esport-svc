package outobjs

type OutQRCodeResult struct {
	Mod    string          `json:"mod"`
	Action string          `json:"action"`
	Result string          `json:"result"`
	Msg    string          `json:"msg"`
	Args   [][]interface{} `json:"args"`
}
