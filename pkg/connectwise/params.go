package connectwise

import "fmt"

type QueryParams struct {
	Conditions            string
	ChildConditions       string
	CustomFieldConditions string
	OrderBy               string // must be asc or desc
	Fields                string
	Columns               string
	Page                  int
	PageSize              int // Default is 25, Max Size is 1000
}

func (q *QueryParams) ToMap() map[string]string {
	m := make(map[string]string)
	if q.Conditions != "" {
		m["conditions"] = q.Conditions
	}

	if q.ChildConditions != "" {
		m["childConditions"] = q.ChildConditions
	}

	if q.CustomFieldConditions != "" {
		m["customFieldConditions"] = q.CustomFieldConditions
	}

	if q.OrderBy != "" {
		m["orderBy"] = q.OrderBy
	}

	if q.Fields != "" {
		m["fields"] = q.Fields
	}

	if q.Columns != "" {
		m["columns"] = q.Columns
	}

	if q.Page != 0 {
		m["page"] = fmt.Sprintf("%d", q.Page)
	}

	if q.PageSize != 0 {
		m["pageSize"] = fmt.Sprintf("%d", q.PageSize)
	}

	return m
}
