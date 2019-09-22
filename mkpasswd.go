package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	length   int
	specials int
	count    int
)

var (
	bangs   bool
	colons  bool
	carrots bool
	pipes   bool
	alphas  bool
	caps    bool
	smalls  bool
)

func main() {
	flag.IntVar(&length, "l", 10, "Length of password")
	//flag.IntVar(&specials, "s", int(length/3), "Number of special chars")
	flag.IntVar(&specials, "s", 3, "Number of special chars")
	flag.IntVar(&count, "c", 1, "Count of passwords shown")
	flag.BoolVar(&bangs, "bangs", true, "Include bangs -  !\"#$%&'()*+,-./]")
	flag.BoolVar(&colons, "colons", true, "Include colons - :;<=>?@")
	flag.BoolVar(&carrots, "carrots", true, "Include carrots - \\[]^_`")
	flag.BoolVar(&pipes, "pipes", true, "Include pipes - {|}~")
	flag.BoolVar(&alphas, "alphas", true, "Include alphas - 0-9")
	flag.BoolVar(&caps, "caps", true, "Include caps - A-Z")
	flag.BoolVar(&smalls, "smalls", true, "Include smalls - a-z")
	flag.Usage = func() {
		fmt.Println("Usage: ", os.Args[0], " options, where options are:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Eg: ", os.Args[0], " -nocarrots=false -nocaps=false")
		fmt.Println("Eg: ", os.Args[0], " -s 8 -l 20")
		fmt.Println("Eg: ", os.Args[0], " -s 8 -l 20 -c 30")
		fmt.Println()
	}
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	schars := make([]byte, 0)
	chars := make([]byte, 0)

	// !"#$%&'()*+,-./
	if bangs != false {
		for i := 32; i < 48; i++ {
			v := fmt.Sprintf("%c", i)
			schars = append(schars, v[0])
		}
	}
	// :;<=>?@
	if colons != false {
		for i := 58; i < 65; i++ {
			v := fmt.Sprintf("%c", i)
			schars = append(schars, v[0])
		}
	}
	//[\]^_`
	if carrots != false {
		for i := 91; i < 97; i++ {
			v := fmt.Sprintf("%c", i)
			schars = append(schars, v[0])
		}
	}
	//{|}~
	if pipes != false {
		for i := 123; i < 127; i++ {
			v := fmt.Sprintf("%c", i)
			schars = append(schars, v[0])
		}
	}
	// 0-9
	if alphas != false {
		for i := 48; i < 58; i++ {
			v := fmt.Sprintf("%c", i)
			chars = append(chars, v[0])
		}
	}
	// A-Z
	if caps != false {
		for i := 65; i < 91; i++ {
			v := fmt.Sprintf("%c", i)
			chars = append(chars, v[0])
		}
	}
	//a-z
	if smalls != false {
		for i := 97; i < 123; i++ {
			v := fmt.Sprintf("%c", i)
			chars = append(chars, v[0])
		}
	}

	/*for _, i := range [][]byte{schars32, schars123, schars91, schars58} {
		for _, j := range i {
			schars = append(schars, j)
		}
	}*/

	for c := 0; c < count; c++ {
		var out []byte
		cs := make(chan byte, 1)
		cc := make(chan byte, 1)

		go func(l int, chars []byte, cc chan byte) {
			for i := 0; i < l; i++ {
				cc <- chars[rand.Intn(len(chars))]
			}
		}(length-specials, chars, cc)
		go func(specials int, schars []byte, cs chan byte) {
			for i := 0; i < specials; i++ {
				cs <- schars[rand.Intn(len(schars))]
			}
		}(specials, schars, cs)

		for i := 0; i < length; i++ {
			select {
			case x := <-cc:
				out = append(out, x)
			case y := <-cs:
				out = append(out, y)

			}
		}

		rand.Shuffle(len(out), func(i, j int) {
			out[i], out[j] = out[j], out[i]
		})
		fmt.Printf("%v\n", string(out))
	}
}
