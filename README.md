# goadapt

The Ck/Return pattern from this package can be used to replace the
common `if err != nil {return err}` boilerplate.  See
examples/ck-return.go for usage.

- Allows shorter, more readable code.
- Provides more descriptive error messages.  
- Includes file name and line number in error strings by default.  
- Allows optional per-error and per-function annotations.  
- Allows optional format strings for annotations.  

Inspired in part by Dave Cheney's
[https://github.com/pkg/errors](https://github.com/pkg/errors) and Harri
Lainio's [https://github.com/lainio/err2](https://github.com/lainio/err2).

This package includes several other possibly-useful functions,
including Assert(), Debug(), and Pprint(), as well as a few shortcut
vars for common `fmt` functions:

```
    Pl  = fmt.Println
    Pf  = fmt.Printf
    Spf = fmt.Sprintf
    Fpf = fmt.Fprintf
```

Docs are admittedly still thin -- see ./examples here and 
[elsewhere on github](https://github.com/search?q=github.com%2Fstevegt%2Fgoadapt+extension%3A.go&type=Code&ref=advsearch&l=&l=)
for now.

## TODO

- add more doc comments
- integrate examples into test cases
- add more examples for full coverage
