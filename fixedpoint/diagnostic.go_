package fixedpoint

import (
	"fmt"
	"hash/fnv"
	"runtime"
	"sync"
)

type DiagnosticInfo struct {
	Function string
	File     string
	Line     int
}

var (
	payloadMap   = make(map[diagnostic]DiagnosticInfo)
	payloadMutex sync.Mutex
)

func getDiagnosticInfo(skip int) DiagnosticInfo {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return DiagnosticInfo{"unknown", "unknown", 0}
	}
	fn := runtime.FuncForPC(pc)
	return DiagnosticInfo{fn.Name(), file, line}
}

func hashDiagnosticInfo(diag DiagnosticInfo) diagnostic {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%s:%s:%d", diag.Function, diag.File, diag.Line)))
	return diagnostic(h.Sum64())
}

func encodeDiagnosticInfo(diag DiagnosticInfo) diagnostic {
	payload := hashDiagnosticInfo(diag)

	payloadMutex.Lock()
	defer payloadMutex.Unlock()

	// Store the diagnostic info in the map if not already present
	if _, exists := payloadMap[payload]; !exists {
		payloadMap[payload] = diag
	}
	return payload
}

func DecodePayload(payload diagnostic) (DiagnosticInfo, bool) {
	payloadMutex.Lock()
	defer payloadMutex.Unlock()

	diag, exists := payloadMap[payload]
	return diag, exists
}
