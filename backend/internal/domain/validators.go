// Package domain provides domain-level validation logic for the DevSync application.
package domain

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Email validation constants
const (
	MinEmailLength = 3
	MaxEmailLength = 255
)

// Password validation constants
const (
	MinPasswordLength = 6
	MaxPasswordLength = 128
)

// Username validation constants
const (
	MinUsernameLength = 2
	MaxUsernameLength = 30
)

// Content validation constants
const (
	MinTitleLength   = 1
	MaxTitleLength   = 200
	MinContentLength = 1
	MaxContentLength = 10000
)

var (
	// emailRegex は基本的なメールアドレスの形式をチェックする正規表現
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// usernameRegex はユーザー名の形式をチェックする正規表現（英数字・アンダースコア・ハイフンのみ）
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
)

// ValidateEmail はメールアドレスのバリデーションを行う
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if len(email) < MinEmailLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("メールアドレスは%d文字以上である必要があります", MinEmailLength), nil)
	}

	if len(email) > MaxEmailLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("メールアドレスは%d文字以下である必要があります", MaxEmailLength), nil)
	}

	if !emailRegex.MatchString(email) {
		return NewError(ErrCodeValidation, "有効なメールアドレスを入力してください", nil)
	}

	return nil
}

// ValidatePassword はパスワードのバリデーションを行う
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("パスワードは%d文字以上である必要があります", MinPasswordLength), nil)
	}

	if len(password) > MaxPasswordLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("パスワードは%d文字以下である必要があります", MaxPasswordLength), nil)
	}

	// 少なくとも1つの数字または記号を含むことを推奨（任意）
	hasNumberOrSymbol := false
	for _, char := range password {
		if unicode.IsDigit(char) || unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasNumberOrSymbol = true
			break
		}
	}

	if !hasNumberOrSymbol {
		return NewError(ErrCodeValidation, "パスワードは数字または記号を少なくとも1つ含める必要があります", nil)
	}

	return nil
}

// ValidateUsername はユーザー名のバリデーションを行う
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)

	if len(username) < MinUsernameLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("ユーザー名は%d文字以上である必要があります", MinUsernameLength), nil)
	}

	if len(username) > MaxUsernameLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("ユーザー名は%d文字以下である必要があります", MaxUsernameLength), nil)
	}

	if !usernameRegex.MatchString(username) {
		return NewError(ErrCodeValidation, "ユーザー名は英数字、アンダースコア、ハイフンのみ使用できます", nil)
	}

	return nil
}

// ValidateTitle はタイトルのバリデーションを行う
func ValidateTitle(title string) error {
	title = strings.TrimSpace(title)

	if len(title) < MinTitleLength {
		return NewError(ErrCodeValidation, "タイトルを入力してください", nil)
	}

	if len(title) > MaxTitleLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("タイトルは%d文字以下である必要があります", MaxTitleLength), nil)
	}

	return nil
}

// ValidateContent はコンテンツのバリデーションを行う
func ValidateContent(content string) error {
	content = strings.TrimSpace(content)

	if len(content) < MinContentLength {
		return NewError(ErrCodeValidation, "内容を入力してください", nil)
	}

	if len(content) > MaxContentLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("内容は%d文字以下である必要があります", MaxContentLength), nil)
	}

	return nil
}

// ValidateURL はURLの形式をチェックする
func ValidateURL(url string) error {
	if url == "" {
		return nil // 空の場合は許可（オプショナル）
	}

	url = strings.TrimSpace(url)

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return NewError(ErrCodeValidation, "URLはhttp://またはhttps://で始まる必要があります", nil)
	}

	if len(url) > 2048 {
		return NewError(ErrCodeValidation, "URLが長すぎます", nil)
	}

	return nil
}

// ValidateTags はタグのバリデーションを行う
func ValidateTags(tags []string) error {
	if len(tags) > 10 {
		return NewError(ErrCodeValidation, "タグは10個まで設定できます", nil)
	}

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if len(tag) == 0 {
			return NewError(ErrCodeValidation, "空のタグは設定できません", nil)
		}
		if len(tag) > 30 {
			return NewError(ErrCodeValidation, "タグは30文字以下である必要があります", nil)
		}
	}

	return nil
}

// ValidatePagination はページネーションパラメータのバリデーションを行い、
// 正規化された値を返す。limitは最大100に制限される。
func ValidatePagination(limit, offset int) (int, int, error) {
	// offsetは0以上
	if offset < 0 {
		return 0, 0, NewError(ErrCodeValidation, "オフセットは0以上である必要があります", nil)
	}

	// limitが指定されていない場合はデフォルト値（10）を使用
	if limit <= 0 {
		limit = 10
	}

	// limitは最大100
	if limit > 100 {
		limit = 100
	}

	return limit, offset, nil
}

// ValidateStringLength は文字列の長さをバリデーションする汎用関数
func ValidateStringLength(s string, min, max int, fieldName string) error {
	s = strings.TrimSpace(s)
	length := len(s)

	if length < min {
		if min == 1 {
			return NewError(ErrCodeValidation, fmt.Sprintf("%sを入力してください", fieldName), nil)
		}
		return NewError(ErrCodeValidation, fmt.Sprintf("%sは%d文字以上である必要があります", fieldName, min), nil)
	}

	if max > 0 && length > max {
		return NewError(ErrCodeValidation, fmt.Sprintf("%sは%d文字以下である必要があります", fieldName, max), nil)
	}

	return nil
}

// ValidateEnum は値が許可されたリストに含まれるかチェックする
func ValidateEnum(value string, allowedValues []string, fieldName string) error {
	if value == "" {
		return nil // 空の場合は許可（オプショナル）
	}

	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}

	return NewError(ErrCodeValidation, fmt.Sprintf("%sの値が不正です", fieldName), nil)
}
