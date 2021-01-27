# EntrustEntryTest

This repository contains my attempt at solving the entry test of this link as per requested: https://github.com/jig/bench


## Folders
In order to proceed with excercise 2, I decided to proceed methodically. For that reason, each folder contains a small step towards the creation of the goab command.

## SimpleAb
With this small project, I attempt to make a small program that sends one single request to the specified http address and returns Latency, TPS and wether if it was successful or a failure.

For that, I had to use the github.com/tcnksm/go-httpstat package.

For now, it only displays the total time spent performing each task of the http connection.

