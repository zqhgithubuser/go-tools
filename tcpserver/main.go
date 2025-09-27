package main

import (
    "io"
    "log"
    "net"
)

func main() {
    // 监听连接
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    for {
        // 接受连接
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
        }

        go func(c net.Conn) {
            var bytes []byte
            for {
                buf := make([]byte, 32)
                _, err := c.Read(buf)
                if err != nil {
                    if err == io.EOF {
                        break
                    } else {
                        log.Fatal(err)
                    }
                }
                bytes = append(bytes, buf...)
            }
            log.Print(string(bytes))

            _, err = conn.Write([]byte("Hello from TCP server"))
            if err != nil {
                log.Fatal(err)
            }
            // 关闭连接
            c.Close()
        }(conn)
    }
}
