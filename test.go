package main

//package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jasonlvhit/gocron"
)

// RRD struct
type RRD struct {
	Name string
}

// Create an RRD
func (r RRD) Create(start int64, DS string, RRA string) error {
	startTime := fmt.Sprintf("%d", start)
	cmd := exec.Command("rrdtool", "create", r.Name, "--start", startTime, DS, RRA)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Update an RRD
func (r RRD) Update() error {
	// i rrdtool update seconds1.rrd \
	//   920805000:000 920805300:300 920805600:600 920805900:900
	cmd := exec.Command("rrdtool", "update", r.Name, fmt.Sprintf("%d"))
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// readFile returns a line-delimited slice of byteslices
func readFile(filename string) ([][]byte, error) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return bytes.Split(input, []byte(`\n`)), nil
}

func getTime() int32 {
	return int32(time.Now().Unix())
}

func not_main() {
	// rrd := RRD{"test.rrd"}
	// rrd.Create(920804400, "DS:seconds:COUNTER:600:U:U", "RRA:AVERAGE:0.5:1:24")
	// rrd.Update()
	memFilename := "/sys/fs/cgroup/memory/user.slice/user-1000.slice/memory.usage_in_bytes"
	func() {
		gocron.Every(2).Seconds().Do(func() error {
			output, err := readFile(memFilename)
			if err != nil {
				return err
			}
			mem, err := strconv.Atoi(strings.TrimSuffix(string(output[0]), "\n"))
			if err != nil {
				return err
			}
			var stat syscall.Statfs_t
			if err := syscall.Statfs("/home/oscar", &stat); err != nil {
				return err
			}
			disk := stat.Bavail * uint64(stat.Bsize)
			fmt.Println(getTime(), mem, disk)
			return nil
		})
		for item := range gocron.Start() {
			fmt.Println(item)
		}
	}()

}
