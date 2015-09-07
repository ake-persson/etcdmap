
# etcdmap
    import "github.com/mickep76/etcdmap"






## func Map
``` go
func Map(root *etcd.Node) map[string]interface{}
```
Map creates a map[string]interface{} from a Etcd directory.


## func MapCreate
``` go
func MapCreate(client *etcd.Client, dir string, d map[string]interface{}) error
```
MapCreate create Etcd directory structure using a map[string]interface{}.









- - -
