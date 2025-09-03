package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestHuffmanEncoding(t *testing.T) {
	sample := []byte("A_DEAD_DAD_CEDED_A_BAD_BABE_A_BEADED_ABACA_BED")
	expected := []byte("1000011101001000110010011101100111001001000111110010011111011111100010001111110100111001001011111011101000111111001")
	data := Compress(sample)

	result := GetContent(data)
	tree := GetHuffmanTree(data)

	xor := make([]byte, len(expected))
	for i := range len(xor) {
		xor[i] = expected[i]^result[i]
		if xor[i] > 0 {
			fmt.Printf("%d -> xor %b (%s) with %b (%s) be %b \n", i, expected[i], string(expected[i]),result[i], string(result[i]), xor[i])
		}
	}

    fmt.Printf("expected (%d)\t '%s'\n result (%d)\t '%s'\n with xor\t '%s'\n using (%d)\t %v\n", len(expected), string(expected), len(result), string(result), string(xor), len(tree), tree)

	if string(result) != string(expected) {
		t.Errorf("should be \n '%s' but got \n '%s'\n", expected, result)
	}
}

func TestHuffmanDecoding(t *testing.T) {
	sample := bytes.Join([][]byte{
		{116},
		[]byte("1000011101001000110010011101100111001001000111110010011111011111100010001111110100111001001011111011101000111111001"),
		{0,46,4,16,0,20,8,12,95,10,0,0,68,10,0,0,0,26,20,24,65,11,0,0,0,15,28,32,69,7,0,0,0,8,36,40,67,2,0,0,66,6,0,0},
	},[]byte{})
	expected := []byte("A_DEAD_DAD_CEDED_A_BAD_BABE_A_BEADED_ABACA_BED")

	result := Decompress(sample)
	tree := GetHuffmanTree(sample)

	xor := make([]byte, len(expected))
	for i := range len(xor) {
		xor[i] = expected[i]^result[i]
		if xor[i] > 0 {
			fmt.Printf("%d -> xor %b (%s) with %b (%s) be %b \n", i, expected[i], string(expected[i]),result[i], string(result[i]), xor[i])
		}
	}

    fmt.Printf("expected (%d)\t '%s'\n result (%d)\t '%s'\n with xor\t '%s'\n using (%d)\t %v\n", len(expected), string(expected), len(result), string(result), string(xor), len(tree), tree)
	
	if string(result) != string(expected) {
		t.Errorf("should be \n '%s' but got \n '%s'\n", expected, result)
	}
}
