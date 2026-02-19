package service

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)


// RoadmapService は学習ロードマップのビジネスロジックを提供する。
// ロードマップとステップのCRUD操作、可視性制御、コピー機能を担当する。
type RoadmapService struct {
	repo repository.RoadmapRepositoryInterface
}

// NewRoadmapService は新しいRoadmapServiceインスタンスを生成する。
func NewRoadmapService(repo repository.RoadmapRepositoryInterface) *RoadmapService {
	return &RoadmapService{repo: repo}
}

// Create は新しいロードマップを作成する。
func (s *RoadmapService) Create(roadmap *model.Roadmap) error {
	return s.repo.Create(roadmap)
}

// GetByID は指定IDのロードマップを可視性チェック付きで取得する。
// 非公開ロードマップはオーナー以外アクセスできない。
func (s *RoadmapService) GetByID(id, userID uint) (*model.Roadmap, error) {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID && !roadmap.IsPublic {
		return nil, ErrForbidden
	}
	return roadmap, nil
}

// GetByUserID は指定ユーザーの全ロードマップを取得する。
func (s *RoadmapService) GetByUserID(userID uint) ([]model.Roadmap, error) {
	return s.repo.GetByUserID(userID)
}

// GetPublicRoadmaps は公開ロードマップをページネーション付きで取得する。
func (s *RoadmapService) GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error) {
	return s.repo.GetPublicRoadmaps(limit, offset)
}

// GetStats は指定ユーザーのロードマップ統計情報を取得する。
func (s *RoadmapService) GetStats(userID uint) (*model.RoadmapStats, error) {
	return s.repo.GetStats(userID)
}

// findAndCheckOwnership は指定IDのロードマップを取得し、所有権を検証する。
func (s *RoadmapService) findAndCheckOwnership(id, userID uint) (*model.Roadmap, error) {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID {
		return nil, ErrForbidden
	}
	return roadmap, nil
}

// findAndCheckStepOwnership はロードマップの所有権とステップの所属を検証する。
func (s *RoadmapService) findAndCheckStepOwnership(roadmapID, stepID, userID uint) (*model.RoadmapStep, error) {
	if _, err := s.findAndCheckOwnership(roadmapID, userID); err != nil {
		return nil, err
	}
	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return nil, err
	}
	if step.RoadmapID != roadmapID {
		return nil, ErrBadRequest
	}
	return step, nil
}

// Update は所有権を検証した後、ロードマップを更新する。
// ステータスが「完了」に変更された場合、完了日時を自動設定する。
func (s *RoadmapService) Update(id, userID uint, updates *model.Roadmap) (*model.Roadmap, error) {
	roadmap, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if updates.Title != "" {
		roadmap.Title = updates.Title
	}
	if updates.Description != "" {
		roadmap.Description = updates.Description
	}
	if updates.Category != "" {
		roadmap.Category = updates.Category
	}
	if updates.Status != "" {
		roadmap.Status = updates.Status
		if roadmap.Status == model.RoadmapStatusCompleted && roadmap.CompletedAt == nil {
			now := time.Now()
			roadmap.CompletedAt = &now
		}
	}

	if err := s.repo.Update(roadmap); err != nil {
		return nil, err
	}
	return roadmap, nil
}

// UpdateVisibility は所有権を検証した後、ロードマップの公開/非公開状態を更新する。
func (s *RoadmapService) UpdateVisibility(id, userID uint, isPublic bool) (*model.Roadmap, error) {
	roadmap, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	roadmap.IsPublic = isPublic

	if err := s.repo.Update(roadmap); err != nil {
		return nil, err
	}
	return roadmap, nil
}

// Delete は所有権を検証した後、ロードマップを削除する。
func (s *RoadmapService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// CopyRoadmap は公開ロードマップをテンプレートとしてコピーする。
// 非公開かつ自分のものでないロードマップはコピーできない。
func (s *RoadmapService) CopyRoadmap(roadmapID, userID uint) (*model.Roadmap, error) {
	original, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return nil, err
	}
	if !original.IsPublic && original.UserID != userID {
		return nil, ErrForbidden
	}
	return s.repo.CopyRoadmap(roadmapID, userID)
}

// CreateStep は所有権を検証した後、ロードマップにステップを追加する。
func (s *RoadmapService) CreateStep(roadmapID, userID uint, step *model.RoadmapStep) error {
	if _, err := s.findAndCheckOwnership(roadmapID, userID); err != nil {
		return err
	}

	step.RoadmapID = roadmapID
	return s.repo.CreateStep(step)
}

