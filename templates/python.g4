@words msg   = hello, world, resrap, generated, code, python, compiler, grammar
@words digit = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9

program: 'def main():\n' stmt+<0.85> '\n\nmain()\n';
stmt: '    print("' msg '")\n'
    | '    print(' number ')\n'
    | '    for i in range(' number '):\n        print(i)\n'
    | '    if ' number ' > ' number ':\n        print("' msg '")\n';
number: digit digit?;
