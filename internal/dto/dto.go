package dto

type Banner struct {
	Tags      []int64 `json:"tag_ids"`
	FeatureId int64   `json:"feature_id"`
	Content   Content `json:"content"`
	IsActive  bool    `json:"is_active"`
}

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}
