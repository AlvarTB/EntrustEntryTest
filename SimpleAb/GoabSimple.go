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
Returns the mean of the speed of the transaction, the latency and whether it was a success or
a failure. Since this version only performs a single request, the mean equals the natural time
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
   res.Body.Close()

   DnsLookup := float32(result.DNSLookup/time.Millisecond)
   TCPConnection := float32(result.TCPConnection/time.Millisecond)
   TLSHandshake := float32(result.TLSHandshake/time.Millisecond)
   serverProcessing := float32(result.ServerProcessing/time.Millisecond)
   contentTransfer := float32(result.ContentTransfer(time.Now())/time.Millisecond)

   return false, nil, DnsLookup + TCPConnection + TLSHandshake + serverProcessing + contentTransfer, 1000/serverProcessing
}


func main() {
   //Create http request with the url sent as argument
   urlAddress := os.Args [1]
   connectionFailed, err, latency, TPS := connectionHandling(urlAddress)
   if connectionFailed == true{
      // Error log
      log.Fatal(err)
   } else{
      // Show the results
      log.Printf("Total latency: %f ms", latency)
      log.Printf("TPS: %f", TPS)
   }
}

