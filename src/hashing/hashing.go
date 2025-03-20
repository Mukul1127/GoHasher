package hashing

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"strings"

	"hash/crc32"
	"hash/crc64"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"

	"github.com/OneOfOne/xxhash"
	"github.com/daegalus/xxh3"
	"github.com/jzelinskie/whirlpool"
	"github.com/mmcloughlin/md4"

	"lukechampine.com/blake3"
)

func HashFile(filePath string, algorithms []hash.Hash, bufferSize int) ([]string, error) {
	// Input Validation
	if filePath == "" {
		return nil, fmt.Errorf("empty file path provided")
	}

	_, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("file does not exist")
	}

	if len(algorithms) == 0 {
		return nil, fmt.Errorf("no hash algorithms provided")
	}

	if bufferSize < 1 {
		return nil, fmt.Errorf("invalid buffer size provided")
	}

	// Reset hashes
	for _, algorithm := range algorithms {
		algorithm.Reset()
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file) // Error handling is unneeded as we are only reading

	// Make MultiWriter
	writers := make([]io.Writer, len(algorithms))
	for i, algorithm := range algorithms {
		writers[i] = algorithm
	}
	multiWriter := io.MultiWriter(writers...)

	// Copy file into hash using buffer
	buffer := make([]byte, bufferSize)
	if _, err := io.CopyBuffer(multiWriter, file, buffer); err != nil {
		return nil, fmt.Errorf("error while hashing file: %w", err)
	}

	// Finalize hashes and create hex strings
	hexStrings := make([]string, len(algorithms))
	for i, algorithm := range algorithms {
		hexStrings[i] = hex.EncodeToString(algorithm.Sum(nil))
	}

	return hexStrings, nil
}

var hashFunctions = []struct {
	name    string
	newHash func() (hash.Hash, error)
}{
	{"XXH_32", func() (hash.Hash, error) { return xxhash.New32(), nil }},
	{"XXH_64", func() (hash.Hash, error) { return xxhash.New64(), nil }},
	{"XXH3_64", func() (hash.Hash, error) { return xxh3.New(), nil }},
	{"XXH3_128", func() (hash.Hash, error) { return xxh3.New128(), nil }},
	{"CRC_32", func() (hash.Hash, error) { return crc32.New(crc32.MakeTable(crc32.IEEE)), nil }},
	{"CRC_64", func() (hash.Hash, error) { return crc64.New(crc64.MakeTable(crc64.ECMA)), nil }},
	{"MD4", func() (hash.Hash, error) { return md4.New(), nil }},
	{"MD5", func() (hash.Hash, error) { return md5.New(), nil }},
	{"SHA1", func() (hash.Hash, error) { return sha1.New(), nil }},
	{"SHA2_224", func() (hash.Hash, error) { return sha256.New224(), nil }},
	{"SHA2_256", func() (hash.Hash, error) { return sha256.New(), nil }},
	{"SHA2_384", func() (hash.Hash, error) { return sha512.New384(), nil }},
	{"SHA2_512", func() (hash.Hash, error) { return sha512.New(), nil }},
	{"SHA2_512/224", func() (hash.Hash, error) { return sha512.New512_224(), nil }},
	{"SHA2_512/256", func() (hash.Hash, error) { return sha512.New512_256(), nil }},
	{"SHA3_224", func() (hash.Hash, error) { return sha3.New224(), nil }},
	{"SHA3_256", func() (hash.Hash, error) { return sha3.New256(), nil }},
	{"SHA3_384", func() (hash.Hash, error) { return sha3.New384(), nil }},
	{"SHA3_512", func() (hash.Hash, error) { return sha3.New512(), nil }},
	{"BLAKE2B_224", func() (hash.Hash, error) { return blake2b.New(16, nil) }},
	{"BLAKE2B_256", func() (hash.Hash, error) { return blake2b.New(32, nil) }},
	{"BLAKE2B_384", func() (hash.Hash, error) { return blake2b.New(48, nil) }},
	{"BLAKE2B_512", func() (hash.Hash, error) { return blake2b.New(64, nil) }},
	{"BLAKE3_224", func() (hash.Hash, error) { return blake3.New(16, nil), nil }},
	{"BLAKE3_256", func() (hash.Hash, error) { return blake3.New(32, nil), nil }},
	{"BLAKE3_384", func() (hash.Hash, error) { return blake3.New(48, nil), nil }},
	{"BLAKE3_512", func() (hash.Hash, error) { return blake3.New(64, nil), nil }},
	{"WHIRLPOOL", func() (hash.Hash, error) { return whirlpool.New(), nil }},
	{"RIPEMD160", func() (hash.Hash, error) { return ripemd160.New(), nil }},
}

func GetHashFunction(name string) (hash.Hash, error) {
	name = strings.ToUpper(name)
	for _, hf := range hashFunctions {
		if hf.name == name {
			return hf.newHash()
		}
	}
	return nil, fmt.Errorf("unsupported hash algorithm: %s", name)
}

func GetHashFunctionNames() []string {
	names := make([]string, len(hashFunctions))
	for i, hf := range hashFunctions {
		names[i] = hf.name
	}
	return names
}
