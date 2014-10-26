/;/               { return SEMI; }
/\(/              { return LPAREN; }
/\)/              { return RPAREN; }
/\{/              { return LBRACE; }
/\}/              { return RBRACE; }
/=/               { return EQUALS; }
/,/               { return COMMA; }
/true/            { lval.bool = true; return BOOLEAN; }
/false/           { lval.bool = false; return BOOLEAN; }
/[0-9][0-9]*/     { lval.number,_ = strconv.Atoi(yylex.Text()); return NUMBER }
/'[^']*'/         { lval.identifier = strings.Trim(yylex.Text(), "'"); return IDENTIFIER }
/"[^"]*"/         { lval.string = strings.Trim(yylex.Text(), "\""); return STRING }
/[A-Za-z_][A-Za-z_0-9]*/ { lval.identifier = yylex.Text(); return IDENTIFIER; }
/[ \t\n]+/        { /* eat up whitespace */ }
/./               { println("Unrecognized character:", yylex.Text()) }
/{[^\{\}\n]*}/    { /* eat up one-line comments */ }
//
package denada

import "fmt"
import "strconv"

func ystream(r io.Reader) {
  lexer := NewLexer(r);
  for {
	var lval yySymType;
    tok := lexer.Lex(&lval);
	if tok==0 {
		break
	}
	fmt.Printf("Token #%d '%s' (lval=%v)\n", tok, lexer.Text(), lval);
  }
}

func stream(r io.Reader) {
  lexer := NewLexer(r);
  for {
    tok := lexer.next(0);
	if tok==-1 {
		break
	}
	fmt.Printf("Token #%d '%s'\n", tok, lexer.Text());
  }
}
