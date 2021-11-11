package explain

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// callback name prefix
const (
	namePrefix = "gorm-plugin-explain:after:"
	ExplainCMD = "EXPLAIN"
)

// CallBackResult call back result
type CallBackResult struct {
	Err     error
	Results []Result
	SQL     string
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
		if gormDB.Error != nil && gormDB.Error != gorm.ErrRecordNotFound {
			gormDB.Logger.Warn(gormDB.Statement.Context, fmt.Sprintf("Explain call back failed: %s", gormDB.Error.Error()))
			return
		}

		if c.enable != nil && !c.enable() {
			gormDB.Logger.Info(gormDB.Statement.Context, "Explain call back not enable")
			return
		}

		sql := gormDB.Dialector.Explain(gormDB.Statement.SQL.String(), gormDB.Statement.Vars...)
		explainSQL := fmt.Sprintf("%s %s", ExplainCMD, sql)

		conn, err := gormDB.DB()
		if err != nil {
			gormDB.Logger.Error(gormDB.Statement.Context, err.Error())
			return
		}

		rows, err := conn.Query(explainSQL)
		if err != nil {
			gormDB.Logger.Error(gormDB.Statement.Context, err.Error())
			return
		}

		result, recom, err := c.explain.Analyze(rows)

		if err != nil {
			gormDB.Logger.Error(gormDB.Statement.Context, fmt.Sprintf("Query: %s, Error: %s", explainSQL, err.Error()))
			return
		}

		if c.fn != nil {
			var resErr error
			if recom != EmptyRecommendation {
				resErr = errors.New(recom)
			}

			c.fn(CallBackResult{Results: result, Err: resErr, SQL: sql})
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
