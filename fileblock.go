package main

import (
	"fmt"
	"log"
	"strings"

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
	fmt.Println("OSMHeader ------------------------")
	fmt.Println("Writing program =", headerBlock.GetWritingprogram())
	fmt.Println("Source =", headerBlock.GetSource())

	fmt.Println("Required features =")
	for _, reqFeats := range headerBlock.GetRequiredFeatures() {
		fmt.Println(reqFeats)
	}

	fmt.Println("Optional features =")
	for _, optFeats := range headerBlock.GetOptionalFeatures() {
		fmt.Println(optFeats)
	}
}

func parseDenseNodes(primitiveBlock *pb.PrimitiveBlock, denseNodes *pb.DenseNodes) {

	granularity := int64(primitiveBlock.GetGranularity())
	latoffset := primitiveBlock.GetLatOffset()
	lonoffset := primitiveBlock.GetLonOffset()
	stringTable := primitiveBlock.GetStringtable().GetS()

	var id, lat, lon int64

	//log.Print("ids=", len(denseNodes.GetId()))
	//log.Print("keysvals=", len(denseNodes.GetKeysVals()))
	var tagIdx int = 0
	for i, deltaId := range denseNodes.GetId() {

		id += deltaId

		lat += denseNodes.GetLat()[i]
		lon += denseNodes.GetLon()[i]

		latitude := 1e-9 * float64(latoffset+(granularity*lat))
		longitude := 1e-9 * float64(lonoffset+(granularity*lon))

		tagsStr := ""
		ti := tagIdx
		for {
			if ti == len(denseNodes.GetKeysVals()) {
				break
			}
			keyIdx := denseNodes.GetKeysVals()[ti]
			if keyIdx == 0 {
				ti++
				break
			}
			valIdx := denseNodes.GetKeysVals()[ti+1]

			tagName := string(stringTable[keyIdx])
			tagValue := string(stringTable[valIdx])
			tagsStr += fmt.Sprintf("%s=%s|", tagName, tagValue)
			ti += 2
		}
		tagIdx = ti

		tagsStr = strings.TrimSuffix(tagsStr, "|")
		fmt.Printf("Dense node id=%d lat=%f lon=%f tags=[%s]\n", id, latitude, longitude, tagsStr)
	}
}

func parseNodes(primitiveBlock *pb.PrimitiveBlock, nodes []*pb.Node) {
	granularity := int64(primitiveBlock.GetGranularity())
	latoffset := primitiveBlock.GetLatOffset()
	lonoffset := primitiveBlock.GetLonOffset()
	stringTable := primitiveBlock.GetStringtable().GetS()

	for _, node := range nodes {

		id := node.GetId()
		lat := node.GetLat()
		lon := node.GetLon()
		latitude := 1e-9 * float64(latoffset+(granularity*lat))
		longitude := 1e-9 * float64(lonoffset+(granularity*lon))

		tagsStr := ""
		for ti := 0; ti < len(node.GetKeys()); ti++ {
			keyIdx := node.GetKeys()[ti]
			valIdx := node.GetVals()[ti]
			tagName := string(stringTable[keyIdx])
			tagValue := string(stringTable[valIdx])
			tagsStr += fmt.Sprintf("%s=%s|", tagName, tagValue)
		}
		tagsStr = strings.TrimSuffix(tagsStr, "|")
		fmt.Printf("Node id=%d lat=%f lon=%f tags=[%s]\n", id, latitude, longitude, tagsStr)
	}
}

