package compiler

import (
    "fmt"
    "strings"
    "rxgui/util/richtext"
    "rxgui/interpreter/lang/source"
)


const BlockClassError = "error"
const BlockClassErrorContentItem = "error-content-item"

func makeErrorDescBlankBlock() richtext.Block {
    var b richtext.Block
    b.AddClass(BlockClassError)
    b.WriteSpan("Error: ", richtext.TAG_B)
    return b
}
func makeErrorDescBlock(msg ...string) richtext.Block {
    var b = makeErrorDescBlankBlock()
    for i, span := range msg {
        b.WriteSpan(span, (func() string {
            if (i % 2) == 0 {
                if strings.HasPrefix(span, "(") && strings.HasSuffix(span, ")") {
                    return richtext.TAG_ERR_NOTE
                } else {
                    return richtext.TAG_ERR
                }
            } else {
                return richtext.TAG_ERR_INLINE
            }
        })())
    }
    return b
}
func makeEmptyErrorContentItemBlock() richtext.Block {
    var b richtext.Block
    b.AddClass(BlockClassErrorContentItem)
    return b
}

type E_DuplicateAlias struct {
    Name  string
}
func (e E_DuplicateAlias) DescribeError() richtext.Block {
    if e.Name == "" {
        return makeErrorDescBlock (
            "duplicate alias",
        )
    } else {
        return makeErrorDescBlock (
            "duplicate alias:",
            e.Name,
        )
    }
}

type E_InvalidAlias struct {
    Name  string
}
func (e E_InvalidAlias) DescribeError() richtext.Block {
    if e.Name == "" {
        return makeErrorDescBlock (
            "invalid alias",
        )
    } else {
        return makeErrorDescBlock (
            "invalid alias:",
            e.Name,
        )
    }
}

type E_AliasTargetNotFound struct {
    Target  string
}
func (e E_AliasTargetNotFound) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "alias target not found:",
        e.Target,
    )
}

type E_DuplicateTypeDecl struct {
    Name  string
}
func (e E_DuplicateTypeDecl) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "duplicate declaration of type",
        e.Name,
    )
}

type E_DuplicateFunDecl struct {
    Name   string
    Assoc  string
}
func (e E_DuplicateFunDecl) DescribeError() richtext.Block {
    if ((e.Name == "") && (e.Assoc == "")) {
        return makeErrorDescBlock (
            "duplicate declaration of entry point in current namespace",
        )
    } else {
        return makeErrorDescBlock (
            "duplicate declaration of global name",
            fmt.Sprintf("%s (%s)", e.Name, e.Assoc),
            "in current namespace",
        )
    }
}

type E_InvalidFieldName struct {
    FieldKind  string
    FieldName  string
}
func (e E_InvalidFieldName) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        fmt.Sprintf("invalid %s name", e.FieldKind),
        e.FieldName,
    )
}

type E_DuplicateField struct {
    FieldKind  string
    FieldName  string
}
func (e E_DuplicateField) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        fmt.Sprintf("duplicate %s", e.FieldKind),
        e.FieldName,
    )
}

type E_GenericEnum struct {}
func (E_GenericEnum) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "a generic type cannot be an enum type",
    )
}

type E_NoSuchType struct {
    TypeRef  string
}
func (e E_NoSuchType) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "no such type:",
        e.TypeRef,
    )
}

type E_NoSuchTypeParameter struct {
    Name  string
}
func (e E_NoSuchTypeParameter) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "no such type parameter:",
        e.Name,
    )
}

type E_NotRecord struct {
    TypeRef  string
}
func (e E_NotRecord) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "type",
        e.TypeRef,
        "is not a record",
    )
}

type E_InvalidRecordTag struct {
    Tag  string
}
func (e E_InvalidRecordTag) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid record tag:",
        e.Tag,
    )
}

type E_NotInterface struct {
    TypeRef  string
}
func (e E_NotInterface) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "type",
        e.TypeRef,
        "is not an interface",
    )
}

type E_DuplicateInterface struct {
    TypeRef  string
}
func (e E_DuplicateInterface) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "duplicate interface",
        e.TypeRef,
    )
}

type E_TypeParamsNotIdentical struct {
    Concrete   string
    Interface  string
}
func (e E_TypeParamsNotIdentical) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "identical type parameters required for type",
        e.Concrete,
        "to implement interface",
        e.Interface,
    )
}

type E_MissingMethod struct {
    Concrete   string
    Interface  string
    Method     string
}
func (e E_MissingMethod) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "missing method",
        e.Method,
        "for type",
        e.Concrete,
        "to implement interface",
        e.Interface,
        "(note: method should be in the same file)",
    )
}

