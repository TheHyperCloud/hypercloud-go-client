package hypercloud

func (h* hypercloud) PerformanceTierListInstance() (json interface{}, err []string) {
    return h.Request("GET", "/performance_tiers/instances", nil)
}

func (h* hypercloud) PerformanceTierListDisk() (json interface{}, err[]string) {
    return h.Request("GET", "/performance_tiers/disks", nil)
}
