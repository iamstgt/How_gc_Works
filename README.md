# How_gc_Works
This aims to provide Go developers with an overiew of the inner workings of gc. 

The gc code: https://github.com/golang/go/tree/master/src/cmd/compile

## 0. Define a function
We're going to define the function for the sake of simplicity to look at its AST and SSA

```
package main

import "fmt"

func myFunc() int {
	a := 1
	b := 2

	if a < 2 {
		a = b
	}

	return a + b
}

func main() {
	fmt.Println(myFunc())
}
```

```
GOSSAFUNC=myFunc go tool compile main.go
```

This invocation will generate ssa.html. If you open it in your browser, you'll see something like this.
![ssa.html](ssa-html.png "ssa.html")

<br>


## 1. Lexical Analysis and  Syntax Analysis
Source code is tokenized and parsed.

The corresponding code for this step; https://github.com/golang/go/tree/master/src/cmd/compile/internal/syntax

## 2. Abstract Syntax Tree (AST)
The AST is a tree-like data structure that represents the structure of a program's source code in a more abstract way. This step involves constructing the AST from the parsed source code.

AST is something like this.

![ast.png](ast.png "ast.png")

Run `ast.go`, if you want to see the full form of AST.

```
go run ast.go
```

<details><summary>Output</summary>

```rb
     0  *ast.GenDecl {
     1  .  TokPos: main.go:3:1
     2  .  Tok: import
     3  .  Lparen: -
     4  .  Specs: []ast.Spec (len = 1) {
     5  .  .  0: *ast.ImportSpec {
     6  .  .  .  Path: *ast.BasicLit {
     7  .  .  .  .  ValuePos: main.go:3:8
     8  .  .  .  .  Kind: STRING
     9  .  .  .  .  Value: "\"fmt\""
    10  .  .  .  }
    11  .  .  .  EndPos: -
    12  .  .  }
    13  .  }
    14  .  Rparen: -
    15  }
     0  *ast.FuncDecl {
     1  .  Name: *ast.Ident {
     2  .  .  NamePos: main.go:5:6
     3  .  .  Name: "myFunc"
     4  .  .  Obj: *ast.Object {
     5  .  .  .  Kind: func
     6  .  .  .  Name: "myFunc"
     7  .  .  .  Decl: *(obj @ 0)
     8  .  .  }
     9  .  }
    10  .  Type: *ast.FuncType {
    11  .  .  Func: main.go:5:1
    12  .  .  Params: *ast.FieldList {
    13  .  .  .  Opening: main.go:5:12
    14  .  .  .  Closing: main.go:5:13
    15  .  .  }
    16  .  .  Results: *ast.FieldList {
    17  .  .  .  Opening: -
    18  .  .  .  List: []*ast.Field (len = 1) {
    19  .  .  .  .  0: *ast.Field {
    20  .  .  .  .  .  Type: *ast.Ident {
    21  .  .  .  .  .  .  NamePos: main.go:5:15
    22  .  .  .  .  .  .  Name: "int"
    23  .  .  .  .  .  }
    24  .  .  .  .  }
    25  .  .  .  }
    26  .  .  .  Closing: -
    27  .  .  }
    28  .  }
    29  .  Body: *ast.BlockStmt {
    30  .  .  Lbrace: main.go:5:19
    31  .  .  List: []ast.Stmt (len = 4) {
    32  .  .  .  0: *ast.AssignStmt {
    33  .  .  .  .  Lhs: []ast.Expr (len = 1) {
    34  .  .  .  .  .  0: *ast.Ident {
    35  .  .  .  .  .  .  NamePos: main.go:6:2
    36  .  .  .  .  .  .  Name: "a"
    37  .  .  .  .  .  .  Obj: *ast.Object {
    38  .  .  .  .  .  .  .  Kind: var
    39  .  .  .  .  .  .  .  Name: "a"
    40  .  .  .  .  .  .  .  Decl: *(obj @ 32)
    41  .  .  .  .  .  .  }
    42  .  .  .  .  .  }
    43  .  .  .  .  }
    44  .  .  .  .  TokPos: main.go:6:4
    45  .  .  .  .  Tok: :=
    46  .  .  .  .  Rhs: []ast.Expr (len = 1) {
    47  .  .  .  .  .  0: *ast.BasicLit {
    48  .  .  .  .  .  .  ValuePos: main.go:6:7
    49  .  .  .  .  .  .  Kind: INT
    50  .  .  .  .  .  .  Value: "1"
    51  .  .  .  .  .  }
    52  .  .  .  .  }
    53  .  .  .  }
    54  .  .  .  1: *ast.AssignStmt {
    55  .  .  .  .  Lhs: []ast.Expr (len = 1) {
    56  .  .  .  .  .  0: *ast.Ident {
    57  .  .  .  .  .  .  NamePos: main.go:7:2
    58  .  .  .  .  .  .  Name: "b"
    59  .  .  .  .  .  .  Obj: *ast.Object {
    60  .  .  .  .  .  .  .  Kind: var
    61  .  .  .  .  .  .  .  Name: "b"
    62  .  .  .  .  .  .  .  Decl: *(obj @ 54)
    63  .  .  .  .  .  .  }
    64  .  .  .  .  .  }
    65  .  .  .  .  }
    66  .  .  .  .  TokPos: main.go:7:4
    67  .  .  .  .  Tok: :=
    68  .  .  .  .  Rhs: []ast.Expr (len = 1) {
    69  .  .  .  .  .  0: *ast.BasicLit {
    70  .  .  .  .  .  .  ValuePos: main.go:7:7
    71  .  .  .  .  .  .  Kind: INT
    72  .  .  .  .  .  .  Value: "2"
    73  .  .  .  .  .  }
    74  .  .  .  .  }
    75  .  .  .  }
    76  .  .  .  2: *ast.IfStmt {
    77  .  .  .  .  If: main.go:9:2
    78  .  .  .  .  Cond: *ast.BinaryExpr {
    79  .  .  .  .  .  X: *ast.Ident {
    80  .  .  .  .  .  .  NamePos: main.go:9:5
    81  .  .  .  .  .  .  Name: "a"
    82  .  .  .  .  .  .  Obj: *(obj @ 37)
    83  .  .  .  .  .  }
    84  .  .  .  .  .  OpPos: main.go:9:7
    85  .  .  .  .  .  Op: <
    86  .  .  .  .  .  Y: *ast.BasicLit {
    87  .  .  .  .  .  .  ValuePos: main.go:9:9
    88  .  .  .  .  .  .  Kind: INT
    89  .  .  .  .  .  .  Value: "2"
    90  .  .  .  .  .  }
    91  .  .  .  .  }
    92  .  .  .  .  Body: *ast.BlockStmt {
    93  .  .  .  .  .  Lbrace: main.go:9:11
    94  .  .  .  .  .  List: []ast.Stmt (len = 1) {
    95  .  .  .  .  .  .  0: *ast.AssignStmt {
    96  .  .  .  .  .  .  .  Lhs: []ast.Expr (len = 1) {
    97  .  .  .  .  .  .  .  .  0: *ast.Ident {
    98  .  .  .  .  .  .  .  .  .  NamePos: main.go:10:3
    99  .  .  .  .  .  .  .  .  .  Name: "a"
   100  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 37)
   101  .  .  .  .  .  .  .  .  }
   102  .  .  .  .  .  .  .  }
   103  .  .  .  .  .  .  .  TokPos: main.go:10:5
   104  .  .  .  .  .  .  .  Tok: =
   105  .  .  .  .  .  .  .  Rhs: []ast.Expr (len = 1) {
   106  .  .  .  .  .  .  .  .  0: *ast.Ident {
   107  .  .  .  .  .  .  .  .  .  NamePos: main.go:10:7
   108  .  .  .  .  .  .  .  .  .  Name: "b"
   109  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 59)
   110  .  .  .  .  .  .  .  .  }
   111  .  .  .  .  .  .  .  }
   112  .  .  .  .  .  .  }
   113  .  .  .  .  .  }
   114  .  .  .  .  .  Rbrace: main.go:11:2
   115  .  .  .  .  }
   116  .  .  .  }
   117  .  .  .  3: *ast.ReturnStmt {
   118  .  .  .  .  Return: main.go:13:2
   119  .  .  .  .  Results: []ast.Expr (len = 1) {
   120  .  .  .  .  .  0: *ast.BinaryExpr {
   121  .  .  .  .  .  .  X: *ast.Ident {
   122  .  .  .  .  .  .  .  NamePos: main.go:13:9
   123  .  .  .  .  .  .  .  Name: "a"
   124  .  .  .  .  .  .  .  Obj: *(obj @ 37)
   125  .  .  .  .  .  .  }
   126  .  .  .  .  .  .  OpPos: main.go:13:11
   127  .  .  .  .  .  .  Op: +
   128  .  .  .  .  .  .  Y: *ast.Ident {
   129  .  .  .  .  .  .  .  NamePos: main.go:13:13
   130  .  .  .  .  .  .  .  Name: "b"
   131  .  .  .  .  .  .  .  Obj: *(obj @ 59)
   132  .  .  .  .  .  .  }
   133  .  .  .  .  .  }
   134  .  .  .  .  }
   135  .  .  .  }
   136  .  .  }
   137  .  .  Rbrace: main.go:14:1
   138  .  }
   139  }
     0  *ast.FuncDecl {
     1  .  Name: *ast.Ident {
     2  .  .  NamePos: main.go:16:6
     3  .  .  Name: "main"
     4  .  .  Obj: *ast.Object {
     5  .  .  .  Kind: func
     6  .  .  .  Name: "main"
     7  .  .  .  Decl: *(obj @ 0)
     8  .  .  }
     9  .  }
    10  .  Type: *ast.FuncType {
    11  .  .  Func: main.go:16:1
    12  .  .  Params: *ast.FieldList {
    13  .  .  .  Opening: main.go:16:10
    14  .  .  .  Closing: main.go:16:11
    15  .  .  }
    16  .  }
    17  .  Body: *ast.BlockStmt {
    18  .  .  Lbrace: main.go:16:13
    19  .  .  List: []ast.Stmt (len = 1) {
    20  .  .  .  0: *ast.ExprStmt {
    21  .  .  .  .  X: *ast.CallExpr {
    22  .  .  .  .  .  Fun: *ast.SelectorExpr {
    23  .  .  .  .  .  .  X: *ast.Ident {
    24  .  .  .  .  .  .  .  NamePos: main.go:17:2
    25  .  .  .  .  .  .  .  Name: "fmt"
    26  .  .  .  .  .  .  }
    27  .  .  .  .  .  .  Sel: *ast.Ident {
    28  .  .  .  .  .  .  .  NamePos: main.go:17:6
    29  .  .  .  .  .  .  .  Name: "Println"
    30  .  .  .  .  .  .  }
    31  .  .  .  .  .  }
    32  .  .  .  .  .  Lparen: main.go:17:13
    33  .  .  .  .  .  Args: []ast.Expr (len = 1) {
    34  .  .  .  .  .  .  0: *ast.CallExpr {
    35  .  .  .  .  .  .  .  Fun: *ast.Ident {
    36  .  .  .  .  .  .  .  .  NamePos: main.go:17:14
    37  .  .  .  .  .  .  .  .  Name: "myFunc"
    38  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    39  .  .  .  .  .  .  .  .  .  Kind: func
    40  .  .  .  .  .  .  .  .  .  Name: "myFunc"
    41  .  .  .  .  .  .  .  .  .  Decl: *ast.FuncDecl {
    42  .  .  .  .  .  .  .  .  .  .  Name: *ast.Ident {
    43  .  .  .  .  .  .  .  .  .  .  .  NamePos: main.go:5:6
    44  .  .  .  .  .  .  .  .  .  .  .  Name: "myFunc"
    45  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 38)
    46  .  .  .  .  .  .  .  .  .  .  }
    47  .  .  .  .  .  .  .  .  .  .  Type: *ast.FuncType {
    48  .  .  .  .  .  .  .  .  .  .  .  Func: main.go:5:1
    49  .  .  .  .  .  .  .  .  .  .  .  Params: *ast.FieldList {
    50  .  .  .  .  .  .  .  .  .  .  .  .  Opening: main.go:5:12
    51  .  .  .  .  .  .  .  .  .  .  .  .  Closing: main.go:5:13
    52  .  .  .  .  .  .  .  .  .  .  .  }
    53  .  .  .  .  .  .  .  .  .  .  .  Results: *ast.FieldList {
    54  .  .  .  .  .  .  .  .  .  .  .  .  Opening: -
    55  .  .  .  .  .  .  .  .  .  .  .  .  List: []*ast.Field (len = 1) {
    56  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Field {
    57  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    58  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: main.go:5:15
    59  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "int"
    60  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
    61  .  .  .  .  .  .  .  .  .  .  .  .  .  }
    62  .  .  .  .  .  .  .  .  .  .  .  .  }
    63  .  .  .  .  .  .  .  .  .  .  .  .  Closing: -
    64  .  .  .  .  .  .  .  .  .  .  .  }
    65  .  .  .  .  .  .  .  .  .  .  }
    66  .  .  .  .  .  .  .  .  .  .  Body: *ast.BlockStmt {
    67  .  .  .  .  .  .  .  .  .  .  .  Lbrace: main.go:5:19
    68  .  .  .  .  .  .  .  .  .  .  .  List: []ast.Stmt (len = 4) {
    69  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.AssignStmt {
    70  .  .  .  .  .  .  .  .  .  .  .  .  .  Lhs: []ast.Expr (len = 1) {
    71  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    72  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: main.go:6:2
    73  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "a"
    74  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    75  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    76  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "a"
    77  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 69)
    78  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
    79  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
    80  .  .  .  .  .  .  .  .  .  .  .  .  .  }
    81  .  .  .  .  .  .  .  .  .  .  .  .  .  TokPos: main.go:6:4
    82  .  .  .  .  .  .  .  .  .  .  .  .  .  Tok: :=
    83  .  .  .  .  .  .  .  .  .  .  .  .  .  Rhs: []ast.Expr (len = 1) {
    84  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
    85  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  ValuePos: main.go:6:7
    86  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Kind: INT
    87  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Value: "1"
    88  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
    89  .  .  .  .  .  .  .  .  .  .  .  .  .  }
    90  .  .  .  .  .  .  .  .  .  .  .  .  }
    91  .  .  .  .  .  .  .  .  .  .  .  .  1: *ast.AssignStmt {
    92  .  .  .  .  .  .  .  .  .  .  .  .  .  Lhs: []ast.Expr (len = 1) {
    93  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    94  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: main.go:7:2
    95  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "b"
    96  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    97  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    98  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "b"
    99  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 91)
   100  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   101  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   102  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   103  .  .  .  .  .  .  .  .  .  .  .  .  .  TokPos: main.go:7:4
   104  .  .  .  .  .  .  .  .  .  .  .  .  .  Tok: :=
   105  .  .  .  .  .  .  .  .  .  .  .  .  .  Rhs: []ast.Expr (len = 1) {
   106  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
   107  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  ValuePos: main.go:7:7
   108  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Kind: INT
   109  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Value: "2"
   110  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   111  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   112  .  .  .  .  .  .  .  .  .  .  .  .  }
   113  .  .  .  .  .  .  .  .  .  .  .  .  2: *ast.IfStmt {
   114  .  .  .  .  .  .  .  .  .  .  .  .  .  If: main.go:9:2
   115  .  .  .  .  .  .  .  .  .  .  .  .  .  Cond: *ast.BinaryExpr {
   116  .  .  .  .  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
   117  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: main.go:9:5
   118  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "a"
   119  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 74)
   120  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   121  .  .  .  .  .  .  .  .  .  .  .  .  .  .  OpPos: main.go:9:7
   122  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Op: <
   123  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Y: *ast.BasicLit {
   124  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  ValuePos: main.go:9:9
   125  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Kind: INT
   126  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Value: "2"
   127  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   128  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   129  .  .  .  .  .  .  .  .  .  .  .  .  .  Body: *ast.BlockStmt {
   130  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Lbrace: main.go:9:11
   131  .  .  .  .  .  .  .  .  .  .  .  .  .  .  List: []ast.Stmt (len = 1) {
   132  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.AssignStmt {
   133  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Lhs: []ast.Expr (len = 1) {
   134  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
   135  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: main.go:10:3
   136  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "a"
   137  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 74)
   138  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   139  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   140  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  TokPos: main.go:10:5
   141  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Tok: =
   142  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Rhs: []ast.Expr (len = 1) {
   143  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
   144  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: main.go:10:7
   145  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "b"
   146  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 96)
   147  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   148  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   149  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   150  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   151  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Rbrace: main.go:11:2
   152  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   153  .  .  .  .  .  .  .  .  .  .  .  .  }
   154  .  .  .  .  .  .  .  .  .  .  .  .  3: *ast.ReturnStmt {
   155  .  .  .  .  .  .  .  .  .  .  .  .  .  Return: main.go:13:2
   156  .  .  .  .  .  .  .  .  .  .  .  .  .  Results: []ast.Expr (len = 1) {
   157  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.BinaryExpr {
   158  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
   159  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: main.go:13:9
   160  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "a"
   161  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 74)
   162  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   163  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  OpPos: main.go:13:11
   164  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Op: +
   165  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Y: *ast.Ident {
   166  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: main.go:13:13
   167  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "b"
   168  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 96)
   169  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   170  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   171  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   172  .  .  .  .  .  .  .  .  .  .  .  .  }
   173  .  .  .  .  .  .  .  .  .  .  .  }
   174  .  .  .  .  .  .  .  .  .  .  .  Rbrace: main.go:14:1
   175  .  .  .  .  .  .  .  .  .  .  }
   176  .  .  .  .  .  .  .  .  .  }
   177  .  .  .  .  .  .  .  .  }
   178  .  .  .  .  .  .  .  }
   179  .  .  .  .  .  .  .  Lparen: main.go:17:20
   180  .  .  .  .  .  .  .  Ellipsis: -
   181  .  .  .  .  .  .  .  Rparen: main.go:17:21
   182  .  .  .  .  .  .  }
   183  .  .  .  .  .  }
   184  .  .  .  .  .  Ellipsis: -
   185  .  .  .  .  .  Rparen: main.go:17:22
   186  .  .  .  .  }
   187  .  .  .  }
   188  .  .  }
   189  .  .  Rbrace: main.go:18:1
   190  .  }
   191  }
```
</details>
<br>

