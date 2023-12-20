package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

var stringDigits = map[string]int{
    "one": 1,
    "two": 2,
    "three": 3,
    "four": 4,
    "five": 5,
    "six": 6,
    "seven": 7,
    "eight": 8,
    "nine": 9,
}

func main() {
    var filepath string

    args := os.Args
    if len(args) > 1 {
        filepath = args[1]
    } else {
        log.Fatal("Need to pass file path as argument")
    }

    file, err := os.Open(filepath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    sum := 0
    for scanner.Scan() {
        text := scanner.Text()
        num := getFirstAndLastDigits(text)
        fmt.Printf("%s parsed digits are: %d\n", text, num)
        sum += num
    }

    fmt.Printf("Sum of digits is: %d\n", sum)

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func getFirstAndLastDigits(s string) int {
    var digits []int
    var buffer = ""
    for i, r := range s {
        if unicode.IsDigit(r) {
            num, _ := strconv.Atoi(string(s[i]))
            digits = append(digits, num)
            buffer = ""
        } else {
            buffer = buffer + string(s[i])
            check := checkForSpeltDigit(buffer)
            if check != 0 {
                digits = append(digits, check)
                // keep the last buffer character in case of overlap? (e.g. twone can be two then one)
                buffer = buffer[(len(buffer)-1):]
            }
        }
    }

    if len(digits) == 0 {
        return 0
    } else {
        r, _ := strconv.Atoi(fmt.Sprintf("%d%d", digits[0], digits[len(digits) - 1]))
        return r
    }
}

func checkForSpeltDigit(s string) int {
    if len(s) < 3 {
        return 0
    }

    val, ok := stringDigits[s]
    if ok {
        return val
    } else {
        return checkForSpeltDigit(s[1:])
    }
}