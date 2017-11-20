package hypercloud

func (h* hypercloud) TemplateList() (json interface{}, err []string) {
    return h.Request("GET", "/templates", nil)
}

func (h* hypercloud) TemplateInfo(templateId string) (json interface{}, err []string) {
    return h.Request("GET", "/templates/" + templateId, nil)
}
