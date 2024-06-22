package main 

import (
    "os"
    "strconv"
    "log"
    "math/rand"
    "time"
)

func write(filename string, size int) {
    data := make([]byte, size)
    for i := 0; i < size; i++ {
        data[i] = byte(rand.Int())
    }
    f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        log.Fatal(err)
    }
    if _, err := f.Write(data); err != nil {
        log.Fatal(err)
    }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}

func background(filename string, size int) {
    for {
        write(filename, size)     
        time.Sleep(2 * time.Second)
    }
}

func foreground(filename string, size int) {
    for {
        write(filename, size)
        time.Sleep(time.Second)
    }
}

func main() {
    filename := os.Args[1]
    size, err := strconv.Atoi(os.Args[2])
    if err != nil {
        log.Fatal("Usage: ./main filename size")
    }

    go background(filename, size)

    foreground(filename, size)
}