type E_WrongMethodType struct {
    Concrete   string
    Interface  string
    Method     string
    Expected   string
    Actual     string
}
func (e E_WrongMethodType) DescribeError() richtext.Block {
    var b = makeErrorDescBlock (
        "bad method",
        e.Method,
        "for type",
        e.Concrete,
        "to implement interface",
        e.Interface,
    )
    b.Append(makeErrorDescBlock (
        "expect method type to be",
        e.Expected,
        "but got",
        e.Actual,
    ))
    return b
}

type E_MethodNameUnavailable struct {
    Name  string
}
func (e E_MethodNameUnavailable) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "unavailable method name",
        e.Name,
        "(conflicts with record field or abstract method)",
    )
}

type E_TypeArgsWrongQuantity struct {
    Type      string
    Given     int
    Required  int
}
func (e E_TypeArgsWrongQuantity) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        fmt.Sprintf(
            "expect %d type argument(s) but given %d for type",
            e.Required,
            e.Given,
        ),
        e.Type,
    )
}

type E_CannotMatchRecord struct {
    TypeDesc  string
}
func (e E_CannotMatchRecord) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "cannot match record from type",
        e.TypeDesc,
    )
}

type E_RecordSizeNotMatching struct {
    PatternArity  int
    RecordSize    int
    RecordDesc    string
}
func (e E_RecordSizeNotMatching) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "cannot match record",
        e.RecordDesc,
        fmt.Sprintf(
            "(size not matching: pattern(%d) / record(%d))",
            e.PatternArity, e.RecordSize,
        ),
    )
}

type E_DuplicateBinding struct {
    BindingName  string
}
func (e E_DuplicateBinding) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "duplicate binding:",
        e.BindingName,
    )
}

type E_UnusedBinding struct {
    BindingName  string
}
func (e E_UnusedBinding) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "unused binding:",
        e.BindingName,
    )
}

type E_LambdaAssignedToIncompatibleType struct {
    TypeDesc  string
}
func (e E_LambdaAssignedToIncompatibleType) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "lambda cannot be assigned to incompatible type",
        e.TypeDesc,
    )
}

type E_ExpectExplicitTypeCast struct {}
func (E_ExpectExplicitTypeCast) DescribeError() richtext.Block {
    return makeErrorDescBlock("expect explicit type cast")
}

type E_ExpectSufficientTypeArguments struct {}
func (E_ExpectSufficientTypeArguments) DescribeError() richtext.Block {
    return makeErrorDescBlock("expect sufficient type arguments")
}

type E_NotAssignable struct {
    From  string
    To    string
}
func (e E_NotAssignable) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "cannot assign from", e.From,
        "to", e.To,
    )
}

type E_AmbiguousAssignmentToUnion struct {
    Union  string
    Key1   string
    Key2   string
}
func (e E_AmbiguousAssignmentToUnion) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "ambiguous assignment to union type", e.Union,
        "as", e.Key1,
        "and", e.Key2,
        "are both assignable items",
    )
}

type E_NoSuchFieldOrMethod struct {
    FieldName  string
    TypeDesc   string
}
func (e E_NoSuchFieldOrMethod) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "no such field/method",
        e.FieldName,
        "on type",
        e.TypeDesc,
    )
}

type E_InteriorRefUnavailable struct {
    InteriorRef  string
    TypeDesc     string
}
func (e E_InteriorRefUnavailable) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "interior reference",
        e.InteriorRef,
        "is unavailable for type",
        e.TypeDesc,
    )
}

type E_InvalidFloat struct {}
func (E_InvalidFloat) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid floating-point number",
    )
}

type E_InvalidText struct {}
func (e E_InvalidText) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid text",
    )
}

type E_InvalidChar struct {}
func (e E_InvalidChar) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid character",
    )
}

type E_InvalidRegexp struct {
    Detail  string
}
func (e E_InvalidRegexp) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid regular expression",
        fmt.Sprintf("(%s)", e.Detail),
    )
}

type E_TooManyTypeArgs struct {}
func (e E_TooManyTypeArgs) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "too many type arguments",
    )
}

type E_MissingOperatorParameter struct {}
func (E_MissingOperatorParameter) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "operator must have at least one parameter",
    )
}

type E_OperatorFirstParameterHasDefaultValue struct {}
func (E_OperatorFirstParameterHasDefaultValue) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "operator cannot have default value on first parameter",
    )
}

type E_MissingVariadicParameter struct {}
func (E_MissingVariadicParameter) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "variadic function/operator must have at least one parameter",
    )
}

type E_InvalidVariadicParameter struct {}
func (E_InvalidVariadicParameter) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "variadic function/operator must have a final list-parameter",
    )
}

