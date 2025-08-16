package dao

import (
	"context"
	"github.com/ego-component/egorm"
	"time"
)

type BusinessConfig struct {
	ID        int64  `gorm:"primaryKey;autoIncrement;comment:'业务ID'"`
	OwnerID   int64  `gorm:"type:BIGINT;comment:'业务方ID'"`
	OwnerType string `gorm:"type:ENUM('person', 'organization');comment:'业务方类型：person-个人,organization-组织'"`
	Name      string `gorm:"type:VARCHAR(255);NOT NULL;comment:'业务名称'"`
	RateLimit int    `gorm:"type:INT;DEFAULT:1000;comment:'每秒最大请求数'"`
	Token     string `gorm:"type:TEXT;NOT NULL;comment:'业务方Token，内部包含bizID'"`
	Ctime     int64
	Utime     int64
}

func (BusinessConfig) TableName() string {
	return "business_configs"
}

type BusinessConfigDAO interface {
	Create(ctx context.Context, config BusinessConfig) (BusinessConfig, error)
	FindByIDs(ctx context.Context, ids []int64) (map[int64]BusinessConfig, error)
	GetByID(ctx context.Context, id int64) (BusinessConfig, error)
	Find(ctx context.Context, offset, limit int) ([]BusinessConfig, error)
	UpdateToken(ctx context.Context, id int64, token string) error
	Update(ctx context.Context, config BusinessConfig) error
	Delete(ctx context.Context, id int64) error
}

func NewBusinessConfigDAO(db *egorm.Component) BusinessConfigDAO {
	return &businessConfigDao{db: db}
}

type businessConfigDao struct {
	db *egorm.Component
}

func (b *businessConfigDao) Create(ctx context.Context, config BusinessConfig) (BusinessConfig, error) {
	now := time.Now().Unix()
	config.Utime = now
	config.Ctime = now
	err := b.db.WithContext(ctx).Model(&BusinessConfig{}).Create(&config).Error
	return config, err
}

func (b *businessConfigDao) FindByIDs(ctx context.Context, ids []int64) (map[int64]BusinessConfig, error) {
	var configs []BusinessConfig
	err := b.db.Where(ctx).Model(&BusinessConfig{}).Where("id id IN ?", ids).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	configMap := make(map[int64]BusinessConfig, len(ids))
	for index := range configs {
		config := configs[index]
		configMap[config.ID] = config
	}
	return configMap, nil
}

func (b *businessConfigDao) GetByID(ctx context.Context, id int64) (BusinessConfig, error) {
	var config BusinessConfig
	err := b.db.WithContext(ctx).Model(&BusinessConfig{}).Where("id = ?", id).First(&config).Error
	return config, err
}

func (b *businessConfigDao) Find(ctx context.Context, offset, limit int) ([]BusinessConfig, error) {
	var configs []BusinessConfig
	err := b.db.Where(ctx).Model(&BusinessConfig{}).Offset(offset).Limit(limit).Find(&configs).Error
	return configs, err
}

func (b *businessConfigDao) UpdateToken(ctx context.Context, id int64, token string) error {
	return b.db.WithContext(ctx).Model(&BusinessConfig{}).Where("id = ?", id).Updates(map[string]interface{}{
		"token": token,
		"utime": time.Now().Unix(),
	}).Error
}

func (b *businessConfigDao) Update(ctx context.Context, config BusinessConfig) error {
	return b.db.WithContext(ctx).Model(&BusinessConfig{}).Where("id=?", config.ID).Updates(map[string]any{
		"owner_id":   config.OwnerID,
		"owner_type": config.OwnerType,
		"name":       config.Name,
		"rate_limit": config.RateLimit,
		"utime":      config.Utime,
	}).Error
}

func (b *businessConfigDao) Delete(ctx context.Context, id int64) error {
	return b.db.WithContext(ctx).Model(&BusinessConfig{}).Where("id = ?", id).Delete(&BusinessConfig{}).Error
}
