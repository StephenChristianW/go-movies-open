package InviteCodeService

import (
	"context"
	"errors"
	"fmt"
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/collections"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/db"
	"github.com/StephenChristianW/go-movies-open/services"
	"github.com/StephenChristianW/go-movies-open/services/ServiceUtils"
	"github.com/StephenChristianW/go-movies-open/utils/UtilsRandom"
	"github.com/StephenChristianW/go-movies-open/utils/UtilsTime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
	"time"
)

var batchSize = config.MaxCreateBatchSize

// InviteCodeInterface 邀请码服务接口
type InviteCodeInterface interface {
	GenerateAndInsertCode() (string, error)
	ListInviteCodes(filter CreateInviteFilter) (services.Pagination, error)
	GenerateAndInsertCodes(num int) (int, error)
	ReBackInviteCodes(ids []string) ([]InviteCode, error)
	DeleteInviteCode(ids []string) (int, error)
	UpdateInviteCodes(updateCode UpdateInviteCode) (int, error)
	GetInviteCode(code string) (InviteCode, error)
}

// 邀请码服务功能-具体实现

// GenerateAndInsertCode 生成邀请码并插入数据库
func (ci CreateInviteCode) GenerateAndInsertCode() (string, error) {
	collection := collections.GetInviteCodeCollection()
	ctx, cancel := db.GetCtx()
	defer cancel()

	newCode, err := ci.newInviteCode()
	if err != nil {
		return "", err
	}

	doc, err := collection.InsertOne(ctx, newCode)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// 若生成重复码，则递归生成新的
			return ci.GenerateAndInsertCode()
		}
		return "", err
	}

	if oid, ok := doc.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("unexpected InsertedID type: %T", doc.InsertedID)
}

// ListInviteCodes 分页获取邀请码列表
func (ci CreateInviteCode) ListInviteCodes(filter CreateInviteFilter) (services.Pagination, error) {
	collection := collections.GetInviteCodeCollection()
	ctx, cancel := db.GetCtx()
	defer cancel()

	conditions := ci.buildFilter(filter)
	opts := db.CalculatePagination(filter.Page, filter.PageSize, "expiredTime", -1)

	var inviteCodes []InviteCode
	cursor, err := collection.Find(ctx, conditions, opts)
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)
	if err != nil {
		return services.Pagination{}, err
	}

	for cursor.Next(ctx) {
		var dbcode InviteCodeInDB
		if err := cursor.Decode(&dbcode); err != nil {
			return services.Pagination{}, err
		}
		inviteCodes = append(inviteCodes, ci.toDTO(dbcode))
	}

	total, err := collection.CountDocuments(ctx, conditions)
	if err != nil {
		return services.Pagination{}, err
	}

	totalPages := (int(total) + filter.PageSize - 1) / filter.PageSize
	return services.Pagination{
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
		Data:       inviteCodes,
	}, nil
}

// GenerateAndInsertCodes 批量生成邀请码
func (ci CreateInviteCode) GenerateAndInsertCodes(num int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if num > batchSize {
		num = batchSize
		log.Printf("数据创建量大于 %d 条, 最多一次请求只插入 %d 条数据", batchSize, batchSize)
	}
	totalInserted, err := ci.insertBatch(ctx, num)
	if err != nil {
		return totalInserted, err
	}
	return totalInserted, nil
}

// ReBackInviteCodes 根据ID获取邀请码
func (ci CreateInviteCode) ReBackInviteCodes(ids []string) ([]InviteCode, error) {
	var codes []InviteCode
	collection := collections.GetInviteCodeCollection()
	ctx, cancel := db.GetCtx()
	defer cancel()

	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("无效的 ObjectID:", id, err)
			continue
		}
		objectIDs = append(objectIDs, objID)
	}

	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return codes, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var dbcode InviteCodeInDB
		if err = cursor.Decode(&dbcode); err != nil {
			log.Println("解码失败:", err)
			continue
		}
		codes = append(codes, ci.toDTO(dbcode))
	}

	if err = cursor.Err(); err != nil {
		return codes, err
	}

	return codes, nil
}

