package hypercloud

func (h* hypercloud) RegionInfo(regionId string) (json interface{}, err []error) {
    if len(regionId) == 3 { //Region code check (i.e. SY3/SV2 etc.)
        regions, errs := h.Request("GET", "/regions", nil)
        if errs != nil {
            return regions, errs
        }
        for _, r := range regions.([]interface{}){
            if r.(map[string]interface{})["code"].(string) == regionId {
                regionId = r.(map[string]interface{})["id"].(string)
                break
            }
        }
    }
    return h.Request("GET", "/regions/" + regionId, nil)
}

func (h* hypercloud) RegionList() (json interface{}, err []error) {
    return h.Request("GET", "/regions", nil)
}
