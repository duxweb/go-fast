package helper

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/gofrs/uuid"
	"github.com/samber/do/v2"
	"github.com/spf13/afero"
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/rand"
	"os"
	"time"
	"unicode"
)

// HashEncode Ciphertext Encryption
func HashEncode(content string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(content), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// HashVerify Ciphertext Verification
func HashVerify(hashedPwd string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func Encryption(str string, keys ...string) (string, error) {
	var key string

	if len(keys) == 1 {
		key = keys[0]
	}

	if key == "" {
		key = config.Load("use").GetString("app.secret")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(str))
	block.Encrypt(ciphertext, []byte(str))
	return hex.EncodeToString(ciphertext), nil
}

func Decryption(str string, keys ...string) (string, error) {
	var key string

	if len(keys) == 1 {
		key = keys[0]
	}

	if key == "" {
		key = config.Load("use").GetString("app.secret")
	}

	ciphertext, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(ciphertext))
	block.Decrypt(plaintext, ciphertext)

	return string(plaintext), nil
}

// PageLimit Calculation
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

// HumanFileSize Format File Size
func HumanFileSize(fileSize int64) (size string) {
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

func CreateDir(dirs ...string) {
	fs := do.MustInvokeNamed[afero.Fs](global.Injector, "os.fs")
	for _, path := range dirs {
		exists, _ := afero.DirExists(fs, path)
		if exists {
			return
		}
		err := fs.MkdirAll(path, 0777)
		if err != nil {
			panic("failed to create " + path + " directory")
		}
	}
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
