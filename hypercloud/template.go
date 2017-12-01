package hypercloud

func (h *hypercloud) TemplateInfo(templateId string) (json interface{}, err []error) {
	return h.Request("GET", "/templates/"+templateId, nil)
}

func (h *hypercloud) TemplateList() (json interface{}, err []error) {
	return h.Request("GET", "/templates", nil)
}

func (h *hypercloud) TemplateSupersede(body interface{}) (json interface{}, err []error) {
	return h.Request("POST", "/templates", body)
}
