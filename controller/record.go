package controller

import (
	"Basil/module"
	"Basil/structs"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Record(c *gin.Context) {
	var tx structs.TransactionMgoLog

	//var result bson.M
	mtcode := c.Param("mtcode")

	mongoClient, err := module.GetMgoCli()
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		wrapResponse(c, 200, "MongoDB connect error", "1100")
		return
	}

	collection := mongoClient.Database("transaction").Collection("transaction_log")

	if err != nil {
		wrapResponse(c, 200, "MongoDB connect error", "1100")
		return
	}
	opts := options.FindOne().SetSkip(0)

	findMtcode := structs.FindByMtcode{Mtcode: mtcode}

	err = collection.FindOne(
		context.TODO(),
		findMtcode,
		opts,
	).Decode(&tx)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			wrapResponse(c, 200, nil, "1014")
			return
		}
	}

	for _, v := range tx.Event {
		if v.Mtcode == mtcode {
			newSlice := []structs.Event{v}
			tx.Event = newSlice
		}
	}

	wrapResponse(c, 200, tx, "0")
}
