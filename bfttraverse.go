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

//////////////////////
/// Graph Iterator ///
//////////////////////

type BFTKmerFunc func(*BFTKmer)

func TraverseGraph(graphFilePath string, kmerFilePath string, kmerFunc BFTKmerFunc) {
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
		if kmer.Exists() {
			graphStats.getKmerStats(kmer)
			kmerFunc(kmer)
		}
	}

	graphStats.finalize()

	fmt.Println(graphStats)

	if err := kmerScanner.Err(); err != nil {
		fmt.Println(err)
	}
}
