package main

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
	"github.com/gotomicro/ego/server/egovernor"
	"github.com/gotomicro/ego/server/egrpc"
	ioc2 "github.com/permission-dev/cmd/platform/ioc"
	"github.com/permission-dev/internal/ioc"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	//创建ego实例
	egoApp := ego.New()
	// 初始化配置
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//zipkin
	tp := ioc.InitZipkinTracer()
	defer func(tp *trace.TracerProvider, ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			elog.Panic("zipkin关闭失败", elog.FieldErr(err))
		}
	}(tp, ctx)
	app := ioc2.InitApp()
	app.StartTasks(ctx)
	servers := make([]server.Server, 0, len(app.GrpcServers)+1)
	//metrics
	servers = append(servers, egovernor.Load("server.governor").Build())
	servers = append(servers, slice.Map(app.GrpcServers, func(_ int, src *egrpc.Component) server.Server {
		return src
	})...)
	if err := egoApp.Serve(servers...).Run(); err != nil {
		elog.Panic("启动失败", elog.FieldErr(err))
	}
}
