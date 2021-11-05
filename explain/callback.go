package explain

import (
	"fmt"
	"gorm.io/gorm"
)

const namePrefix = "gorm-plugin-explain:after:"

type callback struct {
	explain *Explainer
	enable  func() bool
}

func newCallBack() *callback {
	return &callback{}
}

func (c *callback) Register(db *gorm.DB) error {
	var explainCB func(gormDB *gorm.DB)
	explainCB = func(gormDB *gorm.DB) {
		if c.enable != nil && !c.enable() {
			return
		}

		explainSQL := fmt.Sprintf("EXPLAIN %s",
			gormDB.Dialector.Explain(gormDB.Statement.SQL.String(), gormDB.Statement.Vars...))
		sqlDB, err := gormDB.DB()
		if err != nil {
			gormDB.Logger.Error(gormDB.Statement.Context, err.Error())
			return
		}

		result := sqlDB.QueryRow(explainSQL)

		if err = c.explain.Analyze(result); err != nil {
			gormDB.Logger.Error(gormDB.Statement.Context, err.Error())
			return
		}
	}

	if err := db.Callback().Create().After("gorm:create").Register(namePrefix+"gorm:create", explainCB); err != nil {
		return err
	}

	if err := db.Callback().Delete().After("gorm:delete").Register(namePrefix+"gorm:delete", explainCB); err != nil {
		return err
	}

	if err := db.Callback().Query().After("gorm:query").Register(namePrefix+"gorm:query", explainCB); err != nil {
		return err
	}

	if err := db.Callback().Update().After("gorm:update").Register(namePrefix+"gorm:update", explainCB); err != nil {
		return err
	}

	if err := db.Callback().Row().After("gorm:row").Register(namePrefix+"gorm:row", explainCB); err != nil {
		return err
	}

	if err := db.Callback().Raw().After("gorm:raw").Register(namePrefix+"gorm:raw", explainCB); err != nil {
		return err
	}

	return nil
}