// UpdateStep はロードマップの所有権とステップの所属を検証した後、ステップを更新する。
func (s *RoadmapService) UpdateStep(roadmapID, stepID, userID uint, updates *model.RoadmapStep) (*model.RoadmapStep, error) {
	step, err := s.findAndCheckStepOwnership(roadmapID, stepID, userID)
	if err != nil {
		return nil, err
	}

	if updates.Title != "" {
		step.Title = updates.Title
	}
	if updates.Description != "" {
		step.Description = updates.Description
	}
	if updates.ResourceURL != "" {
		step.ResourceURL = updates.ResourceURL
	}

	if err := s.repo.UpdateStep(step); err != nil {
		return nil, err
	}
	return step, nil
}

// UpdateStepCompletion はステップの完了状態を更新する。
// 完了時にはCompletedAtを設定し、未完了に戻す場合はnilにリセットする。
func (s *RoadmapService) UpdateStepCompletion(roadmapID, stepID, userID uint, isCompleted bool) (*model.RoadmapStep, error) {
	step, err := s.findAndCheckStepOwnership(roadmapID, stepID, userID)
	if err != nil {
		return nil, err
	}

	step.IsCompleted = isCompleted
	if isCompleted && step.CompletedAt == nil {
		now := time.Now()
		step.CompletedAt = &now
	} else if !isCompleted {
		step.CompletedAt = nil
	}

	if err := s.repo.UpdateStep(step); err != nil {
		return nil, err
	}
	return step, nil
}

// DeleteStep はロードマップの所有権とステップの所属を検証した後、ステップを削除する。
func (s *RoadmapService) DeleteStep(roadmapID, stepID, userID uint) error {
	if _, err := s.findAndCheckStepOwnership(roadmapID, stepID, userID); err != nil {
		return err
	}
	return s.repo.DeleteStep(stepID)
}

// GetTemplates はテンプレートロードマップの一覧を取得する。
func (s *RoadmapService) GetTemplates() ([]model.Roadmap, error) {
	return s.repo.GetTemplates()
}

// CreateFromTemplate はテンプレートからユーザー用ロードマップを作成する。
// テンプレートでないロードマップからの作成はErrBadRequestを返す。
func (s *RoadmapService) CreateFromTemplate(templateID, userID uint) (*model.Roadmap, error) {
	template, err := s.repo.FindByID(templateID)
	if err != nil {
		return nil, err
	}
	if !template.IsTemplate {
		return nil, ErrBadRequest
	}
	return s.repo.CopyRoadmap(templateID, userID)
}

// SeedTemplates はプリセットテンプレートを初期登録する。
// 既にテンプレートが存在する場合はスキップする。
// userID にはシステムユーザーのIDを指定する（外部キー制約を満たすため）。
func (s *RoadmapService) SeedTemplates(userID uint) error {
	existing, err := s.repo.GetTemplates()
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	templates := []struct {
		Title       string
		Description string
		Category    model.RoadmapCategory
		Steps       []struct {
			Title       string
			Description string
			ResourceURL string
		}
	}{
		{
			Title:       "Webフロントエンドエンジニア",
			Description: "HTML/CSSの基礎からReact/Next.jsを使ったモダンなWeb開発まで、フロントエンドエンジニアに必要なスキルを段階的に学ぶロードマップです。",
			Category:    model.RoadmapCategorySkill,
			Steps: []struct {
				Title       string
				Description string
				ResourceURL string
			}{
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
			Steps: []struct {
				Title       string
				Description string
				ResourceURL string
			}{
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
			Steps: []struct {
				Title       string
				Description string
				ResourceURL string
			}{
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
			Steps: []struct {
				Title       string
				Description string
				ResourceURL string
			}{
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
			Steps: []struct {
				Title       string
				Description string
				ResourceURL string
			}{
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

	for _, tmpl := range templates {
		roadmap := &model.Roadmap{
			UserID:      userID,
			Title:       tmpl.Title,
			Description: tmpl.Description,
			Category:    tmpl.Category,
			IsPublic:    true,
			IsTemplate:  true,
			StepCount:   len(tmpl.Steps),
			Status:      model.RoadmapStatusActive,
		}
		if err := s.repo.Create(roadmap); err != nil {
			return err
		}
		for i, stepDef := range tmpl.Steps {
			stepModel := &model.RoadmapStep{
				RoadmapID:   roadmap.ID,
				Title:       stepDef.Title,
				Description: stepDef.Description,
				OrderIndex:  i,
				ResourceURL: stepDef.ResourceURL,
			}
			if err := s.repo.CreateStep(stepModel); err != nil {
				return err
			}
		}
	}
	return nil
}

// ReorderSteps は所有権を検証した後、ステップの表示順序を一括更新する。
func (s *RoadmapService) ReorderSteps(roadmapID, userID uint, orders []model.StepOrder) error {
	if _, err := s.findAndCheckOwnership(roadmapID, userID); err != nil {
		return err
	}
	return s.repo.ReorderSteps(roadmapID, orders)
}
