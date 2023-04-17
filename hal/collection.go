package hal

import (
	"net/url"

	"github.com/68696c6c/goat/query"
)

// Collection represents a paginated list of API content.
// https://phlyrestfully.readthedocs.io/en/latest/halprimer.html#collections
type Collection[T any] struct {
	*query.Pagination
	*ResourceEmbeds[T]
	*ResourceLinks
}

func NewCollection[T any](resources []T, q query.Builder, linkUrl *url.URL) Collection[[]T] {
	return Collection[[]T]{
		Pagination:     q.GetPagination(),
		ResourceEmbeds: NewEmbeds[[]T](resources),
		ResourceLinks:  NewCollectionLinks(q, linkUrl),
	}
}

func NewCollectionLinks(q query.Builder, linkUrl *url.URL) *ResourceLinks {
	makeLink := func(p *query.Pagination, o *query.Order) Link {
		u := linkUrl.Query()
		p.ApplyToUrl(u)
		o.ApplyToUrl(u)
		linkUrl.RawQuery = u.Encode()
		return NewLink(linkUrl.String())
	}
	pagination := q.GetPagination()
	order := q.GetOrder()
	return NewLinks().
		AddLink("self", makeLink(pagination, order)).
		AddLink("first", makeLink(pagination.First(), order)).
		AddLink("previous", makeLink(pagination.Previous(), order)).
		AddLink("next", makeLink(pagination.Next(), order)).
		AddLink("last", makeLink(pagination.Last(), order))
}
