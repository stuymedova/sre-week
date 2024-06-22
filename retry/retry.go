package main

import (
    "time"
)

func retry(args Args, packets chan Packet, rps int) {
    response := make(chan error)
    packet := Packet{args: args, response: response}
    timeout := time.Second
    for {
        time.Sleep(timeout)
        if send(packet, packets) {
            break
        }
    }
}
