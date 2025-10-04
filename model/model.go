package model

import (
	"net"
	"strconv"
	"strings"
)

func New(conn Connector, view Viewer, log Logger) *Model {
	return &Model{
		connector: conn,
		viewer:    view,
		logger:    log,
	}
}

type Model struct {
	ip        string
	port      string
	login     string
	password  string
	uriHigh   string
	uriLow    string
	status    Status
	connector Connector
	viewer    Viewer
	logger    Logger
}

func (m *Model) RtspHighUri() string {
	return m.uriHigh
}
func (m *Model) RtspLowUri() string {
	return m.uriLow
}

func (m *Model) ChangeIP(ip string) {
	m.ip = strings.TrimSpace(ip)
}

func (m *Model) ChangePort(port string) {
	m.port = strings.TrimSpace(port)
}

func (m *Model) ChangeLogin(login string) {
	m.login = strings.ToLower(strings.TrimSpace(login))
}

func (m *Model) ChangePassword(password string) {
	m.password = password
}

func (m *Model) ChangeUriHigh(uri string) {
	m.uriHigh = strings.ToLower(strings.TrimSpace(uri))
}

func (m *Model) ChangeUriLow(uri string) {
	m.uriLow = strings.ToLower(strings.TrimSpace(uri))
}

func (m *Model) Validate() {
	if net.ParseIP(m.ip) == nil {
		m.logger.Log(invalidIP)
		return
	}

	portInt, err := strconv.Atoi(m.port)
	if err != nil {
		m.logger.Log(invalidPort)
		return
	}

	if portInt < 1 || portInt > 65535 {
		m.logger.Log(invalidPort)
		return
	}

	if len(m.login) == 0 {
		m.logger.Log(invalidLogin)
		return
	}

	if len(m.password) == 0 {
		m.logger.Log(invalidPassword)
		return
	}

	if !strings.HasPrefix(m.uriHigh, "rtsp://") {
		m.logger.Log(invalidURI)
		return
	}

	if !strings.HasPrefix(m.uriLow, "rtsp://") {
		m.logger.Log(invalidURI)
		return
	}

	if m.uriHigh == m.uriLow {
		m.logger.Log(invalidUriPair)
		return
	}

	m.logger.Log(ok)
}
