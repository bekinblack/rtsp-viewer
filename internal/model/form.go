package model

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	invalidIP       = "Некорректный IPv4-адрес"
	invalidPort     = "Порт должен быть числом 1–65535"
	invalidLogin    = "Логин не должен быть пустым"
	invalidPassword = "Пароль не должен быть пустым"
	invalidUriPair  = "Оба RTSP-URI совпадают после нормализации"
)

type Form struct {
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	PathHigh string `yaml:"pathHigh"`
	PathLow  string `yaml:"pathLow"`
}

func (f *Form) SetIP(ip string) {
	f.IP = strings.TrimSpace(ip)
}

func (f *Form) SetPort(port string) {
	f.Port = strings.TrimSpace(port)
}

func (f *Form) SetLogin(login string) {
	f.Login = strings.ToLower(strings.TrimSpace(login))
}

func (f *Form) SetPassword(password string) {
	f.Password = strings.TrimSpace(password)
}

func (f *Form) SetPathHigh(uri string) {
	f.PathHigh = normalizeUri(uri)
}

func (f *Form) SetPathLow(uri string) {
	f.PathLow = normalizeUri(uri)
}

func (f *Form) Validate() (hi, lo string, err error) {
	if net.ParseIP(f.IP) == nil {
		err = errors.New(invalidIP)
		return

	}

	portInt, err := strconv.Atoi(f.Port)
	if err != nil {
		err = errors.New(invalidPort)
		return
	}

	if portInt < 1 || portInt > 65535 {
		err = errors.New(invalidPort)
		return
	}

	if len(f.Login) == 0 {
		err = errors.New(invalidLogin)
		return
	}

	if len(f.Password) == 0 {
		err = errors.New(invalidPassword)
		return
	}

	if f.PathHigh == f.PathLow {
		err = errors.New(invalidUriPair)
		return
	}

	return f.maskedUriHigh(), f.maskedUriLow(), nil
}

func (f *Form) maskedUriHigh() string {
	return f.buildMaskedUri(f.PathHigh)
}
func (f *Form) maskedUriLow() string {
	return f.buildMaskedUri(f.PathLow)
}

func (f *Form) UriHigh() string {
	return f.buildUri(f.PathHigh)
}

func (f *Form) UriLow() string {
	return f.buildUri(f.PathLow)
}

func (f *Form) buildUri(path string) string {
	return fmt.Sprintf(
		"rtsp://%s:%s@%s:%s/%s",
		f.Login, f.Password, f.IP, f.Port, path,
	)
}

func (f *Form) buildMaskedUri(path string) string {
	asterisk := strings.Repeat("*", len(f.Password))
	return fmt.Sprintf(
		"rtsp://%s:%s@%s:%s/%s",
		f.Login, asterisk, f.IP, f.Port, path,
	)
}

func normalizeUri(uri string) string {
	return strings.Trim(strings.TrimSpace(uri), "/")
}
