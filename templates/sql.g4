// SQL template with DDL, DML, joins, grouping, ordering and transactions.
@words table  = users, items, orders, events, records, accounts, products, sessions, invoices, payments
@words column = id, name, value, status, total, created_at, updated_at, email, amount, quantity
@words text   = hello, world, resrap, generated, code, sql, database, query, active, pending
@words digit  = 1, 2, 3, 4, 5, 6, 7, 8, 9

program: statement^;

statement:
      selectstmt<0.28>
    | joinstmt<0.18>
    | insertstmt<0.16>
    | createstmt<0.14>
    | updatestmt<0.10>
    | deletestmt<0.06>
    | indexstmt<0.04>
    | transaction<0.04>;

selectstmt: 'SELECT ' selectlist ' FROM ' table (' WHERE ' condition)? (' ORDER BY ' column (' ASC' | ' DESC'))? (' LIMIT ' number)? ';\n';

joinstmt: 'SELECT ' table '.' column ', ' table '.' column ' FROM ' table ' JOIN ' table ' ON ' table '.id = ' table '.id WHERE ' condition ';\n';

selectlist: '*' | column (', ' column)*;

insertstmt: 'INSERT INTO ' table ' (' columnlist ') VALUES (' valuelist ')' (';\n' | ' RETURNING id;\n');

createstmt: 'CREATE TABLE ' table ' (\n    id INTEGER PRIMARY KEY,\n    name TEXT NOT NULL,\n    value INTEGER DEFAULT ' number ',\n    status REAL\n);\n';

updatestmt: 'UPDATE ' table ' SET ' column ' = ' value ', ' column ' = ' value ' WHERE ' condition ';\n';

deletestmt: 'DELETE FROM ' table ' WHERE ' condition ';\n';

indexstmt: 'CREATE INDEX idx_' table '_' column ' ON ' table ' (' column ');\n';

transaction: 'BEGIN;\n' | 'COMMIT;\n' | 'ROLLBACK;\n';

columnlist: column (', ' column)*;

valuelist: value (', ' value)*;

value: integer | float | string | boolean | 'NULL';

condition: column comparator value (logicalop column comparator value)*;

comparator: ' = ' | ' < ' | ' > ' | ' <= ' | ' >= ' | ' != ' | ' LIKE ';

logicalop: ' AND ' | ' OR ';

string: '\'' text '\'';

boolean: 'TRUE' | 'FALSE';

integer: digit digit?;

float: integer '.' integer;

number: digit digit?;
