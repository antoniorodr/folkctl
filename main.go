/*
Copyright © 2026 ANTONIO RODRIGUEZ <kontakt@antoniorodriguez.no>
*/
package main

import (
	"flag"
	"github.com/antoniorodr/folkctl/cmd"
	"github.com/earthboundkid/versioninfo/v2"
)

func main() {
	versioninfo.AddFlag(nil)
	flag.Parse()
	cmd.Execute()
}
