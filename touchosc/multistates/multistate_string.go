// generated by stringer --type MultiState; DO NOT EDIT

package multistates

import "fmt"

const _MultiState_name = "UnknownPressedReleased"

var _MultiState_index = [...]uint8{0, 7, 14, 22}

func (i MultiState) String() string {
	if i < 0 || i >= MultiState(len(_MultiState_index)-1) {
		return fmt.Sprintf("MultiState(%d)", i)
	}
	return _MultiState_name[_MultiState_index[i]:_MultiState_index[i+1]]
}
