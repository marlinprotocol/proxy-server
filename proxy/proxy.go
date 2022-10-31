package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/marlin/proxy-server/util"
	log "github.com/sirupsen/logrus"
)

type proxies struct {
	TcpToVsockProxies map[address]string
	VsockToTcpProxies map[address]string
	TcpToVsockInstance int
	VsockToTcpInstance int
}

type address struct {
	VsockAddr string
	TcpAddr string
}

func GetProxyInstance() proxies {
	p := proxies {
		TcpToVsockProxies: make(map[address]string),
		VsockToTcpProxies: make(map[address]string),
		TcpToVsockInstance: 1,
		VsockToTcpInstance: 1,
	}
	return p
}

func (p *proxies) LaunchTcpToVsock(tcpAddr string, vsockAddr string) error {
	if !(util.IsTcp(tcpAddr) && util.IsVsock(vsockAddr)) {
		return errors.New("not valid address")
	}
	currentInstance := strconv.Itoa(p.TcpToVsockInstance)
	status, err := exec.Command("marlinctl", "proxy", "tcptovsock", "create", "-t " + tcpAddr, "-v " + vsockAddr, "-i " + currentInstance).Output()
	if err != nil {
		log.Error(err)
		return err
	} else {
		fmt.Println(string(status))
		p.TcpToVsockInstance++
		addrs := address {
			TcpAddr: tcpAddr,
			VsockAddr: vsockAddr,
		}
		p.TcpToVsockProxies[addrs] = currentInstance
		err = addEntry("tcptovsock", &addrs, currentInstance)
		if err != nil {
			log.Error("Error adding entry to proxies. Could cause problems in case of server restart. ", err)
		}
		return nil
	}
}

func (p *proxies) DestroyTcpToVsock(tcpAddr string, vsockAddr string) error {
	if !(util.IsTcp(tcpAddr) && util.IsVsock(vsockAddr)) {
		return errors.New("not valid address")
	}
	addrs := address {
		TcpAddr: tcpAddr,
		VsockAddr: vsockAddr,
	}
	instance := p.TcpToVsockProxies[addrs]

	status, err := exec.Command("marlinctl", "proxy", "tcptovsock", "destroy", "-i " + instance).Output()
	if err != nil {
		log.Error(err)
		return err
	} else {
		fmt.Println(string(status))
		delete(p.TcpToVsockProxies, addrs)
		err = removeEntry("tcptovsock", instance)
		if err != nil {
			log.Error("Error removing entry to proxies. Could cause problems in case of server restart. ", err)
		}
		return nil
	}
}

func (p *proxies) LaunchVsockToTcp(tcpAddr string, vsockAddr string) error {
	if !(util.IsTcp(tcpAddr) && util.IsVsock(vsockAddr)) {
		return errors.New("not valid address")
	}
	currentInstance := strconv.Itoa(p.VsockToTcpInstance)
	status, err := exec.Command("marlinctl", "proxy", "vsocktotcp", "create", "-t " + tcpAddr, "-v " + vsockAddr, "-i " + currentInstance).Output()
	if err != nil {
		log.Error(err)
		return err
	} else {
		fmt.Println(string(status))
		p.VsockToTcpInstance++
		addrs := address {
			TcpAddr: tcpAddr,
			VsockAddr: vsockAddr,
		}
		p.VsockToTcpProxies[addrs] = currentInstance
		err = addEntry("vsocktotcp", &addrs, currentInstance)
		if err != nil {
			log.Error("Error adding entry to proxies. Could cause problems in case of server restart. ", err)
		}
		return nil
	}
}

func (p *proxies) DestroyVsockToTcp(tcpAddr string, vsockAddr string) error {
	if !(util.IsTcp(tcpAddr) && util.IsVsock(vsockAddr)) {
		return errors.New("not valid address")
	}
	addrs := address {
		TcpAddr: tcpAddr,
		VsockAddr: vsockAddr,
	}
	instance := p.VsockToTcpProxies[addrs]

	status, err := exec.Command("marlinctl", "proxy", "vsocktotcp", "destroy", "-i " + instance).Output()
	if err != nil {
		log.Error(err)
		return err
	} else {
		fmt.Println(string(status))
		delete(p.VsockToTcpProxies, addrs)
		err = removeEntry("vsocktotcp", instance)
		if err != nil {
			log.Error("Error removing entry to proxies. Could cause problems in case of server restart. ", err)
		}
		return nil
	}
}

type entry struct {
	Type string
	VsockAddr string
	TcpAddr string
	Id string
}

func addEntry(pType string, addr *address, id string) error {
	user, err := util.GetUser()
	if err != nil {
		return err
	}

	dir := "/home/" + user.Username + "/.marlin"
	filename := dir + "/proxies.json"
	err = util.CreateDirPathIfNotExists(dir)
	if err != nil {
		return err
	}

	exist, err := util.CheckFile(filename)
	if err != nil {
		return err
	}
	var entries []entry

	if exist {
		fileData, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(fileData), &entries)
		if err != nil {
			return err
		}
	}

	entries = append(entries, entry{Type: pType, VsockAddr: addr.VsockAddr, TcpAddr: addr.TcpAddr, Id: id})

	result, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, result, 0644)
    if err != nil {
        return err
    }

	return nil
}

func removeEntry(pType string, id string) error {
	user, err := util.GetUser()
	if err != nil {
		return err
	}

	filename := "/home/" + user.Username + "/.marlin/proxies.json"
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var entries []entry

	err = json.Unmarshal([]byte(fileData), &entries)
	if err != nil {
		return err
	}
	
	idx := -1

	for k := range entries {
		if entries[k].Type == pType && entries[k].Id == id {
			idx = k
			break
		}
	}

	if idx != -1 {
		entries = append(entries[:idx], entries[idx+1:]...)
	} else {
		return errors.New("could not find entry")
	}

	result, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, result, 0644)
    if err != nil {
        return err
    }

	return nil
}

func (p *proxies) ResetRunningInstances() error {
	var entries []entry

	user, err := util.GetUser()
	if err != nil {
		return err
	}

	filename := "/home/" + user.Username + "/.marlin/proxies.json"

	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(fileData), &entries)
	if err != nil {
		return err
	}

	for _, e := range entries {
		addr := address {
			TcpAddr: e.TcpAddr,
			VsockAddr: e.VsockAddr,
		}
		if e.Type == "tcptovsock" {
			p.TcpToVsockProxies[addr] = e.Id
		} else if e.Type == "vsocktotcp" {
			p.VsockToTcpProxies[addr] = e.Id
		}
	}
	return nil
}

