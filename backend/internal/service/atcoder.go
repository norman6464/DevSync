// Package service はAtCoderレーティング取得のビジネスロジックを提供する。
package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AtCoderRatingEntry はAtCoderのレーティング履歴エントリを表す。
type AtCoderRatingEntry struct {
	IsRated     bool   `json:"IsRated"`
	Place       int    `json:"Place"`
	OldRating   int    `json:"OldRating"`
	NewRating   int    `json:"NewRating"`
	Performance int    `json:"Performance"`
	ContestName string `json:"ContestName"`
	EndTime     string `json:"EndTime"`
}

// AtCoderRatingInfo はAtCoderのレーティング情報を表す。
type AtCoderRatingInfo struct {
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Color    string `json:"color"`
	Rank     string `json:"rank"`
}

// AtCoderService はAtCoderレーティング取得を提供するサービス。
type AtCoderService struct {
	client *http.Client
}

// NewAtCoderService は新しいAtCoderServiceインスタンスを生成する。
func NewAtCoderService() *AtCoderService {
	return &AtCoderService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetRating は指定ユーザーのAtCoderレーティング情報を取得する。
func (s *AtCoderService) GetRating(username string) (*AtCoderRatingInfo, error) {
	url := fmt.Sprintf("https://atcoder.jp/users/%s/history/json", username)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("AtCoder APIリクエストに失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("AtCoderユーザーが見つかりません: %s", username)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AtCoder APIエラー: ステータスコード %d", resp.StatusCode)
	}

	var history []AtCoderRatingEntry
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, fmt.Errorf("AtCoder APIレスポンスのパースに失敗: %w", err)
	}

	info := &AtCoderRatingInfo{
		Username: username,
		Rating:   0,
		Color:    "gray",
		Rank:     "灰",
	}

	if len(history) > 0 {
		info.Rating = history[len(history)-1].NewRating
		info.Color = ratingToColor(info.Rating)
		info.Rank = ratingToRank(info.Rating)
	}

	return info, nil
}

// ValidateUsername はAtCoderユーザー名が有効かどうか検証する。
func (s *AtCoderService) ValidateUsername(username string) bool {
	url := fmt.Sprintf("https://atcoder.jp/users/%s/history/json", username)
	resp, err := s.client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// ratingToColor はレーティング値をAtCoderの色名に変換する。
func ratingToColor(rating int) string {
	switch {
	case rating >= 2800:
		return "red"
	case rating >= 2400:
		return "orange"
	case rating >= 2000:
		return "yellow"
	case rating >= 1600:
		return "blue"
	case rating >= 1200:
		return "cyan"
	case rating >= 800:
		return "green"
	case rating >= 400:
		return "brown"
	default:
		return "gray"
	}
}

// ratingToRank はレーティング値をAtCoderのランク名に変換する。
func ratingToRank(rating int) string {
	switch {
	case rating >= 2800:
		return "赤"
	case rating >= 2400:
		return "橙"
	case rating >= 2000:
		return "黄"
	case rating >= 1600:
		return "青"
	case rating >= 1200:
		return "水色"
	case rating >= 800:
		return "緑"
	case rating >= 400:
		return "茶"
	default:
		return "灰"
	}
}
