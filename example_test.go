package etcdmap_test

import (
	"flag"
	"fmt"
	"log"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
)

type User struct {
	Name      string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Group struct {
	Name  string `json:"groupname"`
	Users []User `json:"users"`
}

func main() {
	verbose := flag.Bool("verbose", false, "Verbose")
	node := flag.String("node", "localhost", "Etcd node")
	port := flag.String("port", "5001", "Etcd port")
	flag.Parse()

	// Define nested structure.
	g := Group{
		Name: "staff",
		Users: []User{
			User{
				Name:      "jdoe",
				FirstName: "John",
				LastName:  "Doe",
			},
			User{
				Name:      "lnemoy",
				FirstName: "Leonard",
				LastName:  "Nimoy",
			},
		},
	}

	// Connect to Etcd.
	dbo := []string{fmt.Sprintf("http://%v:%v", *node, *port)}
	if *verbose {
		log.Printf("Connecting to: %s", dbo)
	}
	client := etcd.NewClient(dbo)

	// Create directory structure based on struct.
	err := etcdmap.CreateStruct(client, "/example", g)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get directory structure from Etcd.
	res, err := client.Get("/example", true, true)
	if err != nil {
		log.Fatal(err.Error())
	}

	j, err2 := etcdmap.JSONIndent(res.Node, "    ")
	if err2 != nil {
		log.Fatal(err2.Error())
	}

	fmt.Println(string(j))

	// Output:
	//{
	//    "groupname": "staff",
	//    "users": {
	//        "0": {
	//            "first_name": "John",
	//            "last_name": "Doe",
	//            "username": "jdoe"
	//        },
	//        "1": {
	//            "first_name": "Leonard",
	//            "last_name": "Nimoy",
	//            "username": "lnemoy"
	//        }
	//    }
}


}
