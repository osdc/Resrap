@words msg   = hello, world, resrap, generated, code, c, compiler, grammar
@words digit = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9

program: '#include <stdio.h>\n\nint main(void) {\n' stmt*<0.8> '    return 0;\n}\n';
stmt: '    printf("' msg '\\n");\n'
    | '    printf("%d\\n", ' number ');\n'
    | '    for (int i = 0; i < ' number '; ++i) {\n        printf("%d\\n", i);\n    }\n'
    | '    if (' number ' > ' number ') {\n        printf("' msg '\\n");\n    }\n';
number: digit digit?;