The corresponding code for this step: https://github.com/golang/go/tree/master/src/cmd/compile/internal/gc
## 3. Static Single Assignment (SSA) Form
SSA form is a way of representing a program's control flow in a more structured way, making it easier to perform optimizations. This step involves converting the AST into SSA form. 

Unoptimized SSA is something like this.

![unoptimized-ssa.png](unoptimized-ssa.png "unoptimized-ssa.png")

In b1 block, `v6` and `v7` is the assignment of 1, 2 to a, b. And `If v8 -> b3 b2` is to decide whether to jump program execution to either `b2`(if true) or `b3`(if false).

The corresponding code for this step; https://github.com/golang/go/tree/master/src/cmd/compile/internal/gc
## 4. SSA Optimization
This step involves performing various optimizations on the code in SSA form to improve its performance. These optimizations can range from simple dead code elimination to more complex techniques like loop unrolling.

Optimized SSA is something like this.

![optimized-ssa.png](optimized-ssa.png "optimized-ssa.png")

In the optimized version, the unnecessary memory allocation, initialization, and copying operations have been removed. Instead, the local variables have been directly referenced by their names. Also, the function arguments have been directly used instead of being stored in local variables. Furthermore, the unnecessary variables `v9` and `v11` has been removed, and the variables `v10` and `v12` have been replaced with constant values.

