



proto:
	protoc --go_out=. --go_opt=Mprotos/fileformat.proto=/pb protos/fileformat.proto 
	protoc --go_out=. --go_opt=Mprotos/osmformat.proto=/pb protos/osmformat.proto 

#--go_opt=paths=source_relative
#option go_package = "github.com/deltadecay/readosmpbf/pb";