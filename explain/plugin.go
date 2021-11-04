package explain

import (
	"gorm.io/gorm"
)

type Plugin interface {
	Name() string
	Initialize(db *gorm.DB) error
}

type plugin struct {

}
