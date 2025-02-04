package main

import (
    "sync"
    "practica_1_Go/server1"
    "practica_1_Go/server2"
)

func main() {
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        server1.StartServer1()
    }()

    go func() {
        defer wg.Done()
        server2.StartServer2()
    }()

    wg.Wait()
}