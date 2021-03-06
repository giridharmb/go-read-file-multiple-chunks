## Concurrent File Reads

> What the intent ?
```
> Basic idea is to read a file.

> One function reads the file in one go.

> Another function reads the file in separate 
  chunks and all the chunks are read in parallel using channels

> To Do

- This is just the first draft
- The code must be optimized
- It must be cleaned up
- Experiment with various file sizes and buffer length
```

### Setup

##### Genereate a random file "random_file.bin" using 'dd' Command

##### Example:
```
dd if=/dev/urandom of=large_file.bin bs=1024 count=16384000
```

##### Above command will genereate (large_file.bin) of size 16GB roughly which our program needs

##### Then run the benchmark

> go test -bench=.

##### Sample output on a Mac Book Pro with 16 cores, SSD and 32 GB of RAM

```
2021/03/05 16:43:06 fileName : large_file.bin
2021/03/05 16:43:06 filesize : 16777216000
goos: darwin
goarch: amd64
pkg: sectionReader
BenchmarkPerformanceRead-16    	       1	26100540780 ns/op
BenchmarkNormaleRead-16        	       1	89724022722 ns/op
PASS
ok  	sectionReader	117.784s
```