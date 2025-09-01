package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gotomicro/ego/core/elog"
	"github.com/permission-dev/internal/repository/dao"
	"github.com/permission-dev/internal/repository/dao/audit"
	"github.com/permission-dev/pkg/canalx"
	"github.com/permission-dev/pkg/mqx"
	"time"
)

type UserRoleBinLogConsumer struct {
	consumer mqx.Consumer
	dao      audit.UserRoleLogDAO
	logger   *elog.Component
}

func NewUserRoleBinLogConsumer(consumer *kafka.Consumer, dao audit.UserRoleLogDAO, topic string) (*UserRoleBinLogConsumer, error) {
	err := consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}
	return &UserRoleBinLogConsumer{
		consumer: consumer,
		dao:      dao,
		logger:   elog.DefaultLogger,
	}, nil
}

func (u *UserRoleBinLogConsumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := u.Consume(ctx); err != nil {
				u.logger.Error("消费用户权限Binlog事件失败", elog.FieldErr(err))
				time.Sleep(time.Second) // 防止错误时无限循环
			}
		}
	}
}
func (u *UserRoleBinLogConsumer) Consume(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	msg, err := u.readMessageWithContext(timeoutCtx, -1)
	if err != nil {
		// 如果是context取消错误，直接返回，不记录为错误
		if ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("获取消息失败: %w", err)
	}
	var evt canalx.Message[dao.UserRole]
	err = json.Unmarshal(msg.Value, &evt)
	if err != nil {
		u.logger.Warn("解析消息失败",
			elog.FieldErr(err),
			elog.Any("msg", msg))
		return err
	}
	if evt.Table != evt.Data[0].TableName() ||
		(evt.Type != "INSERT" && evt.Type != "DELETE") {
		return nil
	}
	err = u.dao.BatchCreate(ctx, slice.Map(evt.Data, func(_ int, src dao.UserRole) audit.UserRoleLog {
		var beforeRoleID, afterRoleID int64
		if evt.Type == "INSERT" {
			afterRoleID = src.RoleID
		} else if evt.Type == "DELETE" {
			beforeRoleID = src.RoleID
		}
		return audit.UserRoleLog{
			Operation:    evt.Type,
			BizID:        src.BizID,
			UserID:       src.UserID,
			BeforeRoleID: beforeRoleID,
			AfterRoleID:  afterRoleID,
		}
	}))
	if err != nil {
		u.logger.Warn("创建用户权限变更操作日志失败",
			elog.FieldErr(err),
			elog.Any("evt", evt))
		return err
	}

	// 消费完成，提交消费进度
	_, err = u.consumer.CommitMessage(msg)
	if err != nil {
		u.logger.Warn("提交消息失败",
			elog.FieldErr(err),
			elog.Any("msg", msg))
		return err
	}

	return nil
}
func (u *UserRoleBinLogConsumer) readMessageWithContext(ctx context.Context, timeout time.Duration) (*kafka.Message, error) {
	resultChan := make(chan struct {
		msg *kafka.Message
		err error
	})
	go func() {
		msg, err := u.consumer.ReadMessage(timeout)
		resultChan <- struct {
			msg *kafka.Message
			err error
		}{msg: msg, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultChan:
		return result.msg, result.err
	}
}
