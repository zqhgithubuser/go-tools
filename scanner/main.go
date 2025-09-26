package main

import (
    "fmt"
    "log"
    "net/netip"
    "os"
    "os/exec"
    "sync"
    "time"
)

const ping = "ping"

func main() {
    _, err := exec.LookPath(ping)
    if err != nil {
        log.Fatalf("The %s not found in the system PATH", ping)
    }

    if len(os.Args) != 2 {
        log.Fatal("Error: exactly one argument is required, the network CIDR to scan (for example: 192.168.1.0/24)")
    }

    cidr := os.Args[1]
    hostList, err := HostList(cidr)
    if err != nil {
        log.Fatalf("Failed to generate host list: %v", err)
    }

    results := ScanHosts(hostList, 1*time.Second)

    for _, res := range results {
        if res.Reachable {
            fmt.Printf("Host is reachable: %s\n", res.IPAddr)
        }
    }
}

type ScanResult struct {
    IPAddr    netip.Addr
    Reachable bool
}

func ScanHosts(hosts []netip.Addr, timeout time.Duration) []ScanResult {
    var wg sync.WaitGroup
    ch := make(chan ScanResult, len(hosts))

    for _, host := range hosts {
        wg.Add(1)
        go func(ip netip.Addr) {
            defer wg.Done()
            res := ScanResult{
                IPAddr:    ip,
                Reachable: Ping(ip, timeout),
            }
            ch <- res
        }(host)
    }

    go func() {
        wg.Wait()
        close(ch)
    }()

    var results []ScanResult
    for res := range ch {
        results = append(results, res)
    }
    return results
}

func HostList(cidr string) ([]netip.Addr, error) {
    prefix, err := netip.ParsePrefix(cidr)
    if err != nil {
        return nil, fmt.Errorf("invalid CIDR notation: %w", err)
    }

    var hosts []netip.Addr
    var last netip.Addr

    for ip := prefix.Masked().Addr().Next(); prefix.Contains(ip); ip = ip.Next() {
        if last.IsValid() {
            hosts = append(hosts, last)
        }
        last = ip
    }
    return hosts, nil
}

func Ping(ip netip.Addr, timeout time.Duration) bool {
    args := []string{"-c", "1", "-W", fmt.Sprintf("%d", int(timeout.Seconds())), ip.String()}
    cmd := exec.Command(ping, args...)
    err := cmd.Run()
    return err == nil
}
