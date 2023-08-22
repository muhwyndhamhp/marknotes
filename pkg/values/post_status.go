package values

type PostStatus string

const (
	None      PostStatus = ""
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
)
