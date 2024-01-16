## --- Day 20: Pulse Propagation ---

link: https://adventofcode.com/2023/day/20

Example part 1 output for the example data:
```
Example 1:

# go run *.go example1.txt
2024/01/16 14:32:13 broadcaster -> a, b, c
2024/01/16 14:32:13 %a -> b
2024/01/16 14:32:13 %b -> c
2024/01/16 14:32:13 %c -> inv
2024/01/16 14:32:13 &inv -> a
2024/01/16 14:32:13 map[a:0xc000014080 b:0xc0000140a0 broadcaster:0xc000010030 c:0xc0000140c0 inv:0xc0000140e0]
2024/01/16 14:32:13 (Part 1) After 1000 button pushes: 32000000

Example 2:

# go run *.go example2.txt
2024/01/16 14:32:37 broadcaster -> a
2024/01/16 14:32:37 %a -> inv, con
2024/01/16 14:32:37 &inv -> b
2024/01/16 14:32:37 %b -> con
2024/01/16 14:32:37 &con -> output
2024/01/16 14:32:37 map[a:0xc000014080 b:0xc0000140c0 broadcaster:0xc000010030 con:0xc0000140e0 inv:0xc0000140a0]
2024/01/16 14:32:37 (Part 1) After 1000 button pushes: 11687500
```

Example part 2 output for the example data:
```
N/A Part 2 didn't have an example
```