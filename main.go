package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type TDictionary = map[string]struct{}

var empty = struct{}{}

const maskLetter string = "*"
const commandExit string = "exit"

var reader = bufio.NewReader(os.Stdin)

func getUserInput() (string, error) {
	inp, err := reader.ReadString('\n')
	if (err != nil) || (len(inp) == 0) {
		err = errors.New("error on user input")
		return "", err
	}
	inp = strings.TrimSuffix(inp, "\n")
	inp = strings.TrimSuffix(inp, "\r")
	return inp, nil
}

func readDictioary(fname string) (TDictionary, error) {
	var wordsList = make(TDictionary)
	file, err := os.Open(fname)
	if err != nil {
		err = errors.New("error while opening file")
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordsList[strings.ToLower(scanner.Text())] = empty
	}
	err = scanner.Err()
	if err != nil {
		err = errors.New("error reading file")
		return nil, err
	}
	return wordsList, err
}

func checkWordAllowability(word string, dict TDictionary) bool {
	_, ok := dict[strings.ToLower(word)]
	return !ok
}

func replaceWordWithAsterisk(word string) string {
	n := utf8.RuneCountInString(word)
	return strings.Repeat(maskLetter, n)
}

func censor(word string, dict TDictionary) string {
	if checkWordAllowability(word, dict) {
		return word
	} else {
		return replaceWordWithAsterisk(word)
	}
}

func main() {
	var fname string
	var err error
	var dictionary TDictionary
	var line string = ""
	var new_line = strings.Builder{}
	var new_sentence = strings.Builder{}

	fname, err = getUserInput()
	if err != nil {
		log.Fatal("error while reading file name")
	}

	dictionary, err = readDictioary(fname)
	if err != nil {
		log.Fatal("error while reading file")
	}

	for line != commandExit {
		line, err = getUserInput()
		if err != nil {
			log.Fatal("error while reading word from user")
		}
		new_line.Reset()
		for _, sentence := range strings.Split(line, ".") {
			new_sentence.Reset()
			for _, word := range strings.Split(sentence, " ") {
				new_sentence.WriteString(censor(word, dictionary))
				new_sentence.WriteString(" ")
			}
			new_line.WriteString(strings.TrimSuffix(new_sentence.String(), " "))
			new_line.WriteString(".")
		}
		fmt.Println(strings.TrimSuffix(new_line.String(), "."))
	}
	fmt.Println("Bye!")
}
