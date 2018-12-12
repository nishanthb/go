package main

import(
	"os/exec"
	"context"
	"strings"
	"fmt"
	"time"
)

func main() { 
res,err := TimedExec(10,"/bin/sleep",30)
if err != nil { 
panic err
}
}

// Timed execute
// Call as 	res, err := timedExec(2, cmd,arg1,arg2...)


func TimedExec(t int, c ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(t)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, c[0], c[1:]...)
	out, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return []byte{}, fmt.Errorf("Timeout: %s", c[0])
	}

	return out, err
}



func TimedExecStr(t int, c string)([]byte,error) {
	fields := strings.Fields(c)
	res,err := TimedExec(t,fields[0], fields[1:])
	return res,err
}
