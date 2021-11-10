package explain

import (
	"database/sql"
	"fmt"
	"strings"
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

func (t ResultType) IsValid() bool {
	return ResultTypePriorityMap[t] != 0
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
	requirement explainerOptions
}

type WhiteList interface {
	IsInWhiteList(string string) error
}

type BlackList interface {
	IsInInBlackList(string string) error
}

type ExtraList []ResultExtra

func (extras ExtraList) IsInWhiteList(ex string) error {
	for _, extra := range extras {
		if strings.ToLower(ex) == strings.ToLower(string(extra)) {
			return nil
		}
	}

	return fmt.Errorf("\"%s\" is not in extra white list", ex)
}

func (extras ExtraList) IsInInBlackList(ex string) error {
	for _, extra := range extras {
		if strings.ToLower(ex) == strings.ToLower(string(extra)) {
			return fmt.Errorf("\"%s\" is in extra black list", ex)
		}
	}

	return nil
}

type SelectTypeList []ResultSelectType

func (selectTypes SelectTypeList) IsInWhiteList(st string) error {
	for _, selectType := range selectTypes {
		if strings.ToLower(st) == strings.ToLower(string(selectType)) {
			return nil
		}
	}

	return fmt.Errorf("\"%s\" is not in select type white list", st)
}

func (selectTypes SelectTypeList) IsInInBlackList(st string) error {
	for _, selectType := range selectTypes {
		if strings.ToLower(st) == strings.ToLower(string(selectType)) {
			return fmt.Errorf("\"%s\" is in select type black list", st)
		}
	}

	return nil
}

// explainerOptions explainer requirement
type explainerOptions struct {
	ExtraWhiteList      WhiteList
	ExtraBlackList      BlackList
	SelectTypeWhiteList WhiteList
	SelectTypeBlackList BlackList
	TypeLevel           ResultType
}

// NewExplainer to check sql explain
func NewExplainer(req explainerOptions) *Explainer {
	return &Explainer{requirement: req}
}

func (e *Explainer) Analyze(rows *sql.Rows) ([]Result, error) {
	var results []Result
	if err := rows.Scan(&results); err != nil {
		return nil, err
	}

	for _, row := range results {
		if err := e.requirement.ExtraBlackList.IsInInBlackList(row.Extra); err != nil {
			return results, err
		}

		if err := e.requirement.ExtraWhiteList.IsInWhiteList(row.Extra); err != nil {
			return results, err
		}

		if err := e.requirement.SelectTypeBlackList.IsInInBlackList(row.SelectType); err != nil {
			return results, err
		}

		if err := e.requirement.SelectTypeWhiteList.IsInWhiteList(row.SelectType); err != nil {
			return results, err
		}

		if e.requirement.TypeLevel != ResultTypeNone {
			if !e.requirement.TypeLevel.IsValid() {
				return results, fmt.Errorf("%s is not valid type", e.requirement.TypeLevel)
			}

			if !ResultType(row.Type).IsValid() {
				return results, fmt.Errorf("%s is not valid type from row reulst", row.Type)
			}
		}
	}

	return results, nil
}
