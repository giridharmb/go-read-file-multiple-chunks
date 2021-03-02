package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

/*
Generate a large file:

dd if=/dev/urandom of=large_file.bin bs=1024 count=512000

go test -bench=.
*/

/*
SectionReaderAndLength ...
*/
type SectionReaderAndLength struct {
	sectionReader *io.SectionReader
	contentLength int64
}

func normalRead(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	reader := bufio.NewReader(f)
	buf := make([]byte, 256)

	returnData := make([]byte, int64(0))

	for {
		_, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		for _, byteData := range buf {
			returnData = append(returnData, byteData)
		}
		//fmt.Printf("%s", hex.Dump(buf))
	}
	return returnData, nil
}

func performanceRead(fileName string) ([]byte, error) {

	f, _ := os.Open(fileName)
	defer func() {
		_ = f.Close()
	}()

	fi, err := os.Stat(fileName)
	if err != nil {
		log.Printf("ERROR : could not stat file : %v", err.Error())
		return nil, err
	}
	// get the size
	fileSize := fi.Size()

	var i int64

	offsetAndLengthList := make([]map[string]int64, 0)

	arrryOfNewSectionReaders := make([]SectionReaderAndLength, int64(0))
	chunkLength := int64(20971520) // bytes

	numberOfParts := fileSize / chunkLength

	lastChunkLength := fileSize % (chunkLength)

	for i = 0; i < numberOfParts; i++ {
		myMap := make(map[string]int64)
		myMap["offset"] = i * chunkLength
		myMap["chunkLength"] = chunkLength
		offsetAndLengthList = append(offsetAndLengthList, myMap)
	}

	if lastChunkLength != 0 {
		myMap := make(map[string]int64)
		offset := (numberOfParts) * chunkLength
		myMap["offset"] = offset
		myMap["chunkLength"] = lastChunkLength
		offsetAndLengthList = append(offsetAndLengthList, myMap)
	}

	for _, offsetAndLength := range offsetAndLengthList {
		offset := offsetAndLength["offset"]
		chunkLength := offsetAndLength["chunkLength"]
		s := io.NewSectionReader(f, offset, chunkLength)
		sectionReaderAndLength := SectionReaderAndLength{
			sectionReader: s,
			contentLength: chunkLength,
		}
		arrryOfNewSectionReaders = append(arrryOfNewSectionReaders, sectionReaderAndLength)
	}

	wait := make(chan bool, 1)
	readBytesChannel := make(chan []byte)

	go func() {
		for _, sectionReader := range arrryOfNewSectionReaders {
			contentLength := sectionReader.contentLength
			secReader := sectionReader.sectionReader
			buf := make([]byte, contentLength)
			_, err := secReader.Read(buf)
			if err == io.EOF {
				break
			} else if err != nil {
				break
			}
			readBytesChannel <- buf
		}
	}()

	returnData := make([]byte, int64(0))

	var j int64
	go func() {
		for j = 0; j < int64(len(arrryOfNewSectionReaders)); j++ {
			data := <-readBytesChannel
			for _, byteData := range data {
				returnData = append(returnData, byteData)
			}
		}
		close(wait)
	}()

	<-wait

	return returnData, nil

}

func main() {

	fileName := "random2.bin"

	//log.Printf("file to read (%v)", fileName)
	//
	//var counter int64
	//
	//f, _ := os.Open(fileName)
	//defer func() {
	//	_ = f.Close()
	//}()
	//
	//fi, err := os.Stat(fileName)
	//if err != nil {
	//	log.Printf("ERROR : could not stat file : %v", err.Error())
	//	return
	//}
	//// get the size
	//fileSize := fi.Size()
	//
	//log.Printf("fileSize : %v", fileSize)
	//
	//var i int64
	//
	//offsetAndLengthList := make([]map[string]int64, 0)
	//arrryOfNewSectionReaders := make([]SectionReaderAndLength, 0)
	//chunkLength := int64(300) // bytes
	//
	//numberOfParts := fileSize / chunkLength
	//
	//lastChunkLength := fileSize % (chunkLength)
	//
	//log.Printf("numberOfParts : %v", numberOfParts)
	//
	//log.Printf("lastChunkLength : %v", lastChunkLength)
	//
	//for i = 0; i < numberOfParts; i++ {
	//	myMap := make(map[string]int64)
	//	myMap["offset"] = i * chunkLength
	//	myMap["chunkLength"] = chunkLength
	//	offsetAndLengthList = append(offsetAndLengthList, myMap)
	//}
	//
	//if lastChunkLength != 0 {
	//	myMap := make(map[string]int64)
	//	offset := (numberOfParts) * chunkLength
	//	myMap["offset"] = offset
	//	myMap["chunkLength"] = lastChunkLength
	//	offsetAndLengthList = append(offsetAndLengthList, myMap)
	//}
	//
	//log.Printf("offsetAndLengthList : %v", offsetAndLengthList)
	//
	//for _, offsetAndLength := range offsetAndLengthList {
	//	offset := offsetAndLength["offset"]
	//	chunkLength := offsetAndLength["chunkLength"]
	//	s := io.NewSectionReader(f, offset, chunkLength)
	//	sectionReaderAndLength := SectionReaderAndLength{
	//		sectionReader: s,
	//		contentLength: chunkLength,
	//	}
	//	arrryOfNewSectionReaders = append(arrryOfNewSectionReaders, sectionReaderAndLength)
	//}
	//
	//log.Printf("len(arrryOfNewSectionReaders): %v", len(arrryOfNewSectionReaders))
	//
	//allBytesRead := make([]byte, 0)
	//
	//counter = 0
	//for _, sectionReader := range arrryOfNewSectionReaders {
	//	contentLength := sectionReader.contentLength
	//	secReader := sectionReader.sectionReader
	//	buf := make([]byte, contentLength)
	//	n, err := secReader.Read(buf)
	//	if err == io.EOF {
	//		log.Printf("reached the end")
	//		break
	//	} else if err != nil {
	//		log.Printf("ERROR : could not read onto buffer : %v", err.Error())
	//		break
	//	}
	//
	//	log.Printf("### of bytes read : %v", n)
	//
	//	for _, data := range buf {
	//		allBytesRead = append(allBytesRead, data)
	//	}
	//	log.Printf("counter (%v) : n => %v", counter, n)
	//	counter++
	//}
	//
	////s := io.NewSectionReader(f, 5, chunkLength)

	allBytesRead, _ := performanceRead(fileName)

	log.Printf("allBytesRead : %v", len(allBytesRead))

	outFile := "output"
	writeFileHandle, _ := os.Create(outFile)
	bytesWritten, _ := writeFileHandle.Write(allBytesRead)
	log.Printf("wrote %d bytes , to file (%v)", bytesWritten, outFile)

}
