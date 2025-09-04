package utils

import (
	"go.mongodb.org/mongo-driver/bson"
)

func MergeRangeCondition(filter bson.M, field string, operator string, value int) {
	filter[field] = bson.M{operator: value}
}