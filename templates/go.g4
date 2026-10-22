@words msg   = hello, world, resrap, generated, code, go, compiler, grammar
@words digit = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9

program: 'package main\n\nimport "fmt"\n\nfunc main() {\n' stmt*<0.8> '}\n';
stmt: '    fmt.Println("' msg '")\n'
    | '    fmt.Println(' number ')\n'
    | '    for i := 0; i < ' number '; i++ {\n        fmt.Println(i)\n    }\n'
    | '    if ' number ' > ' number ' {\n        fmt.Println("' msg '")\n    }\n';
number: digit digit?;
