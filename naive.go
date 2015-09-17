package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func deleteEmoticons(text string) string {
	normalizedString := bytes.NewBufferString("")
	for _, runeValue := range text {
		if runeValue >= rune(0x1F600) && runeValue <= rune(0x1F64F) {
			normalizedString.WriteString("")
		} else {
			normalizedString.WriteRune(runeValue)
		}
	}
	return normalizedString.String()
}

func extractNGrams(text string, steps int) []string {
	var ngrams []string
	for i := 0; i < (len(text) - steps + 1); i++ {
		ngrams = append(ngrams, text[i:i+steps])
	}
	return ngrams
}

func main() {
	awesomeTweet := "Die @bedcon ist die großartigste #Konferenz dieses Jahr. 😍 http://bedcon.org"
	lowerCasedTweet := strings.ToLower(awesomeTweet)
	fmt.Printf("Lower Case: %v \n", lowerCasedTweet)

	withoutHashtag := strings.Replace(lowerCasedTweet, "#", "", -1)
	fmt.Printf("Without Hashtag: %v \n", withoutHashtag)

	withoutMention := strings.Replace(withoutHashtag, "@", "", -1)
	fmt.Printf("Without Mention: %v \n", withoutMention)

	urlRegexp := regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)`)
	withoutURL := urlRegexp.ReplaceAllString(withoutMention, "")
	fmt.Printf("Without URL: %v \n", withoutURL)

	withoutEmoticons := deleteEmoticons(withoutURL)
	fmt.Printf("Without Emoticons: %v \n", withoutEmoticons)

	nGrams := extractNGrams(withoutEmoticons, 2)
	fmt.Printf("NGrams: %v \n", nGrams)
}
