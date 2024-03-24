package redir

import (
    "fmt"
    "os/exec"
    "strconv"
    "strings"
    "net"
    "net/netip"
)

func GetOriginalDestination(conn net.Conn) (destination netip.AddrPort, err error) {
    addr, err := lookup(conn.RemoteAddr().(*net.TCPAddr))
    if err != nil {
        return netip.AddrPort{}, err
    }
    return addr.AddrPort(), nil
}

func lookup(addr *net.TCPAddr) (*net.TCPAddr, error) {
    var new_addr net.TCPAddr
    
    out, err := exec.Command("sudo", "-n", "/sbin/pfctl", "-s", "state").Output()
    for _, line := range strings.Split(string(out), "\n") {
        if strings.Contains(line, "ESTABLISHED:ESTABLISHED") {
            if strings.Contains(line, fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port)) {
                fields := strings.Fields(line)
                if len(fields) > 4 {
                    addr := strings.Split(fields[4], ":")
                    new_addr.IP = net.ParseIP(addr[0])
                    new_addr.Port, _ = strconv.Atoi(addr[1])
                    break
                }
            }
        }
    }

    return &new_addr, err
}
