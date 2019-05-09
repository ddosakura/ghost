language ghsh(go);

lang = "ghsh"
package = "github.com/ddosakura/ghost/sh/grammar"
eventBased = true
eventFields = true
eventAST = true
reportTokens = [SingleLineComment, invalid_token]

:: lexer

%s initial, div;
<*> eoi: /{eoi}/

invalid_token:
error:

# Whitespace
<initial, div> {
    WhiteSpace: /[\t\x0b\x0c\x20\xa0\ufeff\p{Zs}]/ (space)
}
LineTerminatorSequence: /[\n\r\u2028\u2029]|\r\n/ (space)

# Comment
SingleLineComment: /#[^\n\r\u2028\u2029]*/

# Identifier
hex = /[0-9A-Fa-f]/
IDStart = /\p{Lu}|\p{Ll}|\p{Lt}|\p{Lm}|\p{Lo}|\p{Nl}/
IDContinue = /{IDStart}|\p{Mn}|\p{Mc}|\p{Nd}|\p{Pc}/
JoinControl = /\u200c|\u200d/
unicodeEscapeSequence = /u(\{{hex}+\}|{hex}{4})/
brokenEscapeSequence = /\\(u({hex}{0,3}|\{{hex}*))?/
identifierStart = /{IDStart}|_|\\{unicodeEscapeSequence}/
identifierPart =  /{identifierStart}|{IDContinue}|{JoinControl}/
Identifier: /{identifierStart}{identifierPart}*/    (class)
# Note: the following rule disables backtracking for incomplete identifiers.
invalid_token: /({identifierStart}{identifierPart}*)?{brokenEscapeSequence}/

# Keywords.
'func':         /func/
'return':       /return/

'if':           /if/
'else':         /else/

'goto':         /goto/

# Punctuation
'=': /=/
'{': /\{/
'}': /\}/
'(': /\(/
')': /\)/
'[': /\[/
']': /\]/
'...': /\.\.\./
'<': /</
'>': />/
'<=': /<=/
'>=': />=/
'==': /==/
'!=': /!=/
'===': /===/
'!==': /!==/
'@': /@/
'+': /\+/
'-': /-/
'*': /\*/
'**': /\*\*/
'/': /\//
'%': /%/
'++': /\+\+/
'--': /--/
'<<': /<</
'>>': />>/
'>>>': />>>/
'&': /&/
'|': /\|/
'^': /^/
'!': /!/
'~': /~/
'&&': /&&/
'||': /\|\|/
'?': /\?/
':': /:/
',': /,/

# Num
#bin = /[0-1]/
#oct = /[0-7]/
#dec = /[0-9]/
#hex = /[0-9A-Fa-f]/
int = /(0+([0-7]*[89][0-9]*)?|[1-9][0-9]*)/ # 
frac = /\.[0-9]*/
exp = /[eE][+-]?[0-9]+/
bad_exp = /[eE][+-]?/
NumericLiteral: /{int}{frac}?{exp}?/# dec
NumericLiteral: /\.[0-9]+{exp}?/    # dec
NumericLiteral: /0[xX]{hex}+/# hex
NumericLiteral: /0[oO][0-7]+/# oct
NumericLiteral: /0+[0-7]+/ 1 # oct(Takes priority over the float rule above)
NumericLiteral: /0[bB][01]+/ # bin
invalid_token: /0[xXbBoO]/
invalid_token: /{int}{frac}?{bad_exp}/
invalid_token: /\.[0-9]+{bad_exp}/

# Str
# s = unquote(s[1 : len(s)-1])
escape = /\\([^1-9xu\n\r\u2028\u2029]|x{hex}{2}|{unicodeEscapeSequence})/
lineCont = /\\([\n\r\u2028\u2029]|\r\n)/
dsChar = /[^\n\r"\\\u2028\u2029]|{escape}|{lineCont}/
ssChar = /[^\n\r'\\\u2028\u2029]|{escape}|{lineCont}/
# TODO check \0 is valid if [lookahead != DecimalDigit]
StringLiteral: /"{dsChar}*"/
StringLiteral: /'{ssChar}*'/
# 按Go标准，反引号内无法转义反引号
MultiLineChars = /([^`]|{escape}|{lineCont})*/
StringLiteral: /`{MultiLineChars}`/

# for -n
unary:

# ### [ Syntax part ]
:: parser
%input Shell;

SyntaxError -> SyntaxProblem
    : error
;

# === [ Identifier & Literal]

IdentifierName -> IdentifierName
    : Identifier
;

%interface Literal;
Literal -> Literal
    : NumericLiteral -> NumLiteral
    | StringLiteral -> StrLiteral
;

# === [ Shell ]

Shell -> Shell
    : ShellItem*
;
ShellItem -> ShellItem
    : FunctionDeclaration
    | Statement
;

Block -> Block
    : '{' Statement* '}'
;
FunctionDeclaration -> FunctionDeclaration
    : 'func' IdentifierName '(' Params ')' Block
;
Params
    : (Param separator ',')*
    | (Param separator ',')+ ',' VParam
    | VParam
;
Param -> Param
    : Identifier
;
VParam -> VParam
    : '...' Param
;

# === [ Expression ]
%interface Expression;

Expression -> Expression /* interface */
    # 赋值表达式
    : AssignmentExpression
;
PrimaryExpression -> Expression /* interface */
    : IdentifierName
    | Literal
    | Parenthesized
;
Parenthesized -> Parenthesized
    : '(' Expression ')'
    | '(' SyntaxError ')'
;
CallExpression -> Expression /* interface */
    : expr=IdentifierName Arguments -> CallExpression
;
Arguments -> Arguments
    : '(' (AssignmentExpression separator ',')* ')'
;

LeftHandSideExpression -> Expression /* interface */
    : PrimaryExpression
    | CallExpression
;

%left '||';
%left '&&';
%left '|';
%left '^';
%left '&';
%left '==' '!=' '===' '!==';
%left '<' '>' '<=' '>=';
%left '<<' '>>' '>>>';
%left '-' '+';
%left '*' '/' '%';
%left unary;
%right '**';

# 一元表达式
UnaryExpression -> Expression /* interface */
    : LeftHandSideExpression
    | '~' UnaryExpression -> UnaryExpression
    | '!' UnaryExpression -> UnaryExpression
;
# 算术表达式
ArithmeticExpression -> Expression /* interface */
    : UnaryExpression
    | left=ArithmeticExpression '+' right=ArithmeticExpression -> AdditiveExpression
    | left=ArithmeticExpression '-' right=ArithmeticExpression -> AdditiveExpression
    #| '+' UnaryExpression %prec unary -> UnaryAdditiveExpression
    | '-' right=ArithmeticExpression %prec unary -> UnaryAdditiveExpression

    | left=ArithmeticExpression '<<' right=ArithmeticExpression -> ShiftExpression
    | left=ArithmeticExpression '>>' right=ArithmeticExpression -> ShiftExpression
    | left=ArithmeticExpression '>>>' right=ArithmeticExpression -> ShiftExpression
    | left=ArithmeticExpression '*' right=ArithmeticExpression -> MultiplicativeExpression
    | left=ArithmeticExpression '/' right=ArithmeticExpression -> MultiplicativeExpression
    | left=ArithmeticExpression '%' right=ArithmeticExpression -> MultiplicativeExpression
    | left=CallExpression '**' right=ArithmeticExpression -> ExponentiationExpression
;
# 二进制表达式
BinaryExpression -> Expression /* interface */
    : ArithmeticExpression
    | left=BinaryExpression '<' right=BinaryExpression -> RelationalExpression
    | left=BinaryExpression '>' right=BinaryExpression -> RelationalExpression
    | left=BinaryExpression '<=' right=BinaryExpression -> RelationalExpression
    | left=BinaryExpression '>=' right=BinaryExpression -> RelationalExpression
    | left=BinaryExpression '==' right=BinaryExpression -> EqualityExpression
    | left=BinaryExpression '!=' right=BinaryExpression -> EqualityExpression
    | left=BinaryExpression '===' right=BinaryExpression -> EqualityExpression
    | left=BinaryExpression '!==' right=BinaryExpression -> EqualityExpression
    | left=BinaryExpression '&' right=BinaryExpression -> BitwiseANDExpression
    | left=BinaryExpression '^' right=BinaryExpression -> BitwiseXORExpression
    | left=BinaryExpression '|' right=BinaryExpression -> BitwiseORExpression
    | left=BinaryExpression '&&' right=BinaryExpression -> LogicalANDExpression
    | left=BinaryExpression '||' right=BinaryExpression -> LogicalORExpression
;
# 条件表达式
ConditionalExpression -> Expression /* interface */
    : BinaryExpression
    | cond=BinaryExpression '?' then=AssignmentExpression ':' else=AssignmentExpression -> ConditionalExpression
;
# 赋值表达式
AssignmentExpression -> Expression /* interface */
    : ConditionalExpression
    | left=LeftHandSideExpression '=' right=AssignmentExpression -> AssignmentExpression
;

# === [ Statement ]
%interface Statement;

Statement -> Statement /* interface */
    : IfStatement
    | LabelledStatement
    | GotoStatement
    | Expression
    | ReturnStatement
    | SingleLineComment -> CommentStatement
;

%right 'else';
IfStatement -> IfStatement
    : 'if' '(' Expression ')' then=Block 'else' else=Block
    | 'if' '(' Expression ')' then=Block %prec 'else'
;

LabelledStatement -> LabelledStatement
    : ':' IdentifierName
;

GotoStatement -> GotoStatement
    : 'goto' IdentifierName
;

ReturnStatement -> ReturnStatement
    : 'return' Expression
;
