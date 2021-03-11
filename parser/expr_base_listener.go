// Code generated from Expr.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // Expr

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseExprListener is a complete listener for a parse tree produced by ExprParser.
type BaseExprListener struct{}

var _ ExprListener = &BaseExprListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExprListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExprListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExprListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExprListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseExprListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseExprListener) ExitStart(ctx *StartContext) {}

// EnterNot is called when production Not is entered.
func (s *BaseExprListener) EnterNot(ctx *NotContext) {}

// ExitNot is called when production Not is exited.
func (s *BaseExprListener) ExitNot(ctx *NotContext) {}

// EnterBoolCompare is called when production BoolCompare is entered.
func (s *BaseExprListener) EnterBoolCompare(ctx *BoolCompareContext) {}

// ExitBoolCompare is called when production BoolCompare is exited.
func (s *BaseExprListener) ExitBoolCompare(ctx *BoolCompareContext) {}

// EnterOr is called when production Or is entered.
func (s *BaseExprListener) EnterOr(ctx *OrContext) {}

// ExitOr is called when production Or is exited.
func (s *BaseExprListener) ExitOr(ctx *OrContext) {}

// EnterBoolFunction is called when production BoolFunction is entered.
func (s *BaseExprListener) EnterBoolFunction(ctx *BoolFunctionContext) {}

// ExitBoolFunction is called when production BoolFunction is exited.
func (s *BaseExprListener) ExitBoolFunction(ctx *BoolFunctionContext) {}

// EnterIn is called when production In is entered.
func (s *BaseExprListener) EnterIn(ctx *InContext) {}

// ExitIn is called when production In is exited.
func (s *BaseExprListener) ExitIn(ctx *InContext) {}

// EnterAnd is called when production And is entered.
func (s *BaseExprListener) EnterAnd(ctx *AndContext) {}

// ExitAnd is called when production And is exited.
func (s *BaseExprListener) ExitAnd(ctx *AndContext) {}

// EnterBoolIdentifier is called when production BoolIdentifier is entered.
func (s *BaseExprListener) EnterBoolIdentifier(ctx *BoolIdentifierContext) {}

// ExitBoolIdentifier is called when production BoolIdentifier is exited.
func (s *BaseExprListener) ExitBoolIdentifier(ctx *BoolIdentifierContext) {}

// EnterCompare is called when production Compare is entered.
func (s *BaseExprListener) EnterCompare(ctx *CompareContext) {}

// ExitCompare is called when production Compare is exited.
func (s *BaseExprListener) ExitCompare(ctx *CompareContext) {}

// EnterBoolVariable is called when production BoolVariable is entered.
func (s *BaseExprListener) EnterBoolVariable(ctx *BoolVariableContext) {}

// ExitBoolVariable is called when production BoolVariable is exited.
func (s *BaseExprListener) ExitBoolVariable(ctx *BoolVariableContext) {}

// EnterBoolean is called when production Boolean is entered.
func (s *BaseExprListener) EnterBoolean(ctx *BooleanContext) {}

// ExitBoolean is called when production Boolean is exited.
func (s *BaseExprListener) ExitBoolean(ctx *BooleanContext) {}

// EnterBoolBracket is called when production BoolBracket is entered.
func (s *BaseExprListener) EnterBoolBracket(ctx *BoolBracketContext) {}

// ExitBoolBracket is called when production BoolBracket is exited.
func (s *BaseExprListener) ExitBoolBracket(ctx *BoolBracketContext) {}

// EnterBracket is called when production Bracket is entered.
func (s *BaseExprListener) EnterBracket(ctx *BracketContext) {}

// ExitBracket is called when production Bracket is exited.
func (s *BaseExprListener) ExitBracket(ctx *BracketContext) {}

// EnterIdentifier is called when production Identifier is entered.
func (s *BaseExprListener) EnterIdentifier(ctx *IdentifierContext) {}

// ExitIdentifier is called when production Identifier is exited.
func (s *BaseExprListener) ExitIdentifier(ctx *IdentifierContext) {}

// EnterVariable is called when production Variable is entered.
func (s *BaseExprListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production Variable is exited.
func (s *BaseExprListener) ExitVariable(ctx *VariableContext) {}

// EnterMulDiv is called when production MulDiv is entered.
func (s *BaseExprListener) EnterMulDiv(ctx *MulDivContext) {}

// ExitMulDiv is called when production MulDiv is exited.
func (s *BaseExprListener) ExitMulDiv(ctx *MulDivContext) {}

// EnterAddSub is called when production AddSub is entered.
func (s *BaseExprListener) EnterAddSub(ctx *AddSubContext) {}

// ExitAddSub is called when production AddSub is exited.
func (s *BaseExprListener) ExitAddSub(ctx *AddSubContext) {}

// EnterExpressionFunction is called when production ExpressionFunction is entered.
func (s *BaseExprListener) EnterExpressionFunction(ctx *ExpressionFunctionContext) {}

// ExitExpressionFunction is called when production ExpressionFunction is exited.
func (s *BaseExprListener) ExitExpressionFunction(ctx *ExpressionFunctionContext) {}

// EnterExpressionNumber is called when production ExpressionNumber is entered.
func (s *BaseExprListener) EnterExpressionNumber(ctx *ExpressionNumberContext) {}

// ExitExpressionNumber is called when production ExpressionNumber is exited.
func (s *BaseExprListener) ExitExpressionNumber(ctx *ExpressionNumberContext) {}

// EnterString is called when production String is entered.
func (s *BaseExprListener) EnterString(ctx *StringContext) {}

// ExitString is called when production String is exited.
func (s *BaseExprListener) ExitString(ctx *StringContext) {}

// EnterSubExpression is called when production SubExpression is entered.
func (s *BaseExprListener) EnterSubExpression(ctx *SubExpressionContext) {}

// ExitSubExpression is called when production SubExpression is exited.
func (s *BaseExprListener) ExitSubExpression(ctx *SubExpressionContext) {}

// EnterNumber is called when production number is entered.
func (s *BaseExprListener) EnterNumber(ctx *NumberContext) {}

// ExitNumber is called when production number is exited.
func (s *BaseExprListener) ExitNumber(ctx *NumberContext) {}

// EnterStringList is called when production stringList is entered.
func (s *BaseExprListener) EnterStringList(ctx *StringListContext) {}

// ExitStringList is called when production stringList is exited.
func (s *BaseExprListener) ExitStringList(ctx *StringListContext) {}

// EnterNumberList is called when production numberList is entered.
func (s *BaseExprListener) EnterNumberList(ctx *NumberListContext) {}

// ExitNumberList is called when production numberList is exited.
func (s *BaseExprListener) ExitNumberList(ctx *NumberListContext) {}

// EnterFunction is called when production function is entered.
func (s *BaseExprListener) EnterFunction(ctx *FunctionContext) {}

// ExitFunction is called when production function is exited.
func (s *BaseExprListener) ExitFunction(ctx *FunctionContext) {}

// EnterArgs is called when production args is entered.
func (s *BaseExprListener) EnterArgs(ctx *ArgsContext) {}

// ExitArgs is called when production args is exited.
func (s *BaseExprListener) ExitArgs(ctx *ArgsContext) {}

// EnterArg is called when production arg is entered.
func (s *BaseExprListener) EnterArg(ctx *ArgContext) {}

// ExitArg is called when production arg is exited.
func (s *BaseExprListener) ExitArg(ctx *ArgContext) {}
