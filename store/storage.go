package store

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/cybercongress/go-cyber/util"
	"github.com/vasilistefanenko/atomicf"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const DbFileFormat = ".cbdata"
const defaultBufferSize = 32768

type Storage interface {
	Put(bytes []byte) error
	Commit(version uint64) error
	LastVersion() int64
	Iterate(process func(bytes []byte)) error
	IterateTillVersion(process func(bytes []byte), ver int64) error
}

//todo: add failures if buffer is full (may be could flush changes)
//todo: add logging
//todo: store rank offset
// WARNING: NOT CONCURRENT SAFE
type BaseStorage struct {
	elementLen  uint64
	dbFilePath  string        // db file ReadWriteCloser
	buffer      *bytes.Buffer // all changes go here until commit happens
	lastVersion int64

	mu *sync.Mutex
}

func NewBaseStorage(name string, dir string, elementLen uint64) (*BaseStorage, error) {
	return NewBaseStorageBuf(name, dir, elementLen, defaultBufferSize)
}

func NewBaseStorageBuf(name string, dir string, elementLen uint64, bufferSize uint) (*BaseStorage, error) {

	dbFilePath := filepath.Join(dir, name+DbFileFormat)
	file, err := atomicf.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Try to recover file if it was corrupted
	if err = file.Recover(); err != nil {
		return nil, err
	}

	lastVersion, err := getLastVersion(file)
	if err != nil {
		return nil, err
	}

	return &BaseStorage{
		elementLen:  elementLen,
		dbFilePath:  dbFilePath,
		buffer:      bytes.NewBuffer(make([]byte, 0, bufferSize)),
		lastVersion: lastVersion,
		mu:          new(sync.Mutex),
	}, nil
}

func getLastVersion(file *atomicf.AtomicFile) (ver int64, err error) {
	ver = int64(-1)
	fileStat, _ := file.Stat()
	if fileStat.Size() > 8 {
		lastVersionOffset := fileStat.Size() - 8
		lastVersionBytes := make([]byte, 8)
		_, err = file.ReadAt(lastVersionBytes, lastVersionOffset)
		if err != nil {
			return
		}
		ver = int64(binary.LittleEndian.Uint64(lastVersionBytes)) // unsafe casting
	}
	return
}

func (bs *BaseStorage) Remove() error {
	return os.Remove(bs.dbFilePath)
}

func (bs *BaseStorage) Put(bytes []byte) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if uint64(len(bytes)) != bs.elementLen {
		return errors.New("invalid element length")
	}
	bs.buffer.Write(bytes)
	return nil
}

func (bs *BaseStorage) getReader() (io.ReadCloser, error) {
	reader, err := os.OpenFile(bs.dbFilePath, os.O_RDONLY|os.O_SYNC, 0666)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (bs *BaseStorage) getWriter() (io.WriteCloser, error) {
	writer, err := atomicf.OpenFile(bs.dbFilePath, os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		return nil, err
	}
	return writer, nil
}

func (bs *BaseStorage) IterateTillVersion(process func(bytes []byte), ver int64) error {

	if bs.lastVersion == -1 { // means that file is empty
		return nil
	}

	reader, err := bs.getReader()
	if err != nil {
		return err
	}
	defer reader.Close()

	bufr := bufio.NewReader(reader)
	for {

		elementsCountBytes, err := util.ReadExactlyNBytes(bufr, 8)
		if err != nil {
			return err
		}
		elementsCount := binary.LittleEndian.Uint64(elementsCountBytes)

		for i := uint64(0); i < elementsCount; i++ {
			elementBytes, err := util.ReadExactlyNBytes(bufr, bs.elementLen)
			if err != nil {
				return err
			}
			process(elementBytes)
		}

		versionBytes, err := util.ReadExactlyNBytes(bufr, 8)
		if err != nil {
			return err
		}
		version := binary.LittleEndian.Uint64(versionBytes)

		if int64(version) == ver { // unsafe casting
			break
		}

	}

	return nil
}

func (bs *BaseStorage) Iterate(process func(bytes []byte)) error {
	return bs.IterateTillVersion(process, bs.lastVersion)
}

//Atomic write on linux fs https://danluu.com/file-consistency/
//
//creat(/dir/log);
//write(/dir/log, “2, 3, [checksum], foo”);
//fsync(/dir/log);
//fsync(/dir);
//pwrite(/dir/orig, 2, “bar”);
//fsync(/dir/orig);
//unlink(/dir/log);
//fsync(/dir);
func (bs *BaseStorage) Commit(version uint64) error {

	bs.mu.Lock()

	writer, err := bs.getWriter()

	defer func() {
		bs.mu.Unlock()
		writer.Close()
		bs.buffer.Reset()
	}()
	if err != nil {
		return err
	}

	elementsCountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(elementsCountBytes, uint64(bs.buffer.Len())/bs.elementLen)

	versionBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(versionBytes, version)

	bytesToCommit := append(elementsCountBytes, bs.buffer.Bytes()...)
	bytesToCommit = append(bytesToCommit, versionBytes...)

	_, err = writer.Write(bytesToCommit) // todo: should be atomic operation
	if err != nil {
		return err
	}

	bs.lastVersion = int64(version) //unsafe casting

	return nil
}

func (bs *BaseStorage) LastVersion() int64 {
	return bs.lastVersion
}
