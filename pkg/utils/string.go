package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"

	"github.com/alextanhongpin/stringdist"
	"github.com/shopspring/decimal"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func NormalizeString(text string) (result string) {
	// To lower
	result = strings.ToLower(text)

	// Remove dup space
	result = strings.Join(strings.Fields(strings.TrimSpace(result)), " ")
	return
}

func NormalizeVietnamese(text string) (result string) {
	// Document at: https://unicode.org/reports/tr15/
	// Remove all diacritics
	trans := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ = transform.String(trans, text)

	// Replace special characters not in unicode.Mn
	result = strings.ReplaceAll(result, "đ", "d")
	result = strings.ReplaceAll(result, "Đ", "D")

	return NormalizeString(result)
}

func NormalizeVietnameseName(text string) (result string) {
	text = NormalizeVietnamese(text)
	text = ReplaceAcronyms(text)

	prefixes := []string{
		"ong ",
		"ba ",
		"anh ",
		"chi ",
		"ong/ba ",
		"anh/chi ",
		"cong ty ",
	}
	for _, prefix := range prefixes {
		text = strings.TrimPrefix(text, prefix)
	}

	return text
}

func ReplaceAcronyms(text string) string {
	acronyms := map[string]string{
		"&":               " va ",
		"cty":             "cong ty",
		"xd":              "xay dung",
		"pt":              "phat trien",
		"gd":              "giam doc",
		"tmcp":            "thuong mai co phan",
		"tnhh":            "trach nhiem huu han",
		"ctcp":            "cong ty co phan",
		"cong ty cp":      "cong ty co phan",
		"cthd":            "cong ty hop danh",
		"cong ty hd":      "cong ty hop danh",
		"dnnn":            "doanh nghiep nha nuoc",
		"doanh nghiep nn": "doanh nghiep nha nuoc",
		"dntn":            "doanh nghiep tu nhan",
		"doanh nghiep tn": "doanh nghiep tu nhan",
		"dhdcd":           "dai hoi dong co dong",
		"gcndkdn":         "giay chung nhan dang ky doanh nghiep",
		"hdqt":            "hoi dong quan tri",
		"hdtv":            "hoi dong thanh vien",
	}

	for acronym, fullForm := range acronyms {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(acronym) + `\b`)
		text = re.ReplaceAllString(text, fullForm)
	}
	return text
}

func IsWordInText(text string, words ...string) bool {
	text = NormalizeVietnamese(text)

	// Escape special characters in each word and join them with the OR operator
	var wordPatterns []string
	for _, word := range words {
		word = NormalizeVietnamese(word)
		word = regexp.QuoteMeta(word)
		wordPatterns = append(wordPatterns, word)
	}
	pattern := `\b(` + strings.Join(wordPatterns, `|`) + `)\b`
	regex := regexp.MustCompile(pattern)
	// Check if any of the words match the regex pattern in the text
	return regex.MatchString(text)
}

func NormalizeBalance(text string) (balance decimal.Decimal, err error) {
	text = RemoveCurrencyIndicators(text)

	// All non-numeric characters
	reg := regexp.MustCompile("[^0-9.,+-]")
	if reg.FindStringIndex(text) != nil {
		err = comerr.WrapMessage(comerr.ErrInvalidData, "invalid char set balance")
		return
	}

	decimalPlace, intPart, _, err := FindBalanceDecimalPlace(text)
	if err != nil {
		return
	}

	if intPart != "" {
		pattern := `^[+-]?\d{1,3}([.,]\d{3})*$`
		re := regexp.MustCompile(pattern)
		if !re.MatchString(intPart) {
			return balance, comerr.WrapMessage(comerr.ErrInvalidData, "invalid thousand seperator")
		}
	}

	// Replace "," "."
	text = strings.ReplaceAll(text, ",", "")
	text = strings.ReplaceAll(text, ".", "")

	balance, err = decimal.NewFromString(text)
	if err != nil {
		err = comerr.ErrInvalidData
		return
	}

	balance = balance.Shift(int32(-decimalPlace))
	return
}

func RemoveCurrencyIndicators(text string) string {
	text = NormalizeVietnamese(text)

	indicators := []string{
		"d",
		"vnd",
		"usd",
		"dollars",
		"dollar",
		"$",
		"eur",
		"€",
	}

	indicatorsRegexp := "(" + strings.Join(indicators, "|") + ")" // Regular expression to match any of the indicators
	// Compile the regular expression
	re := regexp.MustCompile(indicatorsRegexp)

	// Replace indicators with an empty string
	text = re.ReplaceAllString(text, "")

	// Return the modified text
	return NormalizeVietnamese(text)
}

func FindBalanceDecimalPlace(text string) (decimalPlace int, intPart string, fractionalPart string, err error) {
	// Return decimal place, part seperator, decimal place indicator, error
	// If only have "." or only have ",", not have decimal place
	// 1,000,000.00 -> 2, "," "." nil
	indexDot := strings.LastIndex(text, ".")
	indexComma := strings.LastIndex(text, ",")

	if indexDot == -1 && indexComma == -1 {
		decimalPlace = 0
		intPart = text
		return
	}

	fractionalIndex := -1
	if indexDot > indexComma {
		fractionalIndex = indexDot
	} else {
		fractionalIndex = indexComma
	}

	lastParts := text[fractionalIndex+1:]
	if len(lastParts) != 3 {
		decimalPlace = len(lastParts)
		intPart = text[:fractionalIndex]
		fractionalPart = lastParts
	} else {
		intPart = text
	}
	return
}

func LevenshteinPercentage(s1, s2 string) float64 {
	// Calculate the maximum possible distance
	maxLen := max(len(s1), len(s2))
	distance := stringdist.NewDamerauLevenshtein(maxLen+1).Calculate(s1, s2)

	return (float64(distance) / float64(maxLen)) * 100
}
