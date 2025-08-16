package main

import "testing"

func TestHuffman(t *testing.T) {
	const sample = "A_DEAD_DAD_CEDED_A_BAD_BABE_A_BEADED_ABACA_BED"
	const compressed = "1001001101000010010000111101100011000011001111110000111111011111100110011111110100011000011011111011101001111111000"
	data := Compress(sample)
	if data != compressed {
		t.Errorf("should be '%s' but got '%s'", compressed, data)
	}
	result := Decompress(data)
	if result != sample {
		t.Errorf("should be '%s' but got '%s'", sample, result)
	}
}
