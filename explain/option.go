package explain

// options option data
type options struct {
	enable      func() bool
	fn          func(CallBackResult)
	explainOpts explainerOptions
}

// Option interface to apply changes on options
type Option interface {
	apply(*options)
}

// optFunc option function
type optFunc func(*options)

// apply implements Option
func (f optFunc) apply(opts *options) {
	f(opts)
}

// newOptions create a new options
func newOptions() *options {
	return &options{
		explainOpts: explainerOptions{},
	}
}

// EnableFuncOption for controlling this callback whether it is enable,
// it will be called before the callback runs
func EnableFuncOption(f func() bool) Option {
	return optFunc(func(opt *options) {
		opt.enable = f
	})
}

// CallBackFuncOption this function will be called for every explain result
func CallBackFuncOption(cb func(CallBackResult)) Option {
	return optFunc(func(opt *options) {
		opt.fn = cb
	})
}

// ExtraWhiteListOption extra filed white list, if it has a white list,
// the explain won't pass when the extra data is not in white list.
func ExtraWhiteListOption(list []ResultExtra) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.ExtraWhiteList = ExtraList(list)
	})
}

// ExtraBlackListOption extra filed black list, if it has a black list,
// the explain won't pass when the extra data is in black list.
func ExtraBlackListOption(list []ResultExtra) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.ExtraBlackList = ExtraList(list)
	})
}

// SelectTypeWhiteListOption select type filed white list, if it has a white list,
// the explain won't pass when the select type  is not in white list.
func SelectTypeWhiteListOption(list []ResultSelectType) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.SelectTypeWhiteList = SelectTypeList(list)
	})
}

// SelectTypeBlackListOption select type filed black list, if it has a black list,
// the explain won't pass when the select type is in black list.
func SelectTypeBlackListOption(list []ResultSelectType) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.SelectTypeBlackList = SelectTypeList(list)
	})
}

// TypeLevelOption type field requirement. It won't pass if the actually type field is
// worse than it.
func TypeLevelOption(typeLevel ResultType) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.TypeLevel = typeLevel
	})
}
