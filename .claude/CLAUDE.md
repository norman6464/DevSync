# DevSync — Claude Code プロジェクト規約

> DevSync（エンジニアの活動可視化プラットフォーム）で Claude Code が作業するときの規約。
> アーキテクチャの参照実装は姉妹プロダクト FreStyle（`../` の `.claude/CLAUDE.md`）に揃える。

日本語で会話をしてください。

---

## 1. プロジェクト基本情報

- **プロジェクト名**: DevSync — GitHub / AtCoder 等の活動を定量化し、研修生同士で学習量を可視化するプラットフォーム
- **バックエンド**: Go / Gin / GORM（`backend/`）
- **フロントエンド**: React / TypeScript / Vite（`frontend/`）
- **RDB**: PostgreSQL
- **認証**: 現状は自前 JWT。将来 FreStyle の Cognito に統合予定（移行フェーズ 2）
- **リポジトリ**: `github.com/norman6464/DevSync`

---

## 2. 作業管理は Jira のみ（GitHub Issue は禁止）

**DevSync の作業管理は Jira の `DEVSYNC` プロジェクトだけで行う。GitHub Issue は絶対に作成しない・使わない。**

- ❌ **`gh issue create` を実行しない。** GitHub Issue を新規に立てることは禁止する
- ❌ 既存の GitHub Issue を作業管理に使わない（過去はブランチ名末尾の番号 `-2343` 等が GitHub Issue 番号だったが、この運用は廃止。今後は Jira に一本化する）
- ⭕ **作業のチケットは必ず Jira の `DEVSYNC` プロジェクトに起票する**
- **プロジェクト**: `https://frestyle.atlassian.net`（key = `DEVSYNC`）。Atlassian Rovo コネクタ（MCP）経由で操作する
- **チケットは着手前に作る**: 「作業を始める前に必ずチケットを起票する」。PR を先に作らない（順序: チケット → ブランチ → 実装 → PR）
- **起票は未割り当てで行う**: `assignee` を指定しない。着手時に自分にアサインする（[[jira-create-unassigned]] と同じ運用）
- **課題タイプ**: `Design Doc` / `リファクタリング` / `開発タスク` / `hotfix` / `ドキュメント整備` / `エピック`
  - 注意: `リファクタリング` と `開発タスク` は Jira 上の名前の**先頭に半角スペース 2 つ**が入っている（`  リファクタリング` / `  開発タスク`）。createJiraIssue の `issueTypeName` にはこの空白込みで渡す
- **設計判断を伴うものは `Design Doc`** で起票し、背景・選択肢・推奨案・承認記録を残す
- **PR とチケットの相互リンク必須**: チケットに PR URL を書き、PR 本文にチケット番号（`DEVSYNC-NN`）を書いて双方向に辿れるようにする
- **description はテンプレートに沿う**（概要 / ゴール / スコープ外 / 背景・目的 / テスト・検証 / セキュリティ影響 / 影響範囲 / ロールバック方針 / 参考リンク）

---

## 3. クリーンアーキテクチャ（移行完了）

DevSync は旧レイヤード構成（`service` → `repository`、interface は `repository` 側で宣言）から
FreStyle と同じ**クリーンアーキテクチャ（依存性逆転・DIP）**へ移行済み。`internal/service` と
`internal/repository` は撤去済みで、**新規コードをこの構造から外さない**ことがこの節の目的。

### 構造

```
handler → usecase → usecase/repository(port) ← adapter/persistence（DB 実装）
                                               adapter/external（外部 API 実装）
                                               infra（ws / scheduler など実行基盤）
```

配線は `internal/di`（コンテナ）と `internal/router`（ルーティング）が担う。この 2 つだけが全層を参照してよい。

- **依存性逆転（DIP）**: interface（port）は**使う側の `usecase/repository/` で宣言**し、実装は `adapter/persistence/` に置く。依存の向きは usecase ← persistence
- **1 usecase = 1 責務**: `struct` + `NewXxxUseCase` コンストラクタ + `Execute(ctx, ...)`。複数操作を 1 つに詰め込まない
- **ctx を通す**: 全 `Execute` の第一引数に `context.Context` を取り、DB 操作（`db.WithContext(ctx)`）へ伝播する
- **port 充足の保証**: 実装ファイルに `var _ repository.XxxRepository = (*xxxRepository)(nil)` を置き、メソッド追加漏れをビルドで検出する
- **テスト**: usecase は testify/mock（port モック）、handler は「本物の usecase + port モック」で組む（handler は usecase を具象型で受け取る）

