package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/demdxx/gocast/v2"
	"github.com/duxweb/go-fast/config"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/do"
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/rand"
	"net/url"
	"os"
	"time"
	"unicode"
)

// HashEncode Ciphertext Encryption
func HashEncode(content []byte) string {
	hash, err := bcrypt.GenerateFromPassword(content, bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// HashVerify Ciphertext Verification
func HashVerify(hashedPwd string, password []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, password)
	if err != nil {
		return false
	}
	return true
}

// PageLimit Pagination Calculation
func PageLimit(page int, total int, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		return 0, 0
	}
	totalPage := total / limit
	if total%limit != 0 {
		totalPage++
	}
	offset := (page - 1) * limit
	return offset, totalPage
}

// UcFirst Uppercase the First Letter
func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// LcFirst Lowercase the First Letter
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// Url Compile URL
func Url(urlString string, params map[string]any, absolutes ...bool) string {
	var uri url.URL
	q := uri.Query()
	for k, v := range params {
		q.Add(k, gocast.Str(v))
	}
	urlBuild := urlString + "?" + q.Encode()
	var absolute bool
	if len(absolutes) > 0 {
		absolute = absolutes[0]
	}
	if absolute {
		urlBuild = config.Load("app").GetString("app.baseUrl") + urlBuild
	}
	return urlBuild
}

// FormatFileSize Format File Size
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// RandString Random Characters
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// Md5 Generate 32-bit MD5
func Md5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// FileMd5 File MD5
func FileMd5(data []byte) string {
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

// IsExist Determine if Directory/File Exists
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

// CreateDir Create Directory
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 0777)
	if err != nil {
		do.MustInvoke[*zerolog.Logger](nil).Error().Err(err).Msg(dirName)
		return false
	}
	return true
}

// GetUuid Get UUID
func GetUuid() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// Round Keep Decimal Places
func Round(val float64, precision int) float64 {
	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p+0.5) * math.Pow10(-precision)
	}

	return math.Floor(val*p+0.5) / p
}

// InTimeSpan Range Time Query
func InTimeSpan(start, end, check time.Time, includeStart, includeEnd bool) bool {
	_start := start
	_end := end
	_check := check
	if end.Before(start) {
		_end = end.Add(24 * time.Hour)
		if check.Before(start) {
			_check = check.Add(24 * time.Hour)
		}
	}
	if includeStart {
		_start = _start.Add(-1 * time.Nanosecond)
	}
	if includeEnd {
		_end = _end.Add(1 * time.Nanosecond)
	}
	return _check.After(_start) && _check.Before(_end)
}
