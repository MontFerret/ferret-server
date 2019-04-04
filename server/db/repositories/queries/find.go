package queries

var (
	FindAll = `
		FOR i IN %s
			FILTER @` + ParamPageCursor + ` != NULL ? i.created_at < @` + ParamPageCursor + ` : TRUE == TRUE
			LIMIT @` + ParamPageCount + `
			RETURN i
`
	FindAllByScriptID = `
		FOR i IN %s
			FILTER i.script_id == @` + ParamFilterByScriptId + `
			FILTER @` + ParamPageCursor + ` != NULL ? i.created_at < @` + ParamPageCursor + ` : TRUE == TRUE
			LIMIT @` + ParamPageCount + `
			RETURN i
`

	FindOneByName = `
		FOR i IN %s
			FILTER i.name == @` + ParamFilterByName + `
			LIMIT 1
			RETURN i
	`
)
