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

// Result mysql fields
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

// Explainer for sql explain
type Explainer struct {
	requirement explainerOptions
}

// WhiteList check it is in white list
type WhiteList interface {
	IsInWhiteList(string string) error
}

// BlackList check it is in black list
type BlackList interface {
	IsInInBlackList(string string) error
}

// ExtraList a list of extra
type ExtraList []ResultExtra

// IsInWhiteList check if ex is in white list
func (extras ExtraList) IsInWhiteList(ex string) error {
	if len(extras) == 0 {
		return nil
	}

	for _, exDes := range strings.Split(ex, ",") {
		exDes = strings.TrimSpace(exDes)
		for _, extra := range extras {
			if strings.ToLower(exDes) == strings.ToLower(string(extra)) {
				return nil
			}
		}
	}

	return fmt.Errorf("\"%s\" is not in extra white list", ex)
}

// IsInInBlackList check if ex in black list
func (extras ExtraList) IsInInBlackList(ex string) error {
	if len(extras) == 0 {
		return nil
	}

	for _, exDes := range strings.Split(ex, ",") {
		exDes = strings.TrimSpace(exDes)
		for _, extra := range extras {
			if strings.ToLower(exDes) == strings.ToLower(string(extra)) {
				return fmt.Errorf("\"%s\" is in extra black list", ex)
			}
		}
	}

	return nil
}

// SelectTypeList a group of select type
type SelectTypeList []ResultSelectType

// IsInWhiteList check it is in white list
func (selectTypes SelectTypeList) IsInWhiteList(st string) error {
	if len(selectTypes) == 0 {
		return nil
	}

	for _, selectType := range selectTypes {
		if strings.ToLower(st) == strings.ToLower(string(selectType)) {
			return nil
		}
	}

	return fmt.Errorf("\"%s\" is not in select type white list", st)
}

// IsInInBlackList check it is in black list
func (selectTypes SelectTypeList) IsInInBlackList(st string) error {
	if len(selectTypes) == 0 {
		return nil
	}

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

// Analyze every fields
func (e *Explainer) Analyze(rows *sql.Rows) ([]Result, error) {
	results, err := e.extractResults(rows)
	if err != nil {
		return nil, err
	}

	// TODO: refactor analyze process by interface by interface
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

// extractResults extract sql.Rows to result set
func (e *Explainer) extractResults(rows *sql.Rows) ([]Result, error) {
	var results []Result

	for rows.Next() {
		var id, keyLen, rowNum sql.NullInt64
		var selectType, table, typeField, posKey, key, ref, extra sql.NullString

		if err := rows.Scan(
			&id, &selectType, &table, &typeField,
			&posKey, &key, &keyLen, &ref, &rowNum, &extra); err != nil {
			return nil, err
		}

		results = append(results, Result{
			Id:          int(id.Int64),
			SelectType:  selectType.String,
			Table:       table.String,
			Type:        typeField.String,
			PossibleKey: posKey.String,
			Key:         key.String,
			KeyLen:      int(keyLen.Int64),
			Ref:         ref.String,
			Rows:        int(rowNum.Int64),
			Extra:       extra.String,
		})
	}

	return results, nil
}
