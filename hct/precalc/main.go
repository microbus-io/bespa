/*
Copyright (c) 2023-2026 Microbus LLC and various contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/microbus-io/bespa/hct"
)

// main precalculates the HCT tonal map at the 10,20,30,40,50,60,70,80,90,95 and 99 levels
// for all colors at a fidelity of 0x08.
// 66 bytes are generated per color, starting with RGB #000000, #000008, #000010, etc.
// and through #FFFFFF.
// The generated file "tonalmap.bin" is 33^4 bytes long.
func main() {
	err := mainErr()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func mainErr() error {
	tones := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 95, 99}

	f, err := os.Create("../../css/tonalmap.bin")
	if err != nil {
		return err
	}
	buf := bufio.NewWriter(f)
	defer f.Close()

	percent := 0
	cursor := 0
	for r := 0; r <= 0x100; r += 0x08 {
		if r == 0x100 {
			r = 0xff
		}
		for g := 0; g <= 0x100; g += 0x08 {
			if g == 0x100 {
				g = 0xff
			}
			for b := 0; b <= 0x100; b += 0x08 {
				if b == 0x100 {
					b = 0xff
				}
				cursor++
				for _, tone := range tones {
					h := hct.FromInteger(0xff<<24 + r<<16 + g<<8 + b) // No transparency
					h.SetTone(float64(tone))
					t := h.ToInt() & 0xffffff
					buf.WriteByte(byte(t >> 16))
					buf.WriteByte(byte((t >> 8) & 0xff))
					buf.WriteByte(byte(t & 0xff))
				}
				pct := cursor * 100 / (33 * 33 * 33)
				if pct != percent {
					fmt.Printf("%d%%...", pct)
					percent = pct
				}
			}
		}
	}
	fmt.Printf("done\n")
	return nil
}
