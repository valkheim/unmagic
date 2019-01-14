package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const BUF_SIZE = 42

type file struct {
	contents []byte
	header   []byte
	trailer  []byte
}

type signature struct {
	Header      string `json:"header"`
	Trailer     string `json:"trailer"`
	Description string `json:"description"`
}

type signatures struct {
	Signatures     []signature
	headerMatches  []int
	trailerMatches []int
}

func dump(header string, dump []byte) {
	fmt.Println(strings.Repeat("#", 78))
	fmt.Println(strings.Repeat(" ", 78/2-len(header)/2-1), header)
	fmt.Println(strings.Repeat("#", 78))
	fmt.Println(hex.Dump(dump))
}

func loadSignatures() []signature {
	if j, err := os.Open("signatures.json"); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	} else {
		defer j.Close()
		if jb, err := ioutil.ReadAll(j); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		} else {
			signatures := make([]signature, 0)
			if err := json.Unmarshal(jb, &signatures); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			} else {
				return signatures
			}
		}
	}
	return nil
}

func loadFile(filename string) *file {
	if b, err := ioutil.ReadFile(filename); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	} else {
		if stat, err := os.Stat(filename); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		} else {
			var bound int64 = BUF_SIZE
			size := stat.Size()
			if bound > size {
				bound = size
			}
			file := file{b, b[0:bound], b[size-bound:]}
			return &file
		}
	}
	return nil
}

func readHeader(signature signature, fH []byte) bool {
	var bC string
	for i, c := range signature.Header {
		if c != 'x' && c != 'X' {
			bC += string(c)
			if i%2 == 1 {
				b, _ := hex.DecodeString(bC)
				bC = ""
				if fH[i/2] != b[0] {
					return false
				}
			}
		}
	}
	return true
}

func readTrailer(signature signature, fT []byte) bool {
	var bC string
	j := len(fT) - 1
	for i := len(signature.Trailer) - 1; i >= 0; i-- {
		c := signature.Trailer[i]
		if c != 'x' && c != 'X' {
			bC = string(c) + bC
			if i%2 == 0 {
				b, _ := hex.DecodeString(bC)
				if fT[j] != b[0] {
					return false
				}
				j--
				bC = ""
			}
		}
	}
	return true
}

func findSignature(signatures *signatures, header []byte, trailer []byte) {
	for i, signature := range signatures.Signatures {
		if readHeader(signature, header) {
			signatures.headerMatches = append(signatures.headerMatches, i)
		}
		if readTrailer(signature, trailer) {
			signatures.trailerMatches = append(signatures.trailerMatches, i)
		}
	}
}

func showResults(signatures *signatures) {
	fmt.Println("Header matches :")
	for _, i := range signatures.headerMatches {
		fmt.Println(signatures.Signatures[i].Description)
	}
	fmt.Println("\nTrailer matches :")
	for _, i := range signatures.trailerMatches {
		fmt.Println(signatures.Signatures[i].Description)
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file>\n", os.Args[0])
		os.Exit(2)
	}

	/******************************/

	signatures := signatures{Signatures: loadSignatures()}
	if signatures.Signatures == nil {
		os.Exit(1)
	}

	/******************************/

	file := loadFile(os.Args[1])
	if file == nil {
		os.Exit(1)
	}

	dump("header", file.header)
	dump("trailer", file.trailer)

	findSignature(&signatures, file.header, file.trailer)

	showResults(&signatures)
}
