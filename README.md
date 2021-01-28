# EntrustEntryTest

This repository contains my attempt at solving the entry test of this link as per requested: https://github.com/jig/bench

## Exercise 1
In this excercise I had to execute the ab command on a simple server and analyze the results. In order to do a proper study, I tried to emulate the same conditions for all cases.
  1) Each ab was executed on the same machine.
  2) Each ab was executed towards a local nginx server I had set up in the same machine.
  3) When it was needed, always set the -n option to 5000.

Lastly, all of the obsertations here have been performend with, at least, on option enabled since enabling no option equals to sending one request, which is not enough data for a proper conclusion.

### Understanding each option
Each option allows for different things:
  1) -n: This option enables the program to send "n" number of requests to the specified addres. The default value is 1. As we have said, it is important to set this option high enough so we can get observable results.
  2) -c: This option enables concurrency at different leves, that is, the number of requests the server will be asked to handle at the same time. As we will see, we will get some optimization from that, but there will be a point where the server will get overloaded with too much work and it becomes ineffective.
  3) -k: This option activates the Keep-Alive feature, that is, the sockets will not be closed once the request is finished, instead, they are recylced for next requests.

### Experimenting with concurrency
#### No concurrency vs concurrency
Concurrency Level:      1

Time taken for tests:   0.590 seconds

Complete requests:      5000

Failed requests:        0

Total transferred:      4270000 bytes

HTML transferred:       3060000 bytes

Requests per second:    8481.38 [#/sec] (mean)

Time per request:       0.118 [ms] (mean)

Time per request:       0.118 [ms] (mean, across all concurrent requests)

Transfer rate:          7073.33 [Kbytes/sec] received

vs

Concurrency Level:      2

Time taken for tests:   0.269 seconds

Complete requests:      5000

Failed requests:        0

Total transferred:      4270000 bytes

HTML transferred:       3060000 bytes

Requests per second:    18562.79 [#/sec] (mean)

Time per request:       0.108 [ms] (mean)

Time per request:       0.054 [ms] (mean, across all concurrent requests)

Transfer rate:          15481.08 [Kbytes/sec] received

Adding just a bit of concurrency implies a great improvement on the efficiency: more requests are processed at the same time and the ratio of transmission improves too.

#### Concurrency 2 vs concurrency 3
Concurrency Level:      3

Time taken for tests:   0.311 seconds

Complete requests:      5000

Failed requests:        0

Total transferred:      4270000 bytes

HTML transferred:       3060000 bytes

Requests per second:    16061.88 [#/sec] (mean)

Time per request:       0.187 [ms] (mean)

Time per request:       0.062 [ms] (mean, across all concurrent requests)

Transfer rate:          13395.36 [Kbytes/sec] received

Level 3 of efficiency gets much worse here, though it is still slighty better than no concurrency.

### Experimenting with the -k option
#### No concurrency
Concurrency Level:      1

Time taken for tests:   0.264 seconds

Complete requests:      5000

Failed requests:        0

Keep-Alive requests:    4950

Total transferred:      4294750 bytes

HTML transferred:       3060000 bytes

Requests per second:    18916.25 [#/sec] (mean)

Time per request:       0.053 [ms] (mean)

Time per request:       0.053 [ms] (mean, across all concurrent requests)

Transfer rate:          15867.30 [Kbytes/sec] received

If we compare this with the previous situation where we hadn't enabled the k option, we will realize that enabling the K option provides a great improvement. Since there's no need to open and close sockets (except for the first ones) and we keep recycling those, we save a lot of time that we can devote onto performing actual tasks.

#### Adding levels of concurrency
As we keep adding leves of concurrency, we will be glad to realize that the executions with the -k option enabled yield better results than those without the -k option enabled. As expected, though, we can not expect infinite improvement: there's a point were we get stuck in the same levels of efficiency.

For me the point was 6, where we kept getting similar results to:
Concurrency Level:      6

Time taken for tests:   0.085 seconds

Complete requests:      5000

Failed requests:        0

Keep-Alive requests:    4953

Total transferred:      4294765 bytes

HTML transferred:       3060000 bytes

Requests per second:    59061.16 [#/sec] (mean)

Time per request:       0.102 [ms] (mean)

Time per request:       0.017 [ms] (mean, across all concurrent requests)

Transfer rate:          49541.76 [Kbytes/sec] received

## Exercise 2

### Folders
In order to proceed with excercise 2, I decided to proceed methodically. For that reason, each folder contains a small step towards the creation of the goab command.

#### SimpleAb
With this small project, I attempt to make a small program that sends one single request to the specified http address and returns Latency, TPS and wether if it was successful or a failure.

For that, I had to use the github.com/tcnksm/go-httpstat package.

##### Usage (on the same folder as the script)

-run go . <url>

#### Goab
This is the actual goab project. It attempts to emulate the behavior of the ab command. It is based on the SimpleAB small project.
For now, it only registers the -n option. Unfortunately, the results don't seem to be exactly as the ones from the ab implementation.

##### Usage (on the same folder as the script)

-run go . <-n number_of_requests> <url>
