// Code generated by "stringer -type Format topic.go"; DO NOT EDIT.

package channel

import "fmt"

const _Format_name = "SlackMailError"

var _Format_index = [...]uint8{0, 5, 9, 14}

func (i Format) String() string {
	if i < 0 || i >= Format(len(_Format_index)-1) {
		return fmt.Sprintf("Format(%d)", i)
	}
	return _Format_name[_Format_index[i]:_Format_index[i+1]]
}
