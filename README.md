
# etcdmap
    import "github.com/mickep76/etcdmap"






## func CreateMap
``` go
func CreateMap(client *etcd.Client, dir string, d map[string]interface{}) error
```
CreateMap creates a Etcd directory based on map[string]interface{}.


## func CreateStruct
``` go
func CreateStruct(client *etcd.Client, dir string, s interface{}) error
```
CreateStruct creates a Etcd directory based on a struct.


## func Map
``` go
func Map(root *etcd.Node) map[string]interface{}
```
Map returns a map[string]interface{} from a Etcd directory.


## func Struct
``` go
func Struct(root *etcd.Node, s interface{}) error
```
Struct returns a struct from a Etcd directory.









- - -
