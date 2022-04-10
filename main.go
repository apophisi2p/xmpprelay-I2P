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
)

var localAddr *string = flag.String("l", "127.0.0.1:9626", "local XMPP outbound address")
var remoteAddr *string = flag.String("r", "127.0.0.1:4447", "remote I2P proxy address")



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



        go func() {
            defer conn.Close()

            buff := make([]byte, 1024)
            n, err := conn.Read(buff[0:])

            if err != nil {
            	return
            }

            request := string(buff[0:n])


            log.Println("BEGIN Captured Request\n", request)
            log.Println("END Captured Request\n")

///////////////////////////////// See no reason to define a separate function for this

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

            conn2.Write(buff[0:n])

            closer := make(chan struct{}, 2)
            go copy(closer, conn2, conn)
            go copy(closer, conn, conn2)
            <-closer
            log.Println("Connection complete", conn.RemoteAddr())
        }()
    }
}

func copy(closer chan struct{}, dst io.Writer, src io.Reader) {
/// This is ONLY needed to DEBUG. Sends stanzas to Stdout. 
    io.Copy(os.Stdout, io.TeeReader(src, dst))

/// This is for Performance. NO DEBUG Output. Will incorporate into LOGGING system
///    _, _ = io.Copy(dst, src)


    closer <- struct{}{} // connection is closed, send signal to stop proxy
}
