package runner

import (
	"bytes"
	"io"
	"os"
	"unicode"
)

func OutputJudge(exp_out *os.File, code_out *bytes.Buffer) bool {
	const bufSize = 64 * 1024 // 64 KB

	bufA := make([]byte, bufSize)
	bufB := make([]byte, bufSize)

	it_A := bufSize + 1
	it_B := bufSize + 1
	var errA, errB error
	var nA, nB int
	for {
		if it_A >= nA {
			nA, errA = exp_out.Read(bufA)
			it_A = 0
		}
		if it_B >= nB {
			nB, errB = code_out.Read(bufB)
			it_B = 0
		}

		if errA == io.EOF {
			is_whitespace := true
			for errB != io.EOF {
				for it_B < nB {
					if unicode.IsSpace(rune(bufB[it_B])) {
						it_B++
					} else {
						is_whitespace = false
						break
					}
				}
				if is_whitespace == false {
					break
				}
				nB, errB = code_out.Read(bufB)
				it_B = 0
			}
			return is_whitespace
		} else if errB == io.EOF {
			is_whitespace := true
			for errA != io.EOF {
				for it_A < nA {
					if unicode.IsSpace(rune(bufA[it_A])) {
						it_A++
					} else {
						is_whitespace = false
						break
					}
				}
				if is_whitespace == false {
					break
				}
				nA, errA = code_out.Read(bufA)
				it_A = 0
			}
			return is_whitespace
		}

		if unicode.IsSpace(rune(bufA[it_A])) {
			it_A++
		} else if unicode.IsSpace(rune(bufB[it_B])) {
			it_B++
		} else {
			if bufA[it_A] != bufB[it_B] {
				return false
			} else {
				it_A++
				it_B++
			}
		}
	}
}
