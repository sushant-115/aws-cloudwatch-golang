package set

//Report coming from results
type Report interface{}

//ReportSet a custom hashset
type ReportSet struct {
	set map[Report]bool
}

//Add method for adding elements in set if not present
func (set *ReportSet) add(r Report) bool {
	_, found := set.set[r]
	set.set[r] = true
	return !found //False if it existed already
}

func createArray(m map[Report]bool) []Report {
	arr := []Report{}
	for k := range m {
		arr = append(arr, k)
	}
	return arr
}

//MakeSet from reports array
func MakeSet(reports []Report) []Report {

	set := make(map[Report]bool)
	x := ReportSet{set}
	for i := 0; i < len(reports); i++ {
		x.add(reports[i])
	}
	return createArray(x.set)
}
