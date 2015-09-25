//-----------------------------------------------------------------------------
// This volume driver is meant to be used by docker >= 1.8.x
//
// 1- run the driver:
// sudo docker-volume-rbd
//
// 2- run the container:
// docker run -it --volume-driver rbd -v foo:/foo alpine sh
//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------
// Package membership:
//-----------------------------------------------------------------------------

package main

//-----------------------------------------------------------------------------
// Package factored import statement:
//-----------------------------------------------------------------------------

import (

	// Standard library:
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// Community:
	"github.com/calavera/dkvolume"
)

//-----------------------------------------------------------------------------
// Package constant declarations factored into a block:
//-----------------------------------------------------------------------------

const (
	id            = "rbd"
	socketAddress = "/var/run/docker/plugins/rbd.sock"
)

//-----------------------------------------------------------------------------
// Package variable declarations factored into a block:
//-----------------------------------------------------------------------------

var (

	// Predefined defaults:
	defVolRoot = filepath.Join(dkvolume.DefaultDockerRootDirectory, id)

	// Flags:
	volRoot   = flag.String("volroot", defVolRoot, "Docker volumes root directory")
	defPool   = flag.String("pool", "rbd", "Default Ceph pool for RBD operations")
	defSize   = flag.Int("size", 2048, "Default block device image size")
	defFsType = flag.String("fsType", "xfs", "Default file system type for new images")
)

//-----------------------------------------------------------------------------
// func init() is called after all the variable declarations in the package
// have evaluated their initializers, and those are evaluated only after all
// the imported packages have been initialized:
//-----------------------------------------------------------------------------

func init() {

	// Check for mandatory argc:
	if len(os.Args) < 1 {
		usage()
	}

	// Parse commandline flags:
	flag.Usage = usage
	flag.Parse()
}

//-----------------------------------------------------------------------------
// func usage() reports the correct commandline usage for this driver:
//-----------------------------------------------------------------------------

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

//-----------------------------------------------------------------------------
// Function main of package main:
//-----------------------------------------------------------------------------

func main() {

	// Request handler with a driver implementation
	d := initDriver(*volRoot, *defPool, *defFsType, *defSize)
	h := dkvolume.NewHandler(&d)

	// Listen for requests in a unix socket:
	log.Printf("Listening on %s\n", socketAddress)
	fmt.Println(h.ServeUnix("", socketAddress))
}
