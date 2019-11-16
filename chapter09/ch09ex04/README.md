# Chapter 09, exercise 4

Results:

```console
$ for w in 1 10 100 1000 10000 100000 1000000 {5,6,7}000000; do ./ch09ex04 $w; done
Creating 1 goroutines/channels
time: 36.176µs
Creating 10 goroutines/channels
time: 56.75µs
Creating 100 goroutines/channels
time: 62.205µs
Creating 1000 goroutines/channels
time: 330.294µs
Creating 10000 goroutines/channels
time: 4.477175ms
Creating 100000 goroutines/channels
time: 78.811422ms
Creating 1000000 goroutines/channels
time: 795.925556ms
Creating 5000000 goroutines/channels
time: 4.015683577s
Creating 6000000 goroutines/channels
time: 2.202176334s
Creating 7000000 goroutines/channels
Killed
```
