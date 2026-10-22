// C++ template with classes, containers, lambdas, exceptions and templates.
@words message = hello, world, resrap, generated, code, cpp, compiler, grammar, ready, done
@words prefix  = compute, process, update, handle, build, check, render, collect, parse, encode, decode, merge, split, find, apply, visit
@words label   = alpha, beta, gamma, delta, omega
@words tail    = a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p
@words digit   = 1, 2, 3, 4, 5

program: header classdecl helper main function^;

header: '#include <iostream>\n#include <string>\n#include <vector>\n#include <map>\n#include <memory>\n#include <algorithm>\n#include <numeric>\n#include <sstream>\n#include <stdexcept>\n\n';

classdecl: 'class Box {\npublic:\n    int value;\n    std::string name;\n    explicit Box(int v, std::string n) : value(v), name(std::move(n)) {}\n    int get() const { return value; }\n    void add(int n) { value += n; }\n    std::string describe() const { return name + ":" + std::to_string(value); }\n};\n\n';

helper: 'int helper(int x) {\n    return x * 2 + 1;\n}\n\ntemplate <typename T>\nT clamp_value(T value, T low, T high) {\n    return value < low ? low : (value > high ? high : value);\n}\n\n';

main: 'int main() {\n    Box box(helper(3), "' label '");\n    std::vector<int> values = {3, 1, 4, 1, 5};\n    std::sort(values.begin(), values.end());\n    int total = clamp_value(std::accumulate(values.begin(), values.end(), 0), 0, 100);\n    std::cout << "start: " << box.describe() << " " << total << std::endl;\n    std::cout << "done: " << "' message '" << std::endl;\n    return 0;\n}\n\n';

function: 'int ' fname '(int seed) {\n' functionlocals statement*<0.9> '    return total;\n}\n\n';

functionlocals: '    int count = seed;\n    int total = 0;\n    double value = 1.5;\n    std::string label = "' label '";\n    std::vector<int> items;\n    std::map<std::string, int> lookup;\n    std::unique_ptr<int> held = std::make_unique<int>(seed);\n    bool ready = true;\n    int guard = 0;\n';

statement:
      assignment<0.16>
    | printstmt<0.12>
    | ifstmt<0.12>
    | forstmt<0.10>
    | whilestmt<0.08>
    | pushstmt<0.10>
    | rangeforstmt<0.08>
    | callstmt<0.08>
    | lambdastmt<0.07>
    | switchstmt<0.07>
    | stringstmt<0.07>
    | trycatch<0.05>;

simple:
      assignment<0.35>
    | printstmt<0.25>
    | pushstmt<0.20>
    | callstmt<0.20>;

assignment: '    count = count + ' number ';\n' | '    value = value * ' float ';\n' | '    ready = (count % 2 == 0) ? true : false;\n' | '    total += *held + ' number ';\n';

printstmt: '    std::cout << label << ": " << count << " " << value << std::endl;\n';

ifstmt: '    if (count > ' number ') {\n' simple+<0.5> '        count += 1;\n    } else {\n        count -= 1;\n    }\n';

forstmt: '    for (int i = 0; i < ' number '; ++i) {\n        count += i;\n' simple*<0.5> '    }\n';

whilestmt: '    guard = ' number ';\n    while (guard > 0) {\n        --guard;\n        count += 1;\n' simple*<0.5> '    }\n';

pushstmt: '    items.push_back(' number ');\n' | '    lookup[label] = count;\n';

rangeforstmt: '    for (const auto& item : items) {\n        count += item;\n    }\n';

callstmt: '    std::sort(items.begin(), items.end());\n    value = std::max(value, ' float ');\n' | '    total = std::accumulate(items.begin(), items.end(), 0);\n';

lambdastmt: '    {\n        auto add = [](int a, int b) { return a + b; };\n        count = add(count, ' number ');\n    }\n';

switchstmt: '    switch (count % 3) {\n    case 0:\n        count += 1;\n        break;\n    case 1:\n        count += 2;\n        break;\n    default:\n        count += 3;\n        break;\n    }\n';

stringstmt: '    label += "' message '";\n    if (label.find("' message '") != std::string::npos) {\n        count += 1;\n    }\n' | '    std::ostringstream out;\n    out << label << ":" << count;\n    label = out.str();\n';

trycatch: '    try {\n        if (count < 0) {\n            throw std::runtime_error("negative count");\n        }\n        total += count;\n    } catch (const std::exception& e) {\n        total += ' number ';\n    }\n';

fname: prefix '_' tail tail tail tail;
number: digit digit?;
float: number '.' number;
