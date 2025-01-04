//go:build !production

package consts

// IsProduction возвращает true если это production сборка
var IsProduction bool

func init() {
	IsProduction = false
}
