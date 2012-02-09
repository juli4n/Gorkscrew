package main

import (
  "net"
  "os"
  "bufio"
  "strings"
  "io"
  "fmt"
)

func main() {

  if len(os.Args) != 3 {
    fmt.Printf("Usage: %s <proxyhost>:<proxyport> <desthost>:<destport>\n", os.Args[0])
    os.Exit(-1)
  }

  proxyLocation := os.Args[1]
  destLocation := os.Args[2]

  var proxyAddress *net.TCPAddr
  var proxyConn *net.TCPConn
  var err os.Error

  if proxyAddress, err = net.ResolveTCPAddr("tcp4", proxyLocation); err != nil {
    fmt.Printf("Error resolving proxy address: %s\n", proxyLocation)
    os.Exit(-1)
  }

  if proxyConn, err = net.DialTCP("tcp", nil, proxyAddress); err != nil {
    fmt.Printf("Error connecting to proxy server: %s\n", proxyAddress)
    os.Exit(-1)
  }

  proxyConn.Write([]byte("CONNECT " + destLocation + " HTTP/1.0\n\n"))

  bufReader := bufio.NewReader(proxyConn)

  // Read HTTP CONNECT response
  httpResponse, _, _ := bufReader.ReadLine()

  if strings.Contains(string(httpResponse), "200") {
    // Drop empty line
    bufReader.ReadLine()
    go writer(bufReader, os.Stdout, 1024)
    writer(os.Stdin, proxyConn, 1024)
  } else {
    fmt.Printf("Proxy could not open connection to %s\n", destLocation)
    os.Exit(-1)
  }

  os.Exit(0)
}

func writer(r io.Reader, to io.Writer, bufSize int) {
  for {
    b := make([]byte, bufSize)
    n, err := r.Read(b)
    if err != nil {
      return
    }
    if n > 0 {
      to.Write(b[0:n])
    }
  }
}
