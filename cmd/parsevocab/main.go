package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kegsay/gorbslam/internal/bow"
)

var (
	flagFilename = flag.String("filename", "", "Path to the ORB vocabulary")
)

func main() {
	flag.Parse()
	if *flagFilename == "" {
		flag.Usage()
		os.Exit(1)
	}
	vocab, err := bow.NewVocabularyFromFile(*flagFilename)
	if err != nil {
		fmt.Printf("NewVocabularyFromFile: %s\n", err)
		os.Exit(2)
	}
	fmt.Println("Branching Factor:", vocab.BranchingFactor)
	fmt.Println("Depth Levels:", vocab.DepthLevels)
	fmt.Println("Num Nodes:", len(vocab.Nodes))
	fmt.Println("First 5:")
	for i := 0; i < 5; i++ {
		fmt.Printf("%+v\n", vocab.Nodes[i])
	}
	fmt.Println("Num Words:", len(vocab.Words))
	fmt.Println("First 5:")
	for i := 0; i < 5; i++ {
		fmt.Printf("%+v\n", vocab.Words[i])
	}
}
