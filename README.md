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

### 🤖 AI機能
- 投稿時のAI補完
- スキルに基づく学習アドバイス
- おすすめユーザーのマッチング

### 📝 外部サービス連携
- GitHub
- Zenn
- Qiita

---

## 🎬 Demo

### ホーム画面
<p align="center">
  <img src="docs/images/screenshots/home.png" alt="Home" width="800">
</p>

### プロフィール画面
<p align="center">
  <img src="docs/images/screenshots/profile.png" alt="Profile" width="800">
</p>

### ランキング画面
<p align="center">
  <img src="docs/images/screenshots/ranking.png" alt="Ranking" width="800">
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
# Backend
make run-user        # User Serviceを起動
make run-github      # GitHub Serviceを起動
make test            # テスト実行
make lint            # Linter実行
make proto           # Protocol Buffers生成

# Frontend
npm run dev          # 開発サーバー起動
npm run build        # ビルド
npm run test         # テスト実行
npm run lint         # Linter実行

# Infrastructure
make tf-plan         # Terraform plan
make tf-apply        # Terraform apply
make k8s-deploy      # Kubernetesデプロイ
```

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

### Phase 2（拡張）🚧
- [ ] Zenn/Qiita連携
- [ ] 投稿機能
- [ ] チャット機能
- [ ] AI補完

### Phase 3（将来）📋
- [ ] モバイルアプリ
- [ ] チーム機能
- [ ] バッジ・実績システム
- [ ] 週間レポート自動生成

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
