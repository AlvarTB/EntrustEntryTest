package main

import (
   //"io/ioutil"
   //"log"
   "net/http"
   _ "net/http/httptrace"
   "os"
   "log"
   "io/ioutil"
)

/*Descr: Sends a GET request to a web address specified by command argument
Returns speed of transaction (how long it lasted), latency, and if it was
failure or success
*/

//https://jsonplaceholder.typicode.com/posts

func main() {
   urlAddress := os.Args [1]
   resp, err := http.Get(urlAddress)
   if err != nil {
      log.Fatalln(err)
   }
   //We Read the response body on the line below.
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
     //Convert the body to type string
     sb := string(body)
     log.Printf(sb)
}


