# 生成gRPC相关文件
.PHONY: grpc
grpc:
	@buf format -w api/proto
	@buf lint api/proto
	@buf generate api/proto
# 生成go代码
.PHONY: gen
gen:
	@wire ./cmd/platform/ioc/...

.PHONY: run
run:
	@cd cmd/platform && export EGO_DEBUG=true && go run main.go --config=../../config/config.yaml

#查看grpc接口列表
.PHONY: grpclist
grpclist:
	grpcurl -plaintext localhost:9002 list $$(grpcurl -plaintext localhost:9002 list|tail -n 1)

.PHONY: mysqldump
mysqldump:
	@docker exec -i mysql-mysql-1 mysqldump -u root -proot --databases permission > ./scripts/mysql/data.sql
	@echo ""                          >> ./scripts/mysql/init.sql
	@cat ./scripts/mysql/data.sql    >> ./scripts/mysql/init.sql