package testhelper

import (
	"bytes"
	"github.com/jensneuse/graphql-go-tools/pkg/lookup"
	"github.com/jensneuse/graphql-go-tools/pkg/parser"
	"github.com/jensneuse/graphql-go-tools/pkg/printer"
)

func UglifyRequestString(request string) string {
	return string(UglifyRequestBytes([]byte(request)))
}

func UglifyRequestBytes(request []byte) []byte {
	parse := parser.NewParser()
	if err := parse.ParseExecutableDefinition(request); err != nil {
		panic(err)
	}
	astPrint := printer.New()
	look := lookup.New(parse, 512)
	walk := lookup.NewWalker(1024, 8)
	walk.SetLookup(look)
	walk.WalkExecutable()
	astPrint.SetInput(parse, look, walk)
	buff := bytes.Buffer{}
	if err := astPrint.PrintExecutableSchema(&buff); err != nil {
		panic(err)
	}
	return buff.Bytes()
}

func PutLiteralString(p *parser.Parser, literal string) int {
	mod := parser.NewManualAstMod(p)
	ref, _, err := mod.PutLiteralString(literal)
	if err != nil {
		panic(err)
	}
	return ref
}

func LiteralString(p *parser.Parser, cachedName int) string {
	return string(p.CachedByteSlice(cachedName))
}
