package helper

import (
	"github.com/sony/sonyflake"
)

var showflake *sonyflake.Sonyflake

func Showflake() *sonyflake.Sonyflake {
	if showflake == nil {
		showflake = sonyflake.NewSonyflake(sonyflake.Settings{})
	}
	return showflake
}
