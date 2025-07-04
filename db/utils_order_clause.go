package db

type OrderBy string

const (
	OrderASC  OrderBy = "ASC"
	OrderDESC OrderBy = "DESC"
)

func (o OrderBy) IsAsc() bool {
	return o == OrderASC
}

func (o OrderBy) IsDesc() bool {
	return !o.IsAsc()
}

func (o OrderBy) Sign() string {
	if o.IsDesc() {
		return "<"
	}

	return ">"
}

func (o OrderBy) String() string {
	return string(o)
}

func (o OrderBy) FullClause(field string) string {
	return " ORDER BY " + o.Clause(field)
}

func (o OrderBy) Clause(field string) string {
	return field + " " + o.String()
}
