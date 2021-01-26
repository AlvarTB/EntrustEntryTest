package main

import (
   "flag"
   "github.com/tcnksm/go-httpstat"
   "io"
   "io/ioutil"
   "log"
   "net/http"
   _ "net/http/httptrace"
   "os"
   "time"
   "fmt"
)

/*Descr: Sends a GET request to a web address specified by command argument
Returns the mean of the speed of the transaction, the latency and wheter it was a success or
a failure.
If the flag -n is enabled, it sends the number of requests specified by it
*/

func connectionHandling(urlAddress string)(failure bool, err error, Latency float32, TPS float32){
   req, err := http.NewRequest("GET", urlAddress, nil)
   if err != nil {
      return true, err, 0.0, 0.0
   }

   //create httpstat context and add it to the request
   var result httpstat.Result
   ctx := httpstat.WithHTTPStat(req.Context(), &result)
   req = req.WithContext(ctx)

   //send request by default HTTP Client
   client := http.DefaultClient
   res, err := client.Do(req)
   if err != nil {
      return true, err, 0.0, 0.0
   }

   if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
      return true, err, 0.0, 0.0
   }
   //client closes socket
   res.Body.Close()

   DnsLookup := float32(result.DNSLookup/time.Millisecond)
   TCPConnection := float32(result.TCPConnection/time.Millisecond)
   TLSHandshake := float32(result.TLSHandshake/time.Millisecond)
   serverProcessing := float32(result.ServerProcessing/time.Millisecond)
   contentTransfer := float32(result.ContentTransfer(time.Now())/time.Millisecond)

   // Show the results

   return false, nil, DnsLookup + TCPConnection + TLSHandshake + serverProcessing + contentTransfer, 1000/serverProcessing
}


func main() {
   //flag management
   numberReqFlag := flag.Int("n", 1, "Total number of requests")
   flag.Parse()

   fmt.Println(*numberReqFlag)

   //create http request with the url sent as argument
   urlAddress := os.Args [1]
   connectionFailed, err, latency, TPS := connectionHandling(urlAddress)
   if connectionFailed == true{
      log.Fatal(err)
   } else{
      log.Printf("Total latency: %f ms", latency)
      log.Printf("TPS: %f", TPS)
   }
}


//log.Printf("Content transfer: %d ms", int(result.ContentTransfer(time.Now())/time.Millisecond))

