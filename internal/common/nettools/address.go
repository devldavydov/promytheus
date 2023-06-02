package nettools

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrAddressFormat = errors.New("wrong address format")

// Address - utility struct for holding address pair.
type Address struct {
	Host string
	Port int
}

func NewAddress(addr string) (Address, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return Address{}, ErrAddressFormat
	}

	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return Address{}, ErrAddressFormat
	}

	return Address{Host: host, Port: port}, nil
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
