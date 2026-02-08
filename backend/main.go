package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/router"
	"github.com/norman6464/devsync/backend/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.GitHubContribution{},
		&model.GitHubLanguageStat{},
		&model.GitHubRepository{},
		&model.Post{},
		&model.Like{},
		&model.Comment{},
		&model.Message{},
		&model.Notification{},
		&model.PasswordResetToken{},
		&model.ZennArticle{},
		&model.QiitaArticle{},
		&model.LearningGoal{},
		&model.Project{},
		&model.LearningResource{},
		&model.ResourceLike{},
		&model.ResourceSave{},
		&model.BookReview{},
		&model.Question{},
		&model.QuestionVote{},
		&model.Answer{},
		&model.AnswerVote{},
		&model.Roadmap{},
		&model.RoadmapStep{},
		&model.ChatRoom{},
		&model.ChatRoomMember{},
		&model.GroupMessage{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Mark existing users as onboarding completed
	db.Model(&model.User{}).Where("onboarding_completed = ?", false).Update("onboarding_completed", true)

	// Start WebSocket hub
	hub := service.NewHub()
	go hub.Run()

	r := router.Setup(db, cfg, hub)

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
