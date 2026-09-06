package model

import (
	"time"
)

// ReviewStatus は書籍レビューの読書状態を表す型。
type ReviewStatus string

const (
	ReviewStatusNotStarted ReviewStatus = "not_started" // 未読
	ReviewStatusReading    ReviewStatus = "reading"     // 読中
	ReviewStatusCompleted  ReviewStatus = "completed"   // 読了
)

// ValidReviewStatuses は有効な読書ステータスのマップ。
var ValidReviewStatuses = map[ReviewStatus]bool{
	ReviewStatusNotStarted: true,
	ReviewStatusReading:    true,
	ReviewStatusCompleted:  true,
}

// BookReview はユーザーが投稿した書籍レビューを表す。
// Rating は1〜5の整数で評価を表し、Review にレビュー本文を格納する。
type BookReview struct {
	ID          uint         `json:"id"`
	UserID      uint         `json:"user_id"`
	User        User         `json:"user,omitempty"`
	Title       string       `json:"title"`        // 書籍タイトル
	Author      string       `json:"author"`       // 書籍の著者名
	ISBN        string       `json:"isbn"`         // ISBNコード
	Rating      int          `json:"rating"`       // 評価（1〜5の整数）
	Review      string       `json:"review"`       // レビュー本文
	TotalPages  int          `json:"total_pages"`  // 総ページ数
	CurrentPage int          `json:"current_page"` // 現在の読書ページ
	ImageURL    string       `json:"image_url"`    // 書籍カバー画像URL
	Status      ReviewStatus `json:"status"`       // 読書状態
	IsArchived  bool         `json:"is_archived"`  // アーカイブ済みフラグ
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