### 依存方向の機械チェック（archlint）

層をまたぐ禁止依存は `backend/cmd/archlint` が静的に検出する。CI（`backend-test.yml`）で実行され、違反があればビルドを落とす。

```
cd backend && go run ./cmd/archlint .
```

- `domain` / `model` / `dto` は他の内部パッケージ・gin・net/http を import しない
- `usecase/repository`(port) は実装（`adapter`）や `gorm` を import しない（DIP の要）
- `usecase` は `handler` / `adapter` / `infra` / `dto` / gin / net/http を import しない
- `adapter` は `handler` / `usecase` 本体 / `dto` / gin を import しない（依存先は port だけ）
- `infra` は `handler` / `usecase` / `usecase/repository` / `adapter` / `dto` / gin を import しない
- `handler` は `adapter` / `usecase/repository` / `gorm` を import しない（usecase 経由にする）

やむを得ず外す場合は import 行末の `//archlint:allow`（1 行）かファイル先頭の `//archlint:ignore-file`（ファイル全体）を使い、**理由をコメントに書く**。抑制は残さないのが原則で、入れたら解消するチケットを起票する。

### 新しいスライスを足すときの型

- port は**使う側**（`usecase/repository/`）で宣言し、実装は `adapter/persistence`（DB）か `adapter/external`（外部 API）に置く
- 実装ファイルに `var _ repository.XxxRepository = (*xxxRepository)(nil)` を置き、メソッド追加漏れをビルドで検出する
- 不在は `(nil, nil)` に正規化する（`gorm.ErrRecordNotFound` を adapter で吸収し、usecase 側で意味づけする）
- 非同期に走らせる処理はリクエスト ctx で打ち切られないよう `context.WithoutCancel` を使う
- `go build` / `go vet` / `go test ./...` / `go run ./cmd/archlint .` が緑であることを確認して PR にする

### モデルの扱い

- DevSync は `model`（DB エンティティ）と `dto`（API 入出力）を分離している。この分離は FreStyle より厳格で、**尊重する**（`model` を `domain` に統合する大規模変更はしない）
- エラーは `domain` パッケージの sentinel（`domain.ErrBadRequest` 等）に寄せる

---

## 4. コーディング・PR 規約（FreStyle に準拠）

- **言語**: 日本語（PR / チケット / コミット / コメント）、英語（識別子）
- **他社プロダクト名を書かない**（コメント・PR・チケット・docs すべて）
- **コード内コメントに Jira 番号を書かない**（番号は PR / コミット / チケットに書く。[[no-jira-keys-in-code-comments]]）
- **コミットメッセージ / PR タイトルにはチケット番号（`DEVSYNC-NN`）を必ず載せる**
  - コミット: `refactor: follow スライスを DIP へ移行する (DEVSYNC-2)` のように末尾に `(DEVSYNC-NN)`
  - PR タイトル: 同様に末尾へ `(DEVSYNC-NN)`
  - チケットを作る前にコミット / PR しない（順序: Jira 起票 → ブランチ → 実装 → PR）
- **コミットメッセージ**: prefix（`feat` / `fix` / `refactor` / `test` / `docs` / `chore`）+ 日本語。末尾に `Co-Authored-By: <使用した Claude モデル名> <noreply@anthropic.com>`。「Generated with Claude Code」は書かない
- **PR 本文**: `## 概要` / `## 変更内容` / `## テスト` / `## 関連` の 4 セクション。関連にチケット番号（`DEVSYNC-NN`）を書く
- **テストなしでマージしない**。新規コードに Go のテストを付ける
- **`main` へ直接コミットしない**。ブランチ → PR → CodeRabbit レビュー → squash merge
- **シングルタスク運用**: PR は 1 つずつ完了させてから次に着手する
