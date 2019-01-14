package queries

const (
	FindAll = `
		FOR i IN %s
			LIMIT @offset, @count
			RETURN i
`
	FindAllByScriptID = `
		FOR i IN %s
			FILTER i.script_id == @script_id
			LIMIT @offset, @count
			RETURN i
`

	FindOneByName = `
		FOR i IN %s
			FILTER i.name == @name
			LIMIT 1
			RETURN i
	`
)
