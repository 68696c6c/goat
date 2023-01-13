package resource

type Resource struct {
	*Embedded[any]
	*Links
}

func MakeResource(selfHref string, embedded ...any) Resource {
	var embeds *Embedded[any]
	if len(embedded) > 0 {
		embeds = MakeEmbedded[any](embedded[0])
	}
	return Resource{
		Embedded: embeds,
		Links:    MakeResourceLinks(selfHref),
	}
}

func MakeResourceLinks(selfHref string) *Links {
	return MakeLinks().AddLink("self", MakeLink(selfHref))
}
