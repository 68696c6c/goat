package resource

type Link struct {
	Href string `json:"href"`
}

type LinkMap map[string]Link

type Links struct {
	Links LinkMap `json:"_links,omitempty" gorm:"-"`
}

func MakeLinks() *Links {
	return &Links{
		Links: make(LinkMap),
	}
}

func (l *Links) AddLink(key string, link Link) *Links {
	l.Links[key] = link
	return l
}

func MakeLink(href string) Link {
	return Link{
		Href: href,
	}
}

type LinkMaker func(string) Link
