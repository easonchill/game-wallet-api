package structs

type BetsReq struct {
	Account    string     `json:"account"`
	Gamehall   string     `json:"gamehall"`
	Gamecode   string     `json:"gamecode"`
	Session    string     `json:"session"`
	Genre      string     `json:"genre"`
	Data       []betsData `json:"data"`
	CreateTime string     `json:"createTime"`
}

type WinsReq struct {
	List []WinsReqList `json:"list" binding:"required"`
}

type WinsEvent struct {
	Mtcode   string  `json:"mtcode"`
	Amount   float64 `json:"amount"`
	Validbet float64 `json:"validbet"`
	Roundid  string  `json:"roundid"`
	Eventime string  `json:"eventtime"`
	Gamecode string  `json:"gamecode"`
	Gamehall string  `json:"gamehall"`
}

type WinsReqList struct {
	Account   string      `json:"account" binding:"required"`
	Event     []WinsEvent `json:"event"`
	Eventtime string      `json:"eventtime"`
	Ucode     string      `json:"ucode" binding:"required"`
}

type RefundsReq struct {
	Mtcode []string `json:"mtcode" binding:"required,min=1"`
}

type CancelReq struct {
	Mtcode []string `json:"mtcode" binding:"required,min=1"`
}

type AmendReq struct {
	Account    string      `json:"account" binding:"required"`
	Gamehall   string      `json:"gamehall" binding:"required"`
	Gamecode   string      `json:"gamecode" binding:"required"`
	Action     string      `json:"action"`
	Amount     float64     `json:"amount"`
	Data       []AmendData `json:"data" binding:"required"`
	CreateTime string      `json:"createTime" binding:"required"`
}

type AmendsReq struct {
	List []AmendsList `json:"list"`
}
type AmendsEvent struct {
	Mtcode    string  `json:"mtcode"`
	Amount    float64 `json:"amount"`
	Validbet  int     `json:"validbet"`
	Action    string  `json:"action"`
	Roundid   string  `json:"roundid"`
	Eventtime string  `json:"eventtime"`
	Gamecode  string  `json:"gamecode"`
}
type AmendsList struct {
	Account   string        `json:"account"`
	Event     []AmendsEvent `json:"event"`
	Eventtime string        `json:"eventtime"`
	Amount    float64       `json:"amount"`
	Action    string        `json:"action"`
	Ucode     string        `json:"ucode"`
}

type PayoffReq struct {
	Account   string  `json:"account" binding:"required,min=1" form:"account"`
	Eventtime string  `json:"eventTime" binding:"required,min=1" form:"eventtime"`
	Amount    float64 `json:"amount" binding:"required" form:"amount"`
	Mtcode    string  `json:"mtcode" binding:"required,min=1" form:"mtcode"`
	Remark    string  `json:"remark" form:"remark"`
}

type CreateUserReq struct {
	Account    string  `json:"account" binding:"required,min=1" form:"account"`
	Password   string  `json:"password" binding:"required,min=1" form:"password"`
	Currency   string  `json:"currency" binding:"required,min=1" form:"currency"`
	Balance    float64 `json:"balance" binding:"required" form:"balance"`
	Status     bool    `json:"status"`
	LastTime   string  `json:"last_time"`
	CreateTime string  `json:"create_time"`
}

type Tabler interface {
	TableName() string
}

// GORM TableName 会将 CreateUserReq 的表名重写为 `users`
func (CreateUserReq) TableName() string {
	return "users"
}
