package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/marlin/proxy-server/proxy"
	vsock "github.com/mdlayher/vsock"
	log "github.com/sirupsen/logrus"
)

func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Panic("Error reading command line args:", err.Error())
	}

	listner, err := vsock.Listen(uint32(port), nil)
	if err != nil {
		log.Error("listner error ", err)
	} else {
		log.Info("Listening for connections on port ...")
	}
	proxy := proxy.GetProxyInstance()
	err = proxy.ResetRunningInstances()
					
	if err != nil {
		log.Error("Error loading running proxies: ", err)
	}
	for {
		con, err := listner.Accept()
		if err != nil {
			log.Error("Error accepting connections ", err)
		} else {
			log.Info("Enclave connected!")
		}
		buffer := make([]byte, 10240)
		err = con.SetReadDeadline(time.Now().Add(2 * time.Minute))
		if err != nil {
			log.Println("SetReadDeadline failed:", err)
			con.Close()
			continue
		}
		n, err := con.Read(buffer)
		if err != nil {
			log.Error("Error reading buffer: ", err.Error())
			con.Close()
			continue
		}
		if n > 0 {
			requestData := string(buffer[:n])
			// log.Info(requestData)
			request, err := parseRequest(requestData)
			if err != nil {
				con.Close()
				log.Error("", err)
			} else {
				if request.Type == "tcptovsock" {
					if request.Method == "create" {
						err := proxy.LaunchTcpToVsock(request.TcpAddress, request.VsockAddress)
						if err != nil {
							log.Error(err)
							if _, e := con.Write([]byte("ERROR : " + err.Error())); e != nil {
								log.Error("Error send response")
							}
						} else {
							if _, e := con.Write([]byte("SUCCESS")); e != nil {
								log.Error("Error send response")
							}
						}
					} else {
						err := proxy.DestroyTcpToVsock(request.TcpAddress, request.VsockAddress)
						if err != nil {
							log.Error(err)
							if _, e := con.Write([]byte("ERROR : " + err.Error())); e != nil {
								log.Error("Error send response")
							}
						} else {
							if _, e := con.Write([]byte("SUCCESS")); e != nil {
								log.Error("Error send response")
							}
						}
					}
				} else {
					if request.Method == "create" {
						err := proxy.LaunchVsockToTcp(request.TcpAddress, request.VsockAddress)
						if err != nil {
							log.Error(err)
							if _, e := con.Write([]byte("ERROR : " + err.Error())); e != nil {
								log.Error("Error send response")
							}
						} else {
							if _, e := con.Write([]byte("SUCCESS")); e != nil {
								log.Error("Error send response")
							}
						}
					} else {
						err := proxy.DestroyVsockToTcp(request.TcpAddress, request.VsockAddress)
						if err != nil {
							log.Error(err)
							if _, e := con.Write([]byte("ERROR : " + err.Error())); e != nil {
								log.Error("Error send response")
							}
						} else {
							if _, e := con.Write([]byte("SUCCESS")); e != nil {
								log.Error("Error send response")
							}
						}
					}
				
				}
			}
		}
		con.Close()
	}

}

func listen(vsockCon *net.Conn, reqs chan request, ec chan error, m *sync.Mutex) error {		
	buffer := make([]byte, 10240)
	for {
		m.Lock()
		n, err := (*vsockCon).Read(buffer)
		m.Unlock()
		if err != nil {
			log.Error("Error reading buffer: ", err.Error())
			return err
		}
		if n > 0 {
			requestData := string(buffer[:n])
			if strings.Contains(requestData, "CLOSE") {
				ec <- errors.New("CONNECTION CLOSED")
				return nil
			}
			req, err := parseRequest(requestData)
			if err != nil {
				log.Error(err)
			} else {
				reqs <- req
			}
		}
	}
}

