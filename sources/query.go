package sources

import "sort"

type queryResponseItem struct {
	rank int
	item Item
}

type queryResponse []queryResponseItem

func (qr queryResponse) Len() int {
	return len(qr)
}

func (qr queryResponse) Less(i, j int) bool {
	iri := qr[i]
	jri := qr[j]

	if iri.rank > jri.rank {
		return true
	} else if iri.rank == jri.rank && iri.item.LessThan(jri.item) {
		return true
	}

	return false
}

func (qr queryResponse) Swap(i, j int) {
	qr[i], qr[j] = qr[j], qr[i]
}

func Query(items []Item, query string) []Item {
	qr := make(queryResponse, 0)
	for _, item := range items {
		if rank := item.GetRank(query); rank > 0 {
			qr = append(qr, queryResponseItem{
				rank: rank,
				item: item,
			})
		}
	}
	sort.Sort(qr)

	result := make([]Item, 0)
	for _, qri := range qr {
		result = append(result, qri.item)
		if len(result) >= 9 {
			break
		}
	}

	return result
}
