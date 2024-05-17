
# Reading OpenStreetMap PBF files

Test code to read an OpenStreetMap pbf file. See [pbf format](https://wiki.openstreetmap.org/wiki/PBF_Format).

## Generate the pb.go files

To generate the pb files from the proto definitions run:

```
make proto 
```

The run code as 
```
go run .
```

Partial output may look like
```
OSMHeader ------------------------
Writing program = osmium/1.14.0
Source =
Required features =
OsmSchema-V0.6
DenseNodes
Optional features =
Sort.Type_then_ID
OSMData ------------------------
densenodes =  8000
nodes =  0
ways =  0
relations =  0
Dense node id=625025 lat=42.514213 lon=1.552724 tags=[]
Dense node id=625026 lat=42.514453 lon=1.552698 tags=[]
Dense node id=625027 lat=42.515237 lon=1.552978 tags=[]
Dense node id=625028 lat=42.516542 lon=1.553585 tags=[]
Dense node id=625029 lat=42.516894 lon=1.553855 tags=[]
Dense node id=625030 lat=42.517845 lon=1.555245 tags=[]
Dense node id=625032 lat=42.521365 lon=1.558768 tags=[]
Dense node id=625033 lat=42.521943 lon=1.559617 tags=[]
Dense node id=625034 lat=42.522644 lon=1.560765 tags=[]
Dense node id=625035 lat=42.524036 lon=1.562236 tags=[]
Dense node id=625036 lat=42.524601 lon=1.562912 tags=[]
Dense node id=625037 lat=42.524934 lon=1.563941 tags=[]
Dense node id=625039 lat=42.525218 lon=1.564898 tags=[]
Dense node id=625040 lat=42.526162 lon=1.566732 tags=[]
Dense node id=625043 lat=42.527606 lon=1.569090 tags=[]
Dense node id=625050 lat=42.529975 lon=1.572106 tags=[highway=crossing]
Dense node id=625051 lat=42.531152 lon=1.572877 tags=[]
Dense node id=625052 lat=42.532069 lon=1.573862 tags=[]
Dense node id=625053 lat=42.532379 lon=1.574583 tags=[]
Dense node id=625054 lat=42.532591 lon=1.575470 tags=[highway=crossing]
Dense node id=625055 lat=42.532784 lon=1.576523 tags=[]
Dense node id=625056 lat=42.533001 lon=1.578423 tags=[bus=yes|highway=bus_stop|public_transport=stop_position]
Dense node id=625057 lat=42.533386 lon=1.579747 tags=[]
Dense node id=625059 lat=42.534982 lon=1.581263 tags=[]
```
