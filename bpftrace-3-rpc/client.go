package main 

import (
    "os"
    "strconv"
    "log"
    "net/rpc"
)

type Args struct {
  Size int64
  User string
}

func main() {
    port := os.Args[1]
    user := os.Args[2]
    size, _ := strconv.ParseInt(os.Args[3], 10, 64)
    args := Args{Size: size, User: user}
    client, err := rpc.DialHTTP("tcp", "localhost:" + port)
    if err != nil {
        log.Fatal("dial error:", err)
    }
    var reply int64
    err = client.Call("Server.Serve", args, &reply)
    if err != nil {
        log.Fatal("Serve error:", err)
    }
    if reply != size {
	log.Fatal("Serve error: expected %d, but got %d", size, reply);
    }
}

