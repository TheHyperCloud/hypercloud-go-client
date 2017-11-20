package hypercloud

func (h* hypercloud) IPAddressCreate(body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("POST", "/ip_addresses", body)
    return
}

func (h* hypercloud) IPAddressDelete(IPAddrID string) (ret interface{}, err []error) {
    ret, err = h.Request("DELETE", "/ip_addresses/" + IPAddrID, nil)
    return
}

func (h* hypercloud) IPAddressList() (ret interface{}, err []error){
    ret, err = h.Request("GET", "/ip_addresses", nil)
    return
}

func (h* hypercloud) IPAddressListPrivate() (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/ip_addresses/private", nil)
    return
}

func (h* hypercloud) IPAddressesListPublic() (ret interface{}, err []error) {
    ret, err = h.Request("GET", "/ip_addresses/public", nil)
    return
}

func (h* hypercloud) IPAddressInfo(IPAddrID string) (ret interface{}, err []error){
    ret, err = h.Request("GET", "/ip_addresses/" + IPAddrID, nil)
    return
}

func (h* hypercloud) IPAddressUpdate(IPAddrID string, body interface{}) (ret interface{}, err []error) {
    ret, err = h.Request("PUT", "/ip_addresses/" + IPAddrID, body)
    return
}
