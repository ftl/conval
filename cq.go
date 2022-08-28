/*
This file contains the implementation of the specific things for contests announced by the CQ magazine.
*/
package conval

func init() {
	PropertyGetters[WPXPrefixProperty] = PropertyGetterFunc(getWPXPrefix)
}

const (
	WPXPrefixProperty Property = "wpx_prefix"
)

func getWPXPrefix(qso QSO) string {
	return "" // TODO implement
}
