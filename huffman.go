/* Package gohuffman
*  Implement huffman code in go
*  Contain 4 function :
*   - "Compress" and "Decompress" for Data Manipulation
*   - "GetContent" and "GetHuffmanTable" for Data Reading
 */
package main

import (
	"container/heap"
)

type HuffmanTable struct {
	encoding [256][]byte
	decoding map[string]byte
}

type Data struct {
	character  byte
	recurring  uint8
	childLeft  *Data
	childRight *Data
}

func (data *Data) isLeafNode() bool {
	return data.childLeft == nil && data.childRight == nil
}

type Frequency struct {
	buffer [256]*Symbol
	symbols []*Symbol
}

type Symbol struct {
 	ascii byte
	count uint8
}

func (f *Frequency) PopulateFrequencyTable(data []byte) {
	for _, b := range data {
		if f.buffer[b] == nil {
			symbol := &Symbol{ascii: b, count: 0}
			f.symbols = append(f.symbols, symbol)
			f.buffer[b] = symbol
		}
		f.buffer[b].count++
	}
} 

func nodeFromSymbol(symbol *Symbol) *Data {
	return &Data{
		character: symbol.ascii,
		recurring: symbol.count,
		childLeft: nil,
		childRight: nil,
	}
}

func joinNodes(a, b *Data) *Data {
	return &Data {
		character: 0,
		recurring: a.recurring + b.recurring,
		childLeft: a,
		childRight: b,
	}
}

// PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Data

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].recurring < pq[j].recurring
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Data)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func newFrequencyTable() *Frequency {
	return &Frequency {
		buffer: [256]*Symbol{},
		symbols: make([]*Symbol, 0),
	}
}

func newHuffmanTable() *HuffmanTable {
	return &HuffmanTable {
		encoding: [256][]byte{},
		decoding: make(map[string]byte),
	}
}

const (
	HUFFMAN_ENCODE_LEN = 4
	MAX_ENCODING_LEN = 8
	CONTENT_LEN_IDX = 0
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
func Compress(data []byte) ([]byte) {
	frequencyTable := newFrequencyTable()
	frequencyTable.PopulateFrequencyTable(data)
	
	queue := make(PriorityQueue, len(frequencyTable.symbols))

	for i, symbol := range frequencyTable.symbols {
		queue[i] = nodeFromSymbol(symbol)
	}

	heap.Init(&queue)

	for len(queue) > 1  {
		nodeL := heap.Pop(&queue).(*Data)
		nodeR := heap.Pop(&queue).(*Data)

		heap.Push(&queue, joinNodes(nodeL, nodeR))
	}

	root := heap.Pop(&queue).(*Data)
	treeLen := len(frequencyTable.symbols) * 2 - 1

	// Array Structure is 4 byte for each node
	// [ value, count, left, right ]
	index := 0
	huffmanTree := make([]byte, treeLen * HUFFMAN_ENCODE_LEN)
	makeHuffmanTree(root, huffmanTree, &index)

	table := newHuffmanTable()

	populateHuffmanTable(huffmanTree, table)

	result, length := encodeData(data, *table)

	join := make([]byte, 1, index + length + 1) 
	join[CONTENT_LEN_IDX] = byte(length + 1) 
	join = append(join, result...)
	join = append(join, huffmanTree...)

	return join
}

func GetContent(data []byte) []byte {
	length := data[CONTENT_LEN_IDX]
	return data[1:length]
}

func GetHuffmanTree(data []byte) []byte {
	start := data[CONTENT_LEN_IDX]
	return data[start:]
}

func encodeData(data []byte, table HuffmanTable) ([]byte, int) {
	encodedData := make([]byte, 0, len(data)*MAX_ENCODING_LEN)
	for _, b := range data {
		encodedData = append(encodedData, table.encoding[b]...)
	}
	return encodedData, len(encodedData)
}

func populateHuffmanTable(huffmanTree []byte, table *HuffmanTable) {
	if huffmanTree == nil {
		return
	}
	readHuffmanTree(huffmanTree, []byte{}, 0, table)
}

func readHuffmanTree(tree []byte, value []byte, index int, table *HuffmanTable) {
	if tree[index] == 0 {
		readHuffmanTree(tree, append(value, '0'), int(tree[index+2]), table)
		readHuffmanTree(tree, append(value, '1'), int(tree[index+3]), table)
	} else {
		leaf := make([]byte, len(value))
		copy(leaf, value)

		table.encoding[tree[index]] = leaf
		table.decoding[string(leaf)] = tree[index]
	}
}

func makeHuffmanTree(data *Data, tree []byte, index *int) {
	if data == nil {
		return  
	}

	base := *index

	tree[base] = data.character
	tree[base+1] = data.recurring

	*index += HUFFMAN_ENCODE_LEN

	if data.isLeafNode() {
		tree[base+2] = 0 
		tree[base+3] = 0 
		return
	}

	tree[base+2] = byte(*index)
	makeHuffmanTree(data.childLeft, tree, index)

	tree[base+3] = byte(*index)
	makeHuffmanTree(data.childRight, tree, index)
}

func Decompress(data []byte) []byte {
	table := newHuffmanTable()
	populateHuffmanTable(GetHuffmanTree(data), table)

	content := GetContent(data)

	decompressed := make([]byte, 0, data[0]) 
	helper := make([]byte, 0, MAX_ENCODING_LEN)
	for _, b := range content {
		helper = append(helper, b)
		value, found := table.decoding[string(helper)]
		if found  {
			decompressed = append(decompressed, value)
			helper = helper[:0]
		}
	}
	return decompressed
}
