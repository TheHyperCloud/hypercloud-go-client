package hypercloud

func (h* hypercloud) ConsoleSessionInfo(consoleSessionIdentity string) (ret interface{}, err []string) {
    ret, err = h.Request("GET", "/console_sessions/" + consoleSessionIdentity, nil)
    return
}
