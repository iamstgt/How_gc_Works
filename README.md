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
