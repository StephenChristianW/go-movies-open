package ServiceUtils

import "go.mongodb.org/mongo-driver/bson/primitive"

func StringsToObjectIds(ids []string) []primitive.ObjectID {
	var ObjIds []primitive.ObjectID
	for _, id := range ids {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		ObjIds = append(ObjIds, objId)
	}
	return ObjIds
}
