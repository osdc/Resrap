// Full C++ demo grammar.
//
// The `^` at the end of `program` repeats complete generated functions until
// the token budget is reached. All generated translation units are expected to
// pass `g++ -std=c++17 -Wall -Wextra -fsyntax-only`.
@words message = hello, world, resrap, generated, code, cpp, compiler, grammar
@words prefix  = compute, process, update, handle, build, check, render, collect
@words label   = alpha, beta, gamma, delta, omega
@words tail    = a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p
@words digit   = 1, 2, 3, 4, 5

// ---------------------------------------------------------------------------
// Top level: one fixed preamble, then as many valid functions as fit the budget.
// ---------------------------------------------------------------------------
program: header classdecl helper main function^;

header: '#include <iostream>\n#include <string>\n#include <vector>\n#include <algorithm>\n#include <stdexcept>\n\n';

classdecl: 'class Box {\npublic:\n    int value;\n    explicit Box(int v) : value(v) {}\n    int get() const { return value; }\n};\n\n';

helper: 'int helper(int x) {\n    return x * 2 + 1;\n}\n\n';

main: 'int main() {\n    Box box(helper(3));\n    std::cout << "start: " << box.get() << std::endl;\n    std::cout << "done: " << "' label '" << std::endl;\n    return 0;\n}\n\n';

// ---------------------------------------------------------------------------
// A complete function. Its body is a probabilistic sequence of statements.
// ---------------------------------------------------------------------------
function: 'int ' fname '(int seed) {\n    int count = seed;\n    double value = 1.5;\n    std::string label = "' label '";\n    std::vector<int> items;\n    bool ready = true;\n    int guard = 0;\n' statement*<0.8> '    return count;\n}\n\n';

fname: prefix '_' tail tail tail tail;

// ---------------------------------------------------------------------------
// Statements: weighted alternatives showing most common C++ constructs.
// ---------------------------------------------------------------------------
statement:
      assignment<0.18>
    | printstmt<0.12>
    | ifstmt<0.12>
    | forstmt<0.10>
    | whilestmt<0.08>
    | pushstmt<0.10>
    | rangeforstmt<0.07>
    | callstmt<0.07>
    | lambdastmt<0.06>
    | switchstmt<0.05>
    | stringstmt<0.05>;

// Statements safe to nest inside a loop or branch.
simple:
      assignment<0.35>
    | printstmt<0.25>
    | pushstmt<0.20>
    | callstmt<0.20>;

assignment: '    count = count + ' number ';\n' | '    value = value * ' float ';\n' | '    ready = (count % 2 == 0) ? true : false;\n';

printstmt: '    std::cout << "' message '" << " " << count << std::endl;\n';

ifstmt: '    if (count > ' number ') {\n' simple+<0.5> '        count += 1;\n    } else {\n        count -= 1;\n    }\n';

forstmt: '    for (int i = 0; i < ' number '; ++i) {\n        count += i;\n' simple*<0.5> '    }\n';

whilestmt: '    guard = ' number ';\n    while (guard > 0) {\n        --guard;\n        count += 1;\n' simple*<0.5> '    }\n';

pushstmt: '    items.push_back(' number ');\n';

rangeforstmt: '    for (const auto& item : items) {\n        count += item;\n    }\n';

callstmt: '    std::sort(items.begin(), items.end());\n    value = std::max(value, ' float ');\n';

lambdastmt: '    {\n        auto add = [](int a, int b) { return a + b; };\n        count = add(count, ' number ');\n    }\n';

switchstmt: '    switch (count % 3) {\n    case 0:\n        count += 1;\n        break;\n    case 1:\n        count += 2;\n        break;\n    default:\n        count += 3;\n        break;\n    }\n';

stringstmt: '    label += "' message '";\n    if (label.find("' message '") != std::string::npos) {\n        count += 1;\n    }\n';

// ---------------------------------------------------------------------------
// Lexical helpers
// ---------------------------------------------------------------------------
number: digit digit?;
float: number '.' number;
