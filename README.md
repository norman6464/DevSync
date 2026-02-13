# DevSync

<p align="center">
  <img src="docs/images/logo.png" alt="DevSync Logo" width="200">
</p>

<p align="center">
  <strong>エンジニアのためのモチベーション共有プラットフォーム</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#demo">Demo</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#tech-stack">Tech Stack</a> •
  <a href="#getting-started">Getting Started</a> •
  <a href="#contributing">Contributing</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/React-18.2-61DAFB?style=flat-square&logo=react" alt="React">
  <img src="https://img.shields.io/badge/TypeScript-5.3-3178C6?style=flat-square&logo=typescript" alt="TypeScript">
  <img src="https://img.shields.io/badge/Kubernetes-1.29-326CE5?style=flat-square&logo=kubernetes" alt="Kubernetes">
  <img src="https://img.shields.io/badge/AWS-EKS-FF9900?style=flat-square&logo=amazonaws" alt="AWS">
</p>

---

## 📖 概要

**DevSync** は、エンジニア同士がGitHub・Zenn・Qiitaの活動を共有し、お互いの成長を可視化してモチベーションを高め合うプラットフォームです。

### 開発背景

未経験〜駆け出しエンジニアとして毎日勉強する中で、以下の課題を感じていました：

- 一人で学習を続けるとモチベーション維持が難しい
- 周りのエンジニアの頑張りを見て刺激を受けたい
- 既存サービスは「自分の草を見る」だけで、共有・競争の要素がない

DevSyncは、これらの課題を解決するために開発しました。

---

## ✨ Features

### 🌱 GitHub活動の可視化
- コントリビューション（草）カレンダーの表示
- リポジトリごとのコミット分析
- 使用言語の統計・可視化

### 📊 ランキング機能
- 言語別ランキング（Go, TypeScript, Python...）
- 週間・月間コントリビューションランキング
- 資格・スキルカテゴリ別ランキング

### 👥 ソーシャル機能
- ユーザーフォロー・フォロワー
- タイムライン（フォロー中の活動）
- ユーザー検索・探索

### 💬 コミュニケーション
- ダイレクトメッセージ（チャット）
- 学習報告の投稿
- いいね・コメント

### 🎯 学習管理
- 学習目標の設定・進捗追跡
- 週間・月間アクティビティレポート
- 技術書レビュー共有

### 🗺 学習ロードマップテンプレート
- 5種類のキャリアパス別おすすめテンプレート
- テンプレートからワンクリックでロードマップ作成
- ステップ詳細のプレビュー機能
- 進捗追跡・ステップ完了管理

### 🗂 プロジェクトショーケース
- GitHubリポジトリの詳細紹介
- デモURL・スクリーンショット
- 技術スタック・役割の記載
- 注目プロジェクトのハイライト

### 📚 学習リソース共有
- おすすめ教材・チュートリアルの共有
- カテゴリ別（書籍、動画、記事、コース等）
- 難易度設定・タグ付け
- いいね・保存機能

### 📋 ポートフォリオ生成
- プロフィール情報から自動でポートフォリオサイト生成
- 3つのテーマから選択（ミニマル、モダン、グラデーション）
- HTMLダウンロード対応

### 👥 学習グループ（スタディサークル）
- 3〜10人の小グループで同じテーマを学習
- 共有ロードマップ＋メンバー別進捗可視化
- 日次チェックイン（今日やったことを一言で報告）
- グループ内ストリークランキング（連続日数で競争）
- メンバー管理（招待・除外・退出）

### 🤖 AI機能
- 投稿時のAI補完
- スキルに基づく学習アドバイス
- おすすめユーザーのマッチング

### 🔥 ストリーク＆デイリーチャレンジ
- 学習ログの連続記録日数（ストリーク）を自動追跡
- 最長記録・合計学習日数の表示
- 日替わりミニチャレンジ（10種類、日付ベースで毎日変化）
- ダッシュボード＆プロフィールにストリーク情報を表示
- GitHubコントリビューション or 学習ログでバッジ（week-streak / month-streak）獲得可能