// DeleteInviteCode 删除邀请码
func (ci CreateInviteCode) DeleteInviteCode(ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	objectIds := ServiceUtils.StringsToObjectIds(ids)
	ctx, cancel := db.GetCtx()
	defer cancel()
	updateResult, err := collections.GetInviteCodeCollection().UpdateMany(ctx, bson.M{"_id": bson.M{"$in": objectIds}},
		bson.M{"$set": bson.M{"deleted": 2}})
	if err != nil {
		return 0, err
	}
	if updateResult.MatchedCount == 0 {
		return 0, errors.New("无数据匹配")
	}
	if updateResult.ModifiedCount == 0 {
		return 0, errors.New("无数据被删除")
	}
	return int(updateResult.ModifiedCount), nil
}

func (ci CreateInviteCode) UpdateInviteCodes(updateCode UpdateInviteCode) (int, error) {
	if len(updateCode.Ids) == 0 {
		return 0, nil
	}

	objectIds := ServiceUtils.StringsToObjectIds(updateCode.Ids)
	ctx, cancel := db.GetCtx()
	defer cancel()

	// 查询要更新的邀请码
	cursor, err := collections.GetInviteCodeCollection().Find(ctx, bson.M{"_id": bson.M{"$in": objectIds}})
	if err != nil {
		return 0, err
	}

	var noUserIDs []primitive.ObjectID
	for cursor.Next(ctx) {
		var dbcode InviteCodeInDB
		if err = cursor.Decode(&dbcode); err != nil {
			log.Println("Decode error:", err)
			continue
		}
		// 只允许未绑定用户且状态为正常的记录修改
		if dbcode.Username == "" && dbcode.Status == 1 {
			noUserIDs = append(noUserIDs, dbcode.ID)
		}
	}

	if len(noUserIDs) == 0 {
		return 0, errors.New("所有邀请码已被用户绑定，无法修改")
	}

	// 构建更新字段
	updateFields := bson.M{}
	if updateCode.Status != 0 {
		updateFields["status"] = 1
	}
	if updateCode.Days != 0 {
		updateFields["expiredTime"] = UtilsTime.AfterNDays000(updateCode.Days)
	}
	if len(updateFields) == 0 {
		return 0, errors.New("没有要更新的字段")
	}
	updateFields["deleted"] = 1

	update := bson.M{"$set": updateFields}
	updateResult, err := collections.GetInviteCodeCollection().UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": noUserIDs}},
		update,
	)
	if err != nil {
		return 0, err
	}

	if updateResult.ModifiedCount == 0 {
		return 0, errors.New("无数据更改")
	}

	return int(updateResult.ModifiedCount), nil
}

// GetInviteCode 根据codeID 或 code 获取验证码
func (ci CreateInviteCode) GetInviteCode(code string) (InviteCode, error) {
	ctx, cancel := db.GetCtx()
	defer cancel()

	var inviteCodeDB InviteCodeInDB
	var inviteCode InviteCode
	condition := bson.M{}

	// 尝试把 code 转为 ObjectID
	objectId, err := primitive.ObjectIDFromHex(code)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			// code 不是 ObjectID，则用 code 字段查询
			condition["code"] = code
		} else {
			// 其他错误直接返回
			return inviteCode, err
		}
	} else {
		// code 是 ObjectID，查 _id
		condition["_id"] = objectId
	}

	// 查询数据库
	err = collections.GetInviteCodeCollection().FindOne(ctx, condition).Decode(&inviteCodeDB)
	if err != nil {
		return inviteCode, err
	}
	// 转 DTO 返回
	inviteCode = ci.toDTO(inviteCodeDB)
	return inviteCode, nil
}

// 邀请码服务功能-工具函数

