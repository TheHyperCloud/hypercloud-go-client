package hypercloud

func (h* hypercloud) NetworkCreate(body interface{}) (json interface{}, err []error) {
    return h.Request("POST", "/networks", body)
}

func (h* hypercloud) NetworkDelete(netId string) (json interface{}, err []error) {
    return h.Request("DELETE", "/networks/" + netId, nil)
}

func (h* hypercloud) NetworkList() (json interface{}, err []error) {
    return h.Request("GET", "/networks", nil)
}

func (h* hypercloud) NetworkListPrivate() (json interface{}, err []error) {
    return h.Request("GET", "/networks/private", nil)
}

func (h* hypercloud) NetworkListPublic() (json interface{}, err []error) {
    return h.Request("GET", "/networks/public", nil)
}

func (h* hypercloud) NetworkInfo(netId string) (json interface{}, err []error) {
    return h.Request("GET", "/networks/" + netId, nil)
}

func (h* hypercloud) NetworkUpdate(netId string, body interface{}) (json interface{}, err []error) {
    return h.Request("PUT", "/networks/" + netId, body)
}
