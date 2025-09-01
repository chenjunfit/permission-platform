package audit

import (
	"context"
	"github.com/ego-component/egorm"
	"time"
)

type OperationLog struct {
	ID       int64  `gorm:"primaryKey;autoIncrement;comment:'操作日志表自增ID'"`
	Operator string `gorm:"type:VARCHAR(255);comment:'操作者ID，通常为用户ID'"`
	Key      string `gorm:"type:VARCHAR(255);comment:'业务方内唯一标识，用于标识这次操作请求'"`
	BizID    int64  `gorm:"type:BIGINT;NOT NULL;comment:'表示在该业务ID下调用Method执行Request，Operator的业务ID与此业务ID无关系'"`
	Method   string `gorm:"type:TEXT;NOT NULL;comment:'调用的接口名称'"`
	Request  string `gorm:"type:TEXT;NOT NULL;comment:'请求参数JSON序列化后的字符串'"`
	Ctime    int64
	Utime    int64
}

func (o OperationLog) TableName() string {
	return "operation_log"
}

type OperationLogDAO interface {
	Create(ctxt context.Context, log OperationLog) (int64, error)
}

type operationLogDao struct {
	db *egorm.Component
}

func NewOperationLogDao(db *egorm.Component) *operationLogDao {
	return &operationLogDao{db: db}
}

func (o *operationLogDao) Create(ctx context.Context, log OperationLog) (int64, error) {
	now := time.Now().UnixMilli()
	log.Utime = now
	log.Ctime = now
	err := o.db.WithContext(ctx).Model(&OperationLog{}).Create(&log).Error
	return log.ID, err
}
