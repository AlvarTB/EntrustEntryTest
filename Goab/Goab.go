package main

import (
   "flag"
   "github.com/tcnksm/go-httpstat"
   "io"
   "io/ioutil"
   "log"
   "net/http"
   _ "net/http/httptrace"
   "time"
)


/** @Description: attempts to fulfil a Get request with the specified url and returns the results
   @param urlAddress: a string containing the full url who will get the connection requests
   @return failure: a boolean specifying if the connection was or was not a success
   @return err: the error message if an error had happened
   @return Latency: latency of said connection as the sum of the httpstat parameters
   @return TPS:  the estimated transactions per second obtained from that page
 */
func connectionHandling(urlAddress string)(failure bool, err error, Latency float64, TPS float64){
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

   DnsLookup := float64(result.DNSLookup/time.Millisecond)
   TCPConnection := float64(result.TCPConnection/time.Millisecond)
   TLSHandshake := float64(result.TLSHandshake/time.Millisecond)
   serverProcessing := float64(result.ServerProcessing/time.Millisecond)
   contentTransfer := float64(result.ContentTransfer(time.Now())/time.Millisecond)

   // Show the results

   return false, nil, DnsLookup + TCPConnection + TLSHandshake + serverProcessing + contentTransfer,
   1000/serverProcessing
}

/** @Description: Performs the ab operation with the option -n
   @param flag -n: set an integer number for this option to perform n number of transactions
   @param urlAddress: send the url address
   @return TPS:  the estimated transactions per second obtained from that page
*/
func main() {
   totalLatency := 0.0
   totalTPS := 0.0
   successfulConnections := 0.0

   //flag management
   numberReqFlag := flag.Int("n", 1, "Total number of requests")
   flag.Parse()
   urlAddress := flag.Args()[0]

   i := 0
   for i < *numberReqFlag {
      connectionFailed, err, latency, TPS := connectionHandling(urlAddress)
      if connectionFailed == true {
         log.Fatal(err)
      } else{
         successfulConnections = successfulConnections + 1.0
         totalLatency = totalLatency + latency
         totalTPS = totalTPS + TPS
      }
      i++
   }
   if totalLatency > 0 {
      log.Printf("Message")
   }
   if successfulConnections != 0.0 {
      log.Printf("Mean latency: %f ms", totalLatency/successfulConnections)
      log.Printf("Mean TPS: %f", totalTPS/successfulConnections)
      log.Printf("Successful connections: %f", successfulConnections)
   } else{
      log.Printf("Mean latency: -- ms")
      log.Printf("Mean TPS: --")
      log.Printf("Successful connections: 0")
   }
}
