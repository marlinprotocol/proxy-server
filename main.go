package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/marlin/proxy-server/proxy"
	vsock "github.com/mdlayher/vsock"
	log "github.com/sirupsen/logrus"
)

func main() {
	listenSTDIN()
	// port, err := strconv.Atoi(os.Args[1])
	// cid, err := strconv.Atoi(os.Args[2])

	// if err != nil {
	// 	fmt.Println("Error reading command line args:", err.Error())
	// 	os.Exit(1)
	// }

	// requests := make(chan request)

	// go listen(uint32(cid), uint32(port), requests)
	// go proxyLauncher(requests)
}

func listen(cid uint32, port uint32, reqs chan request) error {

	vsockCon, err := vsock.Dial(cid, port, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	defer vsockCon.Close()			
	buffer := make([]byte, 10240)
	for {
		n, err := vsockCon.Read(buffer)
		if err != nil {
			fmt.Println("Error reading buffer: ", err.Error())
			return err
		}
		if n > 0 {
			requestData := string(buffer[:n])
			req, err := parseRequest(requestData)
			if err != nil {
				log.Error(err)
			} else {
				reqs <- req
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

func proxyLauncher(reqs chan request) {
	proxy := proxy.GetProxyInstance()
	for {
		select {
		case request := <- reqs:
			if request.Type == "tcptovsock" {
				if request.Method == "create" {
					err := proxy.LaunchTcpToVsock(request.TcpAddress, request.VsockAddress)
					if err != nil {
						log.Error(err)
						//send error on socket
					}
				} else {
					err := proxy.DestroyTcpToVsock(request.TcpAddress, request.VsockAddress)
					if err != nil {
						log.Error(err)
						//send error on socket
					}
				}
			} else {
				if request.Method == "create" {
					err := proxy.LaunchVsockToTcp(request.TcpAddress, request.VsockAddress)
					if err != nil {
						log.Error(err)
						//send error on socket
					}
				} else {
					err := proxy.DestroyVsockToTcp(request.TcpAddress, request.VsockAddress)
					if err != nil {
						log.Error(err)
						//send error on socket
					}
				}
			}
		}
	}
}