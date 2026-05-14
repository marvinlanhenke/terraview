package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/marvinlanhenke/terraview/internal/mapper"
	"github.com/marvinlanhenke/terraview/internal/terraform"
)

func main() {
	data, err := os.ReadFile("/home/mlanhenke/dev/projects/terraview/testdata/plans/mixed.json")
	if err != nil {
		panic(err)
	}

	plan, err := terraform.Parse(data)
	if err != nil {
		panic(err)
	}

	tree, err := mapper.BuildTree(plan)
	if err != nil {
		panic(err)
	}

	data, err = json.MarshalIndent(tree, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
