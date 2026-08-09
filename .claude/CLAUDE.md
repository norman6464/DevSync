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

## 3. クリーンアーキテクチャへの移行（進行中）

DevSync は現状レイヤード構成（`service` → `repository`、interface は `repository` 側で宣言）。
これを FreStyle と同じ**クリーンアーキテクチャ（依存性逆転・DIP）**へ段階移行する。

### 目標構造（FreStyle 準拠）

```
handler → usecase → usecase/repository(port) ← adapter/persistence(実装)
```

- **依存性逆転（DIP）**: interface（port）は**使う側の `usecase/repository/` で宣言**し、実装は `adapter/persistence/` に置く。依存の向きは usecase ← persistence
- **1 usecase = 1 責務**: `struct` + `NewXxxUseCase` コンストラクタ + `Execute(ctx, ...)`。複数操作を 1 つに詰め込まない
- **ctx を通す**: 全 `Execute` の第一引数に `context.Context` を取り、DB 操作（`db.WithContext(ctx)`）へ伝播する
- **port 充足の保証**: 実装ファイルに `var _ repository.XxxRepository = (*xxxRepository)(nil)` を置き、メソッド追加漏れをビルドで検出する
- **テスト**: usecase は testify/mock（port モック）、handler は「本物の usecase + port モック」で組む（handler は usecase を具象型で受け取る）

### 移行の進め方（ストラングラー）

- **1 スライスずつ**移行し、その都度 `go build` / `go vet` / `go test ./...` が緑であることを確認して PR にする
- **振る舞いを変えない**（既存の HTTP 挙動・レスポンスは同一に保つ）。純粋な構造変更に限定する
- 移行済みスライス: **follow**（パイロット・PR #2345）
- 外部連携を持つスライス（`github` / `atcoder` / `ai_*` / `email`）は切り出しやすいが、認証（`auth`）は Cognito 統合と二重作業になるため後回し
- **全スライスの移行が完了した後**に、層依存を機械検出する自作 linter（FreStyle の `archlint` 等）を CI に導入する

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
