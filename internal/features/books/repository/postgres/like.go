package books_repository_postgres

import "strings"

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func likePattern(s *string) *string {
	if s == nil {
		return nil
	}
	pattern := "%" + likeEscaper.Replace(*s) + "%"
	return &pattern
}
