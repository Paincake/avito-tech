package dto

type Banner struct {
	Tags      []int64 `json:"tag_ids" validate:"nonzero"`
	FeatureId int64   `json:"feature_id" validate:"nonzero"`
	Content   Content `json:"content" validate:"nonzero"`
	IsActive  bool    `json:"is_active" validate:"nonzero"`
	CreatedAt string  `json:"created_at" validate:"nonzero"`
	UpdatedAt string  `json:"updated_at" validate:"nonzero"`
}

type Content struct {
	Title string `json:"title" validate:"nonzero"`
	Text  string `json:"text" validate:"nonzero"`
	Url   string `json:"url" validate:"nonzero"`
}

type User struct {
	Username string `json:"username" required:"true" validate:"nonzero"`
	Password string `json:"password" required:"true" validate:"nonzero"`
}
