// Python template with functions, classes, comprehensions, exceptions and loops.
@words msg   = hello, world, resrap, generated, code, python, compiler, grammar, ready, done
@words digit = 1, 2, 3, 4, 5, 6, 7, 8, 9
@words prefix = compute, process, update, handle, build, check, render, collect, parse, encode, decode, merge, split, find, apply, visit
@words tail   = a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p
@words label  = alpha, beta, gamma, delta, omega

program: imports main definition^ footer;

imports: 'import math\nimport random\n\n';

main: 'def main():\n    values = [3, 1, 4, 1, 5]\n    values.sort()\n    total = sum(values)\n    print("start:", values, total)\n    print("sqrt:", round(math.sqrt(total), 2))\n    print("random:", random.randint(1, 9))\n\n';

footer: '\nif __name__ == "__main__":\n    main()\n';

definition:
      function<0.60>
    | classdecl<0.40>;

function: 'def ' fname '(seed):\n' functionlocals statement*<0.9> '    return total\n\n';

classdecl: 'class ' tname ':\n    def __init__(self, seed):\n        self.seed = seed\n        self.label = "alpha"\n\n    def ' methodname '(self, amount):\n        count = self.seed + amount\n        total = 0\n        values = [count, amount, self.seed]\n        for value in values:\n            total += value\n        return total + len(self.label)\n\n';

functionlocals: '    count = seed\n    total = 0\n    values = []\n    lookup = {}\n    label = "alpha"\n';

statement:
      assign<0.18>
    | printstmt<0.12>
    | ifstmt<0.14>
    | forstmt<0.12>
    | foreach<0.10>
    | whilestmt<0.08>
    | trycatch<0.10>
    | liststmt<0.10>
    | dictstmt<0.06>
    | comprehension<0.10>;

nested:
      '        total += count * ' number '\n'<0.35>
    | '        values.append(count + ' number ')\n'<0.30>
    | '        print(label, count, total)\n'<0.35>;

assign: '    count += ' number '\n' | '    total += count * ' number '\n' | '    total = total * 2 + count\n';

printstmt: '    print(label, count, total)\n';

ifstmt: '    if count > ' number ':\n' nested*<0.6> '        total += ' number '\n    else:\n        total -= 1\n';

forstmt: '    for i in range(' number '):\n        total += i\n' nested*<0.5>;

foreach: '    for value in values:\n        total += value\n';

whilestmt: '    while count < ' number ':\n        count += 1\n        total += count\n' nested*<0.4>;

trycatch: '    try:\n        total += int("' number '")\n    except ValueError:\n        total += 1\n';

liststmt: '    values.append(count + ' number ')\n' | '    total += len(values)\n';

dictstmt: '    lookup[label] = count\n' | '    total += lookup.get(label, 0)\n';

comprehension: '    total += sum(x * x for x in range(' number '))\n' | '    total += len([x for x in values if x > ' number '])\n';

fname: prefix '_' tail tail tail tail;
tname: 'Type_' tail tail tail tail;
methodname: 'score_' tail tail;
number: digit digit?;
