package AdminService

import (
	"context"
	"errors"
	"fmt"
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/collections"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/db"
	"github.com/StephenChristianW/go-movies-open/security/SecurityBcrypt"
	"github.com/StephenChristianW/go-movies-open/services"
	"github.com/StephenChristianW/go-movies-open/services/ServiceUtils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// =======================================
//          管理员管理模块
// =======================================

// CreateAdmin 创建新的管理员账号
func (*AdminService) CreateAdmin(username, password, roleID string) error {
	ctx, cancel := db.GetCtx()
	defer cancel()

	// 检查管理员是否已存在
	var existing AdminInDB
	err := collections.GetAdminCollection().FindOne(ctx, bson.M{
		"username": username,
		"deleted":  1,
	}).Decode(&existing)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 不存在则创建新管理员
			hashedPwd, err := SecurityBcrypt.GenerateHashPwd(password)
			if err != nil {
				return err
			}

			newAdmin := AdminCreate{
				Username: username,
				Password: hashedPwd,
				Status:   1,
				Deleted:  1,
				Role:     roleID,
			}
			_, insertErr := collections.GetAdminCollection().InsertOne(ctx, newAdmin)
			return insertErr
		}
		return err
	}

	return fmt.Errorf("管理员 %s 已存在", username)
}

// DeleteAdmin 逻辑删除管理员（不会删除本人或根管理员）
func (*AdminService) DeleteAdmin(currentUsername, targetUsername string) error {
	if currentUsername == targetUsername {
		return errors.New("无法删除本人")
	}
	if _, ok := config.RootAdmins[targetUsername]; ok {
		return errors.New("该管理员无法删除")
	}

	ctx, cancel := db.GetCtx()
	defer cancel()

	var admin AdminInDB
	err := collections.GetAdminCollection().FindOne(ctx, bson.M{
		"username": targetUsername,
		"deleted":  1,
	}).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("未找到管理员: %s", targetUsername)
		}
		return err
	}

	// 逻辑删除
	_, err = collections.GetAdminCollection().UpdateOne(ctx,
		bson.M{"_id": admin.ID, "deleted": 1},
		bson.M{"$set": bson.M{"deleted": 2}},
	)
	if err != nil {
		return err
	}

	// 删除 token
	return deleteAdminToken(admin.ID.Hex())
}

// ChangeAdminPassword 修改管理员密码，并删除其 token
func (*AdminService) ChangeAdminPassword(username, newPassword string) error {
	ctx, cancel := db.GetCtx()
	defer cancel()

	var admin AdminInDB
	err := collections.GetAdminCollection().FindOne(ctx, bson.M{
		"username": username,
		"deleted":  1,
		"status":   1,
	}).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("用户名错误或管理员 %s 已被禁用", username)
		}
		return err
	}

	hashedPwd, err := SecurityBcrypt.GenerateHashPwd(newPassword)
	if err != nil {
		return err
	}

	_, err = collections.GetAdminCollection().UpdateOne(ctx,
		bson.M{"username": username, "deleted": 1},
		bson.M{"$set": bson.M{"password": hashedPwd}},
	)
	if err != nil {
		return err
	}

	log.Printf("管理员 %s 已更改密码", username)
	return deleteAdminToken(admin.ID.Hex())
}

// BannedAdminUser 禁用管理员账号（非根管理员）
func (*AdminService) BannedAdminUser(username string) error {
	if _, ok := config.RootAdmins[username]; ok {
		return errors.New("初始管理员无法禁用")
	}

	ctx, cancel := db.GetCtx()
	defer cancel()

	var admin AdminInDB
	err := collections.GetAdminCollection().FindOne(ctx, bson.M{
		"username": username,
		"deleted":  1,
		"status":   1,
	}).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("管理员 %s 不存在或已被禁用/删除", username)
		}
		return fmt.Errorf("查询管理员失败: %v", err)
	}

	_, err = collections.GetAdminCollection().UpdateOne(ctx,
		bson.M{"_id": admin.ID},
		bson.M{"$set": bson.M{"status": 2}},
	)
	if err != nil {
		return fmt.Errorf("禁用管理员失败: %v", err)
	}

	return deleteAdminToken(admin.ID.Hex())
}
func (*AdminService) ActiveAdminUser(username string) error {
	ctx, cancel := db.GetCtx()
	defer cancel()
	// 启用条件：未删除、当前状态为禁用
	filter := bson.M{
		"username": username,
		"deleted":  1,
		"status":   2,
	}
	update := bson.M{
		"$set": bson.M{"status": 1},
	}

	doc := collections.GetAdminCollection().FindOneAndUpdate(ctx, filter, update)
	if doc.Err() != nil {
		if errors.Is(doc.Err(), mongo.ErrNoDocuments) {
			return errors.New("该管理员已启用或不存在")
		}
		return doc.Err()
	}
	return nil
}

// AdminList 获取管理员列表（支持分页与过滤）
func (*AdminService) AdminList(filter AdminList) (services.Pagination, error) {
	ctx, cancel := db.GetCtx()
	defer cancel()

	var result services.Pagination
	conditions := buildAdminFilter(filter)

	total, err := collections.GetAdminCollection().CountDocuments(ctx, conditions)
	if err != nil {
		return result, err
	}

	findOptions := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.PageSize)).
		SetLimit(int64(filter.PageSize))

	cursor, err := collections.GetAdminCollection().Find(ctx, conditions, findOptions)
	if err != nil {
		return result, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)

	var admins []Admin
	for cursor.Next(ctx) {
		var a AdminInDB
		if err := cursor.Decode(&a); err != nil {
			return result, err
		}
		admins = append(admins, Admin{
			ID:       a.ID.Hex(),
			Username: a.Username,
			Status:   a.Status,
			Deleted:  a.Deleted,
			Role:     a.Role,
		})
	}

	result.Total = int(total)
	result.Data = admins
	result.Page = filter.Page
	result.PageSize = filter.PageSize
	if filter.PageSize > 0 {
		result.TotalPages = int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize))
	} else {
		result.TotalPages = 1
	}
	return result, nil
}

// buildAdminFilter 构建 MongoDB 查询条件
func buildAdminFilter(filter AdminList) bson.M {
	if len(filter.IDs) > 0 {
		objIDs := ServiceUtils.StringsToObjectIds(filter.IDs)
		if len(objIDs) > 0 {
			return bson.M{"_id": bson.M{"$in": objIDs}}
		}
	}

	m := bson.M{}
	if filter.Username != "" {
		m["username"] = bson.M{"$regex": filter.Username, "$options": "i"}
	}
	if filter.Status != 0 {
		m["status"] = filter.Status
	}
	if filter.Deleted != 0 {
		m["deleted"] = filter.Deleted
	}
	if filter.Role != "" {
		m["role"] = filter.Role
	}
	return m
}
