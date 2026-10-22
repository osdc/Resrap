@words msg   = hello, world, resrap, generated, code, javascript, compiler, grammar
@words digit = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9

program: 'function main() {\n' stmt*<0.8> '}\n\nmain();\n';
stmt: '    console.log("' msg '");\n'
    | '    console.log(' number ');\n'
    | '    for (let i = 0; i < ' number '; i++) {\n        console.log(i);\n    }\n'
    | '    if (' number ' > ' number ') {\n        console.log("' msg '");\n    }\n';
number: digit digit?;
