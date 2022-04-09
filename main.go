package main

import (
    "flag"
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "golang.org/x/net/proxy"
    "regexp"
//    "bufio"
//    "bytes"
//    "encoding/xml"
)

var localAddr *string = flag.String("l", "localhost:9999", "local address")
var remoteAddr *string = flag.String("r", "localhost:80", "remote address")



func main() {
    flag.Parse()
    fmt.Printf("Listening: %v\nProxying: %v\n\n", *localAddr, *remoteAddr)

    listener, err := net.Listen("tcp", *localAddr)
    if err != nil {
        panic(err)
    }
    for {
        conn, err := listener.Accept()
        log.Println("New connection", conn.RemoteAddr())
        if err != nil {
            log.Println("error accepting connection", err)
            continue
        }
//        defer conn.Close()
/// begin lalala
//        scanner := bufio.NewScanner(conn)
//        for scanner.Scan() {
//            message = scanner.Text()
//            fmt.Println("Message Received:", message)
//        }
//        fmt.Printf("END Message \n")
////
//        size := 1
//        buff := make([]byte, 256)
//        reader := bufio.NewReader(conn)
//        for {
//            size, err := reader.ReadByte()
//            if err != nil {
//                return
//            }
//
//        }
//
//        fmt.Printf("size: %x\n", size)
//        fmt.Printf("Received: %x\n", buff[:int(size)])
//        fmt.Fprintf(conn, message)





        go func() {
            defer conn.Close()

///            _, _ = io.Copy(os.Stdout, conn)

            buff := make([]byte, 1024)
            n, err := conn.Read(buff[0:])

            if err != nil {
            	return
            }

            request := string(buff[0:n])


            log.Println("BEGIN Captured Request\n", request)
            log.Println("END Captured Request\n")


/////            var buf bytes.Buffer
/////            io.Copy(&buf, conn)

/////////////////////////////////

            xmppReg := regexp.MustCompile(`to='\S+.i2p'`)
            match := xmppReg.FindString(request)

            domainReg := regexp.MustCompile(`\w+.\w+.i2p`)
            domainTo := domainReg.FindString(match)

            log.Println("Got Domain", domainTo)


/////////////////////////////////


            proxyconn, _ := proxy.SOCKS5("tcp", *remoteAddr, nil, proxy.Direct)
            conn2, err := proxyconn.Dial("tcp", domainTo+":5269")
            if err != nil {
                log.Println("error dialing remote addr", err)
                return
            }
            defer conn2.Close()

//            conn2.Write(buf.Bytes())
            conn2.Write(buff[0:n])
//            conn2.Write([]byte(message))
//            fmt.Fprintf(conn2, message)

            closer := make(chan struct{}, 2)
            go copy(closer, conn2, conn)
            go copy(closer, conn, conn2)
            <-closer
            log.Println("Connection complete", conn.RemoteAddr())
        }()
    }
}

func copy(closer chan struct{}, dst io.Writer, src io.Reader) {
//    copybuf := make([]byte, 1024)
    io.Copy(os.Stdout, io.TeeReader(src, dst))
///    _, _ = io.Copy(dst, src)
    closer <- struct{}{} // connection is closed, send signal to stop proxy
}
