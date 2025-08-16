package rbac

import (
	"context"
	"fmt"
	"github.com/gotomicro/ego/core/econf"
	"github.com/permission-dev/internal/api/grpc/interceptor/auth"
	"github.com/permission-dev/internal/pkg/jwt"
	"github.com/permission-dev/internal/service/rbac"
	iocRbac "github.com/permission-dev/internal/test/integration/ioc/rbac"
	ioc2 "github.com/permission-dev/internal/test/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// 清理环境
func TestMain(m *testing.M) {
	_ = ioc2.InitDBAndTables()
	svc := iocRbac.Init()
	ctx := context.Background()
	t := &testing.T{}
	cleanTestEnviroment(t, ctx, svc)
	exitCode := m.Run()
	cleanTestEnviroment(t, ctx, svc)
	os.Exit(exitCode)
}
func TestRBACInit(t *testing.T) {
	//t.Skip("用于演示生成权限平台的SQL脚本的过程")
	// 1. 先将项目根目录下scripts/mysql/data.sql删掉,用scripts/mysql/database.sql中的内容覆盖scripts/mysql/init.sql中的内容。
	// 2. 运行权限服务
	// 3. 注释掉t.Skip()，手动执行当前测试
	// 4. 观察 scripts/mysql/init.sql 变化
	// 5. 开启 t.Skip()
	ctx := context.Background()
	iocSvc := iocRbac.Init()
	dir, jwtAuthKey, jwtIssuer, err := getJWTConfig()
	require.NoError(t, err)
	require.NotEmpty(t, jwtAuthKey)
	jwtToken := jwt.New(jwtAuthKey, jwtIssuer)
	bizID := int64(1)
	svc := rbac.NewInitService(bizID, 999, 3000, jwtToken,
		iocSvc.BusinessConfigRepo,
		iocSvc.ResourceRepo,
		iocSvc.PermissionRepo,
		iocSvc.RoleRepo,
		iocSvc.RolePermissionRepo,
		iocSvc.UserRoleRepo,
	)

	// 执行
	err = svc.Init(ctx)
	assert.NoError(t, err)

	// 验证
	bizConfig, err := iocSvc.BusinessConfigRepo.FindByID(ctx, bizID)
	assert.NoError(t, err)

	mapClaims, err := jwtToken.Decode(bizConfig.Token)
	assert.NoError(t, err)
	assert.Equal(t, float64(bizID), mapClaims[auth.BizIDName])

	// 调用make命令生成SQL脚本
	cmd := exec.Command("make", "mysqldump")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.Stdout, cmd.Stderr)
}

func getJWTConfig() (dir, key, issuer string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return "", "", "", err
	}
	rootDir := filepath.Clean(dir + "/../../../..")
	path := rootDir + "/config/config.yaml"
	f, err := os.Open(path)
	if err != nil {
		return "", "", "", err
	}

	err = econf.LoadFromReader(f, yaml.Unmarshal)
	if err != nil {
		return "", "", "", err
	}
	return rootDir, econf.GetString("jwt.key"), econf.GetString("jwt.issuer"), nil
}
