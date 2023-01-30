package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

type TransactionMgoLog struct {
	//Id       primitive.ObjectID `json:"id"  bson:"_id"`
	Action   string  `json:"action" bson:"action"`
	Target   Target  `json:"target" bson:"target"`
	Status   Status  `json:"status" bson:"status"`
	Before   float64 `json:"before" bson:"before"`
	Balance  float64 `json:"balance" bson:"balance"`
	Currency string  `json:"currency" bson:"currency"`
	Event    []Event `json:"event"  bson:"event"`
}

type Target struct {
	Account string `json:"account" bson:"account"`
}

type Status struct {
	Create_time string `json:"createtime" bson:"createtime"`
	End_time    string `json:"endtime"  bson:"endtime"`
	Status      string `json:"status" bson:"status"`
	Msg         string `json:"message" bson:"message"`
}

type MtcodeBalanceMgoLog struct {
	Id      primitive.ObjectID `json:"id"  bson:"_id"`
	Action  string             `json:"action" bson:"action"`
	Mtcode  string             `json:"mtcode" bson:"mtcode"`
	Balance float64            `json:"balance"`
}
