package main

import (
	"fmt"

	"github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
)

type User struct {
	Name      string `json:"user"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	u := User{
		Name:      "jdoe",
		FirstName: "John",
		LastName:  "Doe",
	}

	node := []string{fmt.Sprintf("http://%v:%v", "192.168.99.100", "5001")}
	client := etcd.NewClient(node)

	err := etcdmap.CreateStruct(client, "/user/jdoe", &u)
	if err != nil {
		fmt.Println(err.Error)
	}

	res, err := client.Get("/test", true, true)
	if err != nil {
		fmt.Println(err.Error())
	}

	r := User{}
	err2 := etcdmap.Struct(res.Node, &r)
	if err2 != nil {
		fmt.Println(err.Error)
	}

	fmt.Println(r)
}
