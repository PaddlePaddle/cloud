package paddlecloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PaddlePaddle/cloud/go/filemanager/pfsmod"
	"os"
	"path/filepath"
)

func RemoteChunkMeta(s *PfsSubmitter, path string, chunkSize int64) ([]pfsmod.ChunkMeta, error) {

	cmd := pfsmod.NewChunkMetaCmd(path, chunkSize)

	ret, err := s.GetChunkMeta(cmd)
	if err != nil {
		return nil, err
	}

	resp := pfsmod.ChunkMetaResponse{}
	if err := json.Unmarshal(ret, &resp); err != nil {
		return nil, err
	}

	if len(resp.Err) == 0 {
		return resp.Results, nil
	}

	return resp.Results, errors.New(resp.Err)
}

func DownloadChunks(s *PfsSubmitter, src string, dest string, diffMeta []pfsmod.ChunkMeta) error {
	if len(diffMeta) == 0 {
		fmt.Printf("srcfile:%s and destfile:%s are already same\n", src, dest)
		return nil
	}

	for _, meta := range diffMeta {
		err := s.GetChunkData(pfsmod.NewChunkCmd(src, meta.Offset, meta.Len))
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadFile(s *PfsSubmitter, src string, srcFileSize int64, dst string) error {
	srcMeta, err := RemoteChunkMeta(s, src, pfsmod.DefaultChunkSize)
	if err != nil {
		return err
	}

	dstMeta, err := pfsmod.GetChunkMeta(dst, pfsmod.DefaultChunkSize)
	if err != nil {
		if os.IsNotExist(err) {
			if err := pfsmod.CreateSizedFile(dst, srcFileSize); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	diffMeta, err := pfsmod.GetDiffChunkMeta(srcMeta, dstMeta)
	if err != nil {
		return err
	}

	err = DownloadChunks(s, src, dst, diffMeta)
	if err != nil {
		return err
	}

	return nil
}

// Download files to dst
func Download(s *PfsSubmitter, src, dst string) ([]pfsmod.CpCmdResult, error) {
	lsRet, err := RemoteLs(s, pfsmod.NewLsCmd(true, src))
	if err != nil {
		return nil, err
	}

	if len(lsRet) > 1 {
		fi, err := os.Stat(dst)
		if err != nil {
			if err == os.ErrNotExist {
				os.MkdirAll(dst, 0755)
			} else {
				return nil, err
			}
		}

		if !fi.IsDir() {
			return nil, errors.New(pfsmod.StatusText(pfsmod.StatusDestShouldBeDirectory))
		}
	}

	results := make([]pfsmod.CpCmdResult, 0, 100)
	for _, attr := range lsRet {
		if attr.IsDir {
			return results, errors.New(pfsmod.StatusText(pfsmod.StatusOnlySupportFiles))
		}

		realSrc := attr.Path
		_, file := filepath.Split(attr.Path)
		realDst := dst + "/" + file

		fmt.Printf("download src_path:%s dst_path:%s\n", realSrc, realDst)
		if err := DownloadFile(s, realSrc, attr.Size, realDst); err != nil {
			return results, err
		}

		m := pfsmod.CpCmdResult{
			Src: realSrc,
			Dst: realDst,
		}

		results = append(results, m)
	}

	return results, nil
}