package main

import (
	"fmt"
	"regexp"
	"strconv"
	"os/exec"
//	"reflect"
	"runtime"
	"sync"
	"gopkg.in/qml.v1"
	"os"
)

type ServerStatus struct {
	Name string
	Address string
	Ping int
	Status string
}

var default_items = []ServerStatus{
	{"Singapore", "103.10.124.1", 0, "Unknown"},
	{"EU East (Vienna)", "146.66.155.1", 0, "Unknown"},
	{"EU West (Luxembourg)", "146.66.152.1", 0, "Unknown"},
	{"US East (Sterling)", "208.78.164.1", 0, "Unknown"},
	{"US West (Washington)", "192.69.96.1", 0, "Unknown"},
	{"Australia (Sydney)", "103.10.125.1", 0, "Unknown"},
	{"Sweden (Stockholm)", "146.66.156.1", 0, "Unknown"},
	{"South America (Brazil)", "209.197.29.1", 0, "Unknown"},
}

var linux_latency_regex = regexp.MustCompile("\\d+\\.{0,1}\\d+\\s{0,1}ms")
var windows_latency_regex = regexp.MustCompile("\\d+ms")

var ping_lock sync.Mutex
var wg sync.WaitGroup

func PingServers() {

	ping_lock.Lock()
	defer ping_lock.Unlock()

	for i := range default_items {
		default_items[i].Ping = 0
		wg.Add(1)
		go default_items[i].pingServer()
	}
	

	wg.Wait()
}

func (m *ServerStatus) pingServer() {

	defer wg.Done()

	var cmd *exec.Cmd = nil
	if runtime.GOOS == "windows" {
		cmd = exec.Command("C:\\Windows\\System32\\ping.exe", "-n", 
			"1", m.Address)
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("/bin/ping", "-c", "1", m.Address)
	}

	all_bytes, err := cmd.Output()
	if err != nil {
		m.Ping = 0
		m.Status = "Failed"
		fmt.Println("Could not start ping", m.Address)
		return
	}

	var match []byte
	if runtime.GOOS == "windows" {
		match = windows_latency_regex.Find(all_bytes)
	} else if runtime.GOOS == "linux" {
		match = linux_latency_regex.Find(all_bytes)
	}
	if match == nil {
		m.Ping = 0
		m.Status = "Failed"
		fmt.Println("Could not find match")
		return
	}

	lat, err := strconv.ParseFloat(string(match[0:len(match)-3]),
		32)
	if err != nil {
		m.Ping = 0
		m.Status = "Failed"
		fmt.Println("Failed to parse latency")
		return
	}
	m.Ping = int(lat)
	m.Status = "Success"
}

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	engine := qml.NewEngine()
	servers := &ServerModel{
		list: default_items,
		Len: len(default_items),
	}
	engine.Context().SetVar("servers", servers)
	engine.Context().SetVar("ctrl", &Control{})

	controls, err := engine.LoadFile("main.qml")
	if err != nil {
		return err
	}

	window := controls.CreateWindow(nil)

	window.Show()
	window.Wait()
	return nil
}

type Control struct {
}

type ServerModel struct {
	list []ServerStatus
	Len int
}

func (ctrl *Control) Ping() {
	PingServers()
}

func (model *ServerModel) Sync() {
	for i := range default_items {
		qml.Changed(&model.list[i], &model.list[i].Ping)
		qml.Changed(&model.list[i], &model.list[i].Status)
	}
}

func (model *ServerModel) Server(index int) *ServerStatus {
	return &model.list[index]
}
