package utils

import (
	"github.com/Microsoft/go-winio"
	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows/registry"
	"strings"
	"syscall"
)

const afHvSock = 34 // AF_HYPERV

func CheckHVService() bool {
	gcs, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Virtualization\GuestCommunicationServices`, registry.READ)
	if err != nil {
		return false
	}
	defer gcs.Close()

	agentSrvGUID := winio.VsockServiceID(ServicePort)
	agentSrv, err := registry.OpenKey(gcs, agentSrvGUID.String(), registry.READ)
	if err != nil {
		return false
	}
	agentSrv.Close()
	return true
}

func GetVMID() []string {
	type Win32_Process struct {
		CommandLine string
	}
	var processes []Win32_Process
	q := wmi.CreateQuery(&processes, "WHERE Name='wslhost.exe'")
	err := wmi.Query(q, &processes)
	if err != nil {
		return nil
	}

	guids := make(map[string]interface{})
	for _, v := range processes {
		args := strings.Split(v.CommandLine, " ")
		for i := len(args) - 1; i >= 0; i-- {
			if strings.Contains(args[i], "{") {
				guids[args[i]] = 0
				break
			}
		}
	}
	results := make([]string, 0)
	for k, _ := range guids {
		results = append(results, k[1:len(k)-1])
	}
	return results
}

func CheckHvSocket() bool {
	fd, err := syscall.Socket(afHvSock, syscall.SOCK_STREAM, 1)
	if err != nil {
		println(err.Error())
		return false
	}
	syscall.Close(fd)
	return true
}
