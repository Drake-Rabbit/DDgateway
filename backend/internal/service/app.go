package service

import (
	"gateway-service/internal/dto"
	"gateway-service/internal/models"
	"gateway-service/pkg/md5_ram"
)

// AppService 应用服务
type AppService struct{}

// NewAppService 创建应用服务
func NewAppService() *AppService {
	return &AppService{}
}

// GetAppList 获取应用列表
func (s *AppService) GetAppList(params *dto.APPListInput) ([]models.App, int64, error) {
	return models.APPList(params)
}

// GetApp 获取应用详情
func (s *AppService) GetApp(appID string) (*models.App, error) {
	return models.GetByAppID(appID)
}

// CreateApp 创建应用
func (s *AppService) CreateApp(input *dto.APPAddHttpInput) (*models.App, error) {
	// 检查AppID是否已存在
	exists, err := models.Exists(input.AppID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, models.ErrAppAlreadyExists
	}

	// 如果没有提供Secret，则生成随机密钥
	if input.Secret == "" {
		randString := md5_ram.GetRandomString(16)
		input.Secret = md5_ram.MD5(randString)
	}

	app := &models.App{
		AppID:    input.AppID,
		Name:     input.Name,
		Secret:   input.Secret,
		WhiteIPS: input.WhiteIPS,
		Qpd:      input.Qpd,
		Qps:      input.Qps,
	}

	if err := models.Save(app); err != nil {
		return nil, err
	}

	return app, nil
}

// UpdateApp 更新应用
func (s *AppService) UpdateApp(input *dto.APPUpdateHttpInput) error {
	// 检查App是否存在
	exists, err := models.Exists(input.AppID)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrAppNotFound
	}

	// 获取现有App
	app, err := models.GetByAppID(input.AppID)
	if err != nil {
		return err
	}

	// 构建更新数据
	updateData := make(map[string]interface{})
	if input.Name != "" {
		updateData["name"] = input.Name
	}
	if input.Secret != "" {
		updateData["secret"] = input.Secret
	}
	if input.WhiteIPS != "" {
		updateData["white_ips"] = input.WhiteIPS
	}
	if input.Qpd > 0 {
		updateData["qpd"] = input.Qpd
	}
	if input.Qps > 0 {
		updateData["qps"] = input.Qps
	}

	// 更新App
	if err := models.Update(app, updateData); err != nil {
		return err
	}

	// 刷新缓存

	return nil
}

// DeleteApp 删除应用
func (s *AppService) DeleteApp(appID string) error {
	// 检查App是否存在
	app, err := models.GetByAppID(appID)
	if err != nil {
		return models.ErrAppNotFound
	}

	// 软删除
	if err := models.Delete(app); err != nil {
		return err
	}

	return nil
}

// GetAppStats 获取App统计信息
func (s *AppService) GetAppStats() (map[string]interface{}, error) {
	total, err := models.GetCount()
	if err != nil {
		return nil, err
	}

	cacheInfo := map[string]interface{}{
		"total_apps":  total,
		"cached_apps": len(models.AppManagerHandler.GetAppList()),
	}

	return cacheInfo, nil
}

// SearchApp 搜索App（使用构建器模式）
func (s *AppService) SearchApp(params *dto.APPListInput) ([]models.App, int64, error) {
	// 使用新的构建器模式
	query := models.NewAppQuery().
		WhereLike("name", params.Info)

	count, err := query.Count()
	if err != nil {
		return nil, 0, err
	}

	query.Offset((params.PageNo - 1) * params.PageSize).
		Limit(params.PageSize)

	list, err := query.Find()
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}
