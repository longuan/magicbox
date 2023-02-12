package mongobox

import (
	"fmt"
	"strings"
)

type HostAndPort struct {
	Host string
	Port uint16
}

func (hp *HostAndPort) Address() string {
	return fmt.Sprintf("%s:%d", hp.Host, hp.Port)
}

func ConnStringForRs(rsName string, servers []HostAndPort) string {
	connStr := rsName + "/"
	for _, config := range servers {
		connStr += config.Address() + ","
	}
	connStr = strings.TrimSuffix(connStr, ",")
	return connStr
}
