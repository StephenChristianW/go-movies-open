package collections

import (
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserCollection() *mongo.Collection {
	return db.GetStackBuilderCollection(config.UserColl)
}
func GetAdminCollection() *mongo.Collection {
	return db.GetStackBuilderCollection(config.AdminColl)
}
func GetAdminRememberCollection() *mongo.Collection {
	return db.GetStackBuilderCollection(config.AdminRemember)
}
func GetInviteCodeCollection() *mongo.Collection {
	return db.GetStackBuilderCollection(config.InviteCodeColl)
}
func GetBlockIps() *mongo.Collection {
	return db.GetStackBuilderCollection(config.BlockIps)
}
