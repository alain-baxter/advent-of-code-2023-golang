## --- Day 8: Haunted Wasteland ---

link: https://adventofcode.com/2023/day/8

Example part 1 output for the two example data:
```
# go run main.go example1.txt
2024/01/02 15:38:30 RL
2024/01/02 15:38:30
2024/01/02 15:38:30 AAA = (BBB, CCC)
2024/01/02 15:38:30 BBB = (DDD, EEE)
2024/01/02 15:38:30 CCC = (ZZZ, GGG)
2024/01/02 15:38:30 DDD = (DDD, DDD)
2024/01/02 15:38:30 EEE = (EEE, EEE)
2024/01/02 15:38:30 GGG = (GGG, GGG)
2024/01/02 15:38:30 ZZZ = (ZZZ, ZZZ)
2024/01/02 15:38:30 Directions: [R L]
2024/01/02 15:38:30 Node Map: map[AAA:{AAA BBB CCC} BBB:{BBB DDD EEE} CCC:{CCC ZZZ GGG} DDD:{DDD DDD DDD} EEE:{EEE EEE EEE} GGG:{GGG GGG GGG} ZZZ:{ZZZ ZZZ ZZZ}]
2024/01/02 15:38:30 Current Node: {AAA BBB CCC}
2024/01/02 15:38:30 Current Node: {CCC ZZZ GGG}
2024/01/02 15:38:30 final Node reached: {ZZZ ZZZ ZZZ}
2024/01/02 15:38:30 Steps: 2

# go run main.go example2.txt
2024/01/02 15:38:47 LLR
2024/01/02 15:38:47
2024/01/02 15:38:47 AAA = (BBB, BBB)
2024/01/02 15:38:47 BBB = (AAA, ZZZ)
2024/01/02 15:38:47 ZZZ = (ZZZ, ZZZ)
2024/01/02 15:38:47 Directions: [L L R]
2024/01/02 15:38:47 Node Map: map[AAA:{AAA BBB BBB} BBB:{BBB AAA ZZZ} ZZZ:{ZZZ ZZZ ZZZ}]
2024/01/02 15:38:47 Current Node: {AAA BBB BBB}
2024/01/02 15:38:47 Current Node: {BBB AAA ZZZ}
2024/01/02 15:38:47 Current Node: {AAA BBB BBB}
2024/01/02 15:38:47 Current Node: {BBB AAA ZZZ}
2024/01/02 15:38:47 Current Node: {AAA BBB BBB}
2024/01/02 15:38:47 Current Node: {BBB AAA ZZZ}
2024/01/02 15:38:47 final Node reached: {ZZZ ZZZ ZZZ}
2024/01/02 15:38:47 Steps: 6
```

Example part 2 output for the example data:
```
# go run main.go example3.txt
2024/01/03 00:52:36 LR
2024/01/03 00:52:36
2024/01/03 00:52:36 11A = (11B, XXX)
2024/01/03 00:52:36 11B = (XXX, 11Z)
2024/01/03 00:52:36 11Z = (11B, XXX)
2024/01/03 00:52:36 22A = (22B, XXX)
2024/01/03 00:52:36 22B = (22C, 22C)
2024/01/03 00:52:36 22C = (22Z, 22Z)
2024/01/03 00:52:36 22Z = (22B, 22B)
2024/01/03 00:52:36 XXX = (XXX, XXX)
2024/01/03 00:52:36 Directions: [L R]
2024/01/03 00:52:36 Node Map: map[11A:{11A 11B XXX} 11B:{11B XXX 11Z} 11Z:{11Z 11B XXX} 22A:{22A 22B XXX} 22B:{22B 22C 22C} 22C:{22C 22Z 22Z} 22Z:{22Z 22B 22B} XXX:{XXX XXX XXX}]
2024/01/03 00:52:36 Starting Nodes: [{11A 11B XXX} {22A 22B XXX}]
2024/01/03 00:52:36 11A reached 11Z at index 2 with cycle 2
2024/01/03 00:52:36 22A reached 22Z at index 1 with cycle 3
2024/01/03 00:52:36 Cycles: map[11A:2 22A:3]
2024/01/03 00:52:36 LCM of map[11A:2 22A:3] is 6
```