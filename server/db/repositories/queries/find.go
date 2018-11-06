package queries

const (
	FindAll = `
		FOR d IN %s
			LIMIT @offset, @count
			RETURN d
`
	FindOneByName = `
		FOR i IN %s
			FILTER i.name == @name
			LIMIT 1
			RETURN i
	`
)