type E_DuplicateArgument struct {
    Name  string
}
func (e E_DuplicateArgument) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "duplicate argument",
        e.Name,
    )
}

type E_MissingArgument struct {
    Name  string
}
func (e E_MissingArgument) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "missing argument:",
        e.Name,
    )
}

type E_SuperfluousArgument struct {}
func (E_SuperfluousArgument) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "superfluous argument",
    )
}

type E_UnableToInferDefaultValueType struct {
    ArgName  string
}
func (e E_UnableToInferDefaultValueType) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "unable to infer the type of default value of argument",
        e.ArgName,
    )
}

type E_UnableToInferVaType struct {}
func (E_UnableToInferVaType) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "unable to infer variadic argument type",
    )
}

type E_NoSuchThing struct {
    Ref  string
}
func (e E_NoSuchThing) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "no such thing:",
        e.Ref,
    )
}

type E_NoSuchFunction struct {
    FunKindDesc  string
    FunNameDesc  string
}
func (e E_NoSuchFunction) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        fmt.Sprintf("no such %s:", e.FunKindDesc),
        e.FunNameDesc,
    )
}

type E_NotCallable struct {
    TypeDesc  string
}
func (e E_NotCallable) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "cannot call a value of type",
        e.TypeDesc,
    )
}

type E_LambdaCallWrongArgsQuantity struct {
    Given     int
    Required  int
}
func (e E_LambdaCallWrongArgsQuantity) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        fmt.Sprintf("expect %d argument(s) but given %d", e.Required, e.Given),
    )
}

type E_LambdaCallUnorderedArgs struct {}
func (E_LambdaCallUnorderedArgs) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "cannot call a lambda using unordered arguments",
    )
}

type E_InvalidConstructorUsage struct {}
func (E_InvalidConstructorUsage) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid constructor usage",
    )
}

type E_SuperfluousTypeArgs struct {}
func (E_SuperfluousTypeArgs) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "superfluous type arguments",
    )
}

type E_CannotAssignFunctionRequiringImplicitInput struct {}
func (E_CannotAssignFunctionRequiringImplicitInput) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "cannot assign a function/operator that requires implicit input",
    )
}

type E_UnableToUseAsLambda struct {
    InOutDesc  string
}
func (e E_UnableToUseAsLambda) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "unable to use as lambda: unqualified signature",
        e.InOutDesc,
    )
}

type E_MultipleAssignable struct {
    Ref  string
}
func (e E_MultipleAssignable) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "multiple things available for assignment of",
        e.Ref,
    )
}

type E_NoneAssignable struct {
    Ref      string
    Details  [] NoneAssignableErrorDetail
}
type NoneAssignableErrorDetail struct {
    ItemName      string
    ErrorContent  source.ErrorContent
}
func (e E_NoneAssignable) DescribeError() richtext.Block {
    var b = makeErrorDescBlock (
        "nothing qualified for assignment of",
        e.Ref,
    )
    for _, detail := range e.Details {
        var desc = detail.ErrorContent.DescribeError()
        desc = desc.WithoutLeadingSpan()
        desc = desc.WithLeadingSpan((detail.ItemName + ":"), richtext.TAG_B)
        b.Append(desc)
    }
    return b
}

type E_InvalidCondType struct {
    TypeDesc  string
}
func (e E_InvalidCondType) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "cannot use type",
        e.TypeDesc,
        "as a condition",
    )
}

type E_InvalidCondPattern struct {}
func (E_InvalidCondPattern) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid pattern",
    )
}

type E_InvalidWhenOperand struct {
    TypeDesc  string
}
func (e E_InvalidWhenOperand) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid when expression operand type",
        e.TypeDesc,
    )
}

type E_InvalidEachOperand struct {
    TypeDesc  string
}
func (e E_InvalidEachOperand) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid each expression operand",
        e.TypeDesc,
    )
}

type E_InvalidCasePattern struct {}
func (E_InvalidCasePattern) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "invalid pattern",
    )
}

type E_NoSuchCase struct {
    CaseName  string
}
func (e E_NoSuchCase) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "no such case:",
        e.CaseName,
    )
}

type E_DuplicateCase struct {
    CaseName  string
}
func (e E_DuplicateCase) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "duplicate case",
        e.CaseName,
    )
}

type E_MissingCase struct {
    CaseName  string
}
func (e E_MissingCase) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "non-exhaustive cases: missing case",
        e.CaseName,
    )
}

type E_SuperfluousDefaultCase struct {}
func (E_SuperfluousDefaultCase) DescribeError() richtext.Block {
    return makeErrorDescBlock (
        "superfluous default case",
    )
}


