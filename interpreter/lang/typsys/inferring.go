package typsys

import "rxgui/standalone/ctn"


func Match(x Type, y Type, s *InferringState) (bool, *InferringState) {
	if x == nil || y == nil {
		panic("something went wrong")
	}
	if s != nil {
		var X, to_inferring = x.(InferringType)
		var Y, from_inferring = y.(InferringType)
		var inferring string
		var certain Type
		var ok = false
		if to_inferring && IsCertain(y) {
			inferring = X.Id
			certain = y
			ok = true
		}
		if IsCertain(x) && from_inferring {
			inferring = Y.Id
			certain = x
			ok = true
		}
		if ok {
			var inferred, exists = s.getInferred(inferring)
			if exists {
				if Equal(certain, inferred) {
					return true, s
				} else {
					return false, nil
				}
			} else {
				return true, s.newStateSetInferred(inferring, certain)
			}
		}
	}
	switch X := x.(type) {
	case InferringType:
		if Y, ok := y.(InferringType); ok {
			return (X.Id == Y.Id), s
		}
	case ParameterType:
		if Y, ok := y.(ParameterType); ok {
			return (X.Name == Y.Name), s
		}
	case RefType:
		if Y, ok := y.(RefType); ok {
			if (X.Def == Y.Def) {
				if len(X.Args) == len(Y.Args) {
					var n = len(X.Args)
					var all_assignable = true
					for i := 0; i < n; i += 1 {
						var assignable, s1 = Match(X.Args[i], Y.Args[i], s)
						if assignable {
							s = s1
						} else {
							all_assignable = false
							break
						}
					}
					if all_assignable {
						return true, s
					}
				}
			}
		}
	}
	return false, nil
}
func MatchAll(x ([] Type), y ([] Type), s *InferringState) (bool, *InferringState) {
	if len(x) != len(y) {
		return false, nil
	}
	var L = len(x)
	for i := 0; i < L; i += 1 {
		var ok, s1 = Match(x[i], y[i], s)
		if ok {
			s = s1
		} else {
			return false, nil
		}
	}
	return true, s
}

func IsCertain(t Type) bool {
	switch T := t.(type) {
	case InferringType:
		return false
	case RefType:
		for _, arg := range T.Args {
			if !(IsCertain(arg)) {
				return false
			}
		}
	}
	return true
}
func ToInferring(t Type, params ([] string)) Type {
	return Transform(t, func(t Type) (Type, bool) {
		switch T := t.(type) {
		case InferringType:
			panic("invalid arguments")
		case ParameterType:
			for _, p := range params {
				if p == T.Name {
					var inferring = InferringType { p }
					return inferring, true
				}
			}
		}
		return nil, false
	})
}
func GetCertainOrInferred(t Type, s *InferringState) (CertainType, bool) {
	if t == nil {
		panic("invalid argument")
	}
	if IsCertain(t) {
		return CertainType { t }, true
	} else {
		if s != nil {
			return s.GetInferredType(t)
		} else {
			return CertainType {}, false
		}
	}
}

type InferringState struct {
	parameters  [] string
	mapping     ctn.Map[string,Type]
}
func Infer(parameters ([] string)) *InferringState {
	return &InferringState {
		parameters: parameters,
		mapping:    ctn.MakeMap[string,Type](ctn.StringCompare),
	}
}
func (s *InferringState) getInferred(name string) (Type, bool) {
	if s == nil {
		panic("invalid operation")
	}
	var value, exists = s.mapping.Lookup(name)
	return value, exists
}
func (s *InferringState) newStateSetInferred(name string, value Type) *InferringState {
	if s == nil {
		panic("invalid operation")
	}
	return &InferringState {
		parameters: s.parameters,
		mapping:    s.mapping.Inserted(name, value),
	}
}
func (s *InferringState) GetInferredType(t Type) (CertainType, bool) {
	if s == nil {
		panic("invalid operation")
	}
	var ok = true
	var inferred_t = Transform(t, func(t Type) (Type, bool) {
		if T, is_inferring := t.(InferringType); is_inferring {
			var value, exists = s.mapping.Lookup(T.Id)
			if exists {
				return value, true
			} else {
				ok = false
				return nil, false
			}
		} else {
			return nil, false
		}
	})
	if ok {
		return CertainType { inferred_t }, true
	} else {
		return CertainType {}, false
	}
}
func (s *InferringState) GetInferredParameterNamespace(name string) (string, bool) {
	if s == nil {
		panic("invalid operation")
	}
	var value, exists = s.mapping.Lookup(name)
	if exists {
		var ref_type, is_ref_type = value.(RefType)
		if is_ref_type {
			var ns = ref_type.Def.Namespace
			return ns, true
		} else {
			return "", false
		}
	} else {
		return "", false
	}
}


