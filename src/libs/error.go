package libs

type Error struct {
	ErrorName        string `json:"error"`
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_description"`
	ErrorUrl         string `json:"error_url"`
}

func NewError(name string, code string, desc string, url string) *Error {
	err := &Error{
		name,
		code,
		desc,
		url,
	}
	return err
}

func NewLogError() {

}
