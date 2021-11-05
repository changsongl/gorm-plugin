package explain

import (
	"gorm.io/gorm"
)

type Plugin interface {
	Name() string
	Initialize(db *gorm.DB) error
}

type plugin struct {
	cb *callback
}

func (p plugin) Name() string {
	return "gorm-plugin:explain"
}

func (p plugin) Initialize(db *gorm.DB) error {
	return p.cb.Register(db)
}

// New a explain plugin
func New() Plugin {
	return plugin{
		cb: newCallBack(),
	}
}
