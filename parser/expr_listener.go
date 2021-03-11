// Code generated from Expr.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // Expr

import "github.com/antlr/antlr4/runtime/Go/antlr"

// ExprListener is a complete listener for a parse tree produced by ExprParser.
type ExprListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterNot is called when entering the Not production.
	EnterNot(c *NotContext)

	// EnterBoolCompare is called when entering the BoolCompare production.
	EnterBoolCompare(c *BoolCompareContext)

	// EnterOr is called when entering the Or production.
	EnterOr(c *OrContext)

	// EnterBoolFunction is called when entering the BoolFunction production.
	EnterBoolFunction(c *BoolFunctionContext)

	// EnterIn is called when entering the In production.
	EnterIn(c *InContext)

	// EnterAnd is called when entering the And production.
	EnterAnd(c *AndContext)

	// EnterBoolIdentifier is called when entering the BoolIdentifier production.
	EnterBoolIdentifier(c *BoolIdentifierContext)

	// EnterCompare is called when entering the Compare production.
	EnterCompare(c *CompareContext)

	// EnterBoolVariable is called when entering the BoolVariable production.
	EnterBoolVariable(c *BoolVariableContext)

	// EnterBoolean is called when entering the Boolean production.
	EnterBoolean(c *BooleanContext)

	// EnterBoolBracket is called when entering the BoolBracket production.
	EnterBoolBracket(c *BoolBracketContext)

	// EnterBracket is called when entering the Bracket production.
	EnterBracket(c *BracketContext)

	// EnterIdentifier is called when entering the Identifier production.
	EnterIdentifier(c *IdentifierContext)

	// EnterVariable is called when entering the Variable production.
	EnterVariable(c *VariableContext)

	// EnterMulDiv is called when entering the MulDiv production.
	EnterMulDiv(c *MulDivContext)

	// EnterAddSub is called when entering the AddSub production.
	EnterAddSub(c *AddSubContext)

	// EnterExpressionFunction is called when entering the ExpressionFunction production.
	EnterExpressionFunction(c *ExpressionFunctionContext)

	// EnterExpressionNumber is called when entering the ExpressionNumber production.
	EnterExpressionNumber(c *ExpressionNumberContext)

	// EnterString is called when entering the String production.
	EnterString(c *StringContext)

	// EnterSubExpression is called when entering the SubExpression production.
	EnterSubExpression(c *SubExpressionContext)

	// EnterNumber is called when entering the number production.
	EnterNumber(c *NumberContext)

	// EnterStringList is called when entering the stringList production.
	EnterStringList(c *StringListContext)

	// EnterNumberList is called when entering the numberList production.
	EnterNumberList(c *NumberListContext)

	// EnterFunction is called when entering the function production.
	EnterFunction(c *FunctionContext)

	// EnterArgs is called when entering the args production.
	EnterArgs(c *ArgsContext)

	// EnterArg is called when entering the arg production.
	EnterArg(c *ArgContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitNot is called when exiting the Not production.
	ExitNot(c *NotContext)

	// ExitBoolCompare is called when exiting the BoolCompare production.
	ExitBoolCompare(c *BoolCompareContext)

	// ExitOr is called when exiting the Or production.
	ExitOr(c *OrContext)

	// ExitBoolFunction is called when exiting the BoolFunction production.
	ExitBoolFunction(c *BoolFunctionContext)

	// ExitIn is called when exiting the In production.
	ExitIn(c *InContext)

	// ExitAnd is called when exiting the And production.
	ExitAnd(c *AndContext)

	// ExitBoolIdentifier is called when exiting the BoolIdentifier production.
	ExitBoolIdentifier(c *BoolIdentifierContext)

	// ExitCompare is called when exiting the Compare production.
	ExitCompare(c *CompareContext)

	// ExitBoolVariable is called when exiting the BoolVariable production.
	ExitBoolVariable(c *BoolVariableContext)

	// ExitBoolean is called when exiting the Boolean production.
	ExitBoolean(c *BooleanContext)

	// ExitBoolBracket is called when exiting the BoolBracket production.
	ExitBoolBracket(c *BoolBracketContext)

	// ExitBracket is called when exiting the Bracket production.
	ExitBracket(c *BracketContext)

	// ExitIdentifier is called when exiting the Identifier production.
	ExitIdentifier(c *IdentifierContext)

	// ExitVariable is called when exiting the Variable production.
	ExitVariable(c *VariableContext)

	// ExitMulDiv is called when exiting the MulDiv production.
	ExitMulDiv(c *MulDivContext)

	// ExitAddSub is called when exiting the AddSub production.
	ExitAddSub(c *AddSubContext)

	// ExitExpressionFunction is called when exiting the ExpressionFunction production.
	ExitExpressionFunction(c *ExpressionFunctionContext)

	// ExitExpressionNumber is called when exiting the ExpressionNumber production.
	ExitExpressionNumber(c *ExpressionNumberContext)

	// ExitString is called when exiting the String production.
	ExitString(c *StringContext)

	// ExitSubExpression is called when exiting the SubExpression production.
	ExitSubExpression(c *SubExpressionContext)

	// ExitNumber is called when exiting the number production.
	ExitNumber(c *NumberContext)

	// ExitStringList is called when exiting the stringList production.
	ExitStringList(c *StringListContext)

	// ExitNumberList is called when exiting the numberList production.
	ExitNumberList(c *NumberListContext)

	// ExitFunction is called when exiting the function production.
	ExitFunction(c *FunctionContext)

	// ExitArgs is called when exiting the args production.
	ExitArgs(c *ArgsContext)

	// ExitArg is called when exiting the arg production.
	ExitArg(c *ArgContext)
}
