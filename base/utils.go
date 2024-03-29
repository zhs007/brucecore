package brucecorebase

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// JSON - make json to field
func JSON(key string, obj interface{}) zap.Field {
	s, err := json.Marshal(obj)
	if err != nil {
		return zap.Error(err)
	}

	return zap.String(key, string(s))
}

// BuildLogSubFilename -
func BuildLogSubFilename(appName string, version string) string {
	tm := time.Now()
	str := tm.Format("2006-01-02_15:04:05")
	return fmt.Sprintf("%v.%v.%v", appName, version, str)
}

// BuildLogFilename -
func BuildLogFilename(logtype string, subname string) string {
	return fmt.Sprintf("%v.%v.log", subname, logtype)
}

// AppendString - append string
func AppendString(strs ...string) string {
	var buffer bytes.Buffer

	for _, str := range strs {
		if len(str) > 0 {
			buffer.WriteString(str)
		}
	}

	return buffer.String()
}

// GetCurTime - append string
func GetCurTime() int64 {
	return time.Now().Unix()
}

// HashBuffer - hash buffer
func HashBuffer(buf []byte) string {
	// hash := sha256.New()
	// hash.Write(buf)
	// out := hash.Sum(nil)
	sum256 := sha256.Sum256(buf)
	return hex.EncodeToString(sum256[:])
	// return fmt.Sprintf("%x", out)
}
