package validate

import (
	kitvalidate "github.com/gookit/validate"
)

func EnableDefaultConfig() {
	kitvalidate.Config(func(opt *kitvalidate.GlobalOption) {
		opt.SkipOnEmpty = false
		opt.StopOnError = false

		opt.ValidatePrivateFields = false
		opt.FieldTag = "json"

		opt.CheckZero = true
	})
}
