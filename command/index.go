package command

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/ManManavadaria/GO-version-control-system/helper"
)

type IndexEntry struct {
	Ctime time.Time
	Mtime time.Time
	Dev   uint32
	Ino   uint32
	Mode  uint32
	Uid   uint32
	Gid   uint32
	Size  uint32
	Sha1  [20]byte
	Flags uint16
	Path  string
}

type Index struct {
	Signature  [4]byte
	Version    uint32
	EntryCount uint32
	Entries    []IndexEntry
	Extensions []byte
	Checksum   [20]byte
}

func NewIndex() *Index {
	return &Index{
		Signature: [4]byte{'D', 'I', 'R', 'C'},
		Version:   2,
		Entries:   make([]IndexEntry, 0),
	}
}

func getCtime(stat os.FileInfo) time.Time {
	winStat := stat.Sys().(*syscall.Win32FileAttributeData)

	creationTime := winStat.CreationTime

	ctime := time.Unix(0, creationTime.Nanoseconds())
	return ctime
}

func (idx *Index) AddEntry(path string, content []byte) error {

	var storedEntry IndexEntry
	var IsStored bool
	for _, idxEntry := range idx.Entries {
		if path == idxEntry.Path {
			storedEntry = idxEntry
			IsStored = true
		}
	}

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	hashStr, _ := HashObjectFunc(path)
	var byteArr [20]byte
	bytes, err := hex.DecodeString(hashStr)
	if err != nil {
		log.Fatal(err)
	}

	copy(byteArr[:], bytes[:20])

	ctime := getCtime(stat)
	mtime := stat.ModTime()

	var flags uint16
	pathLen := len(path)
	if pathLen > 0xFFF {
		flags = 0xFFF
	} else {
		flags = uint16(pathLen)
	}

	entry := IndexEntry{
		Ctime: ctime,
		Mtime: mtime,
		Dev:   0,
		Ino:   0,
		Mode:  uint32(stat.Mode()),
		Uid:   0,
		Gid:   0,
		Size:  uint32(stat.Size()),
		Sha1:  byteArr,
		Flags: flags,
		Path:  path,
	}

	if IsStored && storedEntry.Path == path && storedEntry.Sha1 == entry.Sha1 {
		return nil
	} else if IsStored && storedEntry.Path == path && storedEntry.Sha1 != entry.Sha1 {
		if i := func(enries []IndexEntry, target IndexEntry) int {
			for i, p := range enries {
				if p == target {
					return i
				}
			}
			return -1
		}(idx.Entries, storedEntry); i != -1 {
			idx.Entries[i] = entry
		}
		return nil
	}

	idx.Entries = append(idx.Entries, entry)
	idx.EntryCount = uint32(len(idx.Entries))

	sort.Slice(idx.Entries, func(i, j int) bool {
		return idx.Entries[i].Path < idx.Entries[j].Path
	})

	return nil
}

func (idx *Index) Write(w io.Writer) error {
	// Write header
	if err := binary.Write(w, binary.BigEndian, idx.Signature); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, idx.Version); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, idx.EntryCount); err != nil {
		return err
	}

	for _, entry := range idx.Entries {
		if err := writeIndexEntry(w, &entry); err != nil {
			return err
		}
	}

	if len(idx.Extensions) > 0 {
		if _, err := w.Write(idx.Extensions); err != nil {
			return err
		}
	}

	h := sha1.New()
	if _, err := w.Write(h.Sum(nil)); err != nil {
		return err
	}
	copy(idx.Checksum[:], h.Sum(nil))

	return nil
}

func writeIndexEntry(w io.Writer, entry *IndexEntry) error {
	// Write each field of the entry
	if err := binary.Write(w, binary.BigEndian, entry.Ctime.Unix()); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, entry.Mtime.Unix()); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, entry.Dev); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, entry.Ino); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, entry.Mode); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, entry.Uid); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, entry.Gid); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, entry.Size); err != nil {
		return err
	}
	if _, err := w.Write(entry.Sha1[:]); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, entry.Flags); err != nil {
		return err
	}
	if _, err := w.Write([]byte(entry.Path)); err != nil {
		return err
	}

	if _, err := w.Write([]byte{0x00}); err != nil {
		return err
	}

	entrySize := 62 + len(entry.Path) + 1
	paddingSize := (8 - (entrySize % 8)) % 8
	padding := make([]byte, paddingSize)
	if _, err := w.Write(padding); err != nil {
		return err
	}
	return nil
}

func InitIndex() {
	index := NewIndex()

	file, err := os.Create(".go-vcs/index")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := index.Write(file); err != nil {
		panic(err)
	}
}

func UpdateIndex(activeFiles []string) {
	index := LoadIndex()

	for _, path := range activeFiles {
		content, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			for i := len(index.Entries) - 1; i >= 0; i-- {
				entry := index.Entries[i]
				if entry.Path == path {
					index.Entries = append(index.Entries[:i], index.Entries[i+1:]...)
				}
			}
			continue
		} else if err != nil {
			helper.PrintError(err.Error())
		}
		if err := index.AddEntry(path, content); err != nil {
			helper.PrintError(err.Error())
		}
	}

	file, err := os.OpenFile(".go-vcs/index", os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Truncate(0)

	if err := index.Write(file); err != nil {
		log.Fatal(err)
	}
}
func UnstageFilesFromIndex(files []TreeDataStruct) {
	index := LoadIndex()

	for _, file := range files {
		for i := len(index.Entries) - 1; i >= 0; i-- {
			entry := index.Entries[i]
			if entry.Path == file.Filename {
				var byteArr [20]byte
				bytes, err := hex.DecodeString(file.Hex)
				if err != nil {
					log.Fatal(err)
				}

				copy(byteArr[:], bytes[:20])
				entry.Sha1 = byteArr
				index.Entries[i] = entry
			}
			continue
		}
	}

	file, err := os.OpenFile(".go-vcs/index", os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Truncate(0)

	if err := index.Write(file); err != nil {
		log.Fatal(err)
	}
}
