package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

type FindByMtcode struct {
	Mtcode string `bson:"event.mtcode"`
}

type FindByMtcodeOnlyOne struct {
	Mtcode string `bson:"event.mtcode"`
}

type UpdateByMtcode struct {
	Mtcode string `bson:"event.mtcode"`
}

type betsData struct {
	Mtcode   string  `json:"mtcode"`
	Amount   float64 `json:"amount"`
	Roundid  string  `json:"roundid"`
	Eventime string  `json:"eventtime"`
}

type AmendData struct {
	Mtcode    string  `json:"mtcode"`
	Amount    float64 `json:"amount"`
	Validbet  float64 `json:"validbet"`
	Roundid   string  `json:"roundid"`
	Eventtime string  `json:"eventtime"`
	Action    string  `json:"action"`
}

type AmendEvent struct {
	Mtcode    string  `json:"mtcode"`
	Amount    float64 `json:"amount"`
	Action    string  `json:"action"`
	Eventtime string  `json:"eventtime"`
}

type AmendRecord struct {
	Action   string       `json:"action"`
	Target   AmendTarget  `json:"target"`
	Status   AmendStatus  `json:"status"`
	Before   float64      `json:"before"`
	Balance  float64      `json:"balance"`
	Currency string       `json:"currency"`
	Event    []AmendEvent `json:"event"`
}
type AmendTarget struct {
	Account string `json:"account"`
}
type AmendStatus struct {
	Createtime string `json:"createtime"`
	Endtime    string `json:"endtime"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type Transaction_mgolog struct {
	Id     primitive.ObjectID `json:"id"  bson:"_id"`
	Action string             `json:"action" bson:"action"`
	Target struct {
		Account string `json:"account" bson:"account"`
	} `json:"target"`
	Status struct {
		Create_time string `json:"createtime" bson:"createtime"`
		End_time    string `json:"endtime"  bson:"endtime"`
		Status      string `json:"status" bson:"status"`
		Msg         string `json:"message" bson:"message"`
	} `json:"status"`
	Before   float64 `json:"before"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Event    []event `json:"event"`
}

type event struct {
	Mtcode    string  `json:"mtcode" bson:"mtcode"`
	Amount    float64 `json:"amount" bson:"amount"`
	Even_time string  `json:"eventime" bson:"eventime"`
	Status    string  `json:"status" bson:"status"`
}
