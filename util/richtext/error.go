package richtext


type Error interface {
    error
    Message() Block
}
func ErrorFrom(err error) Error {
    if err == nil {
        return nil
    }
    return StdError { err }
}

type StdError struct { Underlying error }
func Std2Error(err error) StdError {
    return StdError { err }
}
func (e StdError) Error() string {
    return e.Underlying.Error()
}
func (e StdError) Message() Block {
    var b Block
    b.WriteLine(e.Error(), TAG_ERR)
    return b
}

type Errors ([] Error)
func Std2Errors(err error) Errors {
    return Errors { Std2Error(err) }
}
func ErrorsJoin(errs *Errors, err Error) {
    if err != nil {
        *errs = append(*errs, err)
    }
}
func (errs Errors) Message() Block {
    if len(errs) == 0 {
        panic("invalid operation")
    }
    var b Block
    for _, item := range errs {
        b.Append(item.Message())
    }
    return b
}
func (errs Errors) Error() string {
    return errs.Message().RenderConsole()
}


