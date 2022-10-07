package slack

type slack struct{}

func New(token string) *slack {
	return &slack{}
}
