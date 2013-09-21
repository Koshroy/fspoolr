package statemanage

type Resource struct {
	Path string
	Mime string
}

func NewResource(pth string, mime string) *Resource {
	return &Resource{Path: pth, Mime: mime}
}
