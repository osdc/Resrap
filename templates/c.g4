// C template with functions, structs, enums, control flow and arrays.
@words msg   = hello, world, resrap, generated, code, c, compiler, grammar, ready, done
@words digit = 1, 2, 3, 4, 5, 6, 7, 8, 9
@words prefix = compute, process, update, handle, build, check, render, collect, parse, encode, decode, merge, split, find, apply, visit
@words tail   = a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p
@words label  = alpha, beta, gamma, delta, omega

program: headers declarations*<0.75> main function^;

headers: '#include <stdio.h>\n#include <stdlib.h>\n#include <string.h>\n#include <stdint.h>\n#include <stdbool.h>\n\n';

declarations:
      structdecl<0.28>
    | enumdecl<0.18>
    | typedefdecl<0.16>
    | helperfunc<0.38>;

structdecl: 'typedef struct ' tname ' {\n    int id;\n    double value;\n    char label[32];\n    int items[8];\n} ' tname ';\n\n';

enumdecl: 'enum ' ename ' { ' eval ' = 0, ' eval2 ' = 1, ' eval3 ' = 2 };\n\n';

typedefdecl: 'typedef unsigned long ' tname ';\n\n';

helperfunc: 'static int ' fname '(int seed) {\n' functionlocals statement*<0.85> '    return total;\n}\n\n';

main: 'int main(void) {\n    int values[8] = {3, 1, 4, 1, 5, 9, 2, 6};\n    int total = 0;\n    for (int i = 0; i < 8; ++i) {\n        total += values[i];\n    }\n    printf("start: %d\\n", total);\n    printf("done: %s\\n", "' label '");\n    return 0;\n}\n\n';

function: 'int ' fname '(int seed) {\n' functionlocals statement*<0.9> '    return total;\n}\n\n';

functionlocals: '    int count = seed;\n    int total = 0;\n    double value = 1.5;\n    char label[32] = "alpha";\n    int items[8] = {0};\n    for (int i = 0; i < 8; ++i) {\n        items[i] = i + seed;\n    }\n';

statement:
      assignment<0.18>
    | printstmt<0.14>
    | ifstmt<0.14>
    | forstmt<0.12>
    | whilestmt<0.10>
    | dowhile<0.08>
    | switchstmt<0.10>
    | arraystmt<0.14>;

simple:
      assignment<0.40>
    | printstmt<0.30>
    | arraystmt<0.30>;

assignment: '    count += ' number ';\n' | '    total += count * ' number ';\n' | '    value = value * ' float ';\n' | '    count = (count > ' number ') ? count - ' number ' : count + ' number ';\n' | '    snprintf(label, sizeof(label), "%s", "' msg '");\n';

printstmt: '    printf("%s: %d %.2f\\n", label, count, value);\n';

ifstmt: '    if (count > ' number ') {\n' simple+<0.6> '        total += ' number ';\n    } else {\n        total -= 1;\n    }\n';

forstmt: '    for (int i = 0; i < ' number '; ++i) {\n        total += items[i % 8];\n' simple*<0.5> '    }\n';

whilestmt: '    while (count < ' number ') {\n        ++count;\n        total += count;\n' simple*<0.4> '    }\n';

dowhile: '    do {\n        total += count;\n        ++count;\n    } while (count < ' number ');\n';

switchstmt: '    switch (count % 3) {\n    case 0:\n        total += 1;\n        break;\n    case 1:\n        total += 2;\n        break;\n    default:\n        total += 3;\n        break;\n    }\n';

arraystmt: '    items[count % 8] = count + ' number ';\n' | '    total += items[count % 8];\n';

fname: prefix '_' tail tail tail tail;
tname: 'Type_' tail tail tail;
ename: 'Kind_' tail tail tail;
eval: 'VALUE_A';
eval2: 'VALUE_B';
eval3: 'VALUE_C';
number: digit digit?;
float: number '.' number;
