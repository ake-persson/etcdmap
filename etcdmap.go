// Package etcdmap provides methods for interacting with etcd using struct, map or JSON.
package etcdmap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

// Struct returns a struct from a etcd directory.
// !!! This is not supported for nested struct yet.
func Struct(root *client.Node, s interface{}) error {
	// Convert etcd node to map[string]interface{}
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

// JSON returns an etcd directory as JSON []byte.
func JSON(root *client.Node) ([]byte, error) {
	j, err := json.Marshal(Map(root))
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

// JSONIndent returns an etcd directory as indented JSON []byte.
func JSONIndent(root *client.Node, indent string) ([]byte, error) {
	j, err := json.MarshalIndent(Map(root), "", indent)
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

// Map returns a map[string]interface{} from a etcd directory.
func Map(root *client.Node) map[string]interface{} {
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

// Create etcd directory structure from a map, slice or struct.
func Create(kapi client.KeysAPI, path string, val reflect.Value) error {
	switch val.Kind() {
	case reflect.Ptr:
		orig := val.Elem()
		if !orig.IsValid() {
			return nil
		}
		if err := Create(kapi, path, orig); err != nil {
			return err
		}
	case reflect.Interface:
		orig := val.Elem()
		if err := Create(kapi, path, orig); err != nil {
			return err
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			t := val.Type().Field(i)
			k := t.Tag.Get("etcd")
			if err := Create(kapi, path+"/"+k, val.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			if err := Create(kapi, path+"/"+k.String(), v); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			Create(kapi, fmt.Sprintf("%s/%d", path, i), val.Index(i))
		}
	case reflect.String:
		_, err := kapi.Set(context.TODO(), path, val.String(), nil)
		if err != nil {
			return err
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		_, err := kapi.Set(context.TODO(), path, fmt.Sprintf("%v", val.Interface()), nil)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported type: %s for path: %s", val.Kind(), path)
	}

	return nil
}

// CreateJSON etcd directory structure from JSON.
func CreateJSON(kapi client.KeysAPI, dir string, j []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(j, &m); err != nil {
		return err
	}

	return Create(kapi, dir, reflect.ValueOf(m))
}
