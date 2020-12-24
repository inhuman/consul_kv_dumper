package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/inhuman/consul-kv-mapper"
	"os"
)

func main() {

	if err := realMain(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func realMain() error {
	consulPrefix := flag.String("prefix", "", "consul kv prefix")
	consulHttpAddr := flag.String("addr", "", "consul http address")

	flag.Parse()

	client, err := consulapi.NewClient(&consulapi.Config{
		Address: *consulHttpAddr,
	})

	if err != nil {
		return err
	}

	m, err := consul_kv_mapper.BuildMap(client, *consulPrefix)
	if err != nil {
		return err
	}

	b, err := json.Marshal(&MapWrapper{m})

	fmt.Printf("%s", b)

	return nil
}

type MapWrapper struct {
	*consul_kv_mapper.MapType
}

func (m *MapWrapper) MarshalJSON() ([]byte, error) {

	if m.MapType != nil {
		mm := make(map[consul_kv_mapper.KeyType]interface{})

		mm = get(m.Children)

		return json.Marshal(mm)
	}

	return nil, errors.New("empty")
}

func get(sourceMap map[consul_kv_mapper.KeyType]*consul_kv_mapper.MapType) map[consul_kv_mapper.KeyType]interface{} {

	destMap := make(map[consul_kv_mapper.KeyType]interface{})

	for key, value := range sourceMap {

		if value.Children == nil {
			destMap[key] = value.Value
		} else {
			destMap[key] = get(value.Children)
		}
	}

	return destMap
}
