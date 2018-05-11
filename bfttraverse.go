package bft

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//////////////////
/// GraphStats ///
//////////////////

type GraphStats struct {
	numNodes, numEdges        uint
	avgNumEdges, avgNumColors float32
	numComponents             uint
}

func (g *GraphStats) getKmerStats(kmer *BFTKmer) {
	g.numNodes++
	numEdges := uint(len(kmer.GetSuccessors()))
	g.numEdges += numEdges
	g.avgNumEdges += float32(numEdges)
	g.avgNumColors += float32(len(kmer.colorIds))
}

func (g *GraphStats) finalize() {
	g.avgNumEdges /= float32(g.numNodes)
	g.avgNumColors /= float32(g.numNodes)
}

func (g *GraphStats) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Number of nodes: %d\n", g.numNodes))
	sb.WriteString(fmt.Sprintf("Number of edges: %d\n", g.numEdges))
	sb.WriteString(fmt.Sprintf("Average number of edges per node: %f\n", g.avgNumEdges))
	sb.WriteString(fmt.Sprintf("Average number of colors per node: %f\n", g.avgNumColors))
	sb.WriteString(fmt.Sprintf("Number of components: %d\n", g.numComponents))

	return sb.String()
}

///////////////////////
/// Graph Iterators ///
///////////////////////

func TraverseGraph(graphFilePath string, kmerFilePath string) {
	graph := NewBFTGraph(graphFilePath)
	kmerFile, err := os.Open(kmerFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer kmerFile.Close()

	kmerScanner := bufio.NewScanner(kmerFile)

	graphStats := new(GraphStats)

	// Iterate over each kmer found at kmerFilePath
	for kmerScanner.Scan() {
		kmer := graph.GetKmer(kmerScanner.Text())
		graphStats.getKmerStats(kmer)
	}

	graphStats.finalize()

	fmt.Println(graphStats)

	if err := kmerScanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func IterateGraph(graphFilePath string, kmerFilePath string) {
	nextBFTKmer := IterateKmers(graphFilePath, kmerFilePath)

	graphStats := new(GraphStats)

	// Iterate over each kmer found at kmerFilePath
	for kmer := nextBFTKmer(); kmer != nil; {
		graphStats.getKmerStats(kmer)
	}

	// Garbage collection?
	nextBFTKmer = nil

	graphStats.finalize()

	fmt.Println(graphStats)
}

func IterateKmers(graphFilePath string, kmerFilePath string) func() *BFTKmer {
	graph := NewBFTGraph(graphFilePath)
	kmerFile, err := os.Open(kmerFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer kmerFile.Close()

	kmerScanner := bufio.NewScanner(kmerFile)

	return func() *BFTKmer {
		if kmerScanner.Scan() {
			return graph.GetKmer(kmerScanner.Text())
		} else {
			return nil
		}
	}
}
