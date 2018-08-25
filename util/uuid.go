package util

import (
	crand "crypto/rand"
	"fmt"
	mrand "math/rand"
	"strings"
	"sync"
	"time"
)

var lastTime int64
var lastRand int64
var chars = make([]string, 11, 11)
var mu = &sync.Mutex{}

// 64 chars but ordered by ASCII
const base64 string = "0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZ~abcdefghijklmnopqrstuvwxyz"

func toStr(now int64) string {
	// now do the generation (backwards, so we just %64 then /64 along the way)
	for i := 10; i >= 0; i-- {
		index := now % 64
		chars[i] = string(base64[index])
		now = now / 64
	}

	return strings.Join(chars, "")
}

// GenSID is safe to call from different goroutines since it has it's own locking.
func GenSID() string {
	// lock for lastTime, lastRand, and chars
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().UTC().UnixNano()
	var r int64

	// if we have the same time, just inc lastRand, else create a new one
	if now == lastTime {
		lastRand++
	} else {
		lastRand = mrand.Int63()
	}
	r = lastRand

	// remember this for next time
	lastTime = now

	return toStr(now) + toStr(r)
}

// GenUID generate uuid
func GenUID() string {
	// generate 32 bits timestamp
	unix32bits := uint32(time.Now().UTC().Unix())

	buff := make([]byte, 12)

	numRead, err := crand.Read(buff)

	if numRead != len(buff) || err != nil {
		panic(err)
	}

	uuid := fmt.Sprintf("%x%x%x%x%x%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:])
	return uuid
}