func parseWays(primitiveBlock *pb.PrimitiveBlock, ways []*pb.Way) {
	//granularity := int64(primitiveBlock.GetGranularity())
	//latoffset := primitiveBlock.GetLatOffset()
	//lonoffset := primitiveBlock.GetLonOffset()
	stringTable := primitiveBlock.GetStringtable().GetS()

	for _, way := range ways {

		id := way.GetId()
		// References to nodes by their id
		var ref int64
		for _, deltaRef := range way.GetRefs() {
			ref += deltaRef
			// TODO lookup node with id == ref
		}
		// if optional_features has LocationsOnWays then ways have delta encoded coordinates
		//lats := way.GetLat()
		//lons := way.GetLon()
		//latitude := 1e-9 * float64(latoffset+(granularity*lat))
		//longitude := 1e-9 * float64(lonoffset+(granularity*lon))

		tagsStr := ""
		for ti := 0; ti < len(way.GetKeys()); ti++ {
			keyIdx := way.GetKeys()[ti]
			valIdx := way.GetVals()[ti]
			tagName := string(stringTable[keyIdx])
			tagValue := string(stringTable[valIdx])
			tagsStr += fmt.Sprintf("%s=%s|", tagName, tagValue)
		}
		tagsStr = strings.TrimSuffix(tagsStr, "|")
		fmt.Printf("Way id=%d tags=[%s]\n", id, tagsStr)
	}
}

func parseRelations(primitiveBlock *pb.PrimitiveBlock, relations []*pb.Relation) {
	//granularity := int64(primitiveBlock.GetGranularity())
	//latoffset := primitiveBlock.GetLatOffset()
	//lonoffset := primitiveBlock.GetLonOffset()
	stringTable := primitiveBlock.GetStringtable().GetS()

	for _, relation := range relations {

		id := relation.GetId()

		tagsStr := ""
		for ti := 0; ti < len(relation.GetKeys()); ti++ {
			keyIdx := relation.GetKeys()[ti]
			valIdx := relation.GetVals()[ti]
			tagName := string(stringTable[keyIdx])
			tagValue := string(stringTable[valIdx])
			tagsStr += fmt.Sprintf("%s=%s|", tagName, tagValue)
		}
		tagsStr = strings.TrimSuffix(tagsStr, "|")
		fmt.Printf("Relation id=%d tags=[%s]\n", id, tagsStr)

		var memId int64
		for mi, deltaMemId := range relation.GetMemids() {
			memId += deltaMemId
			// TODO lookup node/way/relation with memId
			role := string(stringTable[relation.GetRolesSid()[mi]])
			typ := relation.GetTypes()[mi]
			var typstr = ""
			if typ == pb.Relation_NODE {
				// lookup node with memid
				typstr = "node"
			} else if typ == pb.Relation_WAY {
				// lookup way with memid
				typstr = "way"
			} else if typ == pb.Relation_RELATION {
				// lookup way with memid
				typstr = "relation"
			}
			fmt.Printf("\tmember type=%s, memid=%d, role=%s\n", typstr, memId, role)

		}
	}
}

func (fileBlock *FileBlock) ParseOSMData() {

	primitiveBlock := &pb.PrimitiveBlock{}
	if err := proto.Unmarshal(fileBlock.GetData(), primitiveBlock); err != nil {
		log.Fatalln("Failed to parse primitive block in OSMData:", err)
	}
	fmt.Println("OSMData ------------------------")
	//log.Println("DateGranularity =", primitiveBlock.GetDateGranularity())
	//log.Println("Granularity =", primitiveBlock.GetGranularity())
	//log.Println("LatOffset =", primitiveBlock.GetLatOffset())
	//log.Println("LonOffset =", primitiveBlock.GetLonOffset())

	primitiveGroups := primitiveBlock.GetPrimitivegroup()

	for _, pg := range primitiveGroups {
		nodes := pg.GetNodes()
		ways := pg.GetWays()
		relations := pg.GetRelations()
		denseNodes := pg.GetDense()
		fmt.Println("densenodes = ", len(denseNodes.GetId()))
		fmt.Println("nodes = ", len(nodes))
		fmt.Println("ways = ", len(ways))
		fmt.Println("relations = ", len(relations))

		if len(pg.GetDense().GetId()) > 0 {
			parseDenseNodes(primitiveBlock, pg.GetDense())
		} else if len(pg.GetNodes()) > 0 {
			parseNodes(primitiveBlock, pg.GetNodes())
		} else if len(pg.GetWays()) > 0 {
			parseWays(primitiveBlock, pg.GetWays())
		} else if len(pg.GetRelations()) > 0 {
			parseRelations(primitiveBlock, pg.GetRelations())
		}
	}

}
