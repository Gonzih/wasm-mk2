package dom

var dom *DOM

type DOM struct {
	MockTemplates map[string]string
}

func RegisterMockTemplate(id, content string) {
	New().MockTemplates[id] = content
}

func New() *DOM {
	if dom == nil {
		dom = &DOM{
			MockTemplates: make(map[string]string, 0),
		}
	}

	return dom
}

func (d *DOM) TemplateContent(id string) string {
	return d.MockTemplates[id]
}
