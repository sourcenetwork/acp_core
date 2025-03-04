grammar PermissionExpr

expr: term | term operator expr;
term: relation | ttu_term | '(' expr ')';
ttu_term: relation '->' operator;

relation: IDENTIFIER;
resource: IDENTIFIER;
operator: '+' #union
        | '-' #difference
        | '&' #intersection
        ;

IDENTIFIER: [a-zA-Z] [a-zA-Z0-9_]+;