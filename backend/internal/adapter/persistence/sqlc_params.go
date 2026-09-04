package persistence

import "math"

// int32Param は LIMIT/OFFSET など sqlc(pgx) の int32 パラメータへ渡す値を安全に変換する。
// usecase 層から渡る値は int のため、極端に大きい入力（不正な offset 等）で
// int32 へのキャストがサイレントにオーバーフロー・符号反転しないようにここで丸める。
// 負値は 0 に、int32 の最大値を超える値は int32 の最大値にクランプする。
func int32Param(n int) int32 {
	if n < 0 {
		return 0
	}
	if n > math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(n)
}
