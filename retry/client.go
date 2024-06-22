package main 

import (
    "os"
    "log"
    "net/rpc"
    "time"
    "strconv"
)

type Args struct {
  Instant time.Time
}

type Packet struct {
    args Args
    response chan error
}

func send(packet Packet, packets chan Packet) bool {
    packets <- packet
    err := <- packet.response
    return err == nil
}

func multisend(packets chan Packet, expected, rps int) {
    response := make(chan error)
    for i:= 0; i < expected; i++ {
        args := Args{Instant: time.Now()}
        packet := Packet{args: args, response: response}
        if !send(packet, packets) {
            go retry(args, packets, rps)
        }
        interval := time.Second / time.Duration(rps)
        time.Sleep(interval)
    }
}

func sender(addr string, packets chan Packet, expected int) {
    var client *rpc.Client
    var err error
    var rps int
    var errorps int
    last := time.Now()
    success := 0
    for packet := <-packets;; packet = <-packets {
        if (time.Since(last).Seconds() > 1.0) {
            log.Printf("RPS %d, Errors %d", rps, errorps)
            rps = 0
            errorps = 0
            last = time.Now()
        }

        rps++
        if client == nil {
            client, err = rpc.DialHTTP("tcp", addr)
            if err != nil {
                client = nil
                packet.response <- err
                errorps++
                continue
            }
        }
        
        var reply bool
        err = client.Call("Server.Serve", packet.args, &reply)
        packet.response <- err
        if err != nil {
            client = nil
            errorps++
            continue
        }
        if reply != true {
            log.Fatal("Serve error");
        }
        success += 1
        if (success == expected) {
            break
        }
    }

    log.Printf("RPS %d, Errors %d", rps, errorps)
}

func main() {
    port := os.Args[1]
    expected, e := strconv.Atoi(os.Args[2])
    if e != nil {
        log.Fatal("parse error:", e)
    }
    rps, e := strconv.Atoi(os.Args[3])
    if e != nil {
        log.Fatal("parse error:", e)
    }

    addr := "localhost:" + port
    log.Printf("Will connect to %s", addr)
    packets := make(chan Packet)

    go multisend(packets, expected, rps)

    sender(addr, packets, expected)
}

