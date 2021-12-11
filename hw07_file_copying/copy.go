package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrCopyFileFailed        = errors.New("copy file failed")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !fileStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := fileStat.Size()

	if fileSize <= offset {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer file.Close()

	out, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(file)
	_, err = io.CopyN(out, barReader, limit)
	if err != nil {
		return ErrCopyFileFailed
	}
	defer bar.Finish()

	return nil
}
