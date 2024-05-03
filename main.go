package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/deltadecay/readosmpbf/decompress"
	"github.com/deltadecay/readosmpbf/pb"
	"google.golang.org/protobuf/proto"
)

func GetUncompressedBlobData(blob *pb.Blob) ([]byte, error) {
	if blob.GetRaw() != nil {
		return blob.GetRaw(), nil
	} else if blob.GetZlibData() != nil {
		return decompress.ZlibData(blob.GetZlibData(), blob.GetRawSize())
	} else if blob.GetLzmaData() != nil {
		return decompress.LzmaData(blob.GetLzmaData(), blob.GetRawSize())
	} else if blob.GetLz4Data() != nil {
		return decompress.Lz4Data(blob.GetLz4Data(), blob.GetRawSize())
	} else if blob.GetZstdData() != nil {
		return decompress.ZstdData(blob.GetZstdData(), blob.GetRawSize())
	}
	return nil, errors.New("blob data compressed in unknown/obsolete format")
}

func ReadFileBlock(br *bufio.Reader) (*FileBlock, error) {
	// See https://wiki.openstreetmap.org/wiki/PBF_Format
	blobHeaderSizeBytes := make([]byte, 4)
	n, err := br.Read(blobHeaderSizeBytes)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New("could not read blob header size in the first 4 bytes")
	}
	blobHeaderSize := int(binary.BigEndian.Uint32(blobHeaderSizeBytes))
	if blobHeaderSize > 64*1024 {
		return nil, errors.New("unexpected large blob header size " + fmt.Sprint(blobHeaderSize))
	}

	blobHeaderBytes := make([]byte, blobHeaderSize)

	n, err = io.ReadFull(br, blobHeaderBytes)
	if err != nil {
		return nil, err
	}
	if n != blobHeaderSize {
		return nil, errors.New("could not read complete blob header")
	}
	blobHeader := &pb.BlobHeader{}
	if err := proto.Unmarshal(blobHeaderBytes, blobHeader); err != nil {
		return nil, errors.New("failed to parse blob header: " + err.Error())
	}

	blobSize := int(blobHeader.GetDatasize())
	if blobSize > 32*1024*1024 {
		return nil, errors.New("unexpected large blob size " + fmt.Sprint(blobSize))
	}

	blobBytes := make([]byte, blobSize)
	n, err = io.ReadFull(br, blobBytes)
	if err != nil {
		return nil, err
	}
	if n != blobSize {
		return nil, errors.New("could not read complete blob")
	}
	blob := &pb.Blob{}
	if err := proto.Unmarshal(blobBytes, blob); err != nil {
		return nil, errors.New("failed to parse blob: " + err.Error())
	}
	//log.Println("blob rawsize=", blob.GetRawSize())

	uncompressedData, err := GetUncompressedBlobData(blob)
	if err != nil {
		return nil, err
	}

	fb := &FileBlock{
		Type:      blobHeader.GetType(),
		IndexData: blobHeader.GetIndexdata(),
		Data:      uncompressedData,
	}
	return fb, nil
}

func main() {

	f, err := os.Open("data/andorra-latest.osm.pbf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	br := bufio.NewReader(f)

	for {

		fileBlock, err := ReadFileBlock(br)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
			break
		}
		if fileBlock == nil {
			break
		}
		if fileBlock.GetType() == "OSMHeader" {
			fileBlock.ParseOSMHeader()
		} else if fileBlock.GetType() == "OSMData" {
			fileBlock.ParseOSMData()
		}

	}

}
