package hypercloud

func (h* hypercloud) RegionInfo(regionId string) (json interface{}, err []string) {
    return h.Request("GET", "/regions/" + regionId, nil)
}

func (h* hypercloud) RegionList() (json interface{}, err []string) {
    return h.Request("GET", "/regions", nil)
}
