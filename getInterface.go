package main

import(
  "strings"
  "os/exec"
  "fmt"
  "log"
  )

func getdevice() (string, error) {
        path, err := exec.LookPath("netstat")
        if err != nil {
                return "", err
        }
        cmd, err := exec.Command(path, "-nr").Output()
        if err != nil {
                return "", err
        }

        var device string
        for _, i := range strings.Split(string(cmd), "\n") {
                fields := strings.Fields(i)
                if fields[0] == "0.0.0.0" {
                        device = fields[len(fields)-1]
                        break
                }
        }
        if device == "" {

                return "", fmt.Errorf("No devices found")
        }
        return device, nil

}

func main() { 
  device,err := getdevice()
  if err != nil { 
    log.Fatal("Unable to get interface: ",err)
  }
  fmt.Println(device)
}
