@words msg   = hello, world, resrap, generated, code, java, compiler, grammar
@words digit = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9

program: 'public class Main {\n    public static void main(String[] args) {\n' stmt*<0.8> '    }\n}\n';
stmt: '        System.out.println("' msg '");\n'
    | '        System.out.println(' number ');\n'
    | '        for (int i = 0; i < ' number '; i++) {\n            System.out.println(i);\n        }\n'
    | '        if (' number ' > ' number ') {\n            System.out.println("' msg '");\n        }\n';
number: digit digit?;
