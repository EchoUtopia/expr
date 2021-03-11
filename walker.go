package expr

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type ParseTreeWalker struct {
	stopped bool
}

func (p *ParseTreeWalker) Walk(listener antlr.ParseTreeListener, t antlr.Tree) {
	if p.stopped {
		return
	}
	switch tt := t.(type) {
	case antlr.ErrorNode:
		listener.VisitErrorNode(tt)
	case antlr.TerminalNode:
		listener.VisitTerminal(tt)
	default:
		p.EnterRule(listener, t.(antlr.RuleNode))
		for i := 0; i < t.GetChildCount(); i++ {
			child := t.GetChild(i)
			p.Walk(listener, child)
		}
		p.ExitRule(listener, t.(antlr.RuleNode))
	}
}
func (p *ParseTreeWalker) EnterRule(listener antlr.ParseTreeListener, r antlr.RuleNode) {

	ctx := r.GetRuleContext().(antlr.ParserRuleContext)
	listener.EnterEveryRule(ctx)
	ctx.EnterRule(listener)
}
func (p *ParseTreeWalker) ExitRule(listener antlr.ParseTreeListener, r antlr.RuleNode) {
	if p.stopped {
		return
	}
	ctx := r.GetRuleContext().(antlr.ParserRuleContext)
	ctx.ExitRule(listener)
	listener.ExitEveryRule(ctx)
}
func (p *ParseTreeWalker) Stop() {
	p.stopped = true
}
