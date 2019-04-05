package dto

import (
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/server/http/api/models"
)

func SearchResultFrom(from dal.QueryResult) models.SearchResult {
	count := float64(from.Count)
	return models.SearchResult{
		Paging: &models.SearchResultPaging{
			Cursors: &models.SearchResultPagingCursors{
				After:  from.AfterCursor.String(),
				Before: from.BeforeCursor.String(),
			},
			Count: &count,
		},
	}
}

func PaginationTo(countP *int32, cursorP *string) dal.Pagination {
	var count uint64
	var cursor dal.Cursor

	if countP != nil {
		count = uint64(*countP)
	}

	if cursorP != nil {
		cursor = dal.Cursor(*cursorP)
	}

	return dal.Pagination{
		Cursor: cursor,
		Count:  count,
	}
}
