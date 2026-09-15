package data

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

const dbTimeout  = time.Second * 3

var client  *mongo.Client

func New(m *mongo.Client) Models {
	client = m

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}