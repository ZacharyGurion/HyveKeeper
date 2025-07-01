package main

import (
    "encoding/json"
    "log"
    "net/http"
    "fmt"
    "os/exec"
    "strconv"
    "strings"
)

type vmMemory struct {
    Used  int `json:"used"`
    Total int `json:"total"`
}

type vmCPU struct {
    Cores int `json:"cores"`
    Usage int `json:"usage"`
}

type vmDisk struct {
    Used  int `json:"used"`
    Total int `json:"total"`
}

type VM struct {
    ID     int      `json:"id"`
    Name   string   `json:"name"`
    Status string   `json:"status"`
    Memory vmMemory `json:"memory"`
    CPU    vmCPU    `json:"cpu"`
    Disk   vmDisk   `json:"disk"`
    Uptime string   `json:"uptime"`
    IP     string   `json:"ip"`
}

type Jail struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Status string `json:"status"`
    Memory string `json:"memory"`
    CPU    string `json:"cpu"`
    Disk   string `json:"disk"`
    Uptime string `json:"uptime"`
    IP     string `json:"ip",omitempty`
    Path   string `json:"path"`
}

type JailInfo struct {
    Jails []Jail `json:"jails"`
    Count int    `json:"count"`
}

func getJailsInfo() (*JailInfo, error) {
    // Execute jls command with detailed output
    cmd := exec.Command("jls", "-v")
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("failed to execute jls: %v", err)
    }

    return parseJlsOutput(string(output))
}

func parseJlsOutput(output string) (*JailInfo, error) {
    lines := strings.Split(strings.TrimSpace(output), "\n")
    var jails []Jail
    var currentJail *Jail

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        if strings.Contains(line, "ID:") {
            if currentJail != nil {
                jails = append(jails, *currentJail)
            }
            currentJail = &Jail{}
        }

        if currentJail == nil {
            continue
        }

        parts := strings.SplitN(line, ":", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        switch key {
        case "ID":
            if jid, err := strconv.Atoi(value); err == nil {
                currentJail.ID = jid
            }
        case "Name":
            currentJail.Name = value
        case "IP":
            currentJail.IP = value
        case "Path":
            currentJail.Path = value
        case "Status":
            currentJail.Status = value
        }
    }

    if currentJail != nil {
        jails = append(jails, *currentJail)
    }

    return &JailInfo{
        Jails: jails,
        Count: len(jails),
    }, nil
}

func getJailsSimple() (*JailInfo, error) {
    cmd := exec.Command("jls", "-h", "jid", "name", "ip4.addr", "host.hostname", "path")
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("failed to execute jls: %v", err)
    }

    return parseJlsSimpleOutput(string(output))
}

func parseJlsSimpleOutput(output string) (*JailInfo, error) {
    lines := strings.Split(strings.TrimSpace(output), "\n")
    var jails []Jail

    for i, line := range lines {
        if i == 0 {
            continue
        }

        fields := strings.Fields(line)
        if len(fields) < 5 {
            continue
        }

        jid, err := strconv.Atoi(fields[0])
        if err != nil {
            continue
        }

        jail := Jail{
            ID:       jid,
            Name:     fields[1],
            IP:       fields[2],
            Path:     fields[4],
            Status:    "running",
        }

        jails = append(jails, jail)
    }

    return &JailInfo{
        Jails: jails,
        Count: len(jails),
    }, nil
}

func jailsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    jailInfo, err := getJailsSimple()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting jail information: %v", err),
        http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(jailInfo)
}

var vms = []VM{
    {
        ID:     100,
        Name:   "prod-web-01",
        Status: "running",
        Memory: vmMemory{Used: 4096, Total: 8192},
        CPU:    vmCPU{Cores: 4, Usage: 45},
        Disk:   vmDisk{Used: 50, Total: 100},
        Uptime: "15d 4h",
        IP:     "10.0.0.10",
    },
    {
        ID:     101,
        Name:   "test-db-01",
        Status: "running",
        Memory: vmMemory{Used: 2048, Total: 4096},
        CPU:    vmCPU{Cores: 2, Usage: 22},
        Disk:   vmDisk{Used: 80, Total: 100},
        Uptime: "7d 12h",
        IP:     "10.0.0.11",
    },
}

func getVMs(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(vms)
}

func main() {
    http.HandleFunc("/vms", getVMs)
    //http.HandleFunc("/jails", jailsHandler)

    log.Println("Server starting on http://localhost:4000")
    if err := http.ListenAndServe(":4000", nil); err != nil {
        log.Fatal(err)
    }
}
