package main 

import (
  "net"
  "log"
  "os"
  "fmt"
  "strings"
  "bufio"
)

func main() {
  args := os.Args[1:]
  
  if len(args) < 1 {
    fmt.Println("sage: go run main.go [port]")
    os.Exit(1)
  }

  port := args[0]

  li, err := net.Listen("tcp", ":" + port)

  if err != nil {
    log.Fatalln(err.Error())
  } else {
    fmt.Println("Local server listening on :" + port)
  }
  defer li.Close()

  for {
    connection, err := li.Accept()
    if err != nil {
      log.Println(err.Error())
      continue
    }
    defer connection.Close()

    go handleRequest(connection)
  }
}

func handleRequest(connection net.Conn) {
  defer connection.Close()
  fmt.Println("Handling Request\n")
  request(connection)
}

func request(connection net.Conn) {
  i := 0
  scanner := bufio.NewScanner(connection)

  for scanner.Scan() {
    ln := scanner.Text()
    fmt.Println(ln)

    if i == 0 {
      method := strings.Fields(ln)[0]
      uri := strings.Fields(ln)[1]

      respond(connection, method, uri)
    }

    if ln == "" {
      break
    }

    i++
  }
}

func respond(connection net.Conn, method string, uri string) {
  fmt.Println("Responding to request:", method, uri)

  if method == "GET" && uri == "/" {
    index(connection)
  } else {
    notFound(connection)
  }
}

func index(connection net.Conn) {
  body := `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>Homepage</title></head><body>
  <h1>HOME PAGE</h1> </body></html>`
  fmt.Fprint(connection, "HTTP/1.1 200 OK\r\n")
  fmt.Fprintf(connection, "Content-Length: %d\r\n", len(body))
  fmt.Fprint(connection, "Content-Type: text/html\r\n\r\n")
  fmt.Fprint(connection, body)

}

func notFound(connection net.Conn) {
  body := `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>Homepage</title></head><body>
  <h1>404 Not Found Error</h1></body></html>`
  fmt.Fprint(connection, "HTTP/1.1 404 Not Found\r\n")
  fmt.Fprintf(connection, "Content-Length: %d\r\n", len(body))
  fmt.Fprint(connection, "Content-Type: text/html\r\n\r\n")
  fmt.Fprint(connection, body)
}

