/*
 * Copyright 2017 caicloud authors. All rights reserved.
 */

package main

import (
	"fmt"

	"github.com/caicloud/golang-template-project/pkg/utils/net"
	"github.com/caicloud/golang-template-project/pkg/version"
)

func main() {
	fmt.Println("Hello world admin")
	fmt.Printf("version '%v', commit '%v', repo root '%v'\n",
		version.VERSION, version.COMMIT, version.REPOROOT)

	net.Helper()
}
