package processor

// GuessLicense tries to guess the licence for this content based on whatever heuristics are set on the
// guesser itself, note that this is a multi stage approach where it mixes each one together
func (l *LicenceGuesser) GuessLicense(content []byte) []License {

	// if we find SPDX markers then just return those because its a very high signal of what was there
	// TODO add ability to disable this check
	spdxIdentified := l.SpdxIdentify(string(content))
	if len(spdxIdentified) != 0 {
		return spdxIdentified
	}

	// now it gets tricky, the vector space will always return results
	// but unless its very confident we should always cross check against keywords
	// which allows us to get an averaged rank

	// try keywords first, and if we get anything of high confidence we should assume it was correct
	keyWordGuessLicence := l.KeyWordGuessLicence(content)
	fkeyWordGuessLicence := l.filterByPercentage(keyWordGuessLicence)

	if len(fkeyWordGuessLicence) != 0 {
		// at this point lets get the vector space guesses and see if there is any overlap between the two
		// with a high level of confidence
		vectorSpaceGuessLicence := l.VectorSpaceGuessLicence(content)
		fvectorSpaceGuessLicence := l.filterByPercentage(vectorSpaceGuessLicence)

		// right now lets check both lists to see if anything is common
		var common []License

		for _, x := range fkeyWordGuessLicence {
			for _, y := range fvectorSpaceGuessLicence {
				if x.LicenseId == y.LicenseId {
					x.ScorePercentage = (x.ScorePercentage + y.ScorePercentage) / 2
					x.MatchType = MatchTypeBlended
					common = append(common, x)
				}
			}
		}

		fcommon := l.filterByPercentage(common)
		return fcommon
	}

	return nil
}

func (l *LicenceGuesser) filterByPercentage(keyWordGuessLicence []License) []License {
	// filter out anything below our cutoff %
	var filteredKeywordGuessLicence []License
	for _, x := range keyWordGuessLicence {
		if x.ScorePercentage >= l.cutoffPercentage {
			filteredKeywordGuessLicence = append(filteredKeywordGuessLicence, x)
		}
	}

	return filteredKeywordGuessLicence
}
