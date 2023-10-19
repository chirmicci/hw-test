package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	l := strings.Fields(text)
	if len(l) == 0 {
		return []string{}
	}
	top := make(map[string]int)
	uniqueList := []string{}
	for _, s := range l {
		if _, ok := top[s]; ok {
			top[s]++
		} else {
			top[s] = 1
			uniqueList = append(uniqueList, s)
		}
	}
	sort.Slice(uniqueList, func(i, j int) bool {
		if top[uniqueList[i]] == top[uniqueList[j]] {
			return uniqueList[i] < uniqueList[j]
		}
		return top[uniqueList[i]] > top[uniqueList[j]]
	})

	topList := make([]string, 0, 10)
	for _, s := range uniqueList {
		if _, ok := top[s]; ok {
			topList = append(topList, s)
			delete(top, s)
			if len(topList) == 10 {
				break
			}
		}
	}

	return topList
}
