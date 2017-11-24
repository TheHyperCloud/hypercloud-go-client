package hypercloud

func (h* hypercloud) DiskCreate(body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/disks", body)
    return
}

func (h* hypercloud) DiskDelete(diskId string) (ret interface{}, err []error) {
    ret, err = h.Request("DELETE", "/disks/" + diskId, nil)
    return
}

func (h* hypercloud) DiskInfo(diskId string) (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/disks/" + diskId, nil)
    return
}

func (h* hypercloud) DiskState(diskId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/disks/" + diskId + "/state", body)
    return
}

func (h* hypercloud) DiskList() (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/disks", nil)
    return
}

func (h* hypercloud) DiskUpdate(diskId string, body interface{}) (ret interface{}, err []error) {
    ret, err =  h.Request("PUT", "/disks/" + diskId, body)
    return
}

func (h* hypercloud) DiskResize(diskId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/disks/" + diskId + "/resize", body)
    return
}

func (h* hypercloud) DiskClone(diskId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/disks/" + diskId + "/clone", body)
    return
}
