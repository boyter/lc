// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"sort"
	"strings"
)

// KeyWordGuessLicence will some content try to guess what the licence is based on checking for unique keywords
// using the prebuilt licence library which contains what we hope are unique ngrams for each licence
func (l *LicenceGuesser) KeyWordGuessLicence(content []byte) []License {
	haystack := LcCleanText(string(content))

	var keywordMatches []License
	var jaccardCandidates []License
	var maxKeywordScore float64

	// Swap out the database to the full one if configured to use it
	db := l.CommonDatabase
	if l.UseFullDatabase {
		db = l.Database
	}

	// Build word set once for Jaccard comparisons
	contentWords := wordSet(haystack)

	for _, lic := range db {
		if len(lic.LicenseTexts) == 0 {
			continue
		}

		matched := false

		// Try keyword matching first
		if len(lic.Keywords) > 0 {
			match := lic.Trie.Match([]byte(haystack))
			if len(match) != 0 {
				lic.ScorePercentage = (float64(len(match)) / float64(len(lic.Keywords))) * 100
				lic.MatchType = MatchTypeKeyword
				keywordMatches = append(keywordMatches, lic)
				matched = true

				if lic.ScorePercentage > maxKeywordScore {
					maxKeywordScore = lic.ScorePercentage
				}
			}
		}

		// Also score via Jaccard for licenses that didn't match keywords
		if !matched {
			licWords := lic.WordSet
			if licWords == nil {
				licWords = wordSet(LcCleanText(lic.LicenseTexts[0]))
			}
			jScore := jaccardSimilarity(contentWords, licWords) * 100
			if jScore > 20 {
				lic.ScorePercentage = jScore
				lic.MatchType = MatchTypeKeyword
				jaccardCandidates = append(jaccardCandidates, lic)
			}
		}
	}

	// Filter keyword matches to keep only those close to the max score
	var filtered []License
	for _, lic := range keywordMatches {
		if lic.ScorePercentage >= (maxKeywordScore * 0.5) {
			filtered = append(filtered, lic)
		}
	}
	if len(filtered) != 0 {
		keywordMatches = filtered
	}

	// Sort Jaccard candidates and keep top 3
	sort.Slice(jaccardCandidates, func(i, j int) bool {
		return jaccardCandidates[i].ScorePercentage > jaccardCandidates[j].ScorePercentage
	})
	if len(jaccardCandidates) > 3 {
		jaccardCandidates = jaccardCandidates[:3]
	}

	// Combine both pools
	candidates := append(keywordMatches, jaccardCandidates...)

	// Use Levenshtein on top candidates to get final ranking
	const maxLevRunes = 1500
	if len(candidates) >= 2 {
		// Sort by current score to pick top N for expensive Levenshtein
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].ScorePercentage > candidates[j].ScorePercentage
		})

		topN := 5
		if len(candidates) < topN {
			topN = len(candidates)
		}

		contentRunes := []rune(string(content))
		if len(contentRunes) > maxLevRunes {
			contentRunes = contentRunes[:maxLevRunes]
		}

		for i := 0; i < topN; i++ {
			licText := ""
			if len(candidates[i].LicenseTexts) > 0 {
				licText = candidates[i].LicenseTexts[0]
			}

			licRunes := []rune(licText)
			if len(licRunes) > maxLevRunes {
				licRunes = licRunes[:maxLevRunes]
			}
			distance := levenshtein.DistanceForStrings(contentRunes, licRunes, levenshtein.DefaultOptions)

			maxLen := len(contentRunes)
			if len(licRunes) > maxLen {
				maxLen = len(licRunes)
			}
			if maxLen == 0 {
				candidates[i].ScorePercentage = 100
			} else {
				candidates[i].ScorePercentage = (1.0 - float64(distance)/float64(maxLen)) * 100
			}
		}

		// Re-sort the top candidates by Levenshtein score
		sort.Slice(candidates[:topN], func(i, j int) bool {
			return candidates[i].ScorePercentage > candidates[j].ScorePercentage
		})
	}

	return candidates
}

// wordSet splits text into a set of unique words
func wordSet(text string) map[string]struct{} {
	words := strings.Fields(text)
	set := make(map[string]struct{}, len(words))
	for _, w := range words {
		set[w] = struct{}{}
	}
	return set
}

// jaccardSimilarity computes the Jaccard similarity between two word sets
func jaccardSimilarity(a, b map[string]struct{}) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 0
	}
	intersection := 0
	for w := range a {
		if _, ok := b[w]; ok {
			intersection++
		}
	}
	union := len(a) + len(b) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}
