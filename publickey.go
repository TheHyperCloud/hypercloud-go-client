package hypercloud

func (h* hypercloud) PublicKeyCreate(body interface{}) (json interface{}, err []string) {
    return h.Request("POST", "/public_keys", body)
}

func (h* hypercloud) PublicKeyDelete(pkId string) (json interface{}, err []string) {
    return h.Request("DELETE", "/public_keys/" + pkId, nil)
}

func (h* hypercloud) PublicKeyInfo(pkId string) (json interface{}, err []string) {
    return h.Request("GET", "/public_keys/" + pkId, nil)
}

func (h* hypercloud) PublicKeyList() (json interface{}, err []string) {
    return h.Request("GET", "/public_keys", nil)
}

func (h* hypercloud) PublicKeyUpdate(pkId string, body interface{}) (json interface{}, err []string) {
    return h.Request("PUT", "/public_keys/" + pkId, body)
}
