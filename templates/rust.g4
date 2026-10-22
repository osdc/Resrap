@words msg   = hello, world, resrap, generated, code, rust, compiler, grammar
@words digit = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9

program: 'fn main() {\n' stmt*<0.8> '}\n';
stmt: '    println!("' msg '");\n'
    | '    println!("{}", ' number ');\n'
    | '    for i in 0..' number ' {\n        println!("{}", i);\n    }\n'
    | '    if ' number ' > ' number ' {\n        println!("' msg '");\n    }\n';
number: digit digit?;
