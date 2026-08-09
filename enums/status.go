package enums

type PostStatus string

const (
	StatusPublish PostStatus = "publish"
	StatusDraft   PostStatus = "draft"
	StatusThrash  PostStatus = "thrash"
)
