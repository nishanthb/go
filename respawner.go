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
        syscall.Exec(os.Args[0], os.Args, os.Environ())
    }()

    fmt.Println(os.Getpid())
    fmt.Printf("Args are: %s\n", os.Args)
    time.Sleep(1 * time.Second)
    log.Print("Respawning")
}
