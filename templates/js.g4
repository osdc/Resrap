// JavaScript template with functions, classes, arrays, maps and exceptions.
@words msg   = hello, world, resrap, generated, code, javascript, compiler, grammar, ready, done
@words digit = 1, 2, 3, 4, 5, 6, 7, 8, 9
@words prefix = compute, process, update, handle, build, check, render, collect, parse, encode, decode, merge, split, find, apply, visit
@words tail   = a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p
@words label  = alpha, beta, gamma, delta, omega

program: header main declaration^ '\nmain();\n';

header: '// generated JavaScript\n\n';

main: 'function main() {\n    const values = [5, 2, 8, 1];\n    values.sort((a, b) => a - b);\n    const total = values.reduce((a, b) => a + b, 0);\n    const lookup = new Map();\n    lookup.set("start", total);\n    console.log("start:", values, lookup.get("start"));\n}\n';

declaration:
      function<0.46>
    | asyncfunction<0.14>
    | classdecl<0.24>
    | constdecl<0.16>;

function: 'function ' fname '(seed) {\n' functionlocals statement*<0.9> '    return total;\n}\n\n';

asyncfunction: 'async function ' fname '(seed) {\n' functionlocals statement*<0.8> '    await Promise.resolve(count);\n    return total;\n}\n\n';

classdecl: 'class ' tname ' {\n    constructor(seed) {\n        this.seed = seed;\n        this.label = "alpha";\n    }\n\n    ' methodname '(amount) {\n' classbody '        return this.seed + amount;\n    }\n}\n\n';

constdecl: 'const ' cname ' = {\n    id: ' number ',\n    name: "' label '",\n    value: ' float ',\n    items: [1, 2, 3],\n    describe() {\n        return this.name + ":" + this.id;\n    }\n};\n\n';

functionlocals: '    let count = seed;\n    let total = 0;\n    const values = [];\n    const lookup = new Map();\n    const label = "alpha";\n';

classbody: '        let count = amount;\n        let total = this.seed;\n        const values = [];\n        const lookup = new Map();\n        const label = this.label;\n';

statement:
      assign<0.18>
    | printstmt<0.12>
    | ifstmt<0.14>
    | forstmt<0.12>
    | forof<0.10>
    | whilestmt<0.08>
    | switchstmt<0.10>
    | arraystmt<0.10>
    | mapstmt<0.08>
    | trycatch<0.08>;

simple:
      assign<0.40>
    | printstmt<0.30>
    | arraystmt<0.30>;

assign: '    count += ' number ';\n' | '    total += count * ' number ';\n' | '    total = total * 2 + count;\n';

printstmt: '    console.log(label, count, total);\n';

ifstmt: '    if (count > ' number ') {\n' simple+<0.6> '        total += ' number ';\n    } else {\n        total -= 1;\n    }\n';

forstmt: '    for (let i = 0; i < ' number '; i++) {\n        total += i;\n' simple*<0.5> '    }\n';

forof: '    for (const value of values) {\n        total += value;\n    }\n';

whilestmt: '    while (count < ' number ') {\n        count += 1;\n        total += count;\n    }\n';

switchstmt: '    switch (count % 3) {\n        case 0:\n            total += 1;\n            break;\n        case 1:\n            total += 2;\n            break;\n        default:\n            total += 3;\n            break;\n    }\n';

arraystmt: '    values.push(count + ' number ');\n' | '    total += values.length;\n' | '    values.sort((a, b) => a - b);\n';

mapstmt: '    lookup.set(label, count);\n' | '    total += lookup.get(label) || 0;\n';

trycatch: '    try {\n        if (count < 0) throw new Error("negative");\n        total += count;\n    } catch (error) {\n        total += ' number ';\n    }\n';

fname: prefix '_' tail tail tail tail;
tname: 'Type_' tail tail tail tail;
cname: 'CONFIG_' tail tail tail;
methodname: 'score_' tail tail;
number: digit digit?;
float: number '.' number;
