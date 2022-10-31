package util

import (
	"errors"
	"os"
	"os/user"
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

func CreateDirPathIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		currentUser, err := GetUser()
		if err != nil {
			return err
		}
		uid, err := strconv.Atoi(currentUser.Uid)
		if err != nil {
			return err
		}
		gid, err := strconv.Atoi(currentUser.Gid)
		if err != nil {
			return err
		}
		createErr := os.MkdirAll(dirPath, 0777)
		if createErr != nil {
			return createErr
		}
		return os.Chown(dirPath, uid, gid)
	}
	return nil
}

func GetUser() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	if os.Geteuid() == 0 {
		// Root, try to retrieve SUDO_USER if exists
		if u := os.Getenv("SUDO_USER"); u != "" {
			usr, err = user.Lookup(u)
			if err != nil {
				return nil, err
			}
		}
	}

	return usr, nil
}

func CheckFile(filename string) (bool, error) {
    _, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return false, err
		} else {
			return false, nil
		}
	}
	return true, nil
}