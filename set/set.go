package set

import "../structs"

//ReportSet a custom hashset
type ReportSet struct {
	set map[structs.Report]bool
}

//Add method for adding elements in set if not present
func (set *ReportSet) add(r structs.Report) bool {
	_, found := set.set[r]
	set.set[r] = true
	return !found //False if it existed already
}

func createArray(m map[structs.Report]bool) []structs.Report {
	arr := []structs.Report{}
	for k := range m {
		arr = append(arr, k)
	}
	return arr
}

//MakeSet from reports array
func MakeSet(reports []structs.Report) []structs.Report {

	set := make(map[structs.Report]bool)
	x := ReportSet{set}
	for i := 0; i < len(reports); i++ {
		x.add(reports[i])
	}
	return createArray(x.set)
}
