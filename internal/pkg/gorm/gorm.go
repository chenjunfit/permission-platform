package gorm

import (
	"context"
	"fmt"
	"github.com/gotomicro/ego/core/elog"
	permissionv1 "github.com/permission-dev/api/proto/gen/permission/v1"
	"gorm.io/gorm"
)

type ctxKey string

const (
	bizIDKey    ctxKey = "biz-id"
	uidKey      ctxKey = "uid"
	resourceKey ctxKey = "resource"
)

type GormAccessPlugin struct {
	client          permissionv1.PermissionServiceClient
	stateMap        map[StatementType]string
	logger          *elog.Component
	permissionToken string
}

type GormAccessPluginOption func(plugin *GormAccessPlugin)

func WithStatementMap(stateMap map[StatementType]string) GormAccessPluginOption {
	return func(plugin *GormAccessPlugin) {
		plugin.stateMap = stateMap
	}
}

func NewGormAccessPlugin(
	client permissionv1.PermissionServiceClient,
	permissionToken string,
	opts ...GormAccessPluginOption,
) *GormAccessPlugin {
	plugin := &GormAccessPlugin{
		client: client,
		stateMap: map[StatementType]string{
			SELECT: "read",
			UPDATE: "update",
			DELETE: "delete",
			CREATE: "create",
		},
		logger:          elog.DefaultLogger,
		permissionToken: permissionToken,
	}
	for index := range opts {
		opts[index](plugin)
	}
	return plugin
}
func (p *GormAccessPlugin) Name() string {
	return "GormAccessPlugin"
}

func (p *GormAccessPlugin) Initialize(db *gorm.DB) error {
	// 查询操作
	if err := db.Callback().Query().Before("gorm:query").Register("access:before_query", p.query); err != nil {
		return err
	}
	// 创建操作
	if err := db.Callback().Create().Before("gorm:create").Register("metrics:before_create", p.create); err != nil {
		return err
	}
	// 更新操作
	if err := db.Callback().Update().Before("gorm:update").Register("metrics:before_update", p.update); err != nil {
		return err
	}
	// 删除操作
	if err := db.Callback().Delete().Before("gorm:delete").Register("metrics:before_delete", p.delete); err != nil {
		return err
	}

	return nil
}
func (p *GormAccessPlugin) query(db *gorm.DB) {
	p.accessCheck(SELECT, db)
}

func (p *GormAccessPlugin) update(db *gorm.DB) {
	p.accessCheck(UPDATE, db)
}

func (p *GormAccessPlugin) create(db *gorm.DB) {
	p.accessCheck(CREATE, db)
}

func (p *GormAccessPlugin) delete(db *gorm.DB) {
	p.accessCheck(DELETE, db)
}

func (p *GormAccessPlugin) accessCheck(stmtType StatementType, db *gorm.DB) {
	ctx := db.Statement.Context
	uid, err := getUID(ctx)
	if err != nil {
		_ = db.AddError(fmt.Errorf("获取uid失败 %w", err))
		return
	}
	action, ok := p.stateMap[stmtType]
	if !ok {
		return
	}
	var key, resourceType string
	if val, ok := db.Statement.Model.(AuthRequire); ok {
		key = val.ResourceKey(ctx)
		resourceType = val.ResourceType(ctx)
	}
	//从ctx获取自定义的resourceKey,resourceType
	//可以实现列的权限控制，行的，多个列的
	//行 /resource/table/具体的id
	//列 /resource/table/某个列
	//多个列/resource/table/columus
	//resource key只是一个标志符的作用
	if val, rerr := getResource(ctx); rerr == nil {
		key, resourceType = val.Key, val.Type
	}
	if key != "" {
		ctx = context.WithValue(ctx, "Authorization", p.permissionToken)
		resp, perr := p.client.CheckPermission(ctx, &permissionv1.CheckPermissionRequest{
			Uid: uid,
			Permission: &permissionv1.Permission{
				ResourceKey:  key,
				ResourceType: resourceType,
				Actions:      []string{action},
			},
		})
		if perr != nil {
			_ = db.AddError(fmt.Errorf("权限校验失败 %w", err))
			elog.Error("权限校验失败",
				elog.FieldErr(err),
				elog.Int64("uid", uid),
				elog.String("action", action),
				elog.String("resourceKey", key),
				elog.String("resourceType", resourceType),
			)
			return
		}
		if !resp.Allowed {
			_ = db.AddError(fmt.Errorf("权限校验失败 %w", err))
		}
	}
}
func getUID(ctx context.Context) (int64, error) {
	value := ctx.Value(bizIDKey)
	fmt.Println(value)
	if value == nil {
		return 0, fmt.Errorf("uid not found in context")
	}
	uid, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid uid type, expected int64 got %T", value)
	}

	return uid, nil
}
func getResource(ctx context.Context) (Resource, error) {
	value := ctx.Value("resource")
	if value == nil {
		return Resource{}, fmt.Errorf("resource not found in ctx")
	}
	res, ok := value.(Resource)
	if !ok {
		return Resource{}, fmt.Errorf("invalid resource type, expected permissionv1.Resource")
	}
	return res, nil
}

type Resource struct {
	Key  string
	Type string
}
