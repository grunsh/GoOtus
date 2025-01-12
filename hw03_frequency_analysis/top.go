package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var topWords = 10

type stringsCounter = map[string]int

var (
	WhiteSpaceRegex = regexp.MustCompile(`\s+`)
	WordRegex       = regexp.MustCompile(`\p{L}+\p{P}*\p{L}+|\p{L}+|-{2,}`)
)

func MySplit(s string) []string {
	return WhiteSpaceRegex.Split(s, -1)
}

func initOrCount(sc *stringsCounter, s string) {
	if value, exists := (*sc)[s]; exists {
		(*sc)[s] = value + 1
	} else {
		(*sc)[s] = 1
	}
}

func GetWord(s string) (string, bool) {
	st := WordRegex.FindString(s)
	if st == "" {
		return "", false
	}
	return st, true
}

func calcRates(s string) stringsCounter {
	rateS := make(stringsCounter)
	splitedS := MySplit(s)
	for _, v := range splitedS {
		if s, ok := GetWord(v); ok {
			initOrCount(&rateS, strings.ToLower(s))
		}
	}
	return rateS
}

type counter [][]int

// Реализуем методы интерфейса sort.Interface.
func (c counter) Len() int           { return len(c) }
func (c counter) Less(i, j int) bool { return c[i][0] > c[j][0] }
func (c counter) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func Top10(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	rateS := calcRates(s)
	ind := make([][]int, len(rateS)) // [0] - Количество повторений слова, [1] - индекс слова в words
	words := make([]string, len(rateS))
	i := 0
	for k, v := range rateS {
		words[i] = k
		ind[i] = append(ind[i], v, i)
		i++
	}
	sort.Sort(counter(ind))
	i = 1
	retSlice := []string{}                  // для возврата результата
	tempSlice := []string{words[ind[0][1]]} // для сортировки пачки элементов с одинаковым количеством вхождений
	for {
		if i >= len(ind) {
			if len(tempSlice) > 1 { // Выполним сортировку, но только если у нас более 1 элемента
				sort.Strings(tempSlice)
			}
			retSlice = append(retSlice, tempSlice...)
			break
		}
		if ind[i][0] == ind[i-1][0] {
			tempSlice = append(tempSlice, words[ind[i][1]])
			i++
			continue
		}
		if len(tempSlice) > 1 {
			sort.Strings(tempSlice)
		}
		retSlice = append(retSlice, tempSlice...)
		tempSlice = tempSlice[:0]
		tempSlice = append(tempSlice, words[ind[i][1]])
		i++
		if i >= topWords {
			break
		}
	}
	if len(retSlice) < topWords {
		topWords = len(retSlice)
	}
	return retSlice[0:topWords]
}
