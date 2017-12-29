package icns_test

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/tmc/icns"
)

func Example() {
	icnsFile, _ := ioutil.ReadFile("testdata/AppIcon.icns")

	icns, err := icns.Parse(bytes.NewReader(icnsFile))
	fmt.Println(err)
	fmt.Println(len(icns))
	fmt.Println(icns[12].Type)

	// output:
	// <nil>
	// 13
	// ic10
}
