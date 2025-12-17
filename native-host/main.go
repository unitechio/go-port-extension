package main

import (
  "bytes"
  "encoding/binary"
  "encoding/json"
  "io"
  "os"
  "syscall"

  netstat "github.com/shirou/gopsutil/v3/net"
  "github.com/shirou/gopsutil/v3/process"
)

type Request struct {
  Action string `json:"action"`
  PID int32 `json:"pid,omitempty"`
}

type PortInfo struct {
  Protocol string `json:"protocol"`
  Port uint32 `json:"port"`
  PID int32 `json:"pid"`
  Process string `json:"process"`
  Status string `json:"status"`
}

type Response struct {
  Success bool `json:"success"`
  Data interface{} `json:"data,omitempty"`
  Error string `json:"error,omitempty"`
}

func readMessage() ([]byte, error) {
  var length uint32
  if err := binary.Read(os.Stdin, binary.LittleEndian, &length); err != nil {
    return nil, err
  }
  buf := make([]byte, length)
  _, err := io.ReadFull(os.Stdin, buf)
  return buf, err
}

func writeMessage(v interface{}) error {
  data, _ := json.Marshal(v)
  var buf bytes.Buffer
  binary.Write(&buf, binary.LittleEndian, uint32(len(data)))
  buf.Write(data)
  _, err := os.Stdout.Write(buf.Bytes())
  return err
}

func listPorts() []PortInfo {
  conns, _ := netstat.Connections("inet")
  res := []PortInfo{}
  for _, c := range conns {
    if c.Laddr.Port == 0 || c.Pid == 0 { continue }
    name := "unknown"
    if p, err := process.NewProcess(c.Pid); err == nil {
      if n, err := p.Name(); err == nil { name = n }
    }
    proto := "UNKNOWN"
    if c.Type == 1 { proto = "TCP" }
    if c.Type == 2 { proto = "UDP" }
    res = append(res, PortInfo{
      Protocol: proto,
      Port: c.Laddr.Port,
      PID: c.Pid,
      Process: name,
      Status: c.Status,
    })
  }
  return res
}

func killPID(pid int32) error {
  p, err := os.FindProcess(int(pid))
  if err != nil { return err }
  return p.Signal(syscall.SIGKILL)
}

func main() {
  for {
    msg, err := readMessage()
    if err != nil { return }
    var req Request
    json.Unmarshal(msg, &req)

    switch req.Action {
    case "list_ports":
      writeMessage(Response{Success: true, Data: listPorts()})
    case "kill":
      err := killPID(req.PID)
      if err != nil {
        writeMessage(Response{Success: false, Error: err.Error()})
      } else {
        writeMessage(Response{Success: true})
      }
    }
  }
}