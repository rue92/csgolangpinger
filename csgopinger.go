package main

import (
	"os/exec"
	"fmt"
	"regexp"
	"strconv"
	"sync"

	_ "image/png"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type ServerStatus struct {
	Name string
	Addr string
	Ping int64
}

type ServerStatusModel struct {
	walk.TableModelBase
	items []*ServerStatus
}

var default_items = []*ServerStatus{
	{"Singapore", "103.10.124.1", 0},
	{"EU East (Vienna)", "146.66.155.1", 0},
	{"EU West (Luxembourg)", "146.66.152.1", 0},
	{"US East (Sterling)", "208.78.164.1", 0},
	{"US West (Washington)", "192.69.96.1", 0},
	{"Australia (Sydney)", "103.10.125.1", 0},
	{"Sweden (Stockholm)", "146.66.156.1", 0},
	{"South America (Brazil)", "209.197.29.1", 0},
}

var latency_regex = regexp.MustCompile("\\d+ms")

var ping_lock sync.Mutex
var wg sync.WaitGroup

var status *walk.StatusBarItem
var mw *walk.MainWindow

func NewServerStatusModel() *ServerStatusModel {
	m := new(ServerStatusModel)
	m.items = default_items
	m.ResetRows()
	return m
}

func (m *ServerStatusModel) PingServers() {

	ping_lock.Lock()
	defer ping_lock.Unlock()

	for i := range m.items {
		m.items[i].Ping = 0
		m.PublishRowChanged(i)
		wg.Add(1)
		go m.pingServer(i)
	}

	wg.Wait()
}

func (m *ServerStatusModel) pingServer(i int) {
	stat := m.items[i]

	cmd := exec.Command("C:\\Windows\\System32\\ping.exe", "-n", "1", 
		stat.Addr)

	all_bytes, err := cmd.Output()
	if err != nil {
		stat.Ping = 0
		fmt.Println("Could not start ping", i)
		return
	}

	match := latency_regex.Find(all_bytes)
	if match == nil {
		stat.Ping = 0
		fmt.Println("Could not find match")
		return
	}
	lat, err := strconv.ParseInt(string(match[0:len(match)-2]),
		10, 64)
	stat.Ping = lat

	if err != nil {
		stat.Ping = 0
		fmt.Println("Failed to parse latency")
		return
	}

	wg.Done()
	m.PublishRowChanged(i)
}

func (m *ServerStatusModel) RowCount() int {
	return len(m.items)
}

func (m *ServerStatusModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Name

	case 1:
		return item.Ping
	}

	panic("Unexpected column")
}

func (m *ServerStatusModel) Len() int {
	return len(m.items)
}

func (m *ServerStatusModel) ResetRows() {
	for i := range m.items {
		m.items[i].Ping = 0
	}

	m.PublishRowsReset()
}

func main() {
	
	walk.SetPanicOnError(true)
	model := NewServerStatusModel()
	
	MainWindow{
		AssignTo: &mw,
		Title: "CS:GO Matchmaking Pings",
		Size: Size{213, 260},
		Layout: VBox{},
		Children: []Widget{
			TableView{
				AlternatingRowBGColor: walk.RGB(255, 255, 224),
				Columns: []TableViewColumn{
					{Title: "Server", Format: "%s",
						Width: 125},
					{Title: "Latency", Format: "%d",
						Alignment: AlignFar,
						Width: 50},
				},
				Model: model,
			},
			PushButton{
				Text: "Ping Servers",
				OnClicked: model.PingServers,
			},
		},
	}
	sb := mw.StatusBar()
	status := walk.NewStatusBarItem()
	status.SetText("Latency Check: ")
	sb.Items().Add(status)
	mw.Run()
}
