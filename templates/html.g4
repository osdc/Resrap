@words msg = hello, world, resrap, generated, code, html, markup, grammar

program: '<!DOCTYPE html>\n<html>\n<head>\n    <meta charset="utf-8">\n    <title>' msg '</title>\n</head>\n<body>\n' elem*<0.85> '</body>\n</html>\n';
elem: '    <h1>' msg '</h1>\n'
    | '    <p>' msg '</p>\n'
    | '    <a href="#">' msg '</a>\n'
    | '    <ul>\n' item*<0.75> '    </ul>\n';
item: '        <li>' msg '</li>\n';
