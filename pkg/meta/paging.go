package meta

type Paging struct {
	Limit  int   `json:"limit" validate:"max=1000"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total,omitempty"`
}

func (p *Paging) SetPage(page int) {
	p.Offset = page * p.Limit
}
