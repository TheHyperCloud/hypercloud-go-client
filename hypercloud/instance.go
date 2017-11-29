package hypercloud

func (h* hypercloud) InstanceBasicCreate(body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/instances", body)
    return
}

func (h* hypercloud) InstanceAssemble(body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/instances/assemble", body)
    return
}

func (h* hypercloud) InstanceDelete(instanceId string) (ret interface{}, err []error) {
    ret, err = h.Request("DELETE", "/instances/" + instanceId, nil)
    return
}

func (h* hypercloud) InstanceInfo(instanceId string) (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/instances/" + instanceId, nil)
    return
}

func (h* hypercloud) InstanceList() (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/instances", nil)
    return
}

/* Takes everything listed in the API reference, as well as the following updates:
    - Availability groups
    - Disks
    - Public keys
    - Networking
*/
func (h* hypercloud) InstanceUpdate(instanceId string, body interface{}) (ret interface{}, err []error) {
    /* Time to munch up the following keys and shove them into the correct functions */
    dat := body.(map[string]interface{})
    if val, ok := dat["availability_groups"]; ok {
        _, erro := h.InstanceUpdateHighAvailability(instanceId, map[string]interface{}{"availability_groups" : val})
        if erro != nil {
            err = append(err, erro...)
        }
        delete(dat, "availability_groups")
    }
    if val, ok := dat["disks"]; ok {
        _, erro := h.InstanceUpdateDisks(instanceId, map[string]interface{}{"disks" : val})
        if erro != nil {
            err = append(err, erro...)
        }
        delete(dat, "disks")
    }
    if val, ok := dat["network_adapters"]; ok {
        _, erro := h.InstanceUpdateNetworking(instanceId, map[string]interface{}{"network_adapters" : val})
        if erro != nil {
            err = append(err, erro...)
        }
        delete(dat, "network_adapters")
    }
    if val, ok := dat["public_keys"]; ok {
        _, erro := h.InstanceUpdatePublicKeys(instanceId, map[string]interface{}{"public_keys" : val})
        if erro != nil {
            err = append(err, erro...)
        }
        delete(dat, "public_keys")
    }
    if len(dat) == 0{
        ret, _ = h.InstanceInfo(instanceId)
        return
    }
    ret, err = h.Request("PUT", "/instances/" + instanceId, dat)
    return
}

func (h* hypercloud) InstanceState(instanceId string) (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/instances/" + instanceId + "/state", nil)
    return
}

func (h* hypercloud) InstanceNote(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/instances/" + instanceId + "/note", body)
    return
}

func (h* hypercloud) InstanceStart(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/instances/" + instanceId + "/start", body)
    return
}

func (h* hypercloud) InstanceStop(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/instances/" + instanceId + "/stop", body)
    return
}

func (h* hypercloud) InstanceRemoteAccess(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/instances/" + instanceId + "/remote_access", body)
    return
}

/* Leaving functionality there, but I've merged this all into the "update" function because it makes sense */
func (h* hypercloud) InstanceUpdateDisks(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("PUT", "/instances/" + instanceId + "/disks", body)
    return
}

func (h* hypercloud) InstanceUpdatePublicKeys(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("PUT", "/instances/" + instanceId + "/public_keys", body)
    return
}

func (h* hypercloud) InstanceUpdateNetworking(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("PUT", "/instances/" + instanceId + "/network_adapters", body)
    return
}

func (h* hypercloud) InstanceUpdateHighAvailability(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("PUT", "/instances/" + instanceId + "/availability_group", body)
    return
}

func (h* hypercloud) InstanceGetContext(instanceId string) (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/instances/" + instanceId + "/context", nil)
    return
}

func (h* hypercloud) InstanceSetContext(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/instances/" + instanceId + "/context", body)
    return
}

func (h* hypercloud) InstanceUpdateContext(instanceId string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("PUT", "/instances/" + instanceId + "/context", body)
    return
}

func (h* hypercloud) InstanceDeleteContextKey(instanceId string, instanceContextKey string) (ret interface{}, err []error) {
    ret, err = h.Request("DELETE", "/instances/" + instanceId + "/context/" + instanceContextKey, nil)
    return
}
