package proxy

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/marlin/proxy-server/util"
	log "github.com/sirupsen/logrus"
)

type proxies struct {
	TcpToVsockProxies map[tcptovsock]string
	VsockToTcpProxies map[vsocktotcp]string
	TcpToVsockInstance int
	VsockToTcpInstance int
}

type vsocktotcp struct {
	VsockAddr string
	TcpAddr string
}

type tcptovsock struct {
	TcpAddr string
	VsockAddr string
}

func GetProxyInstance() proxies {
	p := proxies {
		TcpToVsockProxies: make(map[tcptovsock]string),
		VsockToTcpProxies: make(map[vsocktotcp]string),
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
		addrs := tcptovsock {
			TcpAddr: tcpAddr,
			VsockAddr: vsockAddr,
		}
		p.TcpToVsockProxies[addrs] = currentInstance
		return nil
	}
}

func (p *proxies) DestroyTcpToVsock(tcpAddr string, vsockAddr string) error {
	if !(util.IsTcp(tcpAddr) && util.IsVsock(vsockAddr)) {
		return errors.New("not valid address")
	}
	addrs := tcptovsock {
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
		addrs := vsocktotcp {
			TcpAddr: tcpAddr,
			VsockAddr: vsockAddr,
		}
		p.VsockToTcpProxies[addrs] = currentInstance
		return nil
	}
}

func (p *proxies) DestroyVsockToTcp(tcpAddr string, vsockAddr string) error {
	if !(util.IsTcp(tcpAddr) && util.IsVsock(vsockAddr)) {
		return errors.New("not valid address")
	}
	addrs := vsocktotcp {
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
		return nil
	}
}