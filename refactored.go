// This is a hardly refactred and more general version of the naive
// implementation. It is designed for reuse and can be applied as toolchain.

package main

// This example only uses features from Go's stdlib.
import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// Declaring an interface helps us to apply the Normalize function to all
// following normalizers. It's pretty handy for the implemenation of a
// pipeline.
type Normalizer interface {
	Normalize(text string) string
}

// A chain normalizer collects n normalizers which are extecutes one after
// another. Which helps us to reuse it as some kind of pipeline.
type ChainNormalizer struct {
	Normalizers []Normalizer
}

func NewChainNormalizer(normalizers ...Normalizer) *ChainNormalizer {
	return &ChainNormalizer{
		Normalizers: normalizers}
}

func (chainNormalizer *ChainNormalizer) Normalize(text string) string {
	normalizedString := text
	// Does this contruct look weird? It's nothing else like a for each loop
	// with two assignments. The first one the index is ignored in this case.
	// The second one is the normalizer
	for _, normalizer := range chainNormalizer.Normalizers {
		normalizedString = normalizer.Normalize(normalizedString)
	}

	return normalizedString
}

////////////////////////////////////////////////////////////////////////////////

// A RegexpNormalizer is a is a more general construct which helps us to
// declare more concrete normalizers later like the UrlReplacementNormalizer
// which uses nothing else like a special regexp.
type RegexpNormalizer struct {
	Regexp            *regexp.Regexp
	ReplacementString string
}

func NewRegexpNormalizer(regularExpression string, replacementString string) *RegexpNormalizer {
	return &RegexpNormalizer{
		Regexp:            regexp.MustCompile(regularExpression),
		ReplacementString: replacementString}
}

func (normalizer *RegexpNormalizer) Normalize(text string) string {
	return normalizer.Regexp.ReplaceAllString(text, normalizer.ReplacementString)
}

