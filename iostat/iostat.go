package iostat

import (
	"encoding/json"
	"os/exec"
)

const iostat = "iostat"

type Output struct {
	Sysstat Sysstat `json:"sysstat"`
}

type AvgCPU struct {
	User   float64 `json:"user"`
	Nice   float64 `json:"nice"`
	System float64 `json:"system"`
	Iowait float64 `json:"iowait"`
	Steal  float64 `json:"steal"`
	Idle   float64 `json:"idle"`
}

type Disk struct {
	DiskDevice string  `json:"disk_device"`
	RS         float64 `json:"r/s"`
	WS         float64 `json:"w/s"`
	RkBS       float64 `json:"rkB/s"`
	WkBS       float64 `json:"wkB/s"`
	RrqmS      float64 `json:"rrqm/s"`
	WrqmS      float64 `json:"wrqm/s"`
	Rrqm       float64 `json:"rrqm"`
	Wrqm       float64 `json:"wrqm"`
	RAwait     float64 `json:"r_await"`
	WAwait     float64 `json:"w_await"`
	AquSz      float64 `json:"aqu-sz"`
	RareqSz    float64 `json:"rareq-sz"`
	WareqSz    float64 `json:"wareq-sz"`
	Svctm      float64 `json:"svctm"`
	Util       float64 `json:"util"`
}

type Statistics struct {
	Timestamp string `json:"timestamp"`
	AvgCPU    AvgCPU `json:"avg-cpu"`
	Disk      []Disk `json:"disk"`
}

type Hosts struct {
	Nodename     string       `json:"nodename"`
	Sysname      string       `json:"sysname"`
	Release      string       `json:"release"`
	Machine      string       `json:"machine"`
	NumberOfCpus int          `json:"number-of-cpus"`
	Date         string       `json:"date"`
	Statistics   []Statistics `json:"statistics"`
}

type Sysstat struct {
	Hosts []Hosts `json:"hosts"`
}

func New() (*Output, error) {
	if _, err := exec.LookPath(iostat); err != nil {
		return nil, err
	}

	cmd := exec.Command(iostat, []string{"-x", "-t", "-o", "JSON"}...)
	cmd.Env = []string{"S_TIME_FORMAT=ISO"}

	var b, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var out = new(Output)
	err = json.Unmarshal(b, out)
	return out, err
}
