package algorithms

import (
	"container/heap"
	"fmt"
	"github.com/ZayneLiu/go-mad/data_structures"
)

// BuildFrequencyTable builds a frequency table from the input string
func BuildFrequencyTable(input string) map[rune]int {
	freqTable := make(map[rune]int)
	for _, char := range input {
		freqTable[char]++
	}
	return freqTable
}

// BuildBinaryHeap builds a binary heap from the frequency table
func BuildBinaryHeap(freqTable map[rune]int) *data_structures.HuffmanHeap {
	h := &data_structures.HuffmanHeap{}
	heap.Init(h)

	for char, freq := range freqTable {
		heap.Push(h, &data_structures.HuffmanNode{Char: char, Frequency: freq})
	}

	return h
}

// VisualizeTreeBuilding visualizes the process of building the Huffman tree
func VisualizeTreeBuilding(h *data_structures.HuffmanHeap) *data_structures.HuffmanNode {
	for h.Len() > 1 {
		left := heap.Pop(h).(*data_structures.HuffmanNode)
		right := heap.Pop(h).(*data_structures.HuffmanNode)
		merged := &data_structures.HuffmanNode{
			Frequency: left.Frequency + right.Frequency,
			Left:      left,
			Right:     right,
		}
		heap.Push(h, merged)
		data_structures.PrintHuffmanTree(merged, "", false)
	}
	return heap.Pop(h).(*data_structures.HuffmanNode)
}

// VisualizeBinaryHeapOperations visualizes the operations on the binary heap
func VisualizeBinaryHeapOperations(h *data_structures.HuffmanHeap) {
	for h.Len() > 1 {
		left := heap.Pop(h).(*data_structures.HuffmanNode)
		right := heap.Pop(h).(*data_structures.HuffmanNode)
		merged := &data_structures.HuffmanNode{
			Frequency: left.Frequency + right.Frequency,
			Left:      left,
			Right:     right,
		}
		heap.Push(h, merged)
		fmt.Println("Merged nodes:")
		data_structures.PrintHuffmanTree(left, "", true)
		data_structures.PrintHuffmanTree(right, "", false)
		fmt.Println("Resulting node:")
		data_structures.PrintHuffmanTree(merged, "", false)
	}
}
