.PHONY: db-schema-apply db-schema-diff db-schema-lint db-schema-sql sqlc-generate

# Atlas: db/schema/*.hcl（宣言的desired state）を実DBへ差分適用する。
# ATLAS_DB_URL が未設定ならローカルの docker-compose db（localhost:5432/devsync）を既定値にする。
ATLAS_DB_URL ?= postgres://devsync:devsync@localhost:5432/devsync?sslmode=disable

db-schema-apply:
	ATLAS_DB_URL="$(ATLAS_DB_URL)" atlas schema apply --env local

# CI/レビュー用: 適用せずに差分だけ確認する。
db-schema-diff:
	ATLAS_DB_URL="$(ATLAS_DB_URL)" atlas schema apply --env local --dry-run

db-schema-lint:
	ATLAS_DB_URL="$(ATLAS_DB_URL)" atlas schema lint --env local

# db/schema/*.hcl を sqlc が読めるプレーンSQLへ書き出す（sqlc は HCL を直接読めないための橋渡し）。
# 生成物であり手編集しない。db/schema.sql は .gitignore 済み。
db-schema-sql:
	atlas schema inspect -u "file://db/schema" --dev-url "docker://postgres/18.4/dev" --format '{{ sql . }}' > db/schema_generated.sql

# sqlc: db/queries/*.sql から型付きGoコードを internal/adapter/persistence/sqlcgen へ生成する。
sqlc-generate: db-schema-sql
	sqlc generate