### 🏆 競技プログラミング連携
- AtCoderレーティング表示（色付きランク・レーティング値をプロフィールに表示）
- paizaランク表示（S/A/B/C/D/E の自己申告制）
- 設定画面・オンボーディングから連携設定が可能

### 💻 コードスニペット共有 & レビュー
- 投稿にコードスニペットを添付（30以上の言語に対応）
- Prism.jsベースのシンタックスハイライト（VS Code Dark+テーマ）
- Markdownコードブロックの自動ハイライト（エディタプレビュー・カード・詳細ページ）
- GitHub PR風のインラインコメント（行ごとにレビューコメント可能）
- ワンクリックでコードコピー
- 10言語のi18n対応

### 📊 学習分析ダッシュボード
- ヒートマップ（曜日×時間帯の学習密度）
- カテゴリ別学習時間の可視化（比率計算付き）
- 週間トレンド（前週比成長率）
- 生産性スコア（ポモドーロ・目標達成・ストリークの加重平均）
- AIインサイト（学習パターンの自動分析・アドバイス）

### 📝 外部サービス連携
- GitHub
- Zenn
- Qiita
- AtCoder
- paiza

---

## 🎬 Demo

### ホーム画面
<p align="center">
  <img width="1914" height="978" alt="スクリーンショット 2026-02-08 10 46 47" src="https://github.com/user-attachments/assets/e8c018f3-c1ca-4cc3-a382-6afcfa0d1808" />
</p>

### プロフィール画面
<p align="center">
  <img width="1897" height="819" alt="スクリーンショット 2026-02-08 10 47 56" src="https://github.com/user-attachments/assets/b2c582d5-715f-4f10-af4f-a37994ed4a4f" />
</p>

<p align="center">
  <img width="1895" height="977" alt="スクリーンショット 2026-02-08 10 48 33" src="https://github.com/user-attachments/assets/8fa32ea2-f2e6-45ff-bbf8-b691a004de87" />
</p>


### ランキング画面
<p align="center">
  <img width="1895" height="972" alt="スクリーンショット 2026-02-08 10 49 00" src="https://github.com/user-attachments/assets/81701201-335d-4c09-86b8-68c4a2733898" />
</p>

---

## 🏗 Architecture

