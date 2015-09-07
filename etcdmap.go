package etcdmap

import (
	"reflect"
	"strings"

	"encoding/json"
	"github.com/coreos/go-etcd/etcd"
)

// Struct returns a struct from a Etcd directory.
func Struct(root *etcd.Node, s interface{}) error {
	// Convert Etcd node to map[string]interface{}
	m := Map(root)

	// Marshal map[string]interface{} to JSON.
	j, err := json.Marshal(&m)
	if err != nil {
		return err
	}

	// Unmarshal JSON to struct.
	if err := json.Unmarshal(j, &s); err != nil {
		return err
	}

	return nil
}

// CreateStruct creates a Etcd directory based on a struct.
func CreateStruct(client *etcd.Client, dir string, s interface{}) error {
	// Marshal struct to JSON
	j, err := json.Marshal(&s)
	if err != nil {
		return err
	}

	// Unmarshal JSON to map[string]interface{}
	m := make(map[string]interface{})
	if err := json.Unmarshal(j, &m); err != nil {
		return err
	}

	return CreateMap(client, dir, m)
}

// Map returns a map[string]interface{} from a Etcd directory.
func Map(root *etcd.Node) map[string]interface{} {
	v := make(map[string]interface{})

	for _, n := range root.Nodes {
		keys := strings.Split(n.Key, "/")
		k := keys[len(keys)-1]
		if n.Dir {
			v[k] = make(map[string]interface{})
			v[k] = Map(n)
		} else {
			v[k] = n.Value
		}
	}
	return v
}

// CreateMap creates a Etcd directory based on map[string]interface{}.
func CreateMap(client *etcd.Client, dir string, d map[string]interface{}) error {
	for k, v := range d {
		if reflect.ValueOf(v).Kind() == reflect.Map {
			if _, err := client.CreateDir(dir+"/"+k, 0); err != nil {
				return err
			}
			CreateMap(client, dir+"/"+k, v.(map[string]interface{}))
		} else {
			if _, err := client.Set(dir+"/"+k, v.(string), 0); err != nil {
				return err
			}
		}
	}

	return nil
}
