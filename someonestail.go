// from stackoverflow
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	file := os.Args[1]
	tail(file, os.Stdout)
}

func tail(file string, out io.Writer) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	info, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	osize := info.Size()
	for {
		for line, prefix, err := r.ReadLine(); err != io.EOF; line, prefix, err = r.ReadLine() {
			if prefix {
				fmt.Fprint(out, string(line))
			} else {
				fmt.Fprintln(out, string(line))
			}
			pos, err := f.Seek(0, io.SeekCurrent)
			if err != nil {
				log.Fatal(err)
			}
			for {
				time.Sleep(time.Second)
				newinfo, err := f.Stat()
				if err != nil {
					log.Fatal(err)
				}

				nsize := newinfo.Size()
				if nsize != osize {
					if nsize < osize {
						f.Seek(0, 0)

					} else {
						f.Seek(pos, io.SeekStart)
					}
					r = bufio.NewReader(f)
					osize = nsize
					break
				}
			}
		}
	}
}
