package resource

import (
	"net/url"
)

type Collection[T any] struct {
	Pagination
	*Embedded[T]
	*Links
}

func MakeCollection[T any](resources []T, pagination Pagination, linkUrl *url.URL) Collection[[]T] {
	return Collection[[]T]{
		Pagination: pagination,
		Embedded:   MakeEmbedded[[]T](resources),
		Links:      MakeCollectionLinks(pagination, linkUrl),
	}
}

func MakeCollectionLinks(pagination Pagination, linkUrl *url.URL) *Links {
	makeLink := func(p Pagination) Link {
		u := linkUrl.Query()
		p.AddToQuery(u)
		linkUrl.RawQuery = u.Encode()
		return MakeLink(linkUrl.String())
	}
	return MakeLinks().
		AddLink("self", makeLink(pagination)).
		AddLink("first", makeLink(pagination.First())).
		AddLink("previous", makeLink(pagination.Previous())).
		AddLink("next", makeLink(pagination.Next())).
		AddLink("last", makeLink(pagination.Last()))
}