func proxyLauncher(vsockCon *net.Conn ,reqs chan request, ec chan error, m *sync.Mutex) {
	proxy := proxy.GetProxyInstance()
	err := proxy.ResetRunningInstances()
	if err != nil {
		log.Error("Error loading running proxies: ", err)
	}
	for {
		select {
		case request := <- reqs:
			if request.Type == "tcptovsock" {
				if request.Method == "create" {
					err := proxy.LaunchTcpToVsock(request.TcpAddress, request.VsockAddress)
					if err != nil {
						log.Error(err)
						m.Lock()
						if _, err := (*vsockCon).Write([]byte("ERROR : " + err.Error())); err != nil {
							log.Error("Error send response")
						}
						m.Unlock()
					} else {
						m.Lock()
						if _, err := (*vsockCon).Write([]byte("SUCCESS")); err != nil {
							log.Error("Error send response")
						}
						m.Unlock()
					}
				} else {
					err := proxy.DestroyTcpToVsock(request.TcpAddress, request.VsockAddress)
					if err != nil {
						log.Error(err)
						m.Lock()
						if _, err := (*vsockCon).Write([]byte("ERROR : " + err.Error())); err != nil {
							log.Error("Error send response")
						}
						m.Unlock()
					} else {
						m.Lock()
						if _, err := (*vsockCon).Write([]byte("SUCCESS")); err != nil {
							log.Error("Error send response")
						}
						m.Unlock()
					}
				}
			} else {
				if request.Method == "create" {
					err := proxy.LaunchVsockToTcp(request.TcpAddress, request.VsockAddress)
					if err != nil {
						log.Error(err)
						m.Lock()
						if _, err := (*vsockCon).Write([]byte("ERROR : " + err.Error())); err != nil {
							log.Error("Error send response")
						}
						m.Unlock()
					} else {
						m.Lock()
						if _, err := (*vsockCon).Write([]byte("SUCCESS")); err != nil {
							log.Error("Error send response")
						}
						m.Unlock()
					}
				} else {
					err := proxy.DestroyVsockToTcp(request.TcpAddress, request.VsockAddress)
					if err != nil {
						log.Error(err)
						m.Lock()
						if _, err := (*vsockCon).Write([]byte("ERROR : " + err.Error())); err != nil {
							log.Error("Error send response")
						}
						m.Unlock()
					} else {
						m.Lock()
						if _, err := (*vsockCon).Write([]byte("SUCCESS")); err != nil {
							log.Error("Error send response")
						}
						m.Unlock()
					}
				}
			}
		}
	}
}

func listenSTDIN() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter proxy type: ")
	text, _ := reader.ReadString('\n')
	proxy := proxy.GetProxyInstance()
	if text == "tcptovsock\n" {
		fmt.Println("Enter TCP address:")
		tcpAddr, _ := reader.ReadString('\n')
		fmt.Println("Enter Vsock address:")
		vsockAddr, _ := reader.ReadString('\n')
		proxy.LaunchTcpToVsock(tcpAddr, vsockAddr)
	} else if text == "vsocktotcp" {
		fmt.Println("Enter TCP address:")
		tcpAddr, _ := reader.ReadString('\n')
		fmt.Println("Enter Vsock address:")
		vsockAddr, _ := reader.ReadString('\n')
		proxy.LaunchVsockToTcp(tcpAddr, vsockAddr)
	} else {
		fmt.Println("Enter valid command!")
	}
}

type request struct {
	Type string
	Method string
	TcpAddress string
	VsockAddress string
}

func parseRequest(input string) (request, error) {
	data := strings.Split(input, ",")
	if len(data) != 4 || (data[0] != "tcptovsock" && data[0] != "vsocktotcp") || (data[1] != "create" && data[1] != "destroy") {
		return request{}, errors.New("invalid request")
	} else {
		return request{
			Type: data[0],
			Method: data[1],
			TcpAddress: data[2],
			VsockAddress: data[3],
		}, nil
	}
} 

