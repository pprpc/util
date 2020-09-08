package common

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	NUmStr  = "0123456789"
	CharStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	SpecStr = "+=-@#~,.[]()!%^*$"
)

var incrementalID uint64

// WaitCtrlC .
func WaitCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChan chan os.Signal
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		endWaiter.Done()
	}()
	endWaiter.Wait()
}

// Round .
func Round(val float64, places int) float64 {
	var t float64
	f := math.Pow10(places)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}
	return t
}

// Uint16ConvertByte .
func Uint16ConvertByte(v uint16) (out []byte) {
	out = make([]byte, 2)
	binary.BigEndian.PutUint16(out, v)
	//binary.LittleEndian.PutUint16(buf, data)
	return

}

// Uint32ConvertByte .
func Uint32ConvertByte(v uint32) (out []byte) {
	out = make([]byte, 4)
	binary.BigEndian.PutUint32(out, v)
	return
}

// ByteConvertUint16 .
func ByteConvertUint16(in []byte) (n uint16) {
	n = binary.BigEndian.Uint16(in)
	return
}

// ByteConvertUint32 .
func ByteConvertUint32(in []byte) (n uint32) {
	n = binary.BigEndian.Uint32(in)
	return
}

// RemoveDuplicatesInt .
func RemoveDuplicatesInt(in []int) (ret []int) {
	sort.Ints(in)
	for i, v := range in {
		if i > 0 && in[i-1] == in[i] {
			continue
		}
		ret = append(ret, v)
	}
	return
}

// RemoveDuplicatesAndEmptyString .
func RemoveDuplicatesAndEmptyString(in []string) (ret []string) {
	sort.Strings(in)
	for i, v := range in {
		if (i > 0 && in[i-1] == in[i]) || len(in[i]) == 0 {
			continue
		}
		ret = append(ret, v)
	}
	return
}

// ValueIsExist .
func ValueIsExist(array interface{}, lv interface{}) bool {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice {
		return false
	}
	n := v.Len()
	for i := 0; i < n; i++ {
		t := v.Index(i).Interface()
		if reflect.TypeOf(t) != reflect.TypeOf(lv) {
			continue
		}
		if lv == t {
			return true
		}
	}
	return false
}

// CallStack .
func CallStack(pre string, skip int) string {
	_, f, l, ok := runtime.Caller(skip) // 2
	if !ok {
		return ""
	}
	return fmt.Sprintf("CallerStack, %s, File: %v, Line: %v ", pre, path.Base(f), l)
}

// GetPort .
func GetPort(n net.Addr) (port string) {
	_t := strings.Split(n.String(), ":")
	port = _t[len(_t)-1]
	return
}

// GetPortInt32 .
func GetPortInt32(n net.Addr) (port int32) {
	_, port = GetRemoteIPPort(n)
	return
}

// GetIPAddr .
func GetIPAddr(n net.Addr) string {
	return strings.Split(n.String(), ":")[0]
}

// GetRemoteIPPort .
func GetRemoteIPPort(n net.Addr) (ip string, port int32) {
	ra := strings.Split(n.String(), ":")
	ip = ra[0]
	t, _ := strconv.ParseInt(ra[len(ra)-1], 10, 32)
	port = int32(t)
	return
}

// LoadFileToByte .
func LoadFileToByte(filePath string) (d []byte, err error) {
	var fd *os.File
	fd, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer fd.Close()
	fi, e := fd.Stat()
	if e != nil {
		err = e
		return
	}
	d = make([]byte, fi.Size())
	_, err = fd.Read(d)

	return
}

// GetIID Get incremental id.
func GetIID() uint64 {
	if incrementalID > 4294836215 {
		incrementalID = 0
	}
	return atomic.AddUint64(&incrementalID, 1)
}

// GetRandPasswd .
func GetRandPasswd(length int, ct string) string {
	if ct == "" {
		ct = "mix"
	}
	rand.Seed(time.Now().UnixNano())

	var passwd []byte = make([]byte, length, length)
	var srcStr string
	if ct == "num" {
		srcStr = NUmStr
	} else if ct == "char" {
		srcStr = CharStr
	} else if ct == "mix" {
		srcStr = fmt.Sprintf("%s%s", NUmStr, CharStr)
	} else if ct == "advance" {
		srcStr = fmt.Sprintf("%s%s%s", NUmStr, CharStr, SpecStr)
	} else {
		srcStr = NUmStr
	}

	for i := 0; i < length; i++ {
		index := rand.Intn(len(srcStr))
		passwd[i] = srcStr[index]
	}
	return string(passwd)
}

// GetIPAddrByName .
func GetIPAddrByName(dev string) (ips []string, err error) {
	intfs, e := net.Interfaces()
	if err != nil {
		err = fmt.Errorf("net.Interfaces(), %s", e)
		return
	}
	for _, row := range intfs {
		if row.Name == dev {
			addrs, e := row.Addrs()
			if e != nil {
				err = e
				return
			}
			for _, addr := range addrs {
				ips = append(ips, strings.Split(addr.String(), "/")[0])
			}
			return
		}
	}
	err = fmt.Errorf("not find dev: %s", dev)
	return
}

// RandRange rand range max value
func RandRange(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}
