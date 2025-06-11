package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type Memory struct {
    Used  int `json:"used"`
    Total int `json:"total"`
}

type CPU struct {
    Cores int `json:"cores"`
    Usage int `json:"usage"`
}

type Disk struct {
    Used  int `json:"used"`
    Total int `json:"total"`
}

type VM struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Status string `json:"status"`
    Memory Memory `json:"memory"`
    CPU    CPU    `json:"cpu"`
    Disk   Disk   `json:"disk"`
    Uptime string `json:"uptime"`
    IP     string `json:"ip"`
}

var vms = []VM{
    {
        ID:     100,
        Name:   "prod-web-01",
        Status: "running",
        Memory: Memory{Used: 4096, Total: 8192},
        CPU:    CPU{Cores: 4, Usage: 45},
        Disk:   Disk{Used: 50, Total: 100},
        Uptime: "15d 4h",
        IP:     "10.0.0.10",
    },
    {
        ID:     101,
        Name:   "test-db-01",
        Status: "running",
        Memory: Memory{Used: 2048, Total: 4096},
        CPU:    CPU{Cores: 2, Usage: 22},
        Disk:   Disk{Used: 80, Total: 100},
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
    http.HandleFunc("/api/vms", getVMs)

    log.Println("Server starting on http://localhost:4000")
    if err := http.ListenAndServe(":4000", nil); err != nil {
        log.Fatal(err)
    }
}
