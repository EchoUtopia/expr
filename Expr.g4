grammar Expr ;

start: boolExpression EOF ;

boolExpression
        : '!' boolExpression                                            # Not
        | boolExpression op=(EQ | NEQ )  boolExpression                 # BoolCompare
        | expression op=(GT | LT | GTE | LTE | EQ | NEQ ) expression    # Compare
        | expression  op=(IN | NOTIN) (stringList | numberList)         # In
        | boolExpression AND boolExpression                             # And
        | boolExpression OR boolExpression                              # Or
        | BOOLEAN                                                       # Boolean
        | function                                                      # BoolFunction
        | '(' boolExpression ')'                                        # BoolBracket
        | VAR                                                           # BoolVariable
        | IDENTIFIER                                                    # BoolIdentifier
        ;

expression
        : expression op=(MUL | DIV) expression      # MulDiv
        | '-' expression                            # SubExpression
        | expression op=(ADD | SUB) expression      # AddSub
        | STRING                                    # String
        | IDENTIFIER                                # Identifier
        | number                                    # ExpressionNumber
        | VAR                                       # Variable
        | function                                  # ExpressionFunction
        | '(' expression ')'                        # Bracket
        ;

number: INT | FLOAT;


stringList
        : '(' STRING (',' STRING)* ')'
        ;

numberList
        : '(' number (',' number)* ')'
        ;

function: name=IDENTIFIER '(' fnargs=args? ')';

args: arg (',' arg)*;

arg
    : VAR
    | number
    | BOOLEAN
    | STRING
    ;




OR: 'or' ;
AND: 'and' ;

GT: '>' ;
GTE: '>=' ;
LT: '<' ;
LTE: '<=' ;
EQ: '=' ;
NEQ: '!=' ;

IN: 'in' ;
NOTIN: 'not in' ;


MUL: '*' ;
MOD: '%' ;
DIV: '/' ;
ADD: '+' ;
SUB: '-' ;


INT: '-'? DIGIT+ ;

FLOAT: '-'? DIGIT+ '.' DIGIT+ ;

BOOLEAN: TRUE | FALSE ;

STRING: '\'' (ESC|.)*? '\'' ;

VAR: '$' IDENTIFIER ;

IDENTIFIER: [a-zA-Z_]+[a-zA-Z0-9_]* ;

WS: [ \r\n\t]+ -> skip ;


fragment
ESC : '\\\'' | '\\\\' ;    // match char \' and \\

fragment
DIGIT: [0-9] ;

fragment
TRUE: 'true' ;

fragment
FALSE: 'false' ;