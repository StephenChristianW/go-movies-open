package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CalculatePagination(Page, PageSize int, sortField string, sortOrder int) *options.FindOptions {

	skip := (Page - 1) * PageSize
	if skip < 0 {
		skip = 0
	}
	if PageSize == 0 {
		PageSize = 10 // 默认 10 条
	}
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(PageSize))

	if sortField != "" {
		opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})
	}
	return opts
}
