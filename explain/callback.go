package explain

import (
	"fmt"
	"gorm.io/gorm"
)

// callback name prefix
const namePrefix = "gorm-plugin-explain:after:"

// CallBackResult call back result
type CallBackResult struct {
	Err     error
	Results []Result
}

// callback struct
type callback struct {
	explain *Explainer
	enable  func() bool
	fn      func(CallBackResult)
}

// newCallBack new a call back
func newCallBack(opts *options) *callback {
	return &callback{
		enable:  opts.enable,
		fn:      opts.fn,
		explain: NewExplainer(opts.explainOpts),
	}
}

// Register explain function to all callback processes
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

		rows, err := sqlDB.Query(explainSQL)
		if err != nil {
			gormDB.Logger.Error(gormDB.Statement.Context, err.Error())
			return
		}

		result, err := c.explain.Analyze(rows)

		if err != nil {
			gormDB.Logger.Error(gormDB.Statement.Context, fmt.Sprintf("Query: %s, Error: %s", explainSQL, err.Error()))
			return
		}

		if c.fn != nil {
			c.fn(CallBackResult{Results: result, Err: err})
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
