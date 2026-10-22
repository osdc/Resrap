// Go template with functions, methods, structs, interfaces, maps and slices.
@words msg   = hello, world, resrap, generated, code, go, compiler, grammar, ready, done
@words digit = 1, 2, 3, 4, 5, 6, 7, 8, 9
@words prefix = compute, process, update, handle, build, check, render, collect, parse, encode, decode, merge, split, find, apply, visit
@words tail   = a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p
@words label  = alpha, beta, gamma, delta, omega

program: preamble declaration^;

preamble: 'package main\n\nimport (\n    "fmt"\n    "sort"\n    "strings"\n)\n\ntype Counter struct {\n    Value int\n}\n\nfunc (c Counter) Describe() string {\n    return fmt.Sprintf("counter=%d", c.Value)\n}\n\nfunc main() {\n    values := []int{5, 2, 8, 1}\n    sort.Ints(values)\n    fmt.Println(strings.Join([]string{"start", "go"}, " "), values)\n    c := Counter{Value: values[0]}\n    fmt.Println(c.Describe())\n}\n\n';

declaration:
      function<0.42>
    | method<0.20>
    | structdecl<0.18>
    | interfacedecl<0.12>
    | typededecl<0.08>;

function: 'func ' fname '(seed int) int {\n' functionlocals statement*<0.9> '    return total\n}\n\n';

method: 'func (c Counter) ' fname '(seed int) int {\n' functionlocals statement*<0.9> '    return total + c.Value\n}\n\n';

structdecl: 'type ' tname ' struct {\n    ID int\n    Name string\n    Value float64\n    Items []int\n}\n\n';

interfacedecl: 'type ' tname ' interface {\n    Describe() string\n    Score(seed int) int\n}\n\n';

typededecl: 'type ' tname ' []int\n\n';

functionlocals: '    count := seed\n    total := 0\n    values := []int{count}\n    lookup := map[string]int{"start": count}\n    label := "alpha"\n    _ = values\n    _ = lookup\n    _ = label\n';

statement:
      assign<0.18>
    | printstmt<0.12>
    | ifstmt<0.14>
    | forstmt<0.12>
    | rangefor<0.10>
    | switchstmt<0.10>
    | slicestmt<0.10>
    | mapstmt<0.08>
    | stringstmt<0.08>
    | sortstmt<0.08>;

simple:
      assign<0.40>
    | printstmt<0.30>
    | slicestmt<0.30>;

assign: '    count += ' number '\n' | '    total += count * ' number '\n' | '    total = total*2 + count\n';

printstmt: '    fmt.Println(label, count, total)\n';

ifstmt: '    if count > ' number ' {\n' simple+<0.6> '        total += ' number '\n    } else {\n        total -= 1\n    }\n';

forstmt: '    for i := 0; i < ' number '; i++ {\n        total += i\n' simple*<0.5> '    }\n';

rangefor: '    for _, value := range values {\n        total += value\n    }\n';

switchstmt: '    switch count % 3 {\n    case 0:\n        total += 1\n    case 1:\n        total += 2\n    default:\n        total += 3\n    }\n';

slicestmt: '    values = append(values, count + ' number ')\n' | '    total += len(values)\n';

mapstmt: '    lookup[label] = count\n' | '    total += lookup[label]\n' | '    delete(lookup, "missing")\n';

stringstmt: '    label = strings.TrimSpace(label + " ' msg '")\n' | '    total += strings.Count(label, "a")\n';

sortstmt: '    sort.Ints(values)\n    total += values[0]\n';

fname: prefix '_' tail tail tail tail;
tname: 'Type_' tail tail tail tail;
number: digit digit?;
