@words msg   = hello, world, resrap, generated, code, haskell, compiler, grammar
@words digit = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9

program: 'module Main where\n\nmain :: IO ()\nmain = do\n' stmt*<0.8>;
stmt: '    putStrLn "' msg '"\n'
    | '    print (' number ')\n'
    | '    mapM_ print [1 .. ' number ']\n'
    | '    when (' number ' > ' number ') $ putStrLn "' msg '"\n';
number: digit digit?;
