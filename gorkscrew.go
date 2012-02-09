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

  httpResponse, _, _ := bufReader.ReadLine()

  if strings.Contains(string(httpResponse), "200") {
      bufReader.ReadLine()
      go writer(makeReadChan(proxyConn, 1024), os.Stdout)
      writer(makeReadChan(os.Stdin, 1024), proxyConn)
  } else {
    fmt.Printf("Proxy could not open connection to %s\n", destLocation)
    os.Exit(-1)
  }

  os.Exit(0)

}

func makeReadChan(r io.Reader, bufSize int) chan []byte {
  read := make(chan []byte)
  go func() {
    for {
      b := make([]byte, bufSize)
      n, err := r.Read(b)
      if err != nil {
        return
      }
      if n > 0 {
        read <- b[0:n]
      }
    }
  }()
  return read
}

func writer(from chan []byte, to io.Writer) {
  for {
      block := <-from
      to.Write(block)
      /*
      if line, _, err := from.ReadLine(); err != nil {
        break
      } else {
        if _, err = to.Write(line); err != nil {
          break
        } else {
          to.Write([]byte("\n"))
        }
      }*/
  }
}
