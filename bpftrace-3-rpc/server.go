package main

import (
    "log"
    "net"
    "net/http"
    "net/rpc"
    "os"
)

type Args struct{
    Size int64
    User string
}

type Server int64

func (s *Server) Serve(args *Args, reply *int64) error {
    *reply = args.Size
    return nil
}

func main() {
    port := os.Args[1]
    server := new(Server)
    rpc.Register(server)
    rpc.HandleHTTP()
    listen, e := net.Listen("tcp", ":" + port)
    if e != nil {
        log.Fatal("listen error:", e)
    }
    log.Printf("Listening on port %s", port)
    http.Serve(listen, nil)
}
