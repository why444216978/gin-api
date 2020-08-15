package util

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type generator struct {
	sync.Mutex
}

var (
	g   generator
	mac []byte
	seq int64
)

func init() {
	g = generator{}
	mac = getFirstHardwareAddr()
}

const (
	//http://golang.org/pkg/time
	//A decimal point followed by one or more zeros represents a fractional second, printed to the given number of decimal places. A decimal point followed by one or more nines represents a fractional second, printed to the given number of decimal places, with trailing zeros removed.
	format = "20060102150405.000000" // %Y%m%d%H%M%S%f
)

func getFirstHardwareAddr() []byte {
	is, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	if len(is) == 0 {
		panic("no Hardware Interface")
	}

	for _, v := range is {
		if len([]byte(v.HardwareAddr)) > 0 {
			return []byte(v.HardwareAddr)
		}
	}

	panic("empty Hardware MAC Address")
}

func GenTransNoWithENV(prefix []byte, env string) string {
	if len(prefix) > 2 {
		prefix = prefix[:2]
	}

	dst := make([]byte, 8)
	hex.Encode(dst, mac[:4])

	if strings.ToUpper(env) != ENV_PROD {
		dst = append([]byte("env"+strings.ToLower(env)), dst...)
	}

	return fmt.Sprintf("%s%020s%s%03d", prefix, strings.Replace(time.Now().Format(format), ".", "", 1), dst[:9-len(prefix)], atomic.AddInt64(&seq, 1)%1000)
}

func GenTransNo(prefix []byte) string {
	return GenTransNoWithENV(prefix, ENV_PROD)
}

func GenAccountNo(code int64, currency string, isMerchant bool) string {
	//14 + 1 + 10 + 3 + 5 + 3
	//20171114184840 + B/C + 3301001000 + CNY/YTN + a0a5b + 001
	dst := make([]byte, 6)
	hex.Encode(dst, mac[:3])

	CODEWIDTH := 10
	codeStr := strconv.FormatInt(code, 10)
	if len(codeStr) > CODEWIDTH {
		codeStr = codeStr[:CODEWIDTH]
	} else if len(codeStr) < CODEWIDTH {
		codeStr = fmt.Sprintf("%s%s", codeStr, strings.Repeat("0", CODEWIDTH-len(codeStr)))
	}

	accountType := "C"
	if isMerchant {
		accountType = "B"
	}

	return fmt.Sprintf("%s%s%s%s%s%03d", time.Now().Format(format)[:14], accountType, codeStr, currency, dst[:5], atomic.AddInt64(&seq, 1)%1000)
}
