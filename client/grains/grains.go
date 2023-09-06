package grains

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"

	"github.com/zcalusic/sysinfo"
)

const (
	ERROR_FAILED_TO_DECODE_STRUCT_TO_JSON = 1
)

type Grains struct {
	NumberOfCPUs    int
	NumberOfCores   int
	Archetecture    string
	OperatingSystem string
	Sysinfo         sysinfo.SysInfo
	RunningAsUID    int
	RunningAsGID    int
}

func NewGrains() *Grains {
	grains := &Grains{}
	return grains
}
func (g Grains) ToJSON() string {

	tmp, err := json.MarshalIndent(g, "  ", "    ")
	// tmp, err := json.Marshal(g)
	if err != nil {
		os.Exit(ERROR_FAILED_TO_DECODE_STRUCT_TO_JSON)
	}
	return string(tmp)
}

func (g *Grains) Update() {
	g.findKernelVersion()
	g.NumberOfCores = runtime.NumCPU()
	g.Archetecture = runtime.GOARCH
	g.RunningAsGID = os.Getgid()
	g.RunningAsUID = os.Getuid()
	// g.Kernel.KernelVersion, _ = kernel.GetKernelVersion()
	// g.Kernel.KernelVersonFull = fmt.Sprintf("%s",
	// 	strings.TrimRight(g.OperatingSystem, "\n"))
	g.Sysinfo.GetSysInfo()
}

func (g *Grains) findKernelVersion() {
	os := runtime.GOOS
	switch os {
	// case "windows":
	//     g.OperatingSystem = "Windows"
	// case "darwin":
	//     g.OperatingSystem = "MAC operating system"
	// case "linux":
	//     g.OperatingSystem = "Linux"
	default:
		g.OperatingSystem = strings.TrimRight(os, "\n")
	}
}
