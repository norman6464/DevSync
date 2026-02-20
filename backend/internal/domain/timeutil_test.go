package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartOfMonth(t *testing.T) {
	// 2024年3月15日のケース
	input := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
	result := StartOfMonth(input)

	assert.Equal(t, 2024, result.Year())
	assert.Equal(t, time.March, result.Month())
	assert.Equal(t, 1, result.Day())
	assert.Equal(t, 0, result.Hour())
	assert.Equal(t, 0, result.Minute())
	assert.Equal(t, 0, result.Second())
}

func TestStartOfMonth_FirstDay(t *testing.T) {
	// 月初の場合もそのまま月初を返す
	input := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result := StartOfMonth(input)

	assert.Equal(t, 2024, result.Year())
	assert.Equal(t, time.January, result.Month())
	assert.Equal(t, 1, result.Day())
}

func TestStartOfMonth_LastDay(t *testing.T) {
	// 月末の場合
	input := time.Date(2024, 2, 29, 23, 59, 59, 0, time.UTC)
	result := StartOfMonth(input)

	assert.Equal(t, 2024, result.Year())
	assert.Equal(t, time.February, result.Month())
	assert.Equal(t, 1, result.Day())
}

func TestStartOfWeek_Wednesday(t *testing.T) {
	// 水曜日の場合 → 月曜日を返す
	wednesday := time.Date(2024, 3, 13, 14, 30, 0, 0, time.UTC) // 水曜
	result := StartOfWeek(wednesday)

	assert.Equal(t, 2024, result.Year())
	assert.Equal(t, time.March, result.Month())
	assert.Equal(t, 11, result.Day()) // 月曜
	assert.Equal(t, 0, result.Hour())
}

func TestStartOfWeek_Monday(t *testing.T) {
	// 月曜日の場合 → そのまま月曜を返す
	monday := time.Date(2024, 3, 11, 10, 0, 0, 0, time.UTC) // 月曜
	result := StartOfWeek(monday)

	assert.Equal(t, 11, result.Day())
}

func TestStartOfWeek_Sunday(t *testing.T) {
	// 日曜日の場合 → 前の月曜を返す
	sunday := time.Date(2024, 3, 17, 10, 0, 0, 0, time.UTC) // 日曜
	result := StartOfWeek(sunday)

	assert.Equal(t, 11, result.Day()) // 前の月曜
}

func TestDaysAgo(t *testing.T) {
	base := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	result := DaysAgo(base, 7)

	assert.Equal(t, 2024, result.Year())
	assert.Equal(t, time.March, result.Month())
	assert.Equal(t, 8, result.Day())
}

func TestDaysAgo_Zero(t *testing.T) {
	base := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	result := DaysAgo(base, 0)

	assert.Equal(t, 15, result.Day())
}
