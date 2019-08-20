package main

# Credit -> https://stackoverflow.com/questions/23031752/start-a-process-in-go-and-detach-from-it


import "fmt"
import "os"
import "syscall"

const (
        UID  = 1000
        GUID = 1000
)

func main() {
        // The Credential fields are used to set UID, GID and attitional GIDS of the process
        // You need to run the program as  root to do this
        var cred = &syscall.Credential{UID, GUID, []uint32{} ,true}
        // the Noctty flag is used to detach the process from parent tty
        var sysproc = &syscall.SysProcAttr{Credential: cred, Noctty: true}
        var attr = os.ProcAttr{
                Dir: ".",
                Env: os.Environ(),
                Files: []*os.File{
                        os.Stdin,
                        nil,
                        nil,
                },
                Sys: sysproc,
        }
        process, err := os.StartProcess("/bin/sleep", []string{"/bin/sleep", "100"}, &attr)
        if err == nil {
                fmt.Printf("%#v\n", process)

                // It is not clear from docs, but Realease actually detaches the process
                err = process.Release()
                if err != nil {
                        fmt.Println(err.Error())
                }

        } else {
                fmt.Println(err.Error())
        }
}
