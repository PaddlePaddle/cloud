package pfsmod

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	defaultMaxCreateFileSize = int64(4 * 1024 * 1024 * 1024)
)

const (
	touchCmdName = "touch"
)

// TouchResult represents touch-command's result
type TouchResult struct {
	Path string `json:"path"`
}

// TouchCmd is holds touch command's variables
type TouchCmd struct {
	Method   string `json:"method"`
	FileSize int64  `json:"filesize"`
	Path     string `json:"path"`
}

func (p *TouchCmd) checkFileSize() error {
	if p.FileSize < 0 || p.FileSize > defaultMaxCreateFileSize {
		return errors.New(StatusBadFileSize + ":" + fmt.Sprint(p.FileSize))
	}
	return nil
}

// LocalCheck check the conditions when running local
func (p *TouchCmd) LocalCheck() error {
	return p.checkFileSize()
}

// CloudCheck check the conditions when running on cloud
func (p *TouchCmd) CloudCheck() error {
	if !IsCloudPath(p.Path) {
		return errors.New(StatusShouldBePfsPath + ":" + p.Path)
	}

	if !CheckUser(p.Path) {
		return errors.New(StatusShouldBePfsPath + ":" + p.Path)
	}

	return p.checkFileSize()
}

// ToURLParam encodes a TouchCmd to a URL encoding string
func (p *TouchCmd) ToURLParam() string {
	parameters := url.Values{}
	parameters.Add("method", p.Method)
	parameters.Add("path", p.Path)

	str := fmt.Sprint(p.FileSize)
	parameters.Add("path", str)

	return parameters.Encode()
}

// ToJSON encodes a TouchCmd to a JSON string
func (p *TouchCmd) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// NewTouchCmdFromURLParam return a new TouchCmd with specified path
func NewTouchCmdFromURLParam(path string) (*TouchCmd, int32) {
	cmd := TouchCmd{}

	m, err := url.ParseQuery(path)
	if err != nil ||
		len(m["method"]) == 0 ||
		len(m["filesize"]) == 0 ||
		len(m["path"]) == 0 {
		return nil, http.StatusBadRequest
	}

	cmd.Method = m["method"][0]
	if cmd.Method != touchCmdName {
		return nil, http.StatusBadRequest
	}

	cmd.FileSize, err = strconv.ParseInt(m["filesize"][0], 0, 64)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	cmd.Path = m["path"][0]
	if !IsCloudPath(cmd.Path) {
		return nil, http.StatusBadRequest
	}

	return &cmd, http.StatusOK
}

// NewTouchCmd return a new TouchCmd with specified path and fileSize
func NewTouchCmd(path string, fileSize int64) *TouchCmd {
	return &TouchCmd{
		Method:   touchCmdName,
		Path:     path,
		FileSize: fileSize,
	}
}

// CreateSizedFile creates a file with specified size
func CreateSizedFile(path string, size int64) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer Close(fd)

	if size <= 0 {
		return nil
	}

	_, err = fd.Seek(size-1, 0)
	if err != nil {
		return err
	}

	_, err = fd.Write([]byte{0})
	return err
}

// Run is a function runs TouchCmd
func (p *TouchCmd) Run() (interface{}, error) {
	if p.FileSize < 0 || p.FileSize > defaultMaxCreateFileSize {
		return nil, errors.New(StatusBadFileSize)
	}

	fi, err := os.Stat(p.Path)
	if os.IsExist(err) && fi.IsDir() {
		return nil, errors.New(StatusDirectoryAlreadyExist)
	}

	if os.IsNotExist(err) || fi.Size() != p.FileSize {
		if err := CreateSizedFile(p.Path, p.FileSize); err != nil {
			return nil, err
		}
	}

	return &TouchResult{
		Path: p.Path,
	}, nil
}