package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type chunk struct {
	bufsize int
	offset  int64
}

type BufferData struct {
	bytes []byte
}

type ChunkAndBuffer struct {
	myChunk  chunk
	myBuffer []byte
	myIndex  int
}

type Output struct {
	myIndex      int
	myBufferData BufferData
}

var IndexBufferDataMap map[int]BufferData

func performanceRead(fileName string, writeAfterRead bool) {
	start := time.Now()

	const BufferSize = 314572800
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize := int(fileinfo.Size())

	log.Printf("fileName : %v", fileName)
	log.Printf("filesize : %v", filesize)

	// Number of go routines we need to spawn.
	concurrency := filesize / BufferSize
	// buffer sizes that each of the go routine below should use. ReadAt
	// returns an error if the buffer size is larger than the bytes returned
	// from the file.
	chunksizes := make([]chunk, concurrency)

	// All buffer sizes are the same in the normal case. Offsets depend on the
	// index. Second go routine should start at 100, for example, given our
	// buffer size of 100.
	for i := 0; i < concurrency; i++ {
		chunksizes[i].bufsize = BufferSize
		chunksizes[i].offset = int64(BufferSize * i)
	}

	// check for any left over bytes. Add the residual number of bytes as the
	// the last chunk size.
	if remainder := filesize % BufferSize; remainder != 0 {
		c := chunk{bufsize: remainder, offset: int64(concurrency * BufferSize)}
		concurrency++
		chunksizes = append(chunksizes, c)
	}

	var wg sync.WaitGroup

	arrayOfChunkAndBuffer := make([]ChunkAndBuffer, 0)

	chunkAndBufferChannel := make(chan ChunkAndBuffer)

	outputChannel := make(chan Output)

	IndexBufferDataMap = make(map[int]BufferData)

	for i := 0; i < concurrency; i++ {
		chunk := chunksizes[i]
		buffer := make([]byte, chunk.bufsize)
		chunkAndBuffer := ChunkAndBuffer{myBuffer: buffer, myChunk: chunk, myIndex: i}
		arrayOfChunkAndBuffer = append(arrayOfChunkAndBuffer, chunkAndBuffer)
	}

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			chunkAndBufferChannel <- arrayOfChunkAndBuffer[i]
		}(i, &wg)
	}

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()

			myChunkAndBuffer := <-chunkAndBufferChannel
			chunk := myChunkAndBuffer.myChunk
			buffer := myChunkAndBuffer.myBuffer
			index := myChunkAndBuffer.myIndex
			//bytesread, err := file.ReadAt(buffer, chunk.offset)
			_, err := file.ReadAt(buffer, chunk.offset)
			if err != nil {
				log.Printf("ERROR : %v", err.Error())
				return
			}
			bufferData := BufferData{
				bytes: buffer,
			}
			myOutputData := Output{
				myIndex:      index,
				myBufferData: bufferData,
			}
			outputChannel <- myOutputData
		}(i, &wg)
	}

	//wg.Add(1)
	go func() {
		//defer wg.Done()
		for {
			outputData, ok := <-outputChannel
			if !ok {
				return
			}
			index := outputData.myIndex
			bufferData := outputData.myBufferData
			IndexBufferDataMap[index] = bufferData
		}
	}()

	wg.Wait()

	if !writeAfterRead {
		return
	} else {
		outputFile, _ := os.OpenFile("output_file_1", os.O_CREATE|os.O_WRONLY, 0644)
		datawriter := bufio.NewWriter(outputFile)

		for i := 0; i < concurrency; i++ {
			bufferData := IndexBufferDataMap[i]
			dataBytes := bufferData.bytes
			_, _ = datawriter.Write(dataBytes)
		}

		_ = datawriter.Flush()
	}

	duration := time.Since(start)

	log.Printf("time for performance read and write : %v", duration)

}

func normalRead(fileName string, writeAfterRead bool) {
	start := time.Now()
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("ERROR : could not read file : %v", err.Error())
	}

	if !writeAfterRead {
		return
	} else {
		err = ioutil.WriteFile("output_file_2", data, 0644)
		if err != nil {
			// print it out
			log.Printf("ERROR : could not write file : %v", err.Error())
		}
	}

	duration := time.Since(start)
	log.Printf("time for normal read and write : %v", duration)

}
func main() {
	performanceRead("large_file.bin", false)
	normalRead("large_file.bin", false)
}
