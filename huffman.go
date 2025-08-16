/* Package gohuffman
*  Implement huffman code in go
*  Only contain 2 function "Compress" and "Decompress"
 */
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Data struct {
	character  string
	recurring  uint8
	childLeft  *Data
	childRight *Data
}

var (
	huffmanMap = make(map[string]string)
	queue      []Data
)

/*
	Compress ...

1. Create a leaf node for each symbol and add it to the priority queue.
2. While there is more than one node in the queue:
3. Remove the two nodes of highest priority (lowest probability) from the queue
4. Create a new internal node with these two nodes as children and with probability equal to the sum of the two nodes' probabilities.
5. Add the new node to the queue.
6. The remaining node is the root node and the tree is complete.
*/
func Compress(data string) string {
	dataLength := utf8.RuneCountInString(data)
	queue = make([]Data, 0, dataLength)
	for _, c := range data {
		found := false
		for i, queueData := range queue {
			if queueData.character == string(c) {
				queueData.recurring++
				queue[i] = queueData
				found = !found
				break
			}
		}
		if found {
			continue
		}
		queue = append(queue, Data{character: string(c), recurring: 1})
	}

	for {
		limit := len(queue)
		if limit == 1 {
			break
		}
		for i := 1; i < limit; i++ {
			num1 := queue[i]
			for j := i - 1; j >= 0; j-- {
				if num1.recurring > queue[j].recurring {
					queue[j+1] = queue[j]
					queue[j] = num1
				}
			}
		}

		childLeft := queue[limit-1]
		childRight := queue[limit-2]
		lowestTwo := Data{character: queue[limit-1].character + queue[limit-2].character, recurring: queue[limit-1].recurring + queue[limit-2].recurring, childLeft: &childLeft, childRight: &childRight}
		queue[limit-2] = lowestTwo
		queue = queue[:limit-1]
		fmt.Printf(" data is %s and lchild is %s , rchild is %s -> ", lowestTwo.character, lowestTwo.childLeft.character, lowestTwo.childRight.character)

		for _, queueData := range queue {
			fmt.Printf("%s %d ", queueData.character, queueData.recurring)
		}
		fmt.Printf("\n")
	}

	fmt.Printf(" dataLength is %d\n", len(queue))
	traverseTree(queue[0], make([]rune, 0, len(queue[0].character)))

	var compressed strings.Builder
	for _, character := range data {
		compressed.WriteString(huffmanMap[string(character)])
	}

	return compressed.String()
}

func traverseTree(data Data, path []rune) {
	var paths strings.Builder
	for _, i := range path {
		paths.WriteString(string(i))
	}

	fmt.Printf("%s -> %s:%d \n", paths.String(), data.character, data.recurring)
	huffmanMap[data.character] = paths.String()
	if data.childLeft != nil {
		lpath := append(path, '0')
		traverseTree(*data.childLeft, lpath)
	}
	if data.childRight != nil {
		rpath := append(path, '1')
		traverseTree(*data.childRight, rpath)
	}
}

func translateTree(data Data, path rune) (*Data, *string) {
	if path == '0' {
		if data.childLeft != nil {
			return data.childLeft, nil
		}
	} else {
		if data.childRight != nil {
			return data.childRight, nil
		}
	}
	return nil, &data.character
}

func Decompress(data string) string {
	var decompressed string
	copy := queue[0]
	for _, character := range data {
		translateTree(copy, character)
	}

	return decompressed
}
