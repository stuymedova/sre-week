package main

import (
    "log"
    "net"
    "net/http"
    "net/rpc"
    "os"
    "time"
    "errors"
    "strconv"
    "fmt"
)

type Args struct {
    Instant time.Time
}

type Server struct {
    last time.Time
    count int64
    expected int64
    reject bool
    arrived map[time.Time]time.Time
    maxRPS int64
}

func (s *Server) ValidateArrived() {
    maxLag := time.Duration(0)
    for sent, received := range s.arrived {
        if received.Sub(sent) > maxLag {
            maxLag = received.Sub(sent)
        }
    }
    fmt.Printf("Maximum lag: %s\n", maxLag)
    fmt.Printf("Maximum RPS: %d\n", s.maxRPS)
}

func (s *Server) Serve(args *Args, reply *bool) error {
    if s.reject {
        *reply = false
        return errors.New("Serve refused")
    }

    now := time.Now()
    s.arrived[args.Instant] = now

    if time.Since(s.last).Seconds() > 1.0 {
        if s.maxRPS < s.count {
            s.maxRPS = s.count
        }

        log.Printf("%s RPS %d", s.last, s.count)
        s.last = now
        s.count = 1
    } else {
        s.count++
    }

    if int64(len(s.arrived)) == s.expected {
        s.ValidateArrived()
        go func() {
            time.Sleep(5)
            os.Exit(0)
        }()
    }

    *reply = true
    return nil
}

func main() {
    expected, e := strconv.Atoi(os.Args[1])
    if e != nil {
        log.Fatal("parse error:", e)
    }
    before_downtime, e := strconv.Atoi(os.Args[2])
    if e != nil {
        log.Fatal("parse error:", e)
    }
    downtime, e := strconv.Atoi(os.Args[3])
    if e != nil {
        log.Fatal("parse error:", e)
    }
    server := new(Server)
    server.expected = int64(expected)
    server.arrived = make(map[time.Time]time.Time)
    rpc.Register(server)
    rpc.HandleHTTP()
    listen, e := net.Listen("tcp", ":")
    if e != nil {
        log.Fatal("listen error:", e)
    }
    port := listen.Addr().String()
    log.Printf("Listening on addr %s", port)

    go func() {
        time.Sleep(time.Second * time.Duration(before_downtime))
        log.Printf("Close connection")
        server.reject = true

        time.Sleep(time.Second * time.Duration(downtime))
        log.Printf("Reopen connection")
        server.reject = false
    }()

    http.Serve(listen, nil)
}
