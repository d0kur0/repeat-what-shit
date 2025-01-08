//go:build production

package consts

var IsProduction bool

func init() {
	IsProduction = true
}
