package status

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	pNet "github.com/shirou/gopsutil/v3/net"
	"math"
	"strings"
	"time"
	"unsafe"
)

var cachedFs = make(map[string]struct{})
var timer = 0.0
var prevNetIn uint64
var prevNetOut uint64

func GetStatusFormattedString() string {
	memTotal, memUsed, vmemTotal, VmemUsed := Memory()
	uptime := Uptime()
	cpuLoad := Cpu(2)
	diskTotal, diskUsed := Disk(2)
	return fmt.Sprintf("Uptime: %s\nCPU: %.2f%%\nMem: %s/%s\nSwap: %s/%s\nDisk: %s/%s", uptime, cpuLoad, memUsed, memTotal, VmemUsed, vmemTotal, diskUsed, diskTotal)
}
func Uptime() string {
	bootTime, _ := host.BootTime()
	seconds := uint64(time.Now().Unix()) - bootTime
	duration := time.Duration(seconds) * time.Second
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	if days == 0 {
		return fmt.Sprintf("%dH %dM", hours, minutes)
	}
	return fmt.Sprintf("%dD %dHour %dM", days, hours, minutes)
}
func Memory() (string, string, string, string) {
	memory, _ := mem.VirtualMemory()
	swap, _ := mem.SwapMemory()
	mT := fmt.Sprintf("%.2fGB", float64(memory.Total)/1024.0/1024.0/1024.0)
	mU := fmt.Sprintf("%.2fGB", float64(memory.Used)/1024.0/1024.0/1024.0)
	sT := fmt.Sprintf("%.2fGB", float64(swap.Total)/1024.0/1024.0/1024.0)
	sU := fmt.Sprintf("%.2fGB", float64(swap.Used)/1024.0/1024.0/1024.0)
	return mT, mU, sT, sU
}

func Disk(INTERVAL float64) (string, string) {
	var (
		size, used uint64
	)
	if timer <= 0 {
		diskList, _ := disk.Partitions(false)
		devices := make(map[string]struct{})
		for _, d := range diskList {
			_, ok := devices[d.Device]
			if !ok && checkValidFs(d.Fstype) {
				cachedFs[d.Mountpoint] = struct{}{}
				devices[d.Device] = struct{}{}
			}
		}
		timer = 300.0
	}
	timer -= INTERVAL
	for k := range cachedFs {
		usage, err := disk.Usage(k)
		if err != nil {
			delete(cachedFs, k)
			continue
		}
		size += usage.Total / 1024.0 / 1024.0
		used += usage.Used / 1024.0 / 1024.0
	}
	return fmt.Sprintf("%.2fGB", float64(size)/1024.0), fmt.Sprintf("%.2fGB", float64(used)/1024.0)
}

func Cpu(INTERVAL float64) float64 {
	cpuInfo, _ := cpu.Percent(time.Duration(INTERVAL*float64(time.Second)), false)
	return math.Round(cpuInfo[0]*10) / 10
}

func Traffic(INTERVAL float64) (uint64, uint64, uint64, uint64) {
	var (
		netIn, netOut uint64
	)
	netInfo, _ := pNet.IOCounters(true)
	for _, v := range netInfo {
		if checkInterface(v.Name) {
			netIn += v.BytesRecv
			netOut += v.BytesSent
		}
	}
	rx := uint64(float64(netIn-prevNetIn) / INTERVAL)
	tx := uint64(float64(netOut-prevNetOut) / INTERVAL)
	prevNetIn = netIn
	prevNetOut = netOut
	return netIn, netOut, rx, tx
}

func checkValidFs(name string) bool {
	var validFs = []string{"ext4", "ext3", "ext2", "reiserfs", "jfs", "btrfs", "fuseblk", "zfs", "simfs", "ntfs", "fat32", "exfat", "xfs", "apfs"}
	for _, v := range validFs {
		if strings.ToLower(name) == v {
			return true
		}
	}
	return false
}
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func checkInterface(name string) bool {
	var invalidInterface = []string{"lo", "tun", "kube", "docker", "vmbr", "br-", "vnet", "veth"}
	for _, v := range invalidInterface {
		if strings.Contains(name, v) {
			return false
		}
	}
	return true
}
