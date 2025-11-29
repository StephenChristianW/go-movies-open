package IPManage

import (
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/collections"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/db"
	"github.com/StephenChristianW/go-movies-open/services"
	UserService "github.com/StephenChristianW/go-movies-open/services/User"
	"github.com/StephenChristianW/go-movies-open/utils/UtilsTime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockedIPs struct {
	UserID    string `json:"user_id" bson:"user_id"`
	UserName  string `json:"username" bson:"username"`
	IP        string `json:"ip" bson:"ip"`
	Reason    string `json:"reason" bson:"reason"`
	Active    int    `json:"active" bson:"active"`
	BlockedAt string `json:"blocked_at" bson:"blocked_at"`
}
type BlockedIPsInDB struct {
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	UserName  string             `json:"username" bson:"username"`
	IP        string             `json:"ip" bson:"ip"`
	Reason    string             `json:"reason" bson:"reason"`
	Active    int                `json:"active" bson:"active"`
	BlockedAt string             `json:"blocked_at" bson:"blocked_at"`
}

func AddSuspiciousIp(username, ip, reason string) {
	ctx, cancel := db.GetCtx()
	defer cancel()
	var data BlockedIPs
	var userInDB UserService.UserInDB
	_ = collections.GetUserCollection().FindOne(ctx, bson.M{"username": username, "deleted": 1}).Decode(&userInDB)
	data.UserID = userInDB.ID.Hex()
	data.IP = ip
	data.Reason = reason
	data.Active = 1
	data.BlockedAt = UtilsTime.NowTime()
	_, _ = collections.GetBlockIps().InsertOne(ctx, data)
}

type CreateIPFilter struct {
	UserName string `json:"username" bson:"username"`
	IP       string `json:"ip" bson:"ip"`
	Reason   string `json:"reason" bson:"reason"`
	Active   int    `json:"active" bson:"active"`

	BeginTime string `json:"beginTime" bson:"beginTime"`
	EndTime   string `json:"endTime" bson:"endTime"`
	Page      int    `json:"page" bson:"page"`
	PageSize  int    `json:"pageSize" bson:"pageSize"`
}

func buildFilter(find CreateIPFilter) bson.M {
	filter := bson.M{}

	if find.UserName != "" {
		filter["username"] = find.UserName
	}
	if find.IP != "" {
		filter["ip"] = find.IP
	}
	if find.Reason != "" {
		filter["reason"] = find.Reason
	}
	if find.Active != 0 {
		filter["active"] = find.Active
	}

	if find.BeginTime != "" {
		filter["blocked_at"] = bson.M{
			"$gte": find.BeginTime,
			"$lte": find.EndTime,
		}
	}
	return filter

}

func BlockIpList(ip CreateIPFilter) (services.Pagination, error) {
	ctx, cancel := db.GetCtx()
	defer cancel()
	var ips []BlockedIPs
	conditions := buildFilter(ip)
	opts := db.CalculatePagination(ip.Page, ip.PageSize, "blocked_at", -1)
	cursor, err := collections.GetBlockIps().Find(ctx, conditions, opts)
	if err != nil {
		return services.Pagination{}, err
	}
	for cursor.Next(ctx) {
		var dbData BlockedIPsInDB
		if err := cursor.Decode(&dbData); err != nil {
			return services.Pagination{}, err
		}
		ips = append(ips, BlockedIPs{
			UserID:    dbData.UserID.Hex(),
			UserName:  dbData.UserName,
			IP:        dbData.IP,
			Reason:    dbData.Reason,
			Active:    dbData.Active,
			BlockedAt: dbData.BlockedAt,
		})
	}
	total, err := collections.GetBlockIps().CountDocuments(ctx, conditions)
	if err != nil {
		return services.Pagination{}, err
	}

	totalPages := (int(total) + ip.PageSize - 1) / ip.PageSize
	return services.Pagination{
		Page:       ip.Page,
		PageSize:   ip.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
		Data:       nil,
	}, nil
}
