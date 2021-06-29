// BIP39 spec can be found at https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki

package bip39

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"

	//	"github.com/mises-id/sdk/bip39/wordlists/wordlists"
	"golang.org/x/crypto/pbkdf2"
)

var (
	// bitwise operand for big.Ints
	last11Bits  = big.NewInt(2047)
	shift11Bits = big.NewInt(2048)
	big1        = big.NewInt(1)
	big2        = big.NewInt(2)

	// bitwise masks to isolate checksum bits from entropy+checksum
	checksumMasks = map[int]*big.Int{
		12: big.NewInt(15),
		15: big.NewInt(31),
		18: big.NewInt(63),
		21: big.NewInt(127),
		24: big.NewInt(255),
	}

	checksumShifts = map[int]*big.Int{
		12: big.NewInt(16),
		15: big.NewInt(8),
		18: big.NewInt(4),
		21: big.NewInt(2),
	}

	// wordList is the words for mnemonics
	wordList []string

	// wordMap is the map for lookup from words to indexs
	wordMap map[string]int
)

func init() {
	SetWordList(English)
}

func SetWordList(list []string) {
	wordList = list
	wordMap = map[string]int{}

	for i, v := range wordList {
		wordMap[v] = i
	}
}

func GetWordList() []string {
	return wordList
}

func GetWordIndex(word string) (int, bool) {
	idx, ok := wordMap[word]
	return idx, ok
}

// Create random Entropy bytes, its bitNum must be 32n, and within [128, 256]
func NewEntropy(bitNum int) ([]byte, error) {
	if err := isValidEntropyLen(bitNum); err != nil {
		return nil, err
	}

	entropy := make([]byte, bitNum/8)
	_, _ = rand.Read(entropy)

	return entropy, nil
}

// RestoreEntropy
func RestoreEntropy(mnemonic string) ([]byte, error) {
	mnemonicSlice, isValid := splitMnemonic(mnemonic)
	if !isValid {
		return nil, errors.New("invalid mnenomic")
	}

	// decode mnemonic into a big.Int
	var (
		wordBytes [2]byte
		b         = big.NewInt(0)
	)

	for _, v := range mnemonicSlice {
		idx, found := wordMap[v]
		if !found {
			return nil, fmt.Errorf("word `%v` not found in reverse map", v)
		}

		binary.BigEndian.PutUint16(wordBytes[:], uint16(idx))
		// b shift left 11 bits
		b.Mul(b, shift11Bits)
		// new idx add to b's right most 11bits
		b.Or(b, big.NewInt(0).SetBytes(wordBytes[:]))
	}

	// compute checksum in bigint
	checksum := big.NewInt(0)
	checksumMask := checksumMasks[len(mnemonicSlice)]
	checksum = checksum.And(b, checksumMask)

	// the last bits of b restored from mnemonic are checksums, compute entropy in bigint
	b.Div(b, big.NewInt(0).Add(checksumMask, big1))

	//
	entropy := b.Bytes()
	entropy = padByteSlice(entropy, len(mnemonicSlice)/3*4)

	// compute checksum from entropy, it's 64 bits, so we need shift bits for length less than 64
	eChecksumBytes := computeChecksum(entropy)
	eChecksum := big.NewInt(int64(eChecksumBytes[0]))

	if l := len(mnemonicSlice); l != 24 {
		checksumShift := checksumShifts[l]
		eChecksum.Div(eChecksum, checksumShift)
	}

	if checksum.Cmp(eChecksum) != 0 {
		return nil, errors.New("checksum incorrect")
	}

	return entropy, nil
}

// Generate mnemonic from a given entropy
func NewMnemonic(entropy []byte) (string, error) {
	entropyBitLen := len(entropy) * 8
	checksumBitLen := entropyBitLen / 32
	mnemonicNum := (entropyBitLen + checksumBitLen) / 11

	err := isValidEntropyLen(entropyBitLen)
	if err != nil {
		return "", err
	}

	// Add checksum to entropy
	entropy = addChecksum(entropy)

	// every 11 bits of entropy -> a word of mnemonic
	entropyInt := new(big.Int).SetBytes(entropy)
	words := make([]string, mnemonicNum)
	word := big.NewInt(0)

	for i := mnemonicNum - 1; i >= 0; i-- {
		word.And(entropyInt, last11Bits)
		entropyInt.Div(entropyInt, shift11Bits)

		// one mnemonic in 2 bytes, 16 bits
		wordBytes := padByteSlice(word.Bytes(), 2)

		words[i] = wordList[binary.BigEndian.Uint16(wordBytes)]
	}

	return strings.Join(words, " "), nil
}

func Mnemonic2ByteArray(mnemonic string, raw ...bool) ([]byte, error) {
	var (
		mnemonicSlice = strings.Split(mnemonic, " ")
		entropyLen    = len(mnemonicSlice) * 11
		checksumLen   = entropyLen % 32
		fullLen       = (entropyLen-checksumLen)/8 + 1
	)

	rawEntropyByte, err := RestoreEntropy(mnemonic)
	if err != nil {
		return nil, err
	}

	if len(raw) > 0 && raw[0] {
		return rawEntropyByte, nil
	}

	return padByteSlice(addChecksum(rawEntropyByte), fullLen), nil
}

// NewSeed create seed for private key & public chain code
func NewSeed(mnemonic string, password string) ([]byte, error) {
	_, err := Mnemonic2ByteArray(mnemonic)

	if err != nil {
		return nil, err
	}

	return pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"+password), 2048, 64, sha512.New), nil
}

func IsMnemonicValid(mnemonic string) bool {
	_, err := RestoreEntropy(mnemonic)

	return err == nil
}

func addChecksum(data []byte) []byte {
	hash := computeChecksum(data)
	firstByte := hash[0]

	checksumLen := uint(len(data) / 4)
	dataBig := new(big.Int).SetBytes(data)

	for i := uint(0); i < checksumLen; i++ {
		// 1 shift to left
		dataBig.Mul(dataBig, big2)

		if firstByte&(1<<(7-i)) > 0 {
			dataBig.Or(dataBig, big1)
		}
	}

	return dataBig.Bytes()
}

func computeChecksum(data []byte) []byte {
	hasher := sha256.New()
	_, _ = hasher.Write(data)

	return hasher.Sum(nil)
}

func isValidEntropyLen(bitLen int) error {
	if (bitLen%32) != 0 || bitLen < 128 || bitLen > 256 {
		return errors.New("entropy length must be [128, 256] and a multiple of 32")
	}

	return nil
}

func padByteSlice(slice []byte, length int) []byte {
	offset := length - len(slice)
	if offset <= 0 {
		return slice
	}

	newSlice := make([]byte, length)
	copy(newSlice[offset:], slice)

	return newSlice
}

/*
func compSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
*/
func splitMnemonic(mnemonic string) ([]string, bool) {
	words := strings.Fields(mnemonic)
	num := len(words)

	if (num%3) != 0 || num < 12 || num > 24 {
		return nil, false
	}

	return words, true
}
