/*
        k-cliques Perlocation Mothed Algorithm
        Author: Ye Chen <chen1385@purdue.edu>
        Date: 20/10/2016
*/

package main

import (
        "bufio"
        "strings"
        "strconv"
        "fmt"
        "sort"
        "log"
        "os"
        "github.com/user/cpm/mylibs"
)

// Declare structure Clique
type Clique struct{
        clique []string         // vertex names within the clique
        color int               // the 'color' of the clique
}

// Declare global variables
var graph = make(map[string][]string)            // map each vertex with all its neighbors
var vertices []string                            // all vertices in the graph
var maxcliques []Clique                          // all found cliques
var kCliques = make(map[int]Clique)              // all found k-cliques with each k-clique mapping to a 'color'
var colors int = 0                                // initialize a flag 'color' to mark different k-cliques

/*
        Bron–Kerbosch algorithm to find maximal cliques in a graph recursively.
        Params: clique []string - verteices previously found in the clique
                p []string - vertices could be in the clique
                x []string - vertices should not be in the clique
        Append every found clique in array maxcliques.
*/
func findMaxCliques(clique []string, p []string, x []string) {
        // Base case: when p and x are both empty, append found clique and then return
        if len(x) == 0 && len(p) == 0 {
                clique = append(clique, "")
                vclique := Clique{clique[:len(clique)-1], 0}
                maxcliques = append(maxcliques, vclique)
                return
        }

        for len(p) > 0 {
                v := p[0]                                                   // select the first vertex in the candidate vertices
                neighbors := graph[v]
                tempclique := append(clique, v)                             // new clique with v
                tempx := arrayOperations.IntersectString(x, neighbors)      // new x when v is in the clique
                tempp := arrayOperations.IntersectString(p, neighbors)      // new p when v is in the clique
                sort.Strings(tempp)

                findMaxCliques(tempclique, tempp, tempx)

                // Update x and p: x = x U v, p = p / v
                tempv := []string{}
                tempv = append(tempv, v)
                x = arrayOperations.UnionString(x, tempv)
                p = arrayOperations.DifferenceString(p, tempv)
                sort.Strings(p)
        }
}

/*
        Find the k-cliques by iterating pairs of cliques in the graph, union cliques that
        share vertex number >= k - 1. Assign or merge 'colors' based on the color tags to
        differentiate k-cliques.
*/
func findKcliques(k int) {
        for i := 0; i < len(maxcliques); i++ {
                for j := i + 1; j < len(maxcliques); j++ {
                        if findSharedV(maxcliques[i].clique, maxcliques[j].clique) + 1 >= k {
                                if maxcliques[i].color == 0 && maxcliques[j].color == 0 {
                                        // if both cliques have no 'color', merge and assign a new color to them
                                        colors = colors + 1
                                        maxcliques[i].color = colors
                                        maxcliques[j].color = colors
                                        kclique := arrayOperations.UnionString(maxcliques[i].clique, maxcliques[j].clique)
                                        kCliques[colors] = Clique{kclique, colors}
                                } else {
                                        // else merge and assign both cliques with the existed 'larger color'
                                        maxColor := max(maxcliques[i].color, maxcliques[j].color)
                                        maxcliques[i].color = maxColor
                                        maxcliques[j].color = maxColor
                                        kclique := arrayOperations.UnionString(kCliques[maxColor].clique, maxcliques[i].clique, maxcliques[j].clique)
                                        temp := Clique{kclique, maxColor}
                                        kCliques[maxColor] = temp
                                }
                        }
                }
        }
}

/*
        Read in graph from input file.
*/
func readFromFile(filename string) {
        file, err := os.Open(filename)
        if err != nil {
                log.Fatal(err)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)

        for scanner.Scan() {
                vlist := strings.Fields(scanner.Text())
                vertices = append(vertices, vlist[0]);
                graph[vlist[0]] = vlist[1:]
        }

        if err := scanner.Err(); err != nil {
                log.Fatal(err)
        }
        return;
}

/*
        Return the number of vertices shared by two cliques.
*/
func findSharedV(c1 []string, c2 []string) int{
        shared := arrayOperations.IntersectString(c1, c2);
        return len(shared)
}

/*
        As there is no builtin 'max' function for integers, this is a help function
        to return the larger integer in two integers.
*/
func max(a int, b int) int {
        if a > b { return a } else { return b }
}

func main() {
        // Read in input file
        readFromFile("input.txt")

        // Initiate parames and call findMaxCliques to find maximal cliques in input graph.
        sort.Strings(vertices)
        clique := []string{}
        p := vertices
        x := []string{}
        findMaxCliques(clique, p, x)

        // Call findKcliques to find all k-cliques. Input k as the second argument in Command line.
        // default k = 4
        var k int
        if len(os.Args) == 1 { k = 4 } else { k, _ = strconv.Atoi(os.Args[1]) }
        findKcliques(k)

        // Output k-cliques to stdout.
        fmt.Printf("=== Result ===\nk-cliques with k = %d:\n", k)
        for i := 1; i < colors + 1; i++ {
                fmt.Println(kCliques[i].clique)
        }
}
