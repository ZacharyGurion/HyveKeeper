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
  ID      int   		`json:"id"`
  Name    string		`json:"name"`
  Status  string		`json:"status"`
  Memory  vmMemory  `json:"memory"`
  CPU     vmCPU 		`json:"cpu"`
  Disk    vmDisk		`json:"disk"`
  Uptime  string		`json:"uptime"`
  IP      string		`json:"ip"`
}

type Jail struct {
  ID      int    	  `json:"id"`
  Name    string 	  `json:"name"`
  Status  string 	  `json:"status"`
  Memory  int 			`json:"memory"`
  CPU     string 	  `json:"cpu"`
  Disk    string 	  `json:"disk"`
  IP      string 	  `json:"ip"`
  Path    string 	  `json:"path"`
	Rctl	  RctlInfo	`json:"rctl"`
}

type JailsInfo struct {
  Jails []Jail `json:"jails"`
  Count int    `json:"count"`
}

type RctlInfo struct {
  // Core System Resources
  CPUTime       	int64 `json:"cputime"`
  DataSize      	int64 `json:"datasize"`
  StackSize     	int64 `json:"stacksize"`
  CoreDumpSize  	int64 `json:"coredumpsize"`
  
  // Memory Resources
  MemoryUse     	int64 `json:"memoryuse"`
  MemoryLocked  	int64 `json:"memorylocked"`
  VMemoryUse    	int64 `json:"vmemoryuse"`
  SwapUse       	int64 `json:"swapuse"`
  
  // Process and Thread Resources
  MaxProc       	int64 `json:"maxproc"`
  NThr          	int64 `json:"nthr"`
  
  // File System Resources
  OpenFiles     	int64 `json:"openfiles"`
  PseudoTerminals int64 `json:"pseudoterminals"`
  
  // SysV IPC Resources
  MsgQQueued    	int64 `json:"msgqqueued"`
  MsgQSize      	int64 `json:"msgqsize"`
  NMsgQ         	int64 `json:"nmsgq"`
  NSem          	int64 `json:"nsem"`
  NSemOp        	int64 `json:"nsemop"`
  NShm          	int64 `json:"nshm"`
  ShmSize       	int64 `json:"shmsize"`
  
  // Performance Resources
  WallClock     	int64 `json:"wallclock"`
  PCpu          	int64 `json:"pcpu"`
  ReadBps       	int64 `json:"readbps"`
  WriteBps      	int64 `json:"writebps"`
  ReadIOPS      	int64 `json:"readiops"`
  WriteIOPS     	int64 `json:"writeiops"`
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
  
  // Handle preflight requests
  if r.Method == "OPTIONS" {
    w.WriteHeader(http.StatusOK)
    return
  }
}

func getRctl(name string) (*RctlInfo, error){
	data := &RctlInfo{}
	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", fmt.Sprintf("sudo rctl -u jail:%s", name))
	out, err := cmd.Output()
	arr := strings.Split(string(out), "\n")
	if err != nil {
		return nil, fmt.Errorf("failed to execute rctl: %s", err)
	}
	for _, line := range arr {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		valStr := parts[1]

		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("rctl gave invalid value for %s=%s", key, valStr)
		}
		switch key {
			case "cputime":
        data.CPUTime = val
      case "datasize":
        data.DataSize = val
      case "stacksize":
        data.StackSize = val
      case "coredumpsize":
        data.CoreDumpSize = val
      case "memoryuse":
        data.MemoryUse = val
      case "memorylocked":
        data.MemoryLocked = val
      case "vmemoryuse":
        data.VMemoryUse = val
      case "swapuse":
        data.SwapUse = val
      case "maxproc":
        data.MaxProc = val
      case "nthr":
        data.NThr = val
      case "openfiles":
        data.OpenFiles = val
      case "pseudoterminals":
        data.PseudoTerminals = val
      case "msgqqueued":
        data.MsgQQueued = val
      case "msgqsize":
        data.MsgQSize = val
      case "nmsgq":
        data.NMsgQ = val
      case "nsem":
        data.NSem = val
      case "nsemop":
        data.NSemOp = val
      case "nshm":
        data.NShm = val
      case "shmsize":
        data.ShmSize = val
      case "wallclock":
        data.WallClock = val
      case "pcpu":
        data.PCpu = val
      case "readbps":
        data.ReadBps = val
      case "writebps":
        data.WriteBps = val
      case "readiops":
        data.ReadIOPS = val
      case "writeiops":
        data.WriteIOPS = val
      }
		}
	return data, nil
}

func getJailsSimple() (*JailsInfo, error) {
  cmd := exec.Command("jls", "-h", "jid", "name", "ip4.addr", "host.hostname", "path")
  output, err := cmd.Output()
  if err != nil {
    return nil, fmt.Errorf("failed to execute jls: %v", err)
  }

  lines := strings.Split(strings.TrimSpace(string(output)), "\n")
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

		data, err := getRctl(fields[1])

    var cmd *exec.Cmd
		cmd = exec.Command("bash", "-c", fmt.Sprintf("sudo rctl -u jail:%s | grep memoryuse | grep -v vmemory | cut -d'=' -f2 | tr -d '\\n'", fields[1]))
		out, err := cmd.Output()
		mem := -1 
		if err == nil {
			val, err2 := strconv.Atoi(string(out))
			if err2 == nil {
			  mem = val
			}
    }

    jail := Jail{
      ID:	jid,
      Name:	fields[1],
      IP:	fields[2],
      Path:	fields[4],
      Status: "running",
			Memory: mem,
			Rctl:	*data,
    }
    jails = append(jails, jail)
  }

  return &JailsInfo{
    Jails: jails,
    Count: len(jails),
  }, nil
}

var vms = []VM{
  {
    ID: 100,
    Name: "prod-web-01",
    Status: "running",
    Memory: vmMemory{Used: 4096, Total: 8192},
    CPU: vmCPU{Cores: 4, Usage: 45},
    Disk: vmDisk{Used: 50, Total: 100},
    Uptime: "15d 4h",
    IP: "10.0.0.10",
  },
  {
    ID: 101,
    Name: "test-db-01",
    Status: "running",
    Memory: vmMemory{Used: 2048, Total: 4096},
    CPU: vmCPU{Cores: 2, Usage: 22},
    Disk: vmDisk{Used: 80, Total: 100},
    Uptime: "7d 12h",
    IP: "10.0.0.11",
  },
}

func getVMs(w http.ResponseWriter, r *http.Request) {
  enableCORS(w, r)
  if r.Method == "OPTIONS" {
      return
  }
  w.Header().Set("Content-Type", "application/json")
    //w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
  //w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

  json.NewEncoder(w).Encode(vms)
}

func getJails(w http.ResponseWriter, r *http.Request) {
  enableCORS(w, r)
  if r.Method == "OPTIONS" {
    return
  }

  w.Header().Set("Content-Type", "application/json")

  jailInfo, err := getJailsSimple()
  if err != nil {
    http.Error(w, fmt.Sprintf("Error getting jail information: %v", err),
    http.StatusInternalServerError)
    return
  }

  json.NewEncoder(w).Encode(jailInfo)
}

func main() {
  http.HandleFunc("/vms", getVMs)
  http.HandleFunc("/jails", getJails)

  log.Println("Server starting on http://localhost:4000")
  if err := http.ListenAndServe(":4000", nil); err != nil {
    log.Fatal(err)
  }
}
