package models

import (
	"errors"
	"gateway-service/internal/dto"
	"sync"

	"gorm.io/gorm"
)

//- 提供租户级别的访问控制（IP白名单）
//- 租户级别的限流（QPS/QPD）
//- 作为网关认证的租户标识

type App struct {
	gorm.Model
	AppID    string `json:"app_id" gorm:"column:app_id" description:"租户id	"`
	Name     string `json:"name" gorm:"column:name" description:"租户名称"`
	Secret   string `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS string `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qpd      int64  `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps      int64  `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	IsDelete int8   `json:"is_delete" gorm:"column:is_delete;default:0" description:"是否已删除；0：否；1：是"`
}

func (a *App) TableName() string {
	return "gateway_app"
}

// CRUD functions

func GetByID(id uint) (*App, error) {
	model := &App{}
	err := DB.Where("id = ?", id).First(model).Error
	return model, err
}

func GetByAppID(appID string) (*App, error) {
	model := &App{}
	err := DB.Where("app_id = ? AND is_delete = ?", appID, 0).First(model).Error
	return model, err
}

func Exists(appID string) (bool, error) {
	var count int64
	err := DB.Model(&App{}).Where("app_id = ? AND is_delete = ?", appID, 0).Count(&count).Error
	return count > 0, err
}

func Save(a *App) error {
	if a.ID == 0 {
		// 新增
		return DB.Create(a).Error
	} else {
		// 更新
		return DB.Save(a).Error
	}
}

func Update(a *App, updateData map[string]interface{}) error {
	return DB.Model(a).Where("id = ? AND is_delete = ?", a.ID, 0).Updates(updateData).Error
}

func Delete(a *App) error {
	return DB.Model(a).Update("is_delete", 1).Error
}

func GetCount() (int64, error) {
	var count int64
	err := DB.Model(&App{}).Where("is_delete = ?", 0).Count(&count).Error
	return count, err
}

// Query functions
func APPList(params *dto.APPListInput) ([]App, int64, error) {
	var list []App
	var count int64

	// 基础查询
	query := DB.Model(&App{}).Where("is_delete = ?", 0)

	// 搜索条件
	if params.Info != "" {
		query = query.Where("name LIKE ? OR app_id LIKE ?",
			"%"+params.Info+"%", "%"+params.Info+"%")
	}

	// 计算总数
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.PageNo - 1) * params.PageSize
	err = query.Offset(offset).Limit(params.PageSize).
		Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	return list, count, nil
}

// 构建器模式 - 提供更灵活的查询API
type AppQuery struct {
	db *gorm.DB
}

// NewAppQuery 创建新的查询构建器
func NewAppQuery() *AppQuery {
	return &AppQuery{
		db: DB.Model(&App{}).Where("is_delete = ?", 0),
	}
}

// Where 添加查询条件
func (q *AppQuery) Where(clause string, args ...interface{}) *AppQuery {
	q.db = q.db.Where(clause, args...)
	return q
}

// WhereLike 添加模糊查询条件
func (q *AppQuery) WhereLike(field, keyword string) *AppQuery {
	if keyword != "" {
		q.db = q.db.Where(field+" LIKE ?", "%"+keyword+"%")
	}
	return q
}

// Order 添加排序
func (q *AppQuery) Order(order string) *AppQuery {
	q.db = q.db.Order(order)
	return q
}

// Limit 设置条数限制
func (q *AppQuery) Limit(limit int) *AppQuery {
	q.db = q.db.Limit(limit)
	return q
}

// Offset 设置偏移量
func (q *AppQuery) Offset(offset int) *AppQuery {
	q.db = q.db.Offset(offset)
	return q
}

// Count 获取总数
func (q *AppQuery) Count() (int64, error) {
	var count int64
	err := q.db.Count(&count).Error
	return count, err
}

// Find 查询多条记录
func (q *AppQuery) Find() ([]App, error) {
	var list []App
	err := q.db.Find(&list).Error
	return list, err
}

// First 查询单条记录
func (q *AppQuery) First() (*App, error) {
	var app App
	err := q.db.First(&app).Error
	return &app, err
}

// Batch 批量操作
func BatchUpdate(ids []uint, updates map[string]interface{}) error {
	return DB.Model(&App{}).Where("id IN ?", ids).Updates(updates).Error
}

func BatchDelete(ids []uint) error {
	return DB.Model(&App{}).Where("id IN ?", ids).Update("is_delete", 1).Error
}

var AppManagerHandler *AppManager

func init() {
	AppManagerHandler = NewAppManager()
}

type AppManager struct {
	AppMap   map[string]*App
	AppSlice []*App
	Locker   sync.RWMutex
	init     sync.Once
	err      error
}

func NewAppManager() *AppManager {
	return &AppManager{
		AppMap:   make(map[string]*App),
		AppSlice: make([]*App, 0),
		Locker:   sync.RWMutex{},
		init:     sync.Once{},
	}
}

func (am *AppManager) GetAppList() []*App {
	am.Locker.RLock()
	defer am.Locker.RUnlock()
	return am.AppSlice
}

func (am *AppManager) GetApp(appID string) (*App, bool) {
	am.Locker.RLock()
	defer am.Locker.RUnlock()
	app, exists := am.AppMap[appID]
	return app, exists
}

func (am *AppManager) LoadOnce() error {
	return am.RefreshCache()
}

// 便利方法
func (am *AppManager) GetAppSecret(appID string) (string, error) {
	app, exists := am.GetApp(appID)
	if !exists {
		return "", ErrAppNotFound
	}
	return app.Secret, nil
}

func (am *AppManager) ValidateAppSecret(appID, secret string) bool {
	app, exists := am.GetApp(appID)
	if !exists {
		return false
	}
	return app.Secret == secret
}

// 错误定义
var (
	ErrAppNotFound      = errors.New("app not found")
	ErrAppAlreadyExists = errors.New("app already exists")
	ErrInvalidSecret    = errors.New("invalid secret")
)
