package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// InvertedIndex struct
type InvertedIndex struct {
	Index map[string]*ValueNode
	Size  uint
}

// ValueNode struct
type ValueNode struct {
	Value string
	Index int
	Next  *ValueNode
}

// NewIndex returns a pointer to an Inverted Index object.
func NewIndex() *InvertedIndex {
	// Make a new map that can be given to the InvertedIndex.
	ind := make(map[string]*ValueNode)
	return &InvertedIndex{
		Index: ind,
		Size:  0,
	}
}

// IndexFile reads a file and indexes it.
func (i *InvertedIndex) IndexFile(file string) (map[string]*ValueNode, error) {
	// Read the given file into memory. This should be changed in the future.
	fopen, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// Create a Regex object that can match punctuation within a file's text.
	re, err := regexp.Compile("[.,?]")
	if err != nil {
		return nil, err
	}
	// Remove all punctuation from the file.
	premove := re.ReplaceAllString(string(fopen), " ")
	// Split all of the strings within the file.
	str := strings.Split(premove, " ")
	// Now sort all of the words within the file.
	sort.Strings(str)
	// Enter all of the values found in the file into the index.
	for index, word := range str {
		if i.Index[word] == nil {
			val := &ValueNode{
				Value: file,
				Index: index,
				Next:  nil,
			}
			i.Index[word] = val
			// Increase the size count of the index.
			i.Size++
		} else {
			// Create the value node that will be inserted into the chain.
			val := &ValueNode{
				Value: file,
				Next:  nil,
			}
			// Assign the root value to the first value node.
			place := i.Index[word]
			for place.Next != nil {
				place = place.Next
			}
			// Place the value node when .next == nil.
			place.Next = val
			// Increase size of the total inverse index.
			i.Size++
		}
	}
	// Return nil if no error occurs.
	return i.Index, nil
}

// DecodeIndex is a helper function that decodes the written index to be searched.
func DecodeIndex() *InvertedIndex {
	// Read the index file and then decode the index.
	index, err := os.Open("gosearch-cmd/tmp/index.gob")
	if err != nil {
		panic(err)
	}
	decoder := gob.NewDecoder(index)
	inf := &InvertedIndex{}
	if err = decoder.Decode(inf); err != nil {
		panic(err)
	}
	return inf
}

// SearchByKey searches all indexed documents for the provided key and prints where
// the word occurs within that document.
func (i *InvertedIndex) SearchByKey(key string) error {
	index := DecodeIndex()
	if index.Index[key] != nil {
		// Assign the root value to the first value node.
		place := index.Index[key]
		document, err := ioutil.ReadFile(place.Value)
		if err != nil {
			return err
		}
		str := strings.Split(string(document), " ")
		// Create a color print function.
		cyan := color.New(color.FgCyan).PrintfFunc()
		// Print out the document where the word was found.
		fmt.Printf("\"%s\" found in %s.\n", key, place.Value)
		// After Reading in the document, print to STDOUT.
		for _, keyword := range str {
			// If the keyword matches the key, print the word out.
			if key == keyword {
				cyan("%s ", keyword)
			} else {
				fmt.Printf("%s ", keyword)
			}
		}
		fmt.Println()
		place = place.Next

		// Increase size of the total inverse index.
		i.Size++
	} else {
		fmt.Println("This key has no entries in the index.")
	}
	return nil
}

// PrintIndex prints the given key's index.
func (i *InvertedIndex) PrintIndex() {
	// Print the entire index. Will only print keys and ValueNode address.
	fmt.Println(i.Index)
}

// PrintByKey prints a key's entire chain.
func (i *InvertedIndex) PrintByKey(key string) {
	fmt.Printf("%s:\n", key)
	if i.Index[key] != nil {
		place := i.Index[key]
		for place != nil {
			fmt.Printf("%s\n", place.Value)
			place = place.Next
		}
	} else {
		fmt.Println("This key has no entries in the index.")
	}
}
