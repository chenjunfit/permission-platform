package ioc

import (
	"context"
	"database/sql"
	"github.com/ecodeclub/ekit/retry"
	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/econf"
	"github.com/permission-dev/internal/repository/dao"
	"time"
)

func InitDB() *egorm.Component {
	WaitForDBSetUp(econf.GetString("mysql.dsn"))
	db := egorm.Load("mysql").Build()
	err := dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
func WaitForDBSetUp(dsn string) {
	//只验证dsn是否错误，并不会连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	//采用指数退避方式重试
	const MaxRetryInternal = 10 * time.Second
	const MaxtRetries = 10

	strategy, err := retry.NewExponentialBackoffRetryStrategy(time.Second*1, MaxRetryInternal, MaxtRetries)
	if err != nil {
		panic(err)
	}

	const TimeOut = 5 * time.Second
	for {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		err := db.PingContext(ctx)
		cancel()
		if err == nil {
			break
		}
		next, ok := strategy.Next()
		if !ok {
			panic("重试失败")
		}
		time.Sleep(next)
	}
}
