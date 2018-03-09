package parsers

import (
	"encoding/base64"
	"encoding/json"
	vectorspace "github.com/boyter/golangvectorspace"
	"strings"
)

// Caching the database load result reduces processing time by about 3x for this repository
var Database = []License{}

func loadDatabase() []License {
	if len(Database) != 0 {
		return Database
	}

	var database []License
	data, _ := base64.StdEncoding.DecodeString(database_keywords)
	_ = json.Unmarshal(data, &database)

	for i, v := range database {
		database[i].Concordance = vectorspace.BuildConcordance(strings.ToLower(v.LicenseText))
	}

	Database = database

	return database
}