// Extracting URLs is hard, take a look at the following links.
// https://mathiasbynens.be/demo/url-regex
// http://stackoverflow.com/questions/161738/what-is-the-best-regular-expression-to-check-if-a-string-is-a-valid-url
func NewUrlReplacementNormalizer(replacementString string) *RegexpNormalizer {
	return &RegexpNormalizer{
		Regexp: regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)
																					 (?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+|
																					 (?:www.|[-;:&=\+\$,\w]+@)
																					 [A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?
																					 \??(?:[-\+=&;%@.\w_]*)#
																					 ?(?:[\w]*))?)`),
		ReplacementString: replacementString}
}

////////////////////////////////////////////////////////////////////////////////
// Replacing nothing but strings. This helps us to eleminate certaing characters
// like # und @.
type StringReplacementNormalizer struct {
	NormalizationStrings []string
	ReplacementString    string
	replacer             *strings.Replacer
}

func NewStringReplacementNormalizer(normalizationStrings []string,
	replacementString string) *StringReplacementNormalizer {
	replacerStrings := make([]string, 2*len(normalizationStrings))
	for _, character := range normalizationStrings {
		replacerStrings = append(replacerStrings, character)
		replacerStrings = append(replacerStrings, replacementString)
	}

	replacer := strings.NewReplacer(replacerStrings...)
	return &StringReplacementNormalizer{
		NormalizationStrings: normalizationStrings,
		ReplacementString:    replacementString,
		replacer:             replacer}
}

func (normalizer *StringReplacementNormalizer) Normalize(text string) string {
	return normalizer.replacer.Replace(text)
}

////////////////////////////////////////////////////////////////////////

// The purpose of an UnicodeRangeNormalizer is to eliminate certain Uniocde
// ranges like emoticon or transportation and map.
// A range is defined by a staring and end point.
type UnicodeRangeNormalizer struct {
	StartPoint        int
	EndPoint          int
	ReplacementString string
}

type chart struct {
	StartPoint int
	EndPoint   int
}

// see unicode documentation for more details
var EmoticonChart *chart = &chart{0x1F600, 0x1F64F}
var TransportAndMapChart *chart = &chart{0x1F680, 0x1F6FF}

func NewUnicodeRangeNormalizer(startPoint int, endPoint int,
	replacementString string) *UnicodeRangeNormalizer {
	return &UnicodeRangeNormalizer{
		StartPoint:        startPoint,
		EndPoint:          endPoint,
		ReplacementString: replacementString}
}

func NewUnicodeRangeNormalizerFromChart(chart *chart,
	replacementString string) *UnicodeRangeNormalizer {
	return &UnicodeRangeNormalizer{
		StartPoint:        chart.StartPoint,
		EndPoint:          chart.EndPoint,
		ReplacementString: replacementString}
}

func (normalizer *UnicodeRangeNormalizer) Normalize(text string) string {
	normalizedString := bytes.NewBufferString("")
	for _, runeValue := range text {
		if runeValue >= rune(normalizer.StartPoint) &&
			runeValue <= rune(normalizer.EndPoint) {
			normalizedString.WriteString(normalizer.ReplacementString)
		} else {
			normalizedString.WriteRune(runeValue)
		}
	}
	return normalizedString.String()
}

////////////////////////////////////////////////////////////////////////
// This one is pretty trivial
type LowerCaseNormalizer struct {
	ReplacementString string
}

func NewLowerCaseNormalizer() *LowerCaseNormalizer {
	return &LowerCaseNormalizer{}
}

func (normalizer *LowerCaseNormalizer) Normalize(text string) string {
	return strings.ToLower(text)
}

////////////////////////////////////////////////////////////////////////////////

const Unigram int = 1
const Bigram int = 2
const Trigram int = 3

const WordLevel = "word"
const CharacterLevel = "character"

// The tokenizer is able to extract n grams on word or character level
// the step width ca be defined as well.
type NGramTokenizer struct {
	steps int
	level string
}

func NewNGramTokenizer(steps int, level string) *NGramTokenizer {
	return &NGramTokenizer{
		steps: steps,
		level: level}
}

func (tokenizer *NGramTokenizer) TokenizeString(text string) []string {
	switch tokenizer.level {
	case WordLevel:
		return tokenizer.tokenizeWords(text)
	case CharacterLevel:
		return tokenizer.tokenizeCharacters(text)
	}
	return nil
}

func (tokenizer *NGramTokenizer) tokenizeCharacters(text string) []string {
	var ngrams []string
	for i := 0; i < (len(text) - tokenizer.steps + 1); i++ {
		ngrams = append(ngrams, text[i:i+tokenizer.steps])
	}
	return ngrams
}

func (tokenizer *NGramTokenizer) tokenizeWords(text string) []string {
	var ngrams []string
	splittedWords := strings.Split(text, " ")
	for i := 0; i < (len(splittedWords) - tokenizer.steps + 1); i++ {
		ngrams = append(ngrams,
			strings.Join(splittedWords[i:i+tokenizer.steps], " "))
	}
	return ngrams
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// awesomeTweet := "Die @bedcon ist die großartigste #Konferenz des Jahres. 😍 http://bedcon.org"
	awesomeTweet := "Die @bedcon ist die großartigste #Konferenz des Jahres. http://bedcon.org"

	chain := NewChainNormalizer(
		NewLowerCaseNormalizer(),
		NewUnicodeRangeNormalizerFromChart(EmoticonChart, ""),
		NewUrlReplacementNormalizer(""),
		NewStringReplacementNormalizer([]string{"?", ".", ",", "@", "-", "/", ":",
			"#", "!", ")", "(", "[", "]", "¿"}, ""))

	denoizedTweet := chain.Normalize(awesomeTweet)

	bigramTokenizer := NewNGramTokenizer(Bigram, CharacterLevel)
	tokenizedTweet := bigramTokenizer.TokenizeString(denoizedTweet)

	fmt.Printf("Denoized and tokenized String: %v \n", tokenizedTweet)
}
