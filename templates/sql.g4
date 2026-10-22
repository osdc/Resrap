// SQL template: SELECT, CREATE TABLE and INSERT statements.
@words identifier = users, items, orders, events, records, accounts, products, sessions
@words digit      = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9

program: statement^;

statement
    : selectstmt ';\n'
    | createtablestmt ';\n'
    | insertstmt ';\n';

selectstmt
    : 'SELECT ' selectlist ' FROM ' identifier (' WHERE ' condition)?;

selectlist
    : '*'
    | column (', ' column)*;

column: identifier;

createtablestmt
    : 'CREATE TABLE ' identifier '(' columnlist ')';

columnlist
    : columndef (', ' columndef)*;

columndef
    : identifier ' ' datatype;

insertstmt
    : 'INSERT INTO ' identifier '(' columnlistshort ')' ' VALUES(' valuelist ')';

columnlistshort
    : identifier (', ' identifier)*;

valuelist
    : value (', ' value)*;

value
    : integer
    | float
    | string
    | boolean;

condition
    : identifier comparator value (logicalop identifier comparator value)*;

comparator
    : ' = ' | ' < ' | ' > ' | ' <= ' | ' >= ' | ' != ';

logicalop
    : ' AND ' | ' OR ';

datatype
    : 'INT' | 'FLOAT' | 'TEXT' | 'BOOL';

boolean
    : 'TRUE' | 'FALSE';

string
    : '\'' identifier '\'';

integer
    : digit digit?;

float
    : integer '.' integer;
