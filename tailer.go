import (
	"fmt"
	"github.com/hpcloud/tail"
	"log"
	"sync"
	"strconv"
	"os"
	"strings"
	"time"
)

func main() {
	var filename string = os.Args[1]
	t, err := tail.TailFile(filename, tail.Config{
		Follow: true,
		ReOpen: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	var mutex = &sync.Mutex{}
	mapr := make(map[string]string, 0)
	errctr := make(map[int]int,0)
	var lastkey string
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for _= range ticker.C {
		fmt.Println("================ >>>>>>>>>  <<<<<<<<<< ==================")
		fmt.Printf("%#v\n",errctr)
		fmt.Println("================ >>>>>>>>>  <<<<<<<<<< ==================")
		errctr = make(map[int]int,0)
		}
	}()
	for line := range t.Lines {
		//fmt.Printf("%s",".")
		if line.Text == "" && len(mapr) != 0 {
			if mapr["logo"] == `"RAR"` {
				fmt.Println("====" )
				fmt.Println("MYI: ", mapr["myi"])
				fmt.Println("Complete: ",mapr["complete"])

			tmpsid := strings.Trim(mapr["sid"],`"`)
			tsid,_ := strconv.Atoi(string(tmpsid))
//use strings.EqualFold
			errctr[tsid]++
			mutex.Lock()
			mapr = make(map[string]string, 0)
			mutex.Unlock()
			}
	//fmt.Printf("%#v\n",mapr)
		} else {
			elems := strings.Split(line.Text, "=")
			if len(elems) > 1 {
				lastkey = elems[0]
				mapr[elems[0]] = strings.Join(elems[1:], "=")
			} else {
				mapr[lastkey] += "\n" + line.Text
			}

		}
	}

}