### システム構成図

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              AWS Cloud                                   │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                            EKS Cluster                             │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐     │  │
│  │  │  User   │ │ GitHub  │ │ Article │ │ Ranking │ │  Chat   │     │  │
│  │  │ Service │ │ Service │ │ Service │ │ Service │ │ Service │     │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘     │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────────────────────────────────┐ │  │
│  │  │   AI    │ │  Notif  │ │  ArgoCD / Prometheus / Grafana      │ │  │
│  │  │ Service │ │ Service │ │                                     │ │  │
│  │  └─────────┘ └─────────┘ └─────────────────────────────────────┘ │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│         │              │              │              │                   │
│    ┌────▼────┐    ┌────▼────┐    ┌────▼────┐    ┌────▼────┐            │
│    │   RDS   │    │  Redis  │    │   SQS   │    │ Bedrock │            │
│    │(Postgres)│    │         │    │         │    │  (AI)   │            │
│    └─────────┘    └─────────┘    └─────────┘    └─────────┘            │
└─────────────────────────────────────────────────────────────────────────┘
```

### バックエンド内部アーキテクチャ（クリーンアーキテクチャ）

```
Handler (HTTP) → Service (ビジネスロジック) → Repository (Interface) → GORM (DB)
```

| レイヤー | 責務 | パッケージ |
|----------|------|-----------|
| Handler | HTTPリクエスト/レスポンス、バリデーション | `internal/handler` |
| Service | ビジネスロジック、権限チェック、通知 | `internal/service` |
| Repository | データアクセス抽象化（23インターフェース） | `internal/repository` |
| Model | データ構造、ドメインオブジェクト | `internal/model` |
| Router | DI（依存性注入）、ルーティング | `internal/router` |

### マイクロサービス構成

| サービス | 責務 | 技術 |
|----------|------|------|
| User Service | ユーザー認証、プロフィール、フォロー管理 | Go, PostgreSQL |
| GitHub Service | GitHub API連携、草/コミット/言語データ取得 | Go, PostgreSQL |
| Article Service | Zenn/Qiita API連携、記事データ取得 | Go, PostgreSQL |
| Ranking Service | ランキング集計、スコア計算 | Go, Redis |
| Chat Service | ユーザー間チャット、WebSocket | Go, Redis |
| AI Service | AI補完、レコメンド、マッチング | Go, Bedrock |
| Notification Service | 通知管理（プッシュ、メール） | Go, SQS |

---

## 🛠 Tech Stack

### Backend
| 技術 | バージョン | 用途 |
|------|-----------|------|
| Go | 1.21+ | マイクロサービス実装 |
| Gin | 1.10 | HTTPフレームワーク |
| GORM | 1.25 | ORM |
| testify | 1.9 | テストフレームワーク（モック＆アサーション） |
| gRPC | - | サービス間通信 |
| PostgreSQL | 15 | メインデータベース |
| Redis | 7.x | キャッシュ、ランキング |

### Frontend
| 技術 | バージョン | 用途 |
|------|-----------|------|
| React | 18.2 | UIフレームワーク |
| TypeScript | 5.3 | 型安全な開発 |
| Vite | 5.x | ビルドツール |
| Tailwind CSS | 3.4 | スタイリング |
| TanStack Query | 5.x | データフェッチ |
| Zustand | 4.x | 状態管理 |

### Infrastructure
| 技術 | 用途 |
|------|------|
| AWS EKS | Kubernetesクラスター |
| Terraform | Infrastructure as Code |
| ArgoCD | GitOpsによるCD |
| GitHub Actions | CI |

### Observability
| 技術 | 用途 |
|------|------|
| OpenTelemetry | 分散トレーシング |
| Prometheus | メトリクス収集 |
| Grafana | 可視化・ダッシュボード |
| CloudWatch | ログ管理 |

---

## 🚀 Getting Started

### 必要条件

- Go 1.21+
- Node.js 20+
- Docker & Docker Compose
- kubectl
- AWS CLI（デプロイ時）

### ローカル開発環境のセットアップ

#### 1. リポジトリのクローン

```bash
git clone https://github.com/yourusername/devsync.git
cd devsync
```

#### 2. 環境変数の設定

```bash
# Backend
cp backend/.env.example backend/.env

