package main

import (
    "fmt"
    "log"
    "os"
    "syscall"
    "time"
)

func main() {
    defer func() {
        syscall.Exec(os.Args[0], []string{os.Args[0]}, os.Environ())
    }()

    fmt.Println(os.Getpid())
    time.Sleep(1 * time.Second)
    log.Print("Respawning")
}
