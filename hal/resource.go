package hal

// Resource represents an individual piece of content exposed by the API.
// https://phlyrestfully.readthedocs.io/en/latest/halprimer.html#resources
type Resource struct {
	*ResourceEmbeds[any]
	*ResourceLinks
}

func NewResource(selfHref string, embedded ...any) Resource {
	var embeds *ResourceEmbeds[any]
	if len(embedded) > 0 {
		embeds = NewEmbeds[any](embedded[0])
	}
	return Resource{
		ResourceEmbeds: embeds,
		ResourceLinks:  NewResourceLinks(selfHref),
	}
}

func NewResourceLinks(selfHref string) *ResourceLinks {
	return NewLinks().AddLink("self", NewLink(selfHref))
}

type ResourceEmbeds[T any] struct {
	Embeds T `json:"_embedded,omitempty" gorm:"embedded"`
}

func NewEmbeds[T any](embeds T) *ResourceEmbeds[T] {
	return &ResourceEmbeds[T]{
		Embeds: embeds,
	}
}

type Links map[string]Link

type ResourceLinks struct {
	Links `json:"_links,omitempty" gorm:"-"`
}

func NewLinks() *ResourceLinks {
	return &ResourceLinks{
		Links: make(Links),
	}
}

func (l *ResourceLinks) AddLink(key string, link Link) *ResourceLinks {
	l.Links[key] = link
	return l
}

type Link struct {
	Href string `json:"href"`
}

func NewLink(href string) Link {
	return Link{
		Href: href,
	}
}
