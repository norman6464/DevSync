// Package domain provides domain-level validation logic for the DevSync application.
package domain

import (
	"fmt"
	"net"
	"net/url"
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
	MinPasswordLength = 8
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

	if len([]rune(title)) < MinTitleLength {
		return NewError(ErrCodeValidation, "タイトルを入力してください", nil)
	}

	if len([]rune(title)) > MaxTitleLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("タイトルは%d文字以下である必要があります", MaxTitleLength), nil)
	}

	return nil
}

// ValidateContent はコンテンツのバリデーションを行う
func ValidateContent(content string) error {
	content = strings.TrimSpace(content)

	if len([]rune(content)) < MinContentLength {
		return NewError(ErrCodeValidation, "内容を入力してください", nil)
	}

	if len([]rune(content)) > MaxContentLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("内容は%d文字以下である必要があります", MaxContentLength), nil)
	}

	return nil
}

// ValidateURL はURLの形式をチェックする
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return nil // 空の場合は許可（オプショナル）
	}

	rawURL = strings.TrimSpace(rawURL)

	if len(rawURL) > 2048 {
		return NewError(ErrCodeValidation, "URLが長すぎます", nil)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return NewError(ErrCodeValidation, "URLの形式が不正です", nil)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return NewError(ErrCodeValidation, "URLはhttp://またはhttps://で始まる必要があります", nil)
	}

	if parsed.Host == "" {
		return NewError(ErrCodeValidation, "URLにホスト名が必要です", nil)
	}

	// SSRF対策: ローカルホスト・プライベートIPを拒否
	hostname := parsed.Hostname()
	if isBlockedHost(hostname) {
		return NewError(ErrCodeValidation, "内部ネットワークのURLは許可されていません", nil)
	}

	return nil
}

// cloudMetadataHosts はクラウドメタデータサービスのホスト名一覧。
var cloudMetadataHosts = map[string]bool{
	"metadata.google.internal":        true,
	"metadata.google.internal.":       true,
	"instance-data":                   true,
	"169.254.169.254":                 true,
	"fd00:ec2::254":                   true,
}

// isBlockedHost はSSRF対策として内部ネットワークのホストを検出する。
// IPアドレスの直接指定に加え、DNS解決後のIPもチェックする。
func isBlockedHost(hostname string) bool {
	lower := strings.ToLower(hostname)
	if lower == "localhost" || lower == "0.0.0.0" || lower == "[::1]" {
		return true
	}

	// クラウドメタデータサービスのホスト名をブロック
	if cloudMetadataHosts[lower] {
		return true
	}

	// 直接IPアドレスが指定された場合
	ip := net.ParseIP(hostname)
	if ip != nil {
		return isBlockedIP(ip)
	}

	// ドメイン名の場合はDNS解決してIPをチェック
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return false
	}
	for _, resolvedIP := range ips {
		if isBlockedIP(resolvedIP) {
			return true
		}
	}
	return false
}

// isBlockedIP は指定IPが内部ネットワークに属するかチェックする。
func isBlockedIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified()
}

// ValidateTags はタグのバリデーションを行う
func ValidateTags(tags []string) error {
	if len(tags) > 10 {
		return NewError(ErrCodeValidation, "タグは10個まで設定できます", nil)
	}

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if len([]rune(tag)) == 0 {
			return NewError(ErrCodeValidation, "空のタグは設定できません", nil)
		}
		if len([]rune(tag)) > 30 {
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
	length := len([]rune(s))

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

// ValidateRating は評価値が 1〜5 の範囲かをバリデーションする汎用関数。
func ValidateRating(rating int) error {
	if rating < 1 || rating > 5 {
		return NewError(ErrCodeBadRequest, "評価は1〜5の範囲で指定してください", nil)
	}
	return nil
}

// External username validation constants
const (
	MinExternalUsernameLength = 1
	MaxExternalUsernameLength = 50
)

// externalUsernameRegex は外部サービスのユーザー名の形式をチェックする正規表現。
// 英数字・アンダースコア・ハイフン・ドットのみ許可（URLパスやクエリに安全に埋め込める文字のみ）。
var externalUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

// ValidateExternalUsername は外部サービス（Zenn/Qiita/AtCoder等）のユーザー名をバリデーションする。
// URLインジェクション防止のため、安全な文字のみ許可する。
func ValidateExternalUsername(username string) error {
	username = strings.TrimSpace(username)

	if len(username) < MinExternalUsernameLength {
		return NewError(ErrCodeValidation, "ユーザー名を入力してください", nil)
	}

	if len(username) > MaxExternalUsernameLength {
		return NewError(ErrCodeValidation, fmt.Sprintf("ユーザー名は%d文字以下である必要があります", MaxExternalUsernameLength), nil)
	}

	if !externalUsernameRegex.MatchString(username) {
		return NewError(ErrCodeValidation, "ユーザー名は英数字、アンダースコア、ハイフン、ドットのみ使用できます", nil)
	}

	return nil
}

// validLanguageCodes はYouTube検索で許可するISO 639-1言語コード。
var validLanguageCodes = map[string]bool{
	"ja": true, "en": true, "ko": true, "zh": true,
	"es": true, "fr": true, "de": true, "pt": true,
	"ru": true, "it": true, "ar": true, "hi": true,
	"th": true, "vi": true, "id": true, "tr": true,
	"pl": true, "nl": true, "sv": true, "da": true,
}

// ValidateLanguageCode は言語コードがホワイトリストに含まれるかを検証する。
func ValidateLanguageCode(lang string) error {
	if !validLanguageCodes[lang] {
		return NewError(ErrCodeValidation, "サポートされていない言語コードです", nil)
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
