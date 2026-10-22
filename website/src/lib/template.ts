export const samplegrammer = `@words message = hello, world, resrap, generated, code, grammar
@words name = alpha, beta, gamma, delta, omega
@words digit = 1, 2, 3, 4, 5, 6, 7, 8, 9

program: '#include <stdio.h>\\n\\nint main(void) {\\n    int count = 0;\\n' statement^ '    printf("done: %d\\\\n", count);\\n    return 0;\\n}\\n';

statement:
      '    printf("' message ': %d\\\\n", count);\\n'<0.30>
    | '    count += ' number ';\\n'<0.30>
    | '    if (count > ' number ') {\\n        count -= 1;\\n    } else {\\n        count += 1;\\n    }\\n'<0.20>
    | '    for (int i = 0; i < ' number '; ++i) {\\n        count += i;\\n    }\\n'<0.20>;

number: digit digit?;`;
