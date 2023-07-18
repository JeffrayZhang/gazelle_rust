package ext

import (
	"context"
	"os"

	mapset "github.com/deckarep/golang-set/v2"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/rust"
)

func ParseFile(fileName string) []string {
	// Rust code
	dat, _ := os.ReadFile(fileName)

	sourceCode := []byte(dat)
	// Query with predicates
	rootCratePattern := `
		(
			scoped_identifier
			path:
			(
				identifier
			) @crate_name
			name:
			(
				identifier
			)
		)
	`
	endCratePattern := `
		(
			scoped_identifier
			name:
			(
				identifier
			) @crate_name
		)
	`

	// Parse source code into a tree
	lang := rust.GetLanguage()
	rootNode, _ := sitter.ParseCtx(context.Background(), sourceCode, lang)
	includeSet := findIdentifiers(rootCratePattern, rootNode, sourceCode)
	excludeSet := findIdentifiers(endCratePattern, rootNode, sourceCode)
	// take the difference between includeSet and excludeSet
	// any item that exists in includeSet but does not exist in excludeSet
	A := mapset.NewSet[string](includeSet...)
	B := mapset.NewSet[string](excludeSet...)
	difference := A.Difference(B)
	difference.Remove("std")
	// return difference.ToSlice()
	return difference.ToSlice()
	// first query (A): []string{"src_app_app", "axum", "tokio", "std", "dioxus_liveview", "Router"}
	// second query (B): []string{"app", "extract", "WebSocketUpgrade", "response", "Html", "routing", "get", "Router"}
	// desired output (A - B): []string{"src_app_app", "axum", "tokio" "std", "dioxus_liveview"}
}

// function that takes in a patttern, and n, and outputs a list of identifiers (strings)
// we'll call that function twice
func findIdentifiers(pattern string, rootNode *sitter.Node, sourceCode []byte) []string {
	lang := rust.GetLanguage()
	q, _ := sitter.NewQuery([]byte(pattern), lang)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, rootNode)
	
	identifiers := []string{}
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			identifiers = append(identifiers, c.Node.Content(sourceCode))
		}
	}
	return identifiers
}
