package command

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

func LoadIndex() Index {
	var err error
	file, err := os.Open(".go-vcs/index")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var indexHeader Index
	err = binary.Read(file, binary.BigEndian, &indexHeader.Signature)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(indexHeader.Signature[:], []byte("DIRC")) {
		fmt.Println("Not a valid Git index file")
		return Index{}
	}

	err = binary.Read(file, binary.BigEndian, &indexHeader.Version)
	if err != nil {
		panic(err)
	}
	err = binary.Read(file, binary.BigEndian, &indexHeader.EntryCount)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Signature: %s\n", indexHeader.Signature)
	// fmt.Printf("Version: %d\n", indexHeader.Version)
	// fmt.Printf("Entry Count: %d\n", indexHeader.EntryCount)

	indexHeader.Entries = make([]IndexEntry, 0, indexHeader.EntryCount)

	for i := 0; i < int(indexHeader.EntryCount); i++ {
		var entry IndexEntry

		var ctimeSec, ctimeNsec, mtimeSec, mtimeNsec uint32
		binary.Read(file, binary.BigEndian, &ctimeSec)
		binary.Read(file, binary.BigEndian, &ctimeNsec)
		binary.Read(file, binary.BigEndian, &mtimeSec)
		binary.Read(file, binary.BigEndian, &mtimeNsec)

		entry.Ctime = time.Unix(int64(ctimeSec), int64(ctimeNsec))
		entry.Mtime = time.Unix(int64(mtimeSec), int64(mtimeNsec))

		binary.Read(file, binary.BigEndian, &entry.Dev)
		binary.Read(file, binary.BigEndian, &entry.Ino)
		binary.Read(file, binary.BigEndian, &entry.Mode)
		binary.Read(file, binary.BigEndian, &entry.Uid)
		binary.Read(file, binary.BigEndian, &entry.Gid)
		binary.Read(file, binary.BigEndian, &entry.Size)
		binary.Read(file, binary.BigEndian, &entry.Sha1)
		binary.Read(file, binary.BigEndian, &entry.Flags)

		var pathBytes []byte
		for {
			var b [1]byte
			file.Read(b[:])
			if b[0] == 0x00 {
				break
			}
			pathBytes = append(pathBytes, b[0])
		}
		entry.Path = string(pathBytes)

		entryLength := 62 + len(pathBytes) + 1
		padding := (8 - (entryLength % 8)) % 8
		file.Seek(int64(padding), io.SeekCurrent)

		indexHeader.Entries = append(indexHeader.Entries, entry)

		// fmt.Printf("Entry %d: Path: %s, Sha1: %x, Mode: %o\n", i+1, entry.Path, entry.Sha1, entry.Mode)
	}

	file.Seek(-20, io.SeekEnd)
	err = binary.Read(file, binary.BigEndian, &indexHeader.Checksum)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Checksum (SHA-1): %x\n", indexHeader.Checksum)

	return indexHeader
}
