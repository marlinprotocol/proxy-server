package util

import (
	"errors"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func IsTcp(host string) bool {
	addr := strings.Split(host, ":")
	if len(addr) != 2 {
		log.Error(errors.New("not valip tcp address"))
		return false
	}
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		log.Error(err)
		return false
	} else if port > 65535 || port < 1 {
		log.Error(errors.New("not valid port"))
		return false
	}
	parts := strings.Split(addr[0], ".")

	if len(parts) != 4 {
		log.Error(errors.New("not valid ip address"))
		return false
	}
	
	for _,x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				log.Error(errors.New("not valid ip address"))
				return false
			}
		} else {
			log.Error(errors.New("not valid ip address"))
			return false
		}

	}
	return true
}

func IsVsock(host string) bool {
	addr := strings.Split(host, ":")
	if len(addr) != 2 {
		log.Error(errors.New("not valip tcp address"))
		return false
	}
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		log.Error(err)
		return false
	} else if port > 65535 || port < 1 {
		log.Error(errors.New("not valid port"))
		return false
	}

	_, err = strconv.Atoi(addr[0])
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

