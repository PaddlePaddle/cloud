package paddlecloud

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"

	"github.com/PaddlePaddle/cloud/go/filemanager/pfsmod"
	log "github.com/golang/glog"
	"github.com/google/subcommands"
)

// MkdirCommand represents mkdir command
type MkdirCommand struct {
	//cmd pfsmod.MkdirCmd
}

// Name returns name of MkdirComand
func (*MkdirCommand) Name() string { return "mkdir" }

// Synopsis returns synopsis of MkdirCommand
func (*MkdirCommand) Synopsis() string { return "mkdir directoies on PaddlePaddle Cloud" }

// Usage returns usage of MkdirCommand
func (*MkdirCommand) Usage() string {
	return `mkdir <pfspath>:
	mkdir directories on PaddlePaddleCloud
	Options:
`
}

// SetFlags sets MkdirCommand's parameters
func (p *MkdirCommand) SetFlags(f *flag.FlagSet) {
}

func formatMkdirPrint(results []pfsmod.MkdirResult, err error) {
	if err != nil {
		fmt.Println("\t" + err.Error())
	}
}

// RemoteMkdir creat a directory on cloud
func RemoteMkdir(s *pfsSubmitter, cmd *pfsmod.MkdirCmd) ([]pfsmod.MkdirResult, error) {
	body, err := s.PostFiles(cmd)
	if err != nil {
		return nil, err
	}

	log.V(3).Info(string(body[:]))

	type mkdirResponse struct {
		Err     string               `json:"err"`
		Results []pfsmod.MkdirResult `json:"results"`
	}

	resp := mkdirResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return resp.Results, err
	}

	log.V(1).Infof("%#v\n", resp)

	if len(resp.Err) == 0 {
		return resp.Results, nil
	}

	return resp.Results, errors.New(resp.Err)
}

func remoteMkdir(s *pfsSubmitter, cmd *pfsmod.MkdirCmd) error {
	for _, arg := range cmd.Args {
		subcmd := pfsmod.NewMkdirCmd(arg)

		fmt.Printf("mkdir %s\n", arg)
		results, err := RemoteMkdir(s, subcmd)
		formatMkdirPrint(results, err)
	}
	return nil

}

// Execute runs a MkdirCommand
func (p *MkdirCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() < 1 {
		f.Usage()
		return subcommands.ExitFailure
	}

	cmd, err := pfsmod.NewMkdirCmdFromFlag(f)
	if err != nil {
		return subcommands.ExitFailure
	}
	log.V(1).Infof("%#v\n", cmd)

	s := newPfsCmdSubmitter(UserHomeDir() + "/.paddle/config")
	if err := remoteMkdir(s, cmd); err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}