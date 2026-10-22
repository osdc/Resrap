// Rust template with functions, structs, enums, traits and collections.
@words msg   = hello, world, resrap, generated, code, rust, compiler, grammar, ready, done
@words digit = 1, 2, 3, 4, 5, 6, 7, 8, 9
@words prefix = compute, process, update, handle, build, check, render, collect, parse, encode, decode, merge, split, find, apply, visit
@words tail   = a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p
@words label  = alpha, beta, gamma, delta, omega

program: preamble item^;

preamble: 'use std::collections::HashMap;\n\n#[derive(Debug, Clone)]\nstruct Counter {\n    value: i32,\n}\n\nfn main() {\n    let mut values = vec![5, 2, 8, 1];\n    values.sort();\n    let total: i32 = values.iter().sum();\n    println!("start: {:?} {}", values, total);\n    let mut lookup = HashMap::new();\n    lookup.insert(String::from("start"), total);\n    let counter = Counter { value: total };\n    println!("counter: {:?}", counter);\n}\n\n';

item:
      function<0.40>
    | method<0.18>
    | structdecl<0.16>
    | enumdecl<0.12>
    | traitdecl<0.14>;

function: 'fn ' fname '(seed: i32) -> i32 {\n' functionlocals statement*<0.9> '    total\n}\n\n';

method: 'impl Counter {\n    fn ' fname '(&self, seed: i32) -> i32 {\n' functionlocals statement*<0.85> '        total + self.value\n    }\n}\n\n';

structdecl: 'struct ' tname ' {\n    id: i32,\n    name: String,\n    value: f64,\n    items: Vec<i32>,\n}\n\n';

enumdecl: 'enum ' tname ' {\n    Alpha,\n    Beta,\n    Gamma(i32),\n    Delta { value: i32 },\n}\n\n';

traitdecl: 'trait ' tname ' {\n    fn describe(&self) -> String;\n    fn score(&self, seed: i32) -> i32;\n}\n\n';

functionlocals: '    let mut count = seed;\n    let mut total = 0;\n    let mut values = vec![count];\n    let mut lookup: HashMap<String, i32> = HashMap::new();\n    let label = String::from("alpha");\n';

statement:
      assign<0.18>
    | printstmt<0.12>
    | ifstmt<0.14>
    | forstmt<0.12>
    | rangefor<0.10>
    | matchstmt<0.10>
    | vecstmt<0.12>
    | mapstmt<0.08>
    | stringstmt<0.08>
    | sortstmt<0.06>;

simple:
      assign<0.40>
    | printstmt<0.30>
    | vecstmt<0.30>;

assign: '    count += ' number ';\n' | '    total += count * ' number ';\n' | '    total = total * 2 + count;\n';

printstmt: '    println!("{} {} {}", label, count, total);\n';

ifstmt: '    if count > ' number ' {\n' simple+<0.6> '        total += ' number ';\n    } else {\n        total -= 1;\n    }\n';

forstmt: '    for i in 0..' number ' {\n        total += i;\n' simple*<0.5> '    }\n';

rangefor: '    for value in &values {\n        total += *value;\n    }\n';

matchstmt: '    match count % 3 {\n        0 => total += 1,\n        1 => total += 2,\n        _ => total += 3,\n    }\n';

vecstmt: '    values.push(count + ' number ');\n' | '    total += values.len() as i32;\n';

mapstmt: '    lookup.insert(label.clone(), count);\n' | '    total += lookup.get("start").copied().unwrap_or(0);\n';

stringstmt: '    let text = format!("{}:{}", label, count);\n    total += text.len() as i32;\n' | '    total += label.matches("a").count() as i32;\n';

sortstmt: '    values.sort();\n    total += values[0];\n';

fname: prefix '_' tail tail tail tail;
tname: 'Type_' tail tail tail tail;
number: digit digit?;
