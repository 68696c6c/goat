package resource

type Embedded[T any] struct {
	Embeds T `json:"_embedded,omitempty" gorm:"embedded"`
}

func MakeEmbedded[T any](embeds T) *Embedded[T] {
	return &Embedded[T]{
		Embeds: embeds,
	}
}
