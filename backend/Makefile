.PHONY: db-schema-apply db-schema-diff db-schema-lint db-schema-sql sqlc-generate

# Atlas: internal/infra/database/schema/*.hcl（宣言的desired state）を実DBへ差分適用する。
# ATLAS_DB_URL が未設定ならローカルの docker-compose db（localhost:5432/devsync）を既定値にする。
ATLAS_DB_URL ?= postgres://devsync:devsync@localhost:5432/devsync?sslmode=disable

db-schema-apply:
	ATLAS_DB_URL="$(ATLAS_DB_URL)" atlas schema apply --env local

# CI/レビュー用: 適用せずに差分だけ確認する。
db-schema-diff:
	ATLAS_DB_URL="$(ATLAS_DB_URL)" atlas schema apply --env local --dry-run

db-schema-lint:
	ATLAS_DB_URL="$(ATLAS_DB_URL)" atlas schema lint --env local

# internal/infra/database/schema/*.hcl を sqlc が読めるプレーンSQLへ書き出す
# （sqlc は HCL を直接読めないための橋渡し）。生成物であり手編集しない。.gitignore 済み。
# schema/ ディレクトリ自体に書き出すと、Atlas がそのディレクトリをHCLソースとして
# 読む際にSQLファイルと衝突するため、1つ上の internal/infra/database/ 直下に出力する。
db-schema-sql:
	atlas schema inspect -u "file://internal/infra/database/schema" --dev-url "docker://postgres/18.4/dev" --format '{{ sql . }}' > internal/infra/database/schema_generated.sql

# sqlc: internal/adapter/persistence/queries/*.sql から型付きGoコードを
# internal/adapter/persistence/sqlcgen へ生成する。
sqlc-generate: db-schema-sql
	sqlc generate
