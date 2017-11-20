package hypercloud

func (h* hypercloud) InstanceBasicCreate(body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("POST", "/instances", body)
    return
}

func (h* hypercloud) InstanceAssemble(body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("POST", "/instances/assemble", body)
    return
}

func (h* hypercloud) InstanceDelete(instanceId string) (ret interface{}, err []string) {
    ret, err = h.Request("DELETE", "/instances/" + instanceId, nil)
    return
}

func (h* hypercloud) InstanceInfo(instanceId string) (ret interface{}, err []string) {
    ret, err = h.Request("GET", "/instances/" + instanceId, nil)
    return
}

func (h* hypercloud) InstanceList() (ret interface{}, err []string) {
    ret, err = h.Request("GET", "/instances", nil)
    return
}

func (h* hypercloud) InstanceUpdate(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("PUT", "/instances/" + instanceId, body)
    return
}

func (h* hypercloud) InstanceState(instanceId string) (ret interface{}, err []string) {
    ret, err = h.Request("GET", "/instances/" + instanceId + "/state", nil)
    return
}

func (h* hypercloud) InstanceNote(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("GET", "/instances/" + instanceId + "/note", body)
    return
}

func (h* hypercloud) InstanceStart(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("POST", "/instances/" + instanceId + "/start", body)
    return
}

func (h* hypercloud) InstanceStop(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("POST", "/instances/" + instanceId + "/stop", body)
    return
}

func (h* hypercloud) InstanceRemoteAccess(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("POST", "/instances/" + instanceId + "/remote_access", body)
    return
}

func (h* hypercloud) InstanceUpdateDisks(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("PUT", "/instances/" + instanceId + "/disks", body)
    return
}

func (h* hypercloud) InstanceUpdatePublicKeys(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("PUT", "/instances/" + instanceId + "/public_keys", body)
    return
}

func (h* hypercloud) InstanceUpdateNetworking(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("PUT", "/instances/" + instanceId + "/network_adapters", body)
    return
}

func (h* hypercloud) InstanceGetContext(instanceId string) (ret interface{}, err []string) {
    ret, err = h.Request("GET", "/instances/" + instanceId + "/context", nil)
    return
}

func (h* hypercloud) InstanceSetContext(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("POST", "/instances/" + instanceId + "/context", body)
    return
}

func (h* hypercloud) InstanceUpdateContext(instanceId string, body interface{}) (ret interface{}, err []string) {
    ret, err = h.Request("PUT", "/instances/" + instanceId + "/context", body)
    return
}

func (h* hypercloud) InstanceDeleteContextKey(instanceId string, instanceContextKey string) (ret interface{}, err []string) {
    ret, err = h.Request("DELETE", "/instances/" + instanceId + "/context/" + instanceContextKey, nil)
    return
}
