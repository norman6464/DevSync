package usecase

import "github.com/norman6464/devsync/backend/internal/model"

// roadmapTemplateStep はプリセットのロードマップテンプレートに含まれるステップ。
type roadmapTemplateStep struct {
	Title       string
	Description string
	ResourceURL string
}

// roadmapTemplate は初期登録するプリセットのロードマップテンプレート。
type roadmapTemplate struct {
	Title       string
	Description string
	Category    model.RoadmapCategory
	Steps       []roadmapTemplateStep
}

// presetRoadmapTemplates は起動時に初期登録するテンプレートの定義。
// 内容はシードデータなので usecase 本体とは分けている。
var presetRoadmapTemplates = []roadmapTemplate{
	{
		Title:       "Webフロントエンドエンジニア",
		Description: "HTML/CSSの基礎からReact/Next.jsを使ったモダンなWeb開発まで、フロントエンドエンジニアに必要なスキルを段階的に学ぶロードマップです。",
		Category:    model.RoadmapCategorySkill,
		Steps: []roadmapTemplateStep{
			{Title: "HTML/CSSの基礎", Description: "Webページの構造とスタイリングの基本を学ぶ", ResourceURL: "https://developer.mozilla.org/ja/docs/Learn/HTML"},
			{Title: "JavaScriptの基礎", Description: "変数、関数、DOM操作などJSの基本を習得", ResourceURL: "https://developer.mozilla.org/ja/docs/Learn/JavaScript"},
			{Title: "Git/GitHubの使い方", Description: "バージョン管理とチーム開発の基礎", ResourceURL: "https://docs.github.com/ja"},
			{Title: "TypeScript入門", Description: "型安全な開発の基礎を学ぶ", ResourceURL: "https://www.typescriptlang.org/docs/"},
			{Title: "React入門", Description: "コンポーネント、状態管理、Hooksの基本", ResourceURL: "https://ja.react.dev/learn"},
			{Title: "Tailwind CSS", Description: "ユーティリティファーストCSSフレームワーク", ResourceURL: "https://tailwindcss.com/docs"},
			{Title: "Next.js入門", Description: "SSR/SSG、ルーティング、API Routes", ResourceURL: "https://nextjs.org/docs"},
			{Title: "テスト（Jest/Vitest）", Description: "フロントエンドのユニットテスト・E2Eテスト", ResourceURL: "https://vitest.dev/guide/"},
			{Title: "ポートフォリオ作成", Description: "学んだ技術を使ってオリジナルアプリを作成"},
		},
	},
	{
		Title:       "バックエンドエンジニア（Go）",
		Description: "Go言語の基礎からAPI開発、データベース設計、Dockerまで、バックエンドエンジニアに必要なスキルを学ぶロードマップです。",
		Category:    model.RoadmapCategoryLanguage,
		Steps: []roadmapTemplateStep{
			{Title: "Go言語の基礎", Description: "変数、関数、構造体、インターフェースの基本", ResourceURL: "https://go.dev/tour/welcome/1"},
			{Title: "Git/GitHubの使い方", Description: "バージョン管理とチーム開発の基礎", ResourceURL: "https://docs.github.com/ja"},
			{Title: "SQLとデータベース設計", Description: "RDBMS、テーブル設計、クエリ最適化", ResourceURL: "https://www.postgresql.org/docs/"},
			{Title: "REST API設計", Description: "HTTPメソッド、ステータスコード、RESTful設計原則", ResourceURL: ""},
			{Title: "Ginフレームワーク", Description: "Go製Webフレームワークでの開発", ResourceURL: "https://gin-gonic.com/docs/"},
			{Title: "GORM（ORM）", Description: "GoのORMを使ったデータベース操作", ResourceURL: "https://gorm.io/docs/"},
			{Title: "認証・認可（JWT）", Description: "ユーザー認証とアクセス制御の実装", ResourceURL: ""},
			{Title: "Docker入門", Description: "コンテナ技術の基礎と開発環境構築", ResourceURL: "https://docs.docker.com/get-started/"},
			{Title: "テスト（Go testing）", Description: "ユニットテスト、テーブルドリブンテスト", ResourceURL: "https://go.dev/doc/tutorial/add-a-test"},
			{Title: "APIプロジェクト作成", Description: "学んだ技術を統合してAPIサーバーを構築"},
		},
	},
	{
		Title:       "フルスタックエンジニア",
		Description: "フロントエンドからバックエンド、インフラまで幅広くカバーする、フルスタックエンジニア向けの総合ロードマップです。",
		Category:    model.RoadmapCategorySkill,
		Steps: []roadmapTemplateStep{
			{Title: "HTML/CSS/JavaScriptの基礎", Description: "Web開発の土台となる3言語を習得"},
			{Title: "Git/GitHubの使い方", Description: "バージョン管理の基礎"},
			{Title: "TypeScript + React", Description: "モダンフロントエンド開発"},
			{Title: "バックエンド言語（Go or Node.js）", Description: "サーバーサイドの基礎"},
			{Title: "データベース（PostgreSQL）", Description: "RDBMSの設計と運用"},
			{Title: "REST API開発", Description: "フロントとバックエンドの連携"},
			{Title: "認証・セキュリティ", Description: "JWT、CORS、XSS/CSRF対策"},
			{Title: "Docker & Docker Compose", Description: "開発・本番環境のコンテナ化"},
			{Title: "CI/CD（GitHub Actions）", Description: "テスト・デプロイの自動化"},
			{Title: "クラウドデプロイ（AWS/Vercel）", Description: "本番環境へのデプロイ"},
			{Title: "フルスタックプロジェクト", Description: "全技術を活かしたオリジナルアプリ開発"},
		},
	},
	{
		Title:       "インフラ・DevOpsエンジニア",
		Description: "Linux、ネットワーク、Docker、Kubernetes、CI/CDなど、インフラ・DevOpsエンジニアに必要なスキルを学ぶロードマップです。",
		Category:    model.RoadmapCategorySkill,
		Steps: []roadmapTemplateStep{
			{Title: "Linux基礎", Description: "コマンドライン操作、ファイルシステム、権限管理"},
			{Title: "ネットワーク基礎", Description: "TCP/IP、DNS、HTTP/HTTPS、ファイアウォール"},
			{Title: "シェルスクリプト", Description: "Bashスクリプトによる自動化"},
			{Title: "Git/GitHub", Description: "バージョン管理とGitOps"},
			{Title: "Docker", Description: "コンテナ技術の基礎から実践", ResourceURL: "https://docs.docker.com/"},
			{Title: "CI/CD（GitHub Actions）", Description: "テスト・ビルド・デプロイの自動化"},
			{Title: "AWS基礎（EC2/RDS/S3）", Description: "クラウドインフラの基本サービス", ResourceURL: "https://aws.amazon.com/jp/getting-started/"},
			{Title: "Terraform", Description: "Infrastructure as Codeの実践", ResourceURL: "https://developer.hashicorp.com/terraform/tutorials"},
			{Title: "Kubernetes入門", Description: "コンテナオーケストレーションの基礎", ResourceURL: "https://kubernetes.io/ja/docs/tutorials/"},
			{Title: "監視・ログ管理", Description: "Prometheus、Grafana、CloudWatch"},
		},
	},
	{
		Title:       "モバイルアプリ開発（React Native）",
		Description: "JavaScriptの基礎からReact Nativeを使ったクロスプラットフォームモバイルアプリ開発までを学ぶロードマップです。",
		Category:    model.RoadmapCategoryFramework,
		Steps: []roadmapTemplateStep{
			{Title: "JavaScript/TypeScriptの基礎", Description: "モバイル開発の土台となる言語"},
			{Title: "React基礎", Description: "コンポーネント指向UIの基本"},
			{Title: "React Native入門", Description: "環境構築とHello World", ResourceURL: "https://reactnative.dev/docs/getting-started"},
			{Title: "ナビゲーション", Description: "React Navigationによる画面遷移"},
			{Title: "状態管理", Description: "Context API / Zustand / Redux"},
			{Title: "API連携", Description: "REST API / GraphQLとの通信"},
			{Title: "ネイティブ機能", Description: "カメラ、位置情報、プッシュ通知"},
			{Title: "アプリストア公開", Description: "App Store / Google Playへの公開手順"},
		},
	},
}
