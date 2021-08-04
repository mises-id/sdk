package mobile

import "strings"

var _ MStringList = &mStringListWrapper{}
var _ MUserList = &mUserListWrapper{}

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

func mStringListToSlice(l MStringList) []string {
	ret := []string{}
	for i := 0; i < l.Count(); i++ {
		ret = append(ret, l.Get(i))
	}
	return ret
}
func mUserListToSlice(l MUserList) []MUser {
	ret := []MUser{}
	for i := 0; i < l.Count(); i++ {
		ret = append(ret, l.Get(i))
	}
	return ret
}

func NewMStringList(s string, sep string) MStringList {
	return &mStringListWrapper{strings.Split(s, sep)}
}
