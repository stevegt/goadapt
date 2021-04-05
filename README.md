# goadapt

Use the Ck/Return pattern from this package to replace the common 
`if err != nil {return err}` boilerplate.  This shortens code and
provides more descriptive error messages, including file and line
number.  

Inspired in part by Dave Cheney's
[https://github.com/pkg/errors](https://github.com/pkg/errors) and Harri
Lainio's [https://github.com/lainio/err2](https://github.com/lainio/err2).
