package hypercloud

func (h* hypercloud) RegionInfo(regionId string) (json interface{}, err []error) {
    return h.Request("GET", "/regions/" + regionId, nil)
}

func (h* hypercloud) RegionList() (json interface{}, err []error) {
    return h.Request("GET", "/regions", nil)
}
