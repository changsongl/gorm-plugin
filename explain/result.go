package explain

import (
	"database/sql"
)

// ResultSelectType explain result select type
type ResultSelectType string

const (
	ResultSelectTypeSimple      ResultSelectType = "SIMPLE"
	ResultSelectTypePrimary     ResultSelectType = "PRIMARY"
	ResultSelectTypeSubQuery    ResultSelectType = "SUBQUERY"
	ResultSelectTypeDerived     ResultSelectType = "DERIVED"
	ResultSelectTypeUnion       ResultSelectType = "UNION"
	ResultSelectTypeUnionResult ResultSelectType = "UNION RESULT"
)

// ResultType explain result type
type ResultType string

const (
	ResultTypeNone   ResultType = ""
	ResultTypeAll    ResultType = "all"
	ResultTypeIndex  ResultType = "index"
	ResultTypeRange  ResultType = "range"
	ResultTypeRef    ResultType = "ref"
	ResultTypeEQRef  ResultType = "eq_ref"
	ResultTypeConst  ResultType = "const"
	ResultTypeSystem ResultType = "system"
)

// ResultTypePriorityMap priority of result type
var ResultTypePriorityMap = map[ResultType]int{
	ResultTypeNone:   0,
	ResultTypeAll:    1,
	ResultTypeIndex:  2,
	ResultTypeRange:  3,
	ResultTypeRef:    4,
	ResultTypeEQRef:  5,
	ResultTypeConst:  6,
	ResultTypeSystem: 7,
}

// ResultExtra explain result extra
type ResultExtra string

const (
	ResultExtraFileSort        ResultExtra = "using filesort"
	ResultExtraTemporary       ResultExtra = "using temporary"
	ResultExtraIndex           ResultExtra = "using index"
	ResultExtraWhere           ResultExtra = "using where"
	ResultExtraJoinBuffer      ResultExtra = "using join buffer"
	ResultExtraImpossibleWhere ResultExtra = "impossible where"
	ResultExtraOptimizedAway   ResultExtra = "select tables optimized away"
	ResultExtraDistinct        ResultExtra = "distinct"
)

type Result struct {
	Id          int    `gorm:"id" json:"id"`
	SelectType  string `gorm:"select_type" json:"select_type"`
	Table       string `gorm:"table" json:"table"`
	Type        string `gorm:"type" json:"type"`
	PossibleKey string `gorm:"possible_keys" json:"possible_keys"`
	Key         string `gorm:"key" json:"key"`
	KeyLen      int    `gorm:"key_len" json:"key_len"`
	Ref         string `gorm:"ref" json:"ref"`
	Rows        int    `gorm:"rows" json:"rows"`
	Extra       string `gorm:"Extra" json:"Extra"`
}

type Explainer struct {
}

// NewExplainer to check sql explain
func NewExplainer() *Explainer {
	return &Explainer{}
}

func (e *Explainer) Analyze(result *sql.Row) error {
	return nil
}