# Frontend
cp frontend/.env.example frontend/.env
```

#### 3. Docker Composeで起動

```bash
docker-compose up -d
```

#### 4. マイグレーション実行

```bash
make migrate-up
```

#### 5. アプリケーションにアクセス

- Frontend: http://localhost:5173
- API: http://localhost:8080
- Mailhog: http://localhost:8025（開発用メールUI）
- Grafana: http://localhost:3000

### Kubernetesへのデプロイ

#### 1. EKSクラスターの作成

```bash
cd infrastructure/terraform
terraform init
terraform plan
terraform apply
```

#### 2. ArgoCDのセットアップ

```bash
kubectl apply -k infrastructure/argocd
```

#### 3. アプリケーションのデプロイ

```bash
kubectl apply -f infrastructure/argocd/applications/
```

---

## 📁 Project Structure

```
devsync/
├── backend/                    # バックエンドサービス
│   ├── services/
│   │   ├── user/              # User Service
│   │   ├── github/            # GitHub Service
│   │   ├── article/           # Article Service
│   │   ├── ranking/           # Ranking Service
│   │   ├── chat/              # Chat Service
│   │   ├── ai/                # AI Service
│   │   └── notification/      # Notification Service
│   ├── pkg/                   # 共通パッケージ
│   ├── proto/                 # Protocol Buffers
│   └── docker-compose.yml
│
├── frontend/                   # フロントエンド
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── hooks/
│   │   ├── lib/
│   │   ├── stores/
│   │   └── types/
│   └── package.json
│
├── infrastructure/             # インフラ設定
│   ├── terraform/             # Terraform
│   ├── kubernetes/            # K8sマニフェスト
│   └── argocd/                # ArgoCD設定
│
├── docs/                       # ドキュメント
│   ├── api/                   # API仕様
│   ├── architecture/          # 設計ドキュメント
│   └── images/                # 画像
│
└── README.md
```

---

## 📊 Observability

### Grafana ダッシュボード

DevSyncでは以下のダッシュボードを提供しています：

- **Overview**: システム全体の健全性
- **Service Metrics**: 各サービスのリクエスト数、レイテンシ、エラー率
- **Infrastructure**: CPU、メモリ、ネットワーク使用率
- **Business Metrics**: ユーザー数、アクティブ率、ランキング参加者数

<p align="center">
  <img src="docs/images/screenshots/grafana.png" alt="Grafana Dashboard" width="800">
</p>

### 分散トレーシング

OpenTelemetryによる分散トレーシングで、マイクロサービス間のリクエストフローを可視化しています。

---

## 🔧 Development

### コマンド一覧

```bash
# Docker Compose（全サービス起動）
docker compose up -d --build

# Backend テスト
cd backend && go test ./internal/service/... -v        # 全テスト実行
cd backend && go test ./internal/service/... -cover     # カバレッジ付き実行
cd backend && go test ./internal/service/... -run TestAuth  # 特定テストのみ

# Frontend
npm run dev          # 開発サーバー起動
npm run build        # ビルド
npm run lint         # Linter実行

