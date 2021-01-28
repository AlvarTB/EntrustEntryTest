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
   @return Latency: latency of said connection as the sum of the httpstat parameters
   @return TPS:  the estimated transactions per second obtained from that page
 */
func connectionHandling(urlAddress string)(failure bool, Latency float64, TPS float64){

   //create Get request to the urlAddress
   req, err := http.NewRequest("GET", urlAddress, nil)
   if err != nil {
      return true, 0.0, 0.0
   }

   //create httpstat context and add it to the request
   var result httpstat.Result
   ctx := httpstat.WithHTTPStat(req.Context(), &result)
   req = req.WithContext(ctx)

   //send request by default HTTP Client
   client := http.DefaultClient
   res, err := client.Do(req)
   if err != nil {
      return true, 0.0, 0.0
   }
   if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
      return true, 0.0, 0.0
   }

   //client closes socket
   res.Body.Close()
   end := time.Now()

   //get the stats
   DnsLookup := float64(result.DNSLookup/time.Millisecond)
   TCPConnection := float64(result.TCPConnection/time.Millisecond)
   TLSHandshake := float64(result.TLSHandshake/time.Millisecond)
   serverProcessing := float64(result.ServerProcessing/time.Millisecond)
   contentTransfer := float64(result.ContentTransfer(end)/time.Millisecond)

   //return the results
   return false, DnsLookup + TCPConnection + TLSHandshake + serverProcessing + contentTransfer,
   1000/(DnsLookup + TCPConnection + TLSHandshake + serverProcessing + contentTransfer)

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
   minimum := 99999.99
   maximum := 0.0

   //flag management
   numberReqFlag := flag.Int("n", 1, "Total number of requests")
   flag.Parse()

   //get the url, handle in case it is not present
   if len(flag.Args()) == 0{
      log.Fatal("No url was given")
   }
   urlAddress := flag.Args()[0]

   //perform the -n number of requests. Will perform 1 if the option was not enabled
   i := 0
   for i < *numberReqFlag {
      connectionFailed, latency, TPS := connectionHandling(urlAddress)
      if connectionFailed == false {
         successfulConnections = successfulConnections + 1.0
         totalLatency = totalLatency + latency
         totalTPS = totalTPS + TPS
         if minimum > latency{
            minimum = latency
         }
         if maximum < latency {
            maximum = latency
         }
      }
      i++
      log.Printf("Request finished number: %d " , i)
   }

   //print the results
   if successfulConnections != 0.0 {
      log.Printf("MINIMUM: %f ms", minimum)
      log.Printf("MAXIMUM: %f", maximum)
      log.Printf("Mean latency: %f ms", totalLatency/successfulConnections)
      log.Printf("Mean TPS: %f", totalTPS/successfulConnections)
      log.Printf("Successful connections: %f", successfulConnections)
   } else{
      log.Printf("Mean latency: -- ms")
      log.Printf("Mean TPS: --")
      log.Printf("Successful connections: 0")
   }
}
