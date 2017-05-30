package pfsmod

import (
	"encoding/json"
	"errors"
	"flag"
	log "github.com/golang/glog"
	"os"
	"path/filepath"
	"strconv"
)

const (
	rmCmdName = "rm"
)

//RmResult means Rm command's result
type RmResult struct {
	Path string `json:"path"`
}

//RmCmd means Rm command
type RmCmd struct {
	Method string   `json:"method"`
	R      bool     `json:"r"`
	Args   []string `json:"path"`
}

//LocalCheck check the conditions when running local
func (p *RmCmd) LocalCheck() error {
	if len(p.Args) == 0 {
		return errors.New(StatusText(StatusInvalidArgs))
	}
	return nil
}

//CloudCheck check the conditions when running on cloud
func (p *RmCmd) CloudCheck() error {
	if len(p.Args) == 0 {
		return errors.New(StatusText(StatusInvalidArgs))
	}

	for _, path := range p.Args {
		if !IsCloudPath(path) {
			return errors.New(StatusText(StatusShouldBePfsPath) + ":" + path)
		}

		if !CheckUser(path) {
			return errors.New(StatusText(StatusShouldBePfsPath) + ":" + path)
		}
	}

	return nil
}

//ToURLParam needs not to be implemented
func (p *RmCmd) ToURLParam() string {
	return ""
}

//ToJSON encodes RmCmd to JSON string
func (p *RmCmd) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

//NewRmCmd returns a new RmCmd
func NewRmCmd(r bool, path string) *RmCmd {
	return &RmCmd{
		Method: rmCmdName,
		R:      r,
		Args:   []string{path},
	}
}

//NewRmCmdFromFlag returns a new RmCmd from parsed flags
func NewRmCmdFromFlag(f *flag.FlagSet) (*RmCmd, error) {
	cmd := RmCmd{}

	cmd.Method = rmCmdName
	cmd.Args = make([]string, 0, f.NArg())

	var err error
	f.Visit(func(flag *flag.Flag) {
		if flag.Name == "r" {
			cmd.R, err = strconv.ParseBool(flag.Value.String())
			if err != nil {
				panic(err)
			}
		}
	})

	for _, arg := range f.Args() {
		log.V(2).Info(arg)
		cmd.Args = append(cmd.Args, arg)
	}

	return &cmd, nil
}

//Run runs RmCmd
func (p *RmCmd) Run() (interface{}, error) {
	result := make([]RmResult, 0, 100)
	for _, path := range p.Args {
		list, err := filepath.Glob(path)
		if err != nil {
			return result, err
		}

		for _, arg := range list {
			fi, err := os.Stat(arg)
			if err != nil {
				return result, err
			}

			if fi.IsDir() && !p.R {
				return result, errors.New(StatusText(StatusCannotDelDirectory) + ":" + arg)
			}

			if err := os.RemoveAll(arg); err != nil {
				return result, err
			}

			result = append(result, RmResult{Path: arg})
		}
	}

	return result, nil
}