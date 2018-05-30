package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/awalterschulze/gographviz"
	"github.com/figaro-tech/go-figaro/figp2p"
)

var numHosts = 25

func main() {
	ctx := context.Background()

	log.Printf("Creating %d hosts...", numHosts)

	// Create Boot Node
	bootNode := newBootNode(ctx)

	// Create Client Nodes
	nodes := []*figp2p.Node{bootNode}
	for i := 0; i < numHosts; i++ {
		nodes = append(nodes, newNode(ctx, bootNode))
	}

	// Start all the Nodes (this bootstraps/auto-connects)
	for i, node := range nodes {
		go func(n *figp2p.Node, timeout time.Duration) {
			time.Sleep(timeout)
			n.Start(ctx)
		}(node, time.Duration(int64(i*4393))*time.Millisecond)
	}

	clearDemoFolder()
	printInfo(ctx, nodes)
}

func newBootNode(ctx context.Context) *figp2p.Node {
	bootNode, err := figp2p.NewBootstrapNode(ctx)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Bootstrap Addr:", bootNode.Addrs()[0])
	return bootNode
}

func newNode(ctx context.Context, bootNode *figp2p.Node) *figp2p.Node {
	newNode, err := figp2p.NewNode(ctx, bootNode.Addrs()[:1])
	if err != nil {
		log.Panic(err)
	}
	return newNode
}

func clearDemoFolder() {
	exec.Command("rm", "-rf", "_demo").Run()
	exec.Command("mkdir", "_demo").Run()
}

func printInfo(ctx context.Context, nodes []*figp2p.Node) {
	i := 0
	for {
		time.Sleep(1 * time.Second)
		printConnections(nodes)
		if i%5 == 0 {
			saveGraphStructure(nodes, i)
		}
		i++
	}
}

func printConnections(nodes []*figp2p.Node) {
	var conns []string
	for _, node := range nodes {
		conns = append(conns, fmt.Sprintf("%2d", len(node.Host().Network().Conns())))
	}
	log.Println("Connections:", conns)
}

func saveGraphStructure(nodes []*figp2p.Node, i int) {
	// Create a graphviz data structure
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}

	for _, node := range nodes {
		err := graph.AddNode("G", node.ID().Pretty()[:8], nil)
		if err != nil {
			log.Println("err", err)
		}
	}

	addedConns := make(map[string]bool)
	for _, node := range nodes {
		for _, conn := range node.Host().Network().Conns() {
			local := conn.LocalPeer().Pretty()[:8]
			remote := conn.RemotePeer().Pretty()[:8]
			if !addedConns[local+":"+remote] && !addedConns[remote+":"+local] {
				graph.AddEdge(conn.LocalPeer().Pretty()[:8], conn.RemotePeer().Pretty()[:8], true, nil)
				addedConns[local+":"+remote] = true
			}
		}
	}

	// Turn the graphviz data into an image
	var out bytes.Buffer
	cmd := exec.Command("dot", "-Tpng", fmt.Sprintf("-o_demo/structure[%d].png", i))
	cmd.Stdin = bytes.NewBuffer([]byte(graph.String()))
	cmd.Stdout = &out
	cmd.Run()
}
