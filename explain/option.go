package explain

type options struct {
	enable      func() bool
	fn          func(CallBackResult)
	explainOpts explainerOptions
}

type optFunc func(*options)

func (f optFunc) apply(opts *options) {
	f(opts)
}

type Option interface {
	apply(*options)
}

func newOptions() *options {
	return &options{
		explainOpts: explainerOptions{},
	}
}

func EnableFuncOption(f func() bool) Option {
	return optFunc(func(opt *options) {
		opt.enable = f
	})
}

func CallBackFuncOption(cb func(CallBackResult)) Option {
	return optFunc(func(opt *options) {
		opt.fn = cb
	})
}

func ExtraWhiteListOption(list []ResultExtra) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.ExtraWhiteList = ExtraList(list)
	})
}

func ExtraBlackListOption(list []ResultExtra) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.ExtraBlackList = ExtraList(list)
	})
}

func SelectTypeWhiteListOption(list []ResultSelectType) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.SelectTypeWhiteList = SelectTypeList(list)
	})
}

func SelectTypeBlackListOption(list []ResultSelectType) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.SelectTypeBlackList = SelectTypeList(list)
	})
}

func TypeLevelOption(typeLevel ResultType) Option {
	return optFunc(func(opt *options) {
		opt.explainOpts.TypeLevel = typeLevel
	})
}
