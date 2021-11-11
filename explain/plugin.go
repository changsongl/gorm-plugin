package explain

import (
	"gorm.io/gorm"
)

// plugin
type plugin struct {
	cb *callback
}

// Name of plugin
func (p plugin) Name() string {
	return "gorm-plugin:explain"
}

// Initialize plugin
func (p plugin) Initialize(db *gorm.DB) error {
	return p.cb.Register(db)
}

// New a explain plugin
func New(opts ...Option) gorm.Plugin {
	options := newOptions()
	for _, optFunc := range opts {
		optFunc.apply(options)
	}

	return plugin{
		cb: newCallBack(options),
	}
}
