package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const n = 1e3

type Circ struct {
	arr (*[n]int)
	i   int
}

func (c Circ) write(v int) Circ {
	c.arr[c.i] = v
	c.i = (c.i + 1) % n
	return c
}

func not_main() {
	t := int(unsafe.Sizeof(0)) * n

	mapFile, err := os.Create("test.dat")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = mapFile.Seek(int64(t-1), 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = mapFile.Write([]byte(" "))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mmap, err := syscall.Mmap(int(mapFile.Fd()), 0, int(t), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	circ := Circ{arr: (*[n]int)(unsafe.Pointer(&mmap[0])), i: 0}

	for i := 0; i < n*2; i++ {
		circ = circ.write(i * i)
	}

	fmt.Println(*circ.arr)

	err = syscall.Munmap(mmap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = mapFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
