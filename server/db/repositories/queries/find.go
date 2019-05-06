package queries

import (
	"fmt"
)

var (
	findAll = `
		FOR i IN %s
			SORT i.created_at
			FILTER @` + ParamPageCursor + ` != NULL ? DATE_TIMESTAMP(i.created_at) > @` + ParamPageCursor + ` : TRUE == TRUE
			LIMIT @` + ParamPageCount + `
			RETURN i
`
	findAllByScriptID = `
		FOR i IN %s
			SORT i.created_at
			FILTER i.script_id == @` + ParamFilterByScriptID + `
			FILTER @` + ParamPageCursor + ` != NULL ? DATE_TIMESTAMP(i.created_at) > @` + ParamPageCursor + ` : TRUE == TRUE
			LIMIT @` + ParamPageCount + `
			RETURN i
`

	findOneByName = `
		FOR i IN %s
			FILTER i.name == @` + ParamFilterByName + `
			LIMIT 1
			RETURN i
	`
)

func FindAll(collection string) string {
	return fmt.Sprintf(findAll, collection)
}

func FindAllByScriptID(collection string) string {
	return fmt.Sprintf(findAllByScriptID, collection)
}

func FindOneByName(collection string) string {
	return fmt.Sprintf(findOneByName, collection)
}