# Infrastructure
make tf-plan         # Terraform plan
make tf-apply        # Terraform apply
make k8s-deploy      # Kubernetesデプロイ
```

### テスト構成

Service層のユニットテストを `testify/mock` ベースで実装しています。
全リポジトリインターフェースのモックを使用し、DB不要で高速に実行可能です。

| テストファイル | 対象 | テスト数 |
|---------------|------|---------|
| `auth_test.go` | 認証・JWT・パスワードリセット | 26 |
| `post_test.go` | 投稿CRUD・所有権チェック | 7 |
| `follow_test.go` | フォロー・自己フォロー禁止 | 7 |
| `question_test.go` | 質問CRUD・所有権チェック | 8 |
| `answer_test.go` | 回答・ベストアンサー権限 | 12 |
| `learning_goal_test.go` | 目標・進捗自動完了 | 8 |
| `learning_log_test.go` | 学習ログ・ストリーク・バリデーション | 13 |
| `project_test.go` | プロジェクト・Featured | 7 |
| `book_review_test.go` | 書籍レビュー | 5 |
| `learning_resource_test.go` | リソース・可視性制御 | 12 |
| `roadmap_test.go` | ロードマップ・ステップ所属・テンプレート | 28 |
| `chat_room_test.go` | チャット・メンバーシップ | 16 |
| `user_test.go` | ユーザー検索分岐 | 4 |
| `notification_test.go` | 通知・フォロワー通知 | 5 |
| `message_test.go` | メッセージ・既読 | 4 |
| `code_snippet_test.go` | スニペットCRUD・インラインコメント | 10 |
| `ai_advice_test.go` | ルールエンジン・LLMチャット・会話削除 | 17 |
| `email_test.go` | ウィークリーレポートメール送信・テンプレート | 6 |
| `level_test.go` | XP計算・レベル計算・LevelInfo・Service統合 | 28 |
| `learning_analytics_test.go` | ヒートマップ・カテゴリ・トレンド・生産性スコア・インサイト | 22 |
| `recommendation_test.go` | おすすめユーザー・トレンド投稿・トレンドリソース | 7 |
| `study_circle_test.go` | サークルCRUD・メンバー管理・ステップ・進捗・チェックイン・ストリーク | 29 |

### CI/CD

| ワークフロー | トリガー | 内容 |
|-------------|---------|------|
| Backend Tests | PR / mainプッシュ | Service層ユニットテスト実行、カバレッジレポート生成 |

GitHub Actionsでテストを自動実行し、PRマージ前にテスト通過を確認できます。

### コーディング規約

- **Go**: [Effective Go](https://go.dev/doc/effective_go) + [uber-go/guide](https://github.com/uber-go/guide)
- **TypeScript**: ESLint + Prettier
- **コミットメッセージ**: [Conventional Commits](https://www.conventionalcommits.org/)

---

## 📝 API Documentation

API仕様はOpenAPI (Swagger)で管理しています。

- **Swagger UI**: http://localhost:8080/swagger
- **OpenAPI Spec**: [docs/api/openapi.yaml](docs/api/openapi.yaml)

---

## 🗺 Roadmap

### Phase 1（MVP）✅
- [x] GitHub連携
- [x] 草カレンダー表示
- [x] 言語別統計
- [x] ランキング機能
- [x] ユーザーフォロー

### Phase 2（拡張）✅
- [x] Zenn/Qiita連携
- [x] 投稿機能
- [x] チャット機能
- [x] 通知システム
- [x] 学習目標管理
- [x] アクティビティレポート
- [x] プロフィール共有

### Phase 3（ショーケース）✅
- [x] プロジェクトショーケース
- [x] 学習リソース共有
- [x] ポートフォリオ生成
- [x] 技術書レビュー
- [x] 10言語対応（日本語、英語、中国語、韓国語、スペイン語、フランス語、ドイツ語、ポルトガル語、ロシア語）

### Phase 3.5（エンゲージメント）✅
- [x] 学習ログ（毎日の学習日記機能）
- [x] 学習ストリーク（連続記録日数の自動追跡）
- [x] デイリーチャレンジ（日替わりミニタスク）
- [x] バッジシステムとの連携（GitHub + 学習ログストリーク）

### Phase 3.6（品質基盤）✅
- [x] クリーンアーキテクチャ（Handler → Service → Repository Interface）
- [x] Service層ユニットテスト（165テストケース、16サービス対象）
- [x] GitHub Actions CI（テスト自動実行）
- [x] testify/mockによるモックベーステスト

### Phase 3.7（セキュリティ強化）✅
- [x] JWT認証をhttpOnly Cookieに移行（XSS対策）
- [x] ログアウトエンドポイント追加（`POST /api/v1/auth/logout`）
- [x] localStorageからのトークン管理を完全廃止
- [x] SameSite=Lax属性によるCSRF対策
- [x] 認証ミドルウェア・ハンドラーのユニットテスト追加

### Phase 3.8（データベーススキーマ管理）✅
- [x] 全29テーブルのDDL定義（`backend/db/migrations/000001_create_all_tables.up.sql`）
- [x] ロールバック用SQL（`000001_create_all_tables.down.sql`）
- [x] 実DBスキーマと完全一致する`schema.sql`（`backend/db/schema.sql`）
- [x] 外部キー制約24件・インデックス79件・論理削除5テーブルの完全定義

### Phase 3.9（競技プログラミング連携）✅
- [x] AtCoderレーティング表示（レーティング値・色付きランク）
- [x] paizaランク表示（S/A/B/C/D/E 自己申告制）
- [x] AtCoderユーザー名のバリデーション（外部API検証）
- [x] 設定画面・オンボーディングからの連携設定UI
- [x] プロフィールページでの競技プログラミングセクション表示
- [x] 全10言語のi18n対応

### Phase 3.10（認証セッション永続化）✅
- [x] ページリロード時の認証状態復元（httpOnly Cookieベースのセッション維持）
- [x] アプリ起動時に無条件で `/auth/me` を呼び出し、Cookie検証による自動ログイン復元
- [x] 認証チェック中のローディング表示（ProtectedRouteとの連携）
- [x] 401インターセプターの最適化（初回認証チェック時の不要なリダイレクト防止）

### Phase 3.11（UI・認証バグ修正）✅
- [x] paizaランク表示の文字色を通常色（白）に統一
- [x] GitHubログイン時のアバターURL取得ロジックを修正（メール連携・GitHub連携時にも正しく反映）

### Phase 3.12（学習ロードマップテンプレート）✅
- [x] 5種類のおすすめ学習ロードマップテンプレート（Webフロントエンド、バックエンド Go、フルスタック、インフラ/DevOps、モバイル React Native）
- [x] テンプレートからワンクリックでロードマップ作成（ステップも自動コピー）
- [x] テンプレート一覧表示・ステップ詳細のプレビュー機能
- [x] アプリ起動時の自動シード（システムユーザー自動作成、既存テンプレートがある場合はスキップ）
- [x] TDD（テスト駆動開発）で実装（8テストケース追加）
- [x] 全10言語のi18n対応

### Phase 3.13（コードスニペット共有 & レビュー）✅
- [x] コードスニペットモデル定義（CodeSnippet + SnippetComment）
- [x] リポジトリ・サービス・ハンドラー実装（クリーンアーキテクチャ）
- [x] TDDで実装（10テストケース追加、合計165テスト）
- [x] 投稿作成時にスニペット一括添付（PostForm拡張）
- [x] Markdownプレビュー・PostCard・詳細ページのシンタックスハイライト（Prism.js + VS Code Dark+テーマ）
- [x] GitHub PR風インラインコメント（行ごとのレビュー機能）
- [x] コードコピー・言語選択（30言語対応）
- [x] 全10言語のi18n対応

### Phase 3.14（AIラーニングアドバイザー）✅
- [x] ハイブリッド型AIアドバイザー（ルールベース + LLMオンデマンド）
- [x] 12ルールのルールエンジン（ストリーク・ロードマップ・目標・GitHub活動分析）
- [x] OpenAI GPT-4o-mini連携（net/httpのみ、外部依存なし）
- [x] `OPENAI_API_KEY`未設定時の自動縮退（ルールベースのみ、チャットボタン非表示）
- [x] LLMチャット機能（1日5回制限、会話履歴保存）
- [x] 会話削除機能（確認ダイアログ付き、会話リスト・チャットパネルから削除可能）
- [x] ダッシュボードウィジェット + 専用ページ（レスポンシブ2カラム）
- [x] TDDで実装（17テストケース追加、合計182テスト）
- [x] 全10言語のi18n対応

### Phase 3.15（ポモドーロタイマー）✅
- [x] ポモドーロテクニックに基づく学習タイマー（25分集中→5分休憩→15分長い休憩）
- [x] Web Audio APIによる通知音（集中完了/休憩完了で異なる音色）
- [x] フローティングウィジェット（全ページで利用可能、最小化/展開切替）
- [x] SVG円形プログレス表示、フェーズごとの色分け
- [x] 集中セッション完了時の学習ログ自動記録（source='pomodoro'）
- [x] 学習ログモデルにsourceフィールド追加（manual/pomodoro区別）
- [x] Duration・Category・Sourceのバリデーション追加（TDD: 5テストケース）
- [x] 全10言語のi18n対応

### Phase 3.16（セキュリティ修正 — Critical）✅
- [x] パスワードリセットトークンのAPIレスポンス漏洩を修正（アカウント乗っ取り防止）
- [x] JWTシークレットのハードコードデフォルト値を除去（トークン偽造防止）
- [x] ポートフォリオHTML生成のXSS脆弱性を修正（HTMLエスケープ・URLサニタイズ追加）
- [x] JWT_SECRET環境変数を必須化（未設定時にサーバー起動を阻止）

### Phase 3.17（セキュリティ修正 — High）✅
- [x] WebSocketオリジン検証の追加（CSWSH対策、CORS設定と連動）
- [x] WebSocket認証をhttpOnly Cookieベースに変更（トークンURL露出防止）
- [x] WebSocketグループメッセージに認可チェックを追加（ルームメンバーシップ検証）
- [x] ファイルアップロードにマジックバイトMIMEタイプ検証を追加（Stored XSS防止）
- [x] Markdownレンダリングにrehype-sanitizeを導入（明示的HTMLサニタイズ）

### Phase 3.18（ウィークリーレポートメール）✅
- [x] 毎週月曜9:00にアクティビティサマリーメールを自動送信（cronスケジューラ）
- [x] HTMLメールテンプレート（ダークテーマ、レスポンシブ、10言語対応）
- [x] Settings画面からメール配信ON/OFF・配信言語の設定
- [x] SMTP設定による柔軟なメール配送（未設定時は機能自動無効化）
- [x] 開発用Mailhogコンテナ統合（Web UIでメール確認可能）
- [x] TDDで実装（6テストケース追加）
- [x] 全10言語のi18n対応

### Phase 3.19（学習レベルシステム）✅
- [x] XP（経験値）自動計算（学習ログ・投稿・GitHub・目標・コメント・いいね・ストリーク）
- [x] レベル計算式（累計XP = 100 × N × (N+1) / 2）
- [x] ダッシュボードLevelWidget（レベル・プログレスバー・次レベルまでのXP）
- [x] プロフィールLevelDisplay（XP内訳バーチャート付き）
- [x] LevelBadgeコンポーネント（レベル帯別カラー: gray/green/blue/purple/gold）
- [x] ランキングページに「レベル別」タブ追加
- [x] レベルアップ通知機能（CheckAndNotifyLevelUp）
- [x] TDDで実装（28テストケース追加）
- [x] 全10言語のi18n対応

### Phase 3.20（学習分析ダッシュボード）✅
- [x] ヒートマップ（曜日×時間帯の学習密度を色分け表示）
- [x] カテゴリ別学習時間（横棒グラフ＋比率表示）
- [x] 週間トレンドグラフ（棒グラフ＋前週比成長率）
- [x] 生産性スコア（ポモドーロ率30%＋目標達成率40%＋ストリーク安定性30%）
- [x] AIインサイト（ピーク時間帯・注力カテゴリ・週間成長率・ストリーク傾向の4分析）
- [x] TDDで実装（22テストケース追加）
- [x] 全10言語のi18n対応

### Phase 5.1（おすすめユーザー＆コンテンツレコメンド）✅
- [x] スキルマッチングベースのユーザーレコメンド（共通スキル数でスコアリング）
- [x] トレンド投稿（直近7日間のいいね＋コメント数で集計）
- [x] トレンド学習リソース（直近30日間のいいね＋保存数で集計）
- [x] ダッシュボードに「おすすめユーザー」ウィジェット追加（フォローボタン付き）
- [x] ダッシュボードに「トレンド」ウィジェット追加（投稿/リソースタブ切替）
- [x] TDDで実装（7テストケース追加）
- [x] 全10言語のi18n対応

### Phase 5.2（学習グループ — スタディサークル）✅
- [x] スタディサークル機能（3〜10人の小グループで同じテーマを学習）
- [x] 共有ロードマップ（ステップ作成・並替・メンバー別進捗マトリクス）
- [x] 日次チェックイン（1日1回の一言報告、重複防止）
- [x] グループ内ストリークランキング（連続チェックイン日数で競争）
- [x] メンバー管理（招待・除外・退出、上限チェック）
- [x] ダッシュボードウィジェット（参加中サークル概要表示）
- [x] TDDで実装（29テストケース追加）
- [x] 全10言語のi18n対応

### Phase 5.3（ローディングUI統一と品質向上）✅
- [x] 統一されたローディングコンポーネント（LoadingSpinner、PageLoader、InlineLoader）
- [x] サイズ・カラーバリエーション対応（sm/md/lg、primary/secondary/white）
- [x] スケルトン背景色を bg-gray-800 に統一
- [x] アクセシビリティ対応（role、aria-label）
- [x] 主要ページへの適用（ProfilePage、RankingsPage、PostDetailPage、GoalsPage）
- [x] TDDで実装（4テストスイート追加）
- [x] 全10言語のi18n対応

### Phase 5.4（UIのビジネス向けリデザイン）✅
- [x] ゲーミフィケーション要素の削減（派手な色・グラデーション・過剰なアニメーション削除）
- [x] Teams/Notion風のシンプルなデザインシステム導入
- [x] カラーパレット統一（グレー基調、アクセントは blue-600/green-500 のみ）
- [x] rounded-xl/2xl → rounded-md（控えめな丸み）
- [x] shadow-lg/xl/2xl → shadow-sm（控えめな影）
- [x] hover:scale-105 削除、transition-colors のみ
- [x] 主要コンポーネント更新（LevelWidget、BadgeDisplay、StreakWidget、PomodoroTimer）
- [x] 59ファイル、139箇所の一括スタイル変更

### Phase 5.5（レイアウトと視認性の改善）✅
- [x] グローバルスタイルで視認性向上のユーティリティ追加
- [x] セクション見出しの強調（section-heading クラス）
- [x] カード境界の明確化（card-enhanced クラス）
- [x] ボタンスタイルの統一（btn-primary, btn-secondary）
- [x] フォーカス状態の改善（全インタラクティブ要素）
- [x] ダッシュボードのレイアウト改善（セクションごとにグループ化）
- [x] 一貫性のあるスペーシング（space-y-6）
- [x] 全10言語のi18n対応

### Phase 6.1（統合検索機能）✅
- [x] タブ切り替えUI（ユーザー/投稿/スタディサークル）
- [x] 検索APIエンドポイント（`/api/v1/search/posts`, `/api/v1/search/circles`）
- [x] リポジトリ層の部分一致検索（タイトル・本文・名前・トピック・説明）
- [x] SearchBar / SearchTabs コンポーネント
- [x] usePostSearch / useCircleSearch フック
- [x] 全10言語のi18n対応

### Phase 6.3（投稿の下書き・公開機能）✅
- [x] Postモデルに`is_draft`フィールドを追加
- [x] 下書き専用APIエンドポイント（`GET /api/v1/posts/drafts`, `PUT /api/v1/posts/:id/publish`）
- [x] 下書き作成時は通知をスキップ、公開時にフォロワーへ通知
- [x] 下書き一覧ページ（DraftsPage）の実装
- [x] PostFormに「下書き保存」ボタンを追加
- [x] 全10言語のi18n対応（下書き関連キー）

### Phase 4（将来）📋
- [ ] モバイルアプリ
- [ ] チーム機能
- [ ] バッジ・実績システム拡張

---

## 🤝 Contributing

コントリビューションは大歓迎です！

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'feat: add amazing feature'`)
4. ブランチをプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

詳細は [CONTRIBUTING.md](CONTRIBUTING.md) をご覧ください。

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**カワノ**

- GitHub: [@norman6464](https://github.com/norman6464)
- Zenn: [@norman6464](https://zenn.dev/norman6464)

---

## 🙏 Acknowledgments

- [GitHub API](https://docs.github.com/en/rest)
- [Zenn API](https://zenn.dev/api)
- [shadcn/ui](https://ui.shadcn.com/)
- [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts)

---
