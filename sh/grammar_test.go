package sh

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/ddosakura/ghost/sh/grammar"
	"github.com/ddosakura/ghost/sh/grammar/ast"
	"github.com/ddosakura/ghost/sh/grammar/selector"
	"github.com/kr/pretty"
)

var txtShell = `
func abs(a) {
	if (a<0) {
		return -a
	}
	return a
}
#func print(f, ...v) {
#	if (f>0) {
#		echo(v[0])
#	} else if (f<0) {
#		echo(v[1])
#	} else {
#		echo(v...)
#	}
#}
func print(f, a, b) {
	n = abs(f)+1
	i = 0
	:loop
	if (f>0) {
		echo(a)
	} else {
		if (f<0) {
			echo(b)
		} else {
			echo(a, b)
		}
	}
	i = i+1
	if (i<n) {
		goto loop
	}
}
print(-4, "Hello World",-233^2)
return 0
`

func TestGrammar(t *testing.T) {
	tree, err := ast.Parse("test.sh", txtShell, func(err grammar.SyntaxError) bool {
		pretty.Println("parser", err)
		return true
	})
	if err != nil {
		t.Fatal(err)
	}
	// pretty.Println(tree, err)
	tmp(0, tree.Root())
}

func unquote(s string) string {
	a, e := strconv.Unquote(`"` + s + `"`)
	if e != nil {
		panic(e)
	}
	return a
}

//func expr2Value(p *generater.Pkg, e *ast.Node) value.Value {
//	s := e.Text()
//	s = unquote(s[1 : len(s)-1])
//	return p.GlobalStr(s)
//}

//func parserExpr(e *ast.Node) (value.Value, types.Type) {
//	v, _ := strconv.Atoi(e.Text())
//	t := types.I32
//	return constant.NewInt(t, int64(v)), t
//}

func tmp(dep int, n *ast.Node) {
	space(dep)
	fmt.Println(n.LineColumn())
	space(dep)
	fmt.Println(n.Type(), n.IsValid(),
		strings.ReplaceAll(n.Text(), "\n", "\\n"))
	for _, c := range n.Children(selector.Any) {
		tmp(dep+1, c)
	}
}

func space(dep int) {
	for dep > 0 {
		dep--
		fmt.Print("  ")
	}
}
