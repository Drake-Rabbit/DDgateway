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
}

func (a *App) TableName() string {
	return "gateway_app"
}

// CRUD functions

func GetAppByPrimaryID(id uint) (*App, error) {
	model := &App{}
	err := DB.Where("id = ?", id).First(model).Error
	return model, err
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
	return DB.Model(a).Where("id = ? ", a.ID).Updates(updateData).Error
}

func Delete(a *App) error {
	return DB.Model(a).Delete(&App{}).Where("id = ?", a.ID).Error
}

func GetCount() (int64, error) {
	var count int64
	err := DB.Model(&App{}).Count(&count).Error
	return count, err
}

// Query functions
func APPList(params *dto.APPListInput) ([]App, int64, error) {
	var list []App
	var count int64

	// 基础查询
	query := DB.Model(&App{})

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
	am.init.Do(func() {
		params := &dto.APPListInput{PageNo: 1, PageSize: 99999}
		list, _, err := APPList(params)
		if err != nil {
			am.err = err
			return
		}

		am.Locker.Lock()
		defer am.Locker.Unlock()
		for _, listItem := range list {
			tmpItem := listItem
			am.AppMap[listItem.AppID] = &tmpItem
			am.AppSlice = append(am.AppSlice, &tmpItem)
		}
	})
	return am.err
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
