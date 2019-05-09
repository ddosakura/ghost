package sh

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ddosakura/ghost"
	"github.com/ddosakura/ghost/sh/grammar"
	"github.com/ddosakura/ghost/sh/grammar/ast"
	"github.com/ddosakura/ghost/sh/grammar/selector"
)

// VM for ghost-shell
type VM struct {
	// config
	WD string

	// ast
	tree *ast.Node

	// func list
	fds map[string]*ast.Node

	// ret code
	code int
}

// LoadShell call parser
func LoadShell(s Shell) (*VM, error) {
	tree, err := ast.Parse(s.Name, s.Content, func(err grammar.SyntaxError) bool {
		ghost.Warn(s.Name, "parser error:", err)
		return false
	})
	if err != nil {
		return nil, err
	}
	return &VM{
		tree: tree.Root(),
		code: 0,
	}, nil
}

// Run shell
func (v *VM) Run(args ...string) (code int) {
	//pretty.Println(v.tree)
	cs := v.tree.Children(selector.ShellItem)
	stats := make([]*ast.Node, 0)
	for _, f := range cs {
		c := f.Child(selector.Any)
		if c.Type() == grammar.FunctionDeclaration {
			v.fds[c.Child(selector.IdentifierName).Text()] = c
		} else {
			stats = append(stats, c)
		}
	}
	//pretty.Println(v.fds)
	//pretty.Println(stats)

	v.run(stats, args...)
	return v.code
}

//Statement -> Statement /* interface */
//    : IfStatement
//    | LabelledStatement
//    | GotoStatement
//    | Expression
//    | ReturnStatement
//    | SingleLineComment -> CommentStatement
//;
func (v *VM) run(stats []*ast.Node, args ...string) {
	defer func() {
		e := recover()
		if v.code == 0 && e != nil {
			// vm error
			ghost.ErrorInDefer(e)
		}
	}()
	b := &block{
		vm: v,
		up: nil,
		vals: map[string]*Val{
			// TODO: about array
			"_": val(args),
		},
		labels:  make(map[string]int),
		stats:   stats,
		retv:    val0,
		running: true,
	}
	b.load()
	b.run()
	v.code = b.retv.Int()
}

type block struct {
	vm     *VM
	up     *block
	vals   map[string]*Val
	labels map[string]int
	stats  []*ast.Node

	retv    *Val
	running bool
}

func (b *block) load() {
	// TODO:
	for i, stat := range b.stats {
		switch stat.Type() {
		case grammar.IfStatement:
			//
		case grammar.LabelledStatement:
			label := stat.Child(selector.IdentifierName).Text()
			b.labels[label] = i
		case grammar.GotoStatement:
			//
		case grammar.SingleLineComment:
			//
		case grammar.ReturnStatement:
			//
		default:
			// in grammar.Expression
		}
	}
}

func (b *block) run() {
	// TODO:
	for ip := 0; ip < len(b.stats) && b.running; ip++ {
		stat := b.stats[ip]
		switch stat.Type() {
		case grammar.IfStatement:
			//
		case grammar.LabelledStatement:
			//
		case grammar.GotoStatement:
			//
		case grammar.SingleLineComment:
			//
		case grammar.ReturnStatement:
			expr := stat.Child(selector.Any)
			ret := b.expr(expr)
			b.ret(ret)
		default:
			// in grammar.Expression
			//debug(0, stat)
			b.expr(stat)
		}
	}
}

func (b *block) ret(v *Val) {
	if b.up == nil {
		b.retv = v
	} else {
		b.running = false
		b.up.ret(v)
	}
}

//var Expression = []NodeType{
//    AdditiveExpression,
//    AssignmentExpression,
//    BitwiseANDExpression,
//    BitwiseORExpression,
//    BitwiseXORExpression,
//    CallExpression,
//    ConditionalExpression,
//    EqualityExpression,
//    ExponentiationExpression,
//    IdentifierName,
//    LogicalANDExpression,
//    LogicalORExpression,
//    MultiplicativeExpression,
//    NumLiteral,
//    Parenthesized,
//    RelationalExpression,
//    ShiftExpression,
//    StrLiteral,
//    UnaryAdditiveExpression,
//    UnaryExpression,
//}
func (b *block) expr(stat *ast.Node) *Val {
	defer func() {
		e := recover()
		if e != nil {
			if e != ErrUnknowExpr {
				// shell error
				b.vm.code = -1
				ghost.ErrorInDefer(e)
			}
			panic(e)
		}
	}()
	return b.parseExpr(stat)
}

func (b *block) parseExpr(stat *ast.Node) *Val {
	// TODO:
	switch stat.Type() {
	case grammar.AdditiveExpression:
	case grammar.AssignmentExpression:
	case grammar.BitwiseANDExpression:
	case grammar.BitwiseORExpression:
	case grammar.BitwiseXORExpression:
	case grammar.CallExpression:
		fn := stat.Child(selector.IdentifierName).Text()
		argNodes := stat.
			Child(selector.Arguments).
			Children(selector.Any)
		args := make([]*Val, 0, len(argNodes))
		for _, arg := range argNodes {
			args = append(args, b.expr(arg))
		}

		f := b.findFunc(fn)
		if f != nil {
			return f(args)
		}

		argList := make([]string, 0, len(args))
		for _, arg := range args {
			argList = append(argList, arg.String())
		}

		// TODO: 改一下调用方式
		c := exec.Command(fn, argList...)
		c.Dir = RootDirTmp
		e := execCmd(c)
		// TODO: 处理异常及程序返回值情况(分变量接收返回值和无变量接收返回值)
		if e != nil {
			b.vm.code = -1
			ghost.Error(e)
		}
		return val0
	case grammar.ConditionalExpression:
	case grammar.EqualityExpression:
	case grammar.ExponentiationExpression:
	case grammar.IdentifierName:
	case grammar.LogicalANDExpression:
	case grammar.LogicalORExpression:
	case grammar.MultiplicativeExpression:
	case grammar.NumLiteral:
	case grammar.Parenthesized:
	case grammar.RelationalExpression:
	case grammar.ShiftExpression:
	case grammar.StrLiteral:
		// TODO: 去掉 " ' `
		return val(stat.Text())
	case grammar.UnaryAdditiveExpression:
	case grammar.UnaryExpression:
	}
	ghost.Warn("parseExpr", grammar.ExportTypeStr[stat.Type()])
	panic(ErrUnknowExpr)
}

func (b *block) findFunc(fn string) func([]*Val) *Val {
	// TODO:
	return nil
}

// TODO: delete

func debug(dep int, n *ast.Node) {
	_space(dep)
	fmt.Println(n.LineColumn())
	_space(dep)
	fmt.Println(n.Type(), n.IsValid(),
		strings.ReplaceAll(n.Text(), "\n", "\\n"))
	for _, c := range n.Children(selector.Any) {
		debug(dep+1, c)
	}
}

func _space(dep int) {
	for dep > 0 {
		dep--
		fmt.Print("  ")
	}
}
