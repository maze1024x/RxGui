package util

import "errors"

type ErrorWithExtraData interface {
	error
	Serialize(desc string) ([] byte)
	LookupExtraData(key string) (string, bool)
}

type WrappedErrorWithExtraData struct {
	Underlying   ErrorWithExtraData
	Description  string
}
func (e *WrappedErrorWithExtraData) Error() string {
	return e.Underlying.Error()
}
func (e *WrappedErrorWithExtraData) Serialize(desc string) ([] byte) {
	if desc == "" {
		desc = e.Description
	}
	return e.Underlying.Serialize(desc)
}
func (e *WrappedErrorWithExtraData) LookupExtraData(key string) (string, bool) {
	return e.Underlying.LookupExtraData(key)
}

func WrapError(e error, wrap func(string)(string)) error {
	var wrapped_desc = wrap(e.Error())
	var e_with_extra, with_extra = e.(ErrorWithExtraData)
	if with_extra {
		return &WrappedErrorWithExtraData {
			Underlying:  e_with_extra,
			Description: wrapped_desc,
		}
	} else {
		return errors.New(wrapped_desc)
	}
}


