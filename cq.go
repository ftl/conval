/*
This file contains the implementation of the specific things for contests announced by the CQ magazine.
*/
package conval

func init() {
	commonPropertyGetters[WPXPrefixProperty] = PropertyGetterFunc(getWPXPrefix)
}

const (
	WPXPrefixProperty Property = "wpx_prefix"
)

func getWPXPrefix(qso QSO) string {
	return WPXPrefix(qso.TheirCall)
}
