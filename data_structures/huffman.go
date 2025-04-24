package data_structures

import (
	"container/heap"
	"fmt"
)

// HuffmanNode represents a node in the Huffman tree
type HuffmanNode struct {
	Char      rune
	Frequency int
	Left      *HuffmanNode
	Right     *HuffmanNode
}

// HuffmanHeap is a priority queue for Huffman nodes
type HuffmanHeap []*HuffmanNode

func (h HuffmanHeap) Len() int           { return len(h) }
func (h HuffmanHeap) Less(i, j int) bool { return h[i].Frequency < h[j].Frequency }
func (h HuffmanHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *HuffmanHeap) Push(x interface{}) {
	*h = append(*h, x.(*HuffmanNode))
}

func (h *HuffmanHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// BuildHuffmanTree builds the Huffman tree from the given frequency table
func BuildHuffmanTree(freqTable map[rune]int) *HuffmanNode {
	h := &HuffmanHeap{}
	heap.Init(h)

	for char, freq := range freqTable {
		heap.Push(h, &HuffmanNode{Char: char, Frequency: freq})
	}

	for h.Len() > 1 {
		left := heap.Pop(h).(*HuffmanNode)
		right := heap.Pop(h).(*HuffmanNode)
		merged := &HuffmanNode{
			Frequency: left.Frequency + right.Frequency,
			Left:      left,
			Right:     right,
		}
		heap.Push(h, merged)
	}

	return heap.Pop(h).(*HuffmanNode)
}

// GenerateCodes generates the binary codes for each character in the Huffman tree
func GenerateCodes(root *HuffmanNode) map[rune]string {
	codes := make(map[rune]string)
	var generate func(node *HuffmanNode, code string)
	generate = func(node *HuffmanNode, code string) {
		if node == nil {
			return
		}
		if node.Left == nil && node.Right == nil {
			codes[node.Char] = code
		}
		generate(node.Left, code+"0")
		generate(node.Right, code+"1")
	}
	generate(root, "")
	return codes
}

// Encode encodes the input string using the Huffman codes
func Encode(input string, codes map[rune]string) string {
	var encoded string
	for _, char := range input {
		encoded += codes[char]
	}
	return encoded
}

// Decode decodes the encoded string using the Huffman tree
func Decode(encoded string, root *HuffmanNode) string {
	var decoded string
	node := root
	for _, bit := range encoded {
		if bit == '0' {
			node = node.Left
		} else {
			node = node.Right
		}
		if node.Left == nil && node.Right == nil {
			decoded += string(node.Char)
			node = root
		}
	}
	return decoded
}

// PrintHuffmanTree prints the Huffman tree for visualization
func PrintHuffmanTree(node *HuffmanNode, prefix string, isLeft bool) {
	if node != nil {
		fmt.Printf("%s", prefix)
		if isLeft {
			fmt.Printf("├──")
		} else {
			fmt.Printf("└──")
		}
		if node.Char != 0 {
			fmt.Printf("%c:%d\n", node.Char, node.Frequency)
		} else {
			fmt.Printf("%d\n", node.Frequency)
		}
		PrintHuffmanTree(node.Left, prefix+"│   ", true)
		PrintHuffmanTree(node.Right, prefix+"    ", false)
	}
}
