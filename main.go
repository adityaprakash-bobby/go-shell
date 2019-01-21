package main

import (
    "fmt"
    "bufio"
    "os"
    "os/exec"
    "os/user"
    // "os/signal"
    // "syscall"
    "io/ioutil"
    "errors"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    for {
        userName, err := user.Current()
        hostName, err := os.Hostname()
        workingDir, err := os.Getwd()

        defaultDirDepth := 3
        folders := strings.Split(workingDir, string(os.PathSeparator))
        shortWorkingDir := ""

        if len(folders) > defaultDirDepth {
            shortWorkingDir = strings.Join([]string{shortWorkingDir, "..."}, "/")
            for i := defaultDirDepth; i > 0; i-- {
                shortWorkingDir = strings.Join([]string{shortWorkingDir, folders[len(folders) - i]}, "/")
            }
        } else {
            shortWorkingDir = workingDir
        }

        fmt.Print(userName.Username + "@" + hostName + "[" + shortWorkingDir + "]$ ")

        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println(os.Stderr, err)
        }

        updateHistory(input)

        if err = execInput(input); err != nil {
            fmt.Println(os.Stderr, err)
        }

    }
}

func execInput(input string) error {

    input = strings.TrimSuffix(input,"\n")
    input = strings.TrimSpace(input)
    args := strings.Split(input, " ")

    // sigc := make(chan os.Signal, 1)
    // signal.Notify(sigc,
    //     syscall.SIGHUP,
    //     syscall.SIGINT,
    //     syscall.SIGTERM,
    //     syscall.SIGQUIT)
    // if sigc != nil {
    //     fmt.Println("Signal caught.")
    // }

    switch args[0] {

    case "cd":
        if  len(args) < 2 {
            return errors.New("Specify a path to change directory.")
        }
        return os.Chdir(args[1])

    case "history":
        readHistory()

    case "exit":
        os.Exit(0)

    }

    cmd := exec.Command(args[0], args[1:] ...)

    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout

    return cmd.Run()
}

func readHistory() {
    file, err := ioutil.ReadFile(".gohistory")

    if err != nil {
        panic("Error!")
    }

    str := string(file)
    fmt.Println(str)
}

func updateHistory(input string) {

    file, err := os.OpenFile(".gohistory", os.O_APPEND|os.O_WRONLY, 0600)

    if err != nil {
        panic("Error!")
    }

    buf := bufio.NewWriter(file)
    buf.WriteString(input)
    buf.Flush()
    file.Close()

}
