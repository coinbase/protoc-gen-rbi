package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"regexp"
)

var (
	validRubyField = regexp.MustCompile(`\A[a-z][A-Za-z0-9_]*\z`)
)

func main() {
	pgs.Init(pgs.DebugEnv("DEBUG")).RegisterModule(RBS()).RegisterPostProcessor(pgsgo.GoFmt()).Render()
}
