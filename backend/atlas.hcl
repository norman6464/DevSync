// Atlas プロジェクト設定。宣言的スキーマ管理（atlas schema apply）専用。
// バージョン管理された migrate diff は使わず、db/schema/*.hcl を desired state として
// 実DBへ直接差分適用する。

variable "db_url" {
  type    = string
  default = getenv("ATLAS_DB_URL")
}

env "local" {
  src = "file://db/schema"
  url = var.db_url

  // 差分計算のための使い捨てDB。ローカルのDockerでコンテナを都度立ち上げる。
  dev = "docker://postgres/18.4/dev"
}
