package xcpprof

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

// ServiceInfo service info
type ServiceInfo struct {
	HeapSys            uint64
	HeapAlloc          uint64
	Gos                int
	StartLen           int64
	TCPCount, UDPCount int32
}

var startTime int64

// NoticePprofCB notice pprof call back.
type NoticePprofCB func(info ServiceInfo)

//
type HealthCheckCB func(w http.ResponseWriter, r *http.Request)

func init() {
	startTime = time.Now().Unix()
}

// StartPprof .
func StartPprof(addr string) {
	// fmt.Sprintf("%s:%d", ipaddr, port)
	http.ListenAndServe(addr, nil)
}

func StartPprofV2(addr string, infocb, errcb HealthCheckCB) {
	// fmt.Sprintf("%s:%d", ipaddr, port)
	if infocb != nil {
		http.HandleFunc("/fdcheck/getinfo", infocb)
	}
	if errcb != nil {
		http.HandleFunc("/fdcheck/geterror", errcb)
	}
	http.ListenAndServe(addr, nil)
}

// ShowSysInfo .
func ShowSysInfo() (info ServiceInfo) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	//sysinfo = fmt.Sprintf("HeapSys: %d KB, HeapAlloc: %d KB, TotalAlloc: %d KB, Alloc: %d KB, goroutines: %d", m.HeapSys / 1024, m.HeapAlloc / 1024, m.TotalAlloc / 1024, m.Alloc / 1024, runtime.NumGoroutine())
	info.HeapSys = m.HeapSys / 1024
	info.HeapAlloc = m.HeapAlloc / 1024
	info.Gos = runtime.NumGoroutine()
	info.StartLen = time.Now().Unix() - startTime
	return
}

// ReportService  report service.
func ReportService(intervalSec int, fn NoticePprofCB) {
	if fn == nil {
		return
	}
	time.AfterFunc(time.Duration(intervalSec)*time.Second, func() {
		fn(ShowSysInfo())
		ReportService(intervalSec, fn)
	})
}