// 生成随机邀请码（8位字母数字组合）
func generateRandomCode() (string, error) {
	return UtilsRandom.RandomAlphaNumString(8)
}

// 将数据库结构转换为前端结构
func (ci CreateInviteCode) toDTO(dbcode InviteCodeInDB) InviteCode {
	return InviteCode{
		ID:          dbcode.ID.Hex(),
		Code:        dbcode.Code,
		Username:    dbcode.Username,
		Status:      dbcode.Status,
		IsExpired:   dbcode.IsExpired,
		ExpiredTime: dbcode.ExpiredTime,
	}
}

func (ci CreateInviteCode) buildFilter(find CreateInviteFilter) bson.M {
	// 如果传入了 ID 列表，直接返回只筛选这些 ID
	if len(find.IDs) > 0 {
		objIDs := ServiceUtils.StringsToObjectIds(find.IDs)
		if len(objIDs) > 0 {
			return bson.M{"_id": bson.M{"$in": objIDs}}
		}
	}
	filter := bson.M{}

	if find.Code != "" {
		filter["code"] = find.Code
	}
	if find.Username != "" {
		filter["username"] = find.Username
	}
	if find.Status != 0 {
		filter["status"] = find.Status
	}
	if find.BeginTime != "" {
		filter["expiredTime"] = bson.M{
			"$gte": find.BeginTime,
			"$lte": find.EndTime,
		}
	}
	return filter
}

// insertBatch 插入一批数据，并处理重复 key 错误，确保最终成功插入 size 条
func (ci CreateInviteCode) insertBatch(ctx context.Context, size int) (int, error) {
	if size <= 0 {
		return 0, nil
	}

	var docs []interface{}
	for i := 0; i < size; i++ {
		newCode, err := ci.newInviteCode()
		if err != nil {
			return 0, err
		}
		docs = append(docs, newCode)
	}

	res, err := collections.GetInviteCodeCollection().InsertMany(ctx, docs)
	if err != nil {
		// 处理批量写入错误
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			dupCount := 0
			for _, we := range bulkErr.WriteErrors {
				if strings.Contains(we.Message, "E11000") {
					dupCount++
				}
			}

			successCount := size - dupCount

			// 如果还有重复的，再递归补齐
			if dupCount > 0 {
				retryCount, retryErr := ci.insertBatch(ctx, dupCount)
				if retryErr != nil {
					return successCount + retryCount, retryErr
				}
				return successCount + retryCount, nil
			}

			return successCount, nil
		}

		// 其它错误（例如网络超时、权限问题），直接返回
		return 0, fmt.Errorf("批量插入失败: %w", err)
	}

	// 正常插入，返回成功数量
	return len(res.InsertedIDs), nil
}

// 构建新的 InviteCode 对象
func (ci CreateInviteCode) newInviteCode() (CreateInviteCode, error) {
	code, err := generateRandomCode()
	if err != nil {
		return ci, err
	}
	return CreateInviteCode{
		Code:        code,
		Username:    "",                         // 初始未绑定用户名
		Status:      1,                          // 状态 1 表示可用
		IsExpired:   1,                          // 1 表示未过期
		ExpiredTime: UtilsTime.AfterNDays000(3), // 默认三天后过期
		Deleted:     1,
	}, nil
}

func CheckExpiredCodes() error {
	today0 := UtilsTime.Today000()
	collection := collections.GetInviteCodeCollection()
	ctx, cancel := db.GetCtx()
	defer cancel()

	update := bson.M{
		"$set": bson.M{"deleted": 2}, // 标记为已删除
	}

	filter := bson.M{
		"username":     "",                     // 未被领取
		"expired_time": bson.M{"$lte": today0}, // 过期
		"deleted":      1,                      // 当前状态正常
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {

		return err
	}
	log.Printf("CheckExpiredCodes: %d codes expired and marked deleted\n", result.ModifiedCount)
	return nil

}
