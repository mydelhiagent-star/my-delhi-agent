package models

// models/common.go
type BaseQueryParams struct {
    Page    *int    `query:"page"`
    Limit   *int    `query:"limit"`
    Sort    *string `query:"sort"`
    Order   *string `query:"order"`
}

// ✅ Default values
func (b *BaseQueryParams) SetDefaults() {
    if b.Page == nil || *b.Page < 1 {
        b.Page = &[]int{1}[0]
    }
    if b.Limit == nil || *b.Limit < 1 || *b.Limit > 100 {
        b.Limit = &[]int{20}[0]
    }
    if b.Sort == nil || *b.Sort == "" {
        b.Sort = &[]string{"created_at"}[0]
    }
    if b.Order == nil || *b.Order == "" {
        b.Order = &[]string{"desc"}[0]
    }
}