.PHONY: ensure-atlas db-schema-apply db-schema-diff db-schema-lint db-schema-sql sqlc-generate

# Atlas: internal/infra/database/schema/schema.hcl（宣言的desired state）を実DBへ差分適用する。
# プロジェクト設定ファイル（atlas.hcl）は持たず、対象ごとに --url/--to/--dev-url を都度指定する
# （バージョン管理された migrate diff は使わない）。
# ATLAS_DB_URL が未設定ならローカルの docker-compose db（localhost:5432/devsync）を既定値にする。
ATLAS_DB_URL ?= postgres://devsync:devsync@localhost:5432/devsync?sslmode=disable

SCHEMA_DIR := internal/infra/database/schema
SCHEMA_HCL := $(SCHEMA_DIR)/schema.hcl
SCHEMA_GEN_SQL := $(SCHEMA_DIR)/schema.gen.sql

# 差分計算のための使い捨てDB。ローカルのDockerでコンテナを都度立ち上げる。
ATLAS_DEV_URL := docker://postgres/18.4/dev

# Atlas CLI: アーキテクチャ別のバイナリを直接取得する。版は固定（-latestはcanaryで日によって挙動が変わる）。
ATLAS_VERSION := v1.3.3
GO_INSTALL_DIR := $(shell go env GOBIN)
ifeq ($(GO_INSTALL_DIR),)
GO_INSTALL_DIR := $(shell go env GOPATH)/bin
endif
ATLAS := $(GO_INSTALL_DIR)/atlas
ATLAS_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ATLAS_ARCH := $(shell uname -m | sed 's/x86_64/amd64/; s/aarch64/arm64/')

ensure-atlas:
	@if [ ! -x "$(ATLAS)" ] || ! "$(ATLAS)" version 2>/dev/null | grep -q '$(ATLAS_VERSION)'; then \
	  echo ">>> installing atlas ($(ATLAS_VERSION)) for $(ATLAS_OS)/$(ATLAS_ARCH)..."; \
	  mkdir -p "$(GO_INSTALL_DIR)"; \
	  curl -fsSL -o "$(ATLAS)" "https://release.ariga.io/atlas/atlas-$(ATLAS_OS)-$(ATLAS_ARCH)-$(ATLAS_VERSION)"; \
	  chmod +x "$(ATLAS)"; \
	fi

db-schema-apply: ensure-atlas
	@"$(ATLAS)" schema apply --url "$(ATLAS_DB_URL)" --to "file://$(SCHEMA_HCL)" --dev-url "$(ATLAS_DEV_URL)"

# CI/レビュー用: 適用せずに差分だけ確認する。
db-schema-diff: ensure-atlas
	@"$(ATLAS)" schema apply --url "$(ATLAS_DB_URL)" --to "file://$(SCHEMA_HCL)" --dev-url "$(ATLAS_DEV_URL)" --dry-run

db-schema-lint: ensure-atlas
	@"$(ATLAS)" schema lint --url "file://$(SCHEMA_HCL)" --dev-url "$(ATLAS_DEV_URL)"

# schema.hcl を sqlc が読めるプレーンSQLへ書き出す（sqlc は HCL を直接読めないための橋渡し）。
# schema.hcl は単一ファイルなのでAtlasの読み込み先と衝突せず、schema/ 配下に共存できる。
# schema.gen.sql は DO NOT EDIT な生成物だがコミット対象（internal/infra/database/schema.go が
# go:embed してバイナリへ埋め込み、空DBへの起動時自己適用に使うため）。
# CIはこのtargetを再実行し、git diffでdrift（手直し・生成し忘れ）を検査する。
db-schema-sql: ensure-atlas
	@"$(ATLAS)" schema inspect --url "file://$(SCHEMA_HCL)" --dev-url "$(ATLAS_DEV_URL)" --format '{{ sql . }}' > $(SCHEMA_GEN_SQL)

# sqlc: internal/adapter/persistence/queries/*.sql から型付きGoコードを
# internal/adapter/persistence/sqlcgen へ生成する。
sqlc-generate: db-schema-sql
	sqlc generate