The corresponding code for this step; https://github.com/golang/go/tree/master/src/cmd/compile/internal/ssa
## 5. Assembly
The SSA that is architecture-independent is converted into assembly that is architecture-dependent.

```
go tool objdump -S main.o > main.objdump
```
```
TEXT "".myFunc(SB) gofile../Users/iamstgt/~/How_gc_Works/main.go
	return a + b
  0x13cd		b804000000		MOVL $0x4, AX		
  0x13d2		c3			RET			

TEXT "".main(SB) gofile../Users/iamstgt/~/How_gc_Works/main.go
func main() {
  0x13d3		493b6610		CMPQ 0x10(R14), SP	[2:2]R_USEIFACE:type.int [2:2]R_USEIFACE:type.*os.File	
  0x13d7		765b			JBE 0x1434		
  0x13d9		4883ec40		SUBQ $0x40, SP		
  0x13dd		48896c2438		MOVQ BP, 0x38(SP)	
  0x13e2		488d6c2438		LEAQ 0x38(SP), BP	
	fmt.Println(myFunc())
  0x13e7		440f117c2428		MOVUPS X15, 0x28(SP)	
  0x13ed		b804000000		MOVL $0x4, AX		
  0x13f2		90			NOPL			
  0x13f3		e800000000		CALL 0x13f8		[1:5]R_CALL:runtime.convT64<1>	
  0x13f8		488d0d00000000		LEAQ 0(IP), CX		[3:7]R_PCREL:type.int		
  0x13ff		48894c2428		MOVQ CX, 0x28(SP)	
  0x1404		4889442430		MOVQ AX, 0x30(SP)	
	return Fprintln(os.Stdout, a...)
  0x1409		488b1d00000000		MOVQ 0(IP), BX		[3:7]R_PCREL:os.Stdout			
  0x1410		488d0500000000		LEAQ 0(IP), AX		[3:7]R_PCREL:go.itab.*os.File,io.Writer	
  0x1417		488d4c2428		LEAQ 0x28(SP), CX	
  0x141c		bf01000000		MOVL $0x1, DI		
  0x1421		4889fe			MOVQ DI, SI		
  0x1424		e800000000		CALL 0x1429		[1:5]R_CALL:fmt.Fprintln	
}
  0x1429		488b6c2438		MOVQ 0x38(SP), BP	
  0x142e		4883c440		ADDQ $0x40, SP		
  0x1432		90			NOPL			
  0x1433		c3			RET			
func main() {
  0x1434		e800000000		CALL 0x1439		[1:5]R_CALL:runtime.morestack_noctxt	
  0x1439		eb98			JMP "".main(SB)		

```
This invovation will generate the main.objdump. Objdump command disassembles executable files and it printes a disassembly of all text symbols (code) in the binary.


## 6. Machine code for a given target
Finally, the assembly is compiled into machine code that can be executed on a target CPU architecture.

The corresponding code for this step (for the x86 architecture); https://github.com/golang/go/tree/master/src/cmd/compile/internal/x86

There's one more step to generate the ready-to-use output. This is not a role of gc, but rather the linker that combines multiple object files into one and generates an executable binary that is executed by CPU.
