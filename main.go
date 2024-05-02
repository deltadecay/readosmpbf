package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/deltadecay/readosmpbf/pb"
	"google.golang.org/protobuf/proto"
)

func GetBlobData(blob *pb.Blob) []byte {
	var data []byte
	if blob.GetRaw() != nil {
		//log.Println("Raw")
		data = blob.GetRaw()
	} else if blob.GetZlibData() != nil {
		//log.Println("Zlib")
		r, err := zlib.NewReader(bytes.NewReader(blob.GetZlibData()))
		if err != nil {
			log.Fatalln("Failed to zlib decompress data")
		}
		defer r.Close()

		//l := blob.GetRawSize() + bytes.MinRead
		//databuf := make([]byte, 0, l+l/10)
		// NewBuffer wants len=0 but desired cap
		databuf := make([]byte, 0, blob.GetRawSize())

		buf := bytes.NewBuffer(databuf)
		if _, err = buf.ReadFrom(r); err != nil {
			log.Fatalln("Failed to read the decompressed data:", err)
		}
		if buf.Len() != int(blob.GetRawSize()) {
			log.Fatalf("Raw blob size is %d but expected %d", buf.Len(), blob.GetRawSize())
		}
		data = buf.Bytes()

	} else if blob.GetLzmaData() != nil {
		log.Println("Lzma")
	} else if blob.GetLz4Data() != nil {
		log.Println("Lz4")
	} else if blob.GetZstdData() != nil {
		log.Println("Zstd")
	}
	return data
}

func ReadFileBlock(br *bufio.Reader) (*FileBlock, error) {
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
	//log.Println("blobheader size=", blobHeaderSize)

	blobHeaderBytes := make([]byte, blobHeaderSize)

	//n, err = br.Read(blobHeaderBytes)
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
	//log.Println("blobheader type=", blobHeader.GetType())
	//log.Println("blob size=", blobHeader.GetDatasize())

	blobSize := int(blobHeader.GetDatasize())
	if blobSize > 32*1024*1024 {
		return nil, errors.New("unexpected large blob size " + fmt.Sprint(blobSize))
	}

	blobBytes := make([]byte, blobSize)
	//n, err = br.Read(blobBytes)
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

	fb := &FileBlock{
		Type:      blobHeader.GetType(),
		IndexData: blobHeader.GetIndexdata(),
		Data:      GetBlobData(blob),
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
