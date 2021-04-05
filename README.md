# goadapt

Use the Ck/Return pattern from this package to replace the common `if
err != nil {return err}` boilerplate. 

- Allows shorter, more readable code.
- Provides more descriptive error messages.  
- Includes file name and line number in error strings by default.  
- Allows optional per-error and per-function annotations.  
- Allows optional format strings for annotations.  

Inspired in part by Dave Cheney's
[https://github.com/pkg/errors](https://github.com/pkg/errors) and Harri
Lainio's [https://github.com/lainio/err2](https://github.com/lainio/err2).
