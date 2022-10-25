package source

import (
	"rxgui/util/richtext"
)


type Errors ([] *Error)
type Error struct {
	Location  Location
	Content   ErrorContent
}
type ErrorContent interface {
	DescribeError() richtext.Block
}

// ----------------------
func ErrorsFrom(err *Error) Errors {
	if err != nil {
		return Errors([] *Error { err })
	} else {
		return nil
	}
}
func ErrorsJoin(errs *Errors, err *Error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}
func ErrorsJoinAll(errs *Errors, another Errors) {
	for _, err := range another {
		*errs = append(*errs, err)
	}
}
func (errs Errors) Message() richtext.Block {
	if len(errs) == 0 {
		panic("invalid operation")
	}
	var b richtext.Block
	for _, item := range errs {
		b.Append(item.Message())
	}
	return b
}
func (errs Errors) Error() string {
	return errs.Message().RenderConsole()
}

// ----------------------
func MakeError(loc Location, content ErrorContent) *Error {
	return &Error {
		Location: loc,
		Content:  content,
	}
}
func (e *Error) Description() richtext.Block {
	return e.Content.DescribeError()
}
func (e *Error) Message() richtext.Block {
	return e.Location.FormatMessage(e.Description())
}
func (e *Error) Error() string {
	return e.Message().RenderPlainText()
}


