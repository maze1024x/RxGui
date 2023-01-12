package typsys

import (
    "fmt"
    "strings"
)


func Equal(t1 Type, t2 Type) bool {
    switch T1 := t1.(type) {
    case InferringType:
        if T2, ok := t2.(InferringType); ok {
            return (T1.Id == T2.Id)
        }
    case ParameterType:
        if T2, ok := t2.(ParameterType); ok {
            return (T1.Name == T2.Name)
        }
    case RefType:
        if T2, ok := t2.(RefType); ok {
            if T1.Def == T2.Def {
                if len(T1.Args) == len(T2.Args) {
                    var n = len(T1.Args)
                    var all_equal = true
                    for i := 0; i < n; i += 1 {
                        var equal = Equal(T1.Args[i], T2.Args[i])
                        if !(equal) {
                            all_equal = false
                            break
                        }
                    }
                    return all_equal
                }
            }
        }
    }
    return false
}

func Transform(t Type, f func(t Type)(Type,bool)) Type {
    switch T := t.(type) {
    case RefType:
        var mapped_args = make([] Type, len(T.Args))
        for i, arg := range T.Args {
            mapped_args[i] = Transform(arg, f)
        }
        var t = RefType {
            Def:  T.Def,
            Args: mapped_args,
        }
        if u, ok := f(t); ok {
            return u
        } else {
            return t
        }
    default:
        if u, ok := f(t); ok {
            return u
        } else {
            return t
        }
    }
}

func Describe(t Type) string {
    switch T := t.(type) {
    case InferringType:
        return ("(" + T.Id + ")")
    case ParameterType:
        return T.Name
    case RefType:
        if len(T.Args) == 0 {
            return T.Def.String()
        } else {
            var name_desc = T.Def.String()
            var n = len(T.Args)
            var arg_desc = make([] string, n)
            for i := 0; i < n; i += 1 {
                arg_desc[i] = Describe(T.Args[i])
            }
            var args_desc = strings.Join(arg_desc, ",")
            return fmt.Sprintf("%s[%s]", name_desc, args_desc)
        }
    default:
        panic("impossible branch")
    }
}

func Inflate(t Type, params ([] string), args ([] Type)) Type {
    return Transform(t, func(t Type) (Type, bool) {
        switch T := t.(type) {
        case InferringType:
            panic("invalid argument")
        case ParameterType:
            for i := range params {
                if params[i] == T.Name {
                    if i < len(args) {
                        return args[i], true
                    }
                }
            }
        }
        return nil, false
    })
}

func DescribeCertain(t CertainType) string {
    return Describe(t.Type)
}

func DescribeWithInferringState(t Type, s *InferringState) string {
    if s == nil {
        return Describe(t)
    } else {
        return Describe(Transform(t, func(t Type) (Type, bool) {
            switch T := t.(type) {
            case InferringType:
                var current, has_current = s.getInferred(T.Id)
                if has_current {
                    return current, true
                }
            }
            return nil, false
        }))
    }
}


