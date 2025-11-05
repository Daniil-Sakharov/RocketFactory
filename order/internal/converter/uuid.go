package converter

import "github.com/google/uuid"

// stringsToUUIDs конвертирует []string в []uuid.UUID
func StringsToUUIDs(strs []string) []uuid.UUID {
	uuids := make([]uuid.UUID, 0, len(strs))
	for _, str := range strs {
		if u, err := uuid.Parse(str); err == nil {
			uuids = append(uuids, u)
		}
	}
	return uuids
}

// uuidsToStrings конвертирует []uuid.UUID в []string
func UuidsToStrings(uuids []uuid.UUID) []string {
	strs := make([]string, 0, len(uuids))
	for _, u := range uuids {
		strs = append(strs, u.String())
	}
	return strs
}
