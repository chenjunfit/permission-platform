package gorm

import "context"

type AuthRequire interface {
	ResourceType(ctx context.Context) string
	ResourceKey(ctx context.Context) string
}

type StatementType string

const (
	SELECT StatementType = "SELECT"
	UPDATE StatementType = "UPDATE"
	DELETE StatementType = "DELETE"
	CREATE StatementType = "CREATE"
)
