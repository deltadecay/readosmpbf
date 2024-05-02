package main

import (
	"log"

	"github.com/deltadecay/readosmpbf/pb"
	"google.golang.org/protobuf/proto"
)

type FileBlock struct {
	Type      string
	IndexData []byte
	Data      []byte
}

func (fb *FileBlock) GetType() string {
	if fb != nil {
		return fb.Type
	}
	return ""
}
func (fb *FileBlock) GetData() []byte {
	if fb != nil {
		return fb.Data
	}
	return nil
}

func (fileBlock *FileBlock) ParseOSMHeader() {

	headerBlock := &pb.HeaderBlock{}
	if err := proto.Unmarshal(fileBlock.GetData(), headerBlock); err != nil {
		log.Fatalln("Failed to parse header block in OSMHeader:", err)
	}
	log.Println("OSMHeader ------------------------")
	log.Println("Writing program =", headerBlock.GetWritingprogram())
	log.Println("Source =", headerBlock.GetSource())

	log.Print("Required features =")
	for _, reqFeats := range headerBlock.GetRequiredFeatures() {
		log.Print(" ", reqFeats)
	}
	log.Println()

	log.Print("Optional features =")
	for _, optFeats := range headerBlock.GetOptionalFeatures() {
		log.Print(" ", optFeats)
	}
	log.Println()
}

func (fileBlock *FileBlock) ParseOSMData() {

	primitiveBlock := &pb.PrimitiveBlock{}
	if err := proto.Unmarshal(fileBlock.GetData(), primitiveBlock); err != nil {
		log.Fatalln("Failed to parse primitive block in OSMData:", err)
	}
	log.Println("OSMData ------------------------")
	log.Println("DateGranularity =", primitiveBlock.GetDateGranularity())
	log.Println("Granularity =", primitiveBlock.GetGranularity())
	log.Println("LatOffset =", primitiveBlock.GetLatOffset())
	log.Println("LonOffset =", primitiveBlock.GetLonOffset())

	granularity := int64(primitiveBlock.GetGranularity())
	latoffset := primitiveBlock.GetLatOffset()
	lonoffset := primitiveBlock.GetLonOffset()

	primitiveGroups := primitiveBlock.GetPrimitivegroup()

	for _, pg := range primitiveGroups {
		denseNodes := pg.GetDense()
		log.Println("densenodes = ", len(denseNodes.GetId()))
		var id, lat, lon int64
		for i, deltaId := range denseNodes.GetId() {
			id += deltaId
			lat += denseNodes.GetLat()[i]
			lon += denseNodes.GetLon()[i]

			latitude := 1e-9 * float64(latoffset+(granularity*lat))
			longitude := 1e-9 * float64(lonoffset+(granularity*lon))

			log.Printf("Dense node id=%d lat=%f lon=%f", id, latitude, longitude)
		}

	}

	//stringTable := primitiveBlock.GetStringtable()
	/*for _, s := range stringTable.GetS() {
		log.Println(string(s))
	}*/
}
