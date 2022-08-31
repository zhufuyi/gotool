package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zhufuyi/goctl/templates/handler/internal/cache"
	"github.com/zhufuyi/goctl/templates/handler/internal/model"

	"github.com/spf13/cast"
	cacheBase "github.com/zhufuyi/pkg/cache"
	"github.com/zhufuyi/pkg/goredis"
	"github.com/zhufuyi/pkg/mysql/query"
	"gorm.io/gorm"
)

var _ UserExampleDao = (*userExampleDao)(nil)

// UserExampleDao 定义dao接口
type UserExampleDao interface {
	Create(ctx context.Context, table *model.UserExample) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.UserExample) error
	GetByID(ctx context.Context, id uint64) (*model.UserExample, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.UserExample, int64, error)
}

type userExampleDao struct {
	db    *gorm.DB
	cache cache.UserExampleCache
}

// NewUserExampleDao 创建dao接口
func NewUserExampleDao(db *gorm.DB, cache cache.UserExampleCache) UserExampleDao {
	return &userExampleDao{db: db, cache: cache}
}

// Create 创建一条记录，插入记录后，id值被回写到table中
func (d *userExampleDao) Create(ctx context.Context, table *model.UserExample) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID 根据id删除一条记录
func (d *userExampleDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UserExample{}).Error
	if err != nil {
		return nil
	}

	// delete cache
	_ = d.cache.Del(ctx, id)

	return nil
}

// Deletes 根据id删除多条记录
func (d *userExampleDao) Deletes(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.UserExample{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.cache.Del(ctx, id)
	}

	return nil
}

// UpdateByID 根据id更新记录
func (d *userExampleDao) UpdateByID(ctx context.Context, table *model.UserExample) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}
	// todo generate the update fields code to here
	// delete the templates code start
	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Password != "" {
		update["password"] = table.Password
	}
	if table.Email != "" {
		update["email"] = table.Email
	}
	if table.Phone != "" {
		update["phone"] = table.Phone
	}
	if table.Avatar != "" {
		update["avatar"] = table.Avatar
	}
	if table.Age > 0 {
		update["age"] = table.Age
	}
	if table.Gender > 0 {
		update["gender"] = table.Gender
	}
	if table.LoginAt > 0 {
		update["login_at"] = table.LoginAt
	}
	// delete the templates code end
	err := d.db.WithContext(ctx).Model(table).Where("id = ?", table.ID).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.cache.Del(ctx, table.ID)

	return nil
}

// GetByID 根据id获取一条记录
func (d *userExampleDao) GetByID(ctx context.Context, id uint64) (*model.UserExample, error) {
	record, err := d.cache.Get(ctx, id)

	if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, model.ErrRecordNotFound
	}

	// 从mysql获取
	if errors.Is(err, goredis.ErrRedisNotFound) {
		table := &model.UserExample{}
		err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
		if err != nil {
			// if data is empty, set not found cache to prevent cache penetration(防止缓存穿透)
			if err.Error() == model.ErrRecordNotFound.Error() {
				err = d.cache.SetCacheWithNotFound(ctx, id)
				if err != nil {
					return nil, err
				}
				return nil, model.ErrRecordNotFound
			}
			return nil, err
		}

		if table.ID == 0 {
			return nil, model.ErrRecordNotFound
		}

		// set cache
		err = d.cache.Set(ctx, id, table, cacheBase.DefaultExpireTime)
		if err != nil {
			return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
		}

		return table, nil
	}

	if err != nil {
		// fail fast, if cache error return, don't request to db
		return nil, err
	}

	return record, nil
}

// GetByColumns 根据分页和列信息筛选多条记录
// params 包括分页参数和查询参数
// 分页参数(必须):

//	page: 页码，从0开始
//	size: 每页行数
//	sort: 排序字段，默认是id倒叙，可以在字段前添加-号表示倒序，没有-号表示升序，多个字段用逗号分隔
//
// 查询参数(非必须):
//
//	name: 列名
//	exp: 表达式，有=、!=、>、>=、<、<=、like七种类型，值为空时默认是=
//	value: 列值
//	logic: 表示逻辑类型，有&(and)、||(or)两种类型，值为空时默认是and
//
// 示例: 查询年龄大于20的男性
//
//	params = &query.Params{
//	    Page: 0,
//	    Size: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "男",
//		},
//	}
func (d *userExampleDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.UserExample, int64, error) {
	query, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = d.db.WithContext(ctx).Model(&model.UserExample{}).Where(query, args...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, total, nil
	}

	records := []*model.UserExample{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(query, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// GetByIDs 根据id批量获取
func (d *userExampleDao) GetByIDs(ctx context.Context, ids []uint64) ([]*model.UserExample, error) {
	records := []*model.UserExample{}

	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	var missedID []uint64
	for _, id := range ids {
		item, ok := itemMap[cast.ToString(id)]
		if !ok {
			missedID = append(missedID, id)
			continue
		}
		records = append(records, item)
	}

	// get missed data
	if len(missedID) > 0 {
		var missedData []*model.UserExample
		err = d.db.WithContext(ctx).Where("id IN (?)", missedID).Find(&missedData).Error
		if err != nil {
			return nil, err
		}

		if len(missedData) > 0 {
			records = append(records, missedData...)
			err = d.cache.MultiSet(ctx, missedData, 10*time.Minute)
			if err != nil {
				return nil, err
			}
		}
	}

	return records, nil
}
