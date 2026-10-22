// Haskell template with functions, pattern matching, lists, data and type aliases.
@words msg   = hello, world, resrap, generated, code, haskell, compiler, grammar, ready, done
@words digit = 1, 2, 3, 4, 5, 6, 7, 8, 9
@words prefix = compute, process, update, handle, build, check, render, collect, parse, encode, decode, merge, split, find, apply, visit
@words tail   = a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p
@words label  = alpha, beta, gamma, delta, omega

program: module imports declaration^ maindecl;

module: 'module Main where\n\n';

imports: 'import Data.List (sort, nub)\nimport Data.Maybe (fromMaybe)\n\n';

maindecl: 'data Counter = Counter { value :: Int }\n\nclass Describable a where\n    describe :: a -> String\n\ninstance Describable Counter where\n    describe (Counter v) = "counter=" ++ show v\n\nmain :: IO ()\nmain = do\n    let values = [5, 2, 8, 1, 3]\n    print (sort values)\n    print (sum values)\n    print (nub [1, 2, 1, 3, 2])\n    print (fromMaybe 0 (Just (sum values)))\n    print (describe (Counter (sum values)))\n';

declaration:
      function<0.50>
    | casefunction<0.20>
    | listfunction<0.18>
    | typededecl<0.12>;

function: fname ' seed =\n    let count = seed + ' number '\n        total = count * ' number '\n        values = [1 .. ' number ']\n    in if count > ' number '\n        then total + sum values\n        else total - count\n\n';

casefunction: fname ' seed =\n    case seed `mod` 3 of\n        0 -> seed + ' number '\n        1 -> seed * ' number '\n        _ -> seed - ' number '\n\n';

listfunction: fname ' seed = sum (map (\\x -> x * ' number ' + seed) (filter even [1 .. ' number ']))\n\n';

typededecl: 'type ' tname ' = [Int]\n\n';

fname: prefix '_' tail tail tail tail;
tname: 'Type_' tail tail tail tail;
number: digit digit?;
