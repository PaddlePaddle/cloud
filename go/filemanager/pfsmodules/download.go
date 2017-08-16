package pfsmodules

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/golang/glog"
)

func downloadFile(src string, srcFileSize int64, dst string) error {
	w := FileHandle{}
	if err := w.Open(dst, os.O_RDWR, srcFileSize); err != nil {
		return err
	}
	defer w.Close()

	r := RFileHandle{}
	if err := r.Open(src, os.O_RDONLY, 0); err != nil {
		return err
	}
	defer r.Close()

	offset := int64(0)
	size := defaultChunkSize

	for {
		m, err := r.GetChunkMeta(offset, size)
		if err != nil {
			return err
		}
		c, err := w.ReadChunk(offset, size)
		if err != nil {
			return err
		}
		offset += m.Len

		if m.Checksum == c.Checksum {
			log.V(4).Infof("remote chunk is same as local chunk:%s", c.ToString())
			continue
		}

		c, err = r.ReadChunk(offset, size)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := w.WriteChunk(c); err != nil {
			return err
		}
	}

	if offset != srcFileSize {
		return fmt.Errorf("expect %d but read %d", srcFileSize, offset)
	}

	return nil
}

func checkBeforeDownLoad(src []LsResult, dst string) (bool, error) {
	var bDir bool
	fi, err := os.Stat(dst)
	if err == nil {
		bDir = fi.IsDir()
		if !fi.IsDir() && len(src) > 1 {
			return bDir, errors.New(StatusDestShouldBeDirectory)
		}
	} else if os.IsNotExist(err) {
		return false, nil
	}

	return bDir, err
}

func download(src, dst string) error {
	log.V(1).Infof("download %s to %s\n", src, dst)
	lsRet, err := RemoteLs(NewLsCmd(true, src))
	if err != nil {
		return err
	}

	bDir, err := checkBeforeDownLoad(lsRet, dst)
	if err != nil {
		return err
	}

	for _, attr := range lsRet {
		if attr.IsDir {
			ColorError("Download %s error info:%s\n", StatusOnlySupportFiles)
			return errors.New(StatusOnlySupportFiles)
		}

		realSrc := attr.Path
		realDst := dst

		if bDir {
			_, file := filepath.Split(attr.Path)
			realDst = dst + "/" + file
		}

		if err := downloadFile(realSrc, attr.Size, realDst); err != nil {
			ColorError("Download %s to %s error info:%s\n", realSrc, realDst, err)
			return err
		}
		ColorInfo("Downloaded %s\n", realSrc)
	}

	return nil
}
