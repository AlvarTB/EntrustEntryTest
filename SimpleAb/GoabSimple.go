package main

import (
   "github.com/tcnksm/go-httpstat"
   "io"
   "io/ioutil"
   "log"
   "net/http"
   _ "net/http/httptrace"
   "os"
   "time"
)

/*Descr: Sends a GET request to a web address specified by command argument
Returns speed of transaction (how long it lasted), latency, and if it was
failure or success
*/

//Example = https://jsonplaceholder.typicode.com/posts

func main() {
   //Create http request with the url sent as argument
   urlAddress := os.Args [1]
   req, err := http.NewRequest("GET", urlAddress, nil)
   if err != nil {
      log.Fatal(err)
   }

   //create httpstat context and add it to the request
   var result httpstat.Result
   ctx := httpstat.WithHTTPStat(req.Context(), &result)
   req = req.WithContext(ctx)

   //send request by default HTTP Client
   client := http.DefaultClient
   res, err := client.Do(req)
   if err != nil {
      log.Fatal(err)
   }

   if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
      log.Fatal(err)
   }
   res.Body.Close()
//   end := time.Now()
   // Show the results
   log.Printf("DNS lookup: %d ms", int(result.DNSLookup/time.Millisecond))
   log.Printf("TCP connection: %d ms", int(result.TCPConnection/time.Millisecond))
   log.Printf("TLS handshake: %d ms", int(result.TLSHandshake/time.Millisecond))
   log.Printf("Server processing: %d ms", int(result.ServerProcessing/time.Millisecond))
   log.Printf("Content transfer: %d ms", int(result.ContentTransfer(time.Now())/time.Millisecond))
}


