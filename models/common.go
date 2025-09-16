package models

// models/common.go
type BaseQueryParams struct {
    Page    int    `query:"page"`
    Limit   int    `query:"limit"`
    Sort    string `query:"sort"`
    Order   string `query:"order"`
}

// ✅ Default values
func (b *BaseQueryParams) SetDefaults() {
    if b.Page < 1 {
        b.Page = 1
    }
    if b.Limit < 1 || b.Limit > 100 {
        b.Limit = 20
    }
    if b.Sort == "" {
        b.Sort = "created_at"
    }
    if b.Order == "" {
        b.Order = "desc"
    }
}