package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

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
        num, _ := getFirstAndLastDigits(text)
        fmt.Printf("%s parsed digits are: %d\n", text, num)
        sum += num
    }

    fmt.Printf("Sum of digits is: %d\n", sum)

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func getFirstAndLastDigits(s string) (int, error) {
    var digits []byte
    for i, r := range s {
        if unicode.IsDigit(r) {
            digits = append(digits, s[i])
        }
    }

    if len(digits) == 0 {
        return 0, nil
    } else if len(digits) == 1 {
        digits = append(digits, digits[0])
        return strconv.Atoi(string(digits))
    } else {
        return strconv.Atoi(string([]byte{digits[0], digits[len(digits) - 1]}))
    }
}
