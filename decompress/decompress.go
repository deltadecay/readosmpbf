package decompress

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"

	"github.com/klauspost/compress/zstd"
	"github.com/pedroalbanese/lzma"
	"github.com/pierrec/lz4"
)

func ZlibData(compressedData []byte, rawSize int32) ([]byte, error) {
	rc, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, errors.New("failed to decompress zlib data")
	}
	defer rc.Close()
	// NewBuffer wants len=0 but desired cap
	databuf := make([]byte, 0, rawSize)
	buf := bytes.NewBuffer(databuf)
	if _, err = buf.ReadFrom(rc); err != nil {
		return nil, errors.New("failed to decompress data: " + err.Error())
	}
	if buf.Len() != int(rawSize) {
		return nil, fmt.Errorf("raw size is %d but expected %d", buf.Len(), rawSize)
	}
	return buf.Bytes(), nil
}

func LzmaData(compressedData []byte, rawSize int32) ([]byte, error) {
	rc := lzma.NewReader(bytes.NewReader(compressedData))
	/*if err != nil {
		return nil, errors.New("failed to decompress lzma data")
	}*/
	defer rc.Close()
	// NewBuffer wants len=0 but desired cap
	databuf := make([]byte, 0, rawSize)
	buf := bytes.NewBuffer(databuf)
	if _, err := buf.ReadFrom(rc); err != nil {
		return nil, errors.New("failed to decompress data: " + err.Error())
	}
	if buf.Len() != int(rawSize) {
		return nil, fmt.Errorf("raw size is %d but expected %d", buf.Len(), rawSize)
	}
	return buf.Bytes(), nil
}

func Lz4Data(compressedData []byte, rawSize int32) ([]byte, error) {
	r := lz4.NewReader(bytes.NewReader(compressedData))
	/*if err != nil {
		return nil, errors.New("failed to decompress lz4 data")
	}*/
	//defer r.Close()
	// NewBuffer wants len=0 but desired cap
	databuf := make([]byte, 0, rawSize)
	buf := bytes.NewBuffer(databuf)
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, errors.New("failed to decompress data: " + err.Error())
	}
	if buf.Len() != int(rawSize) {
		return nil, fmt.Errorf("raw size is %d but expected %d", buf.Len(), rawSize)
	}
	return buf.Bytes(), nil
}

func ZstdData(compressedData []byte, rawSize int32) ([]byte, error) {
	rc, err := zstd.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, errors.New("failed to decompress zstd data")
	}
	defer rc.Close()
	// NewBuffer wants len=0 but desired cap
	databuf := make([]byte, 0, rawSize)
	buf := bytes.NewBuffer(databuf)
	if _, err = buf.ReadFrom(rc); err != nil {
		return nil, errors.New("failed to decompress data: " + err.Error())
	}
	if buf.Len() != int(rawSize) {
		return nil, fmt.Errorf("raw size is %d but expected %d", buf.Len(), rawSize)
	}
	return buf.Bytes(), nil
}
