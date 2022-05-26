package structs

type WinsSuccessList struct {
	Account  string  `json:"account"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Ucode    string  `json:"ucode"`
}

type WinsFailList struct {
	Account string `json:"account"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Ucode   string `json:"ucode"`
}

type WinsResp struct {
	Success []WinsSuccessList `json:"success"`
	Failed  []WinsFailList    `json:"failed"`
}

type BetsResp struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type RefundsResp struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type CancelResp struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type AmendResp struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

//下面是amends回傳要用的結構

type AmendsSuccessList struct {
	Account  string  `json:"account"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Ucode    string  `json:"ucode"`
}

type AmendsFailList struct {
	Account string `json:"account"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Ucode   string `json:"ucode"`
}

type AmendsResp struct {
	Success []AmendsSuccessList `json:"success"`
	Failed  []AmendsFailList    `json:"failed"`
}
