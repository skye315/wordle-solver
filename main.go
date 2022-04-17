package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("hello skye")
	wordList := getWordList()

	for i := 0; i < 6; i++ {
		// make a guess
		guess := makeAGuess(wordList)

		// get back mask (green yellow grey stuff)
		maskForGuess := getMaskForWord(guess)
		fmt.Println(maskForGuess)

		// if mask is 5 green we won!!!!
		if maskForGuess.didIWin() {
			fmt.Println("Yay you won!!!!!")
			os.Exit(0)
		}

		wordList = filterBadWords(wordList, maskForGuess)
	}
	fmt.Println("done")
}

func getWordList() []word {
	readFile, err := os.Open("word-list.txt")
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var wordList []word
	for fileScanner.Scan() {
		w := word(fileScanner.Text())
		if w.valid() {
			wordList = append(wordList, w)
		}
	}
	readFile.Close()
	return wordList
}

func makeAGuess(wordList []word) word {
	if len(wordList) == 0 {
		fmt.Println("You ran out of words so I guess you lost")
		os.Exit(1)
	}
	return wordList[0]
}

func getMaskForWord(guess word) mask {
	fmt.Println("Try guessing", guess, "and tell me what mask you get back")
	for {
		fmt.Println("Type a mask of 5 letters using g for green, y for yellow and b for black")
		reader := bufio.NewReader(os.Stdin)
		// ReadString will block until the delimiter is entered
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input = strings.TrimSuffix(input, "\n")
		out := mask{
			guess:        guess,
			maskForGuess: input,
		}
		if out.valid() {
			return out
		}
	}
}

func filterBadWords(wordList []word, maskForGuess mask) []word {
	out := make([]word, 0, len(wordList))
	for _, oneWord := range wordList {
		if maskForGuess.matchesWord(oneWord) {
			out = append(out, oneWord)
		}
	}
	return out
}

type mask struct {
	guess        word
	maskForGuess string
}

func (m mask) valid() bool {
	match, _ := regexp.MatchString("^[bgy]{5}$", m.maskForGuess)
	return match
}

func (m mask) didIWin() bool {
	return m.maskForGuess == "ggggg"
}

func (m mask) matchesWord(w word) bool {
	yellows := make(map[byte]int)
	blacks := make(map[byte]int)
	greens := make(map[byte]int)
	for i := 0; i < len(m.maskForGuess); i++ {
		if m.maskForGuess[i] == 'g' {
			if w[i] != m.guess[i] {
				fmt.Println("We failed the green match at pos", i, w)
				return false
			}
			greens[m.guess[i]]++
		}
		if m.maskForGuess[i] == 'b' {
			if w[i] == m.guess[i] {
				fmt.Println("We failed the black match at pos", i, w)
				return false
			}
			blacks[m.guess[i]]++
		}
		if m.maskForGuess[i] == 'y' {
			if w[i] == m.guess[i] {
				fmt.Println("We failed the yellow match at pos", i, w)
				return false
			}
			yellows[m.guess[i]]++
		}
	}
	for letter, letterCnt := range yellows {
		cnt := 0
		for j := 0; j < len(w); j++ {
			if w[j] == letter {
				cnt++
			}
		}
		if cnt < letterCnt {
			fmt.Println("We failed the yellow count ", cnt, letter, letterCnt, w)
			return false
		}
	}
	for letter := range blacks {
		cnt := 0
		for j := 0; j < len(w); j++ {
			if w[j] == letter {
				cnt++
			}
		}
		if cnt > (yellows[letter] + greens[letter]) {
			fmt.Println("We failed the black count ", cnt, letter, yellows[letter], w)
			return false
		}
	}

	return true
}

type word string

func (w word) valid() bool {
	match, _ := regexp.MatchString("^[a-z]{5}$", string(w))
	return match
}
