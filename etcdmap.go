// Package etcdmap provides methods for interacting with Etcd using struct, map or JSON.
package etcdmap

import (
	"fmt"
	"reflect"
	"strings"

	"encoding/json"
	"github.com/coreos/go-etcd/etcd"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))

// Struct returns a struct from a Etcd directory.
// !!! This is not supported for nested struct yet.
func Struct(root *etcd.Node, s interface{}) error {
	// Convert Etcd node to map[string]interface{}
	m := Map(root)

	// Yes this is a hack, so what it works.
	// Marshal map[string]interface{} to JSON.
	j, err := json.Marshal(&m)
	if err != nil {
		return err
	}

	// Yes this is a hack, so what it works.
	// Unmarshal JSON to struct.
	if err := json.Unmarshal(j, &s); err != nil {
		return err
	}

	return nil
}

// JSON returns an Etcd directory as JSON []byte.
func JSON(root *etcd.Node) ([]byte, error) {
	j, err := json.Marshal(Map(root))
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

// JSONIndent returns an Etcd directory as indented JSON []byte.
func JSONIndent(root *etcd.Node, indent string) ([]byte, error) {
	j, err := json.MarshalIndent(Map(root), "", indent)
	if err != nil {
		return []byte{}, err
	}

	return j, nil
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

// Create Etcd directory from a map, slice or struct.
func Create(client *etcd.Client, path string, val reflect.Value) error {

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			t := val.Type().Field(i)
			k := t.Tag.Get("etcd")
			Create(client, path+"/"+k, val.Field(i))
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			Create(client, path+"/"+k.String(), v)
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			Create(client, fmt.Sprintf("%s/%d", path, i), val.Index(i))
		}
	case reflect.String:
		if _, err := client.Set(path, val.String(), 0); err != nil {
			return err
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		if _, err := client.Set(path, fmt.Sprintf("%v", val.Interface()), 0); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported type: %s for path: %s", val.Kind(), path)
	}

	return nil
}
