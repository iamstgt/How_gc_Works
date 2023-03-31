# How_gc_Works
This aims to help Go developers gain a deeper understanding of the inner workings of gc.

## 1. Lexical Analysis and Syntax Analysis
The compiler analyzes the source code and converts it into a form that can be easily processed.

The corresponding code for this step; https://github.com/golang/go/tree/master/src/cmd/compile/internal/syntax

## 2. Abstract Syntax Tree (AST)
The AST is a tree-like data structure that represents the structure of a program's source code in a more abstract way. This step involves constructing the AST from the parsed source code.

The corresponding code for this step: https://github.com/golang/go/tree/master/src/cmd/compile/internal/gc
## 3. Static Single Assignment (SSA) Form
SSA form is a way of representing a program's control flow in a more structured way, making it easier to perform optimizations. This step involves converting the AST into SSA form. 

The corresponding code for this step; https://github.com/golang/go/tree/master/src/cmd/compile/internal/gc
## 4. SSA Optimization
This step involves performing various optimizations on the code in SSA form to improve its performance. These optimizations can range from simple dead code elimination to more complex techniques like loop unrolling.

The corresponding code for this step; https://github.com/golang/go/tree/master/src/cmd/compile/internal/ssa
## 5. Machine code for a given target
Finally, the optimized code is compiled into machine code that can be executed on a specific target platform.

The corresponding code for this step (for the x86 architecture); https://github.com/golang/go/tree/master/src/cmd/compile/internal/x86
