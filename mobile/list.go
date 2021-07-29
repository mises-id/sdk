package mobile

type mStringListWrapper struct {
	slice []string
}
type mUserListWrapper struct {
	slice []MUser
}

func (w *mStringListWrapper) Count() int {
	return len(w.slice)
}
func (w *mStringListWrapper) Get(idx int) string {
	return w.slice[idx]
}

func (w *mUserListWrapper) Count() int {
	return len(w.slice)
}
func (w *mUserListWrapper) Get(idx int) MUser {
	return w.slice[idx]
}

func MStringListToSlice(l MStringList) []string {
	ret := []string{}
	for i := 0; i < l.Count(); i++ {
		ret = append(ret, l.Get(i))
	}
	return ret
}
func MUserListToSlice(l MUserList) []MUser {
	ret := []MUser{}
	for i := 0; i < l.Count(); i++ {
		ret = append(ret, l.Get(i))
	}
	return ret
}
