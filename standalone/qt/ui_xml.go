package qt

import (
	"fmt"
	"errors"
	"unicode"
	"reflect"
)


func LoadWidget(ui_xml string, base_dir string, ctx Pkg, k func(Widget)(error)) error {
	var w, err = loadWidget(ui_xml, base_dir, ctx)
	if err != nil {
		return fmt.Errorf("failed to load widget from ui xml: %w", err)
	}
	return k(w)
}

func (obj Object) FindChild(name string, v reflect.Value) error {
	var child, ok = findChild(obj, name)
	if !(ok) {
		return errors.New(fmt.Sprintf("child Object %s not found", name))
	}
	var err = child.AssignTo(v)
	if err != nil {
		return fmt.Errorf("child %s: %w", name, err)
	}
	return nil
}

func (obj Object) AssignTo(v reflect.Value) error {
	if v.Kind() != reflect.Ptr { panic("invalid argument") }
	var class = obj.ClassName()
	var required = (func() string {
		var type_name = v.Elem().Type().Name()
		return ("Q" + type_name)
	})()
	if class == required {
		var target = v.Elem().Field(0)
		for !(reflect.TypeOf(obj).AssignableTo(target.Type())) {
			target = target.Field(0)
		}
		target.Set(reflect.ValueOf((interface{})(obj)))
		return nil
	} else {
		return errors.New(fmt.Sprintf(
			"class not matching: %s (%s)", class, required))
	}
}

func (obj Object) AssignChildrenTo(v reflect.Value) error {
	if v.Kind() != reflect.Ptr { panic("invalid argument") }
	if v.Elem().Kind() != reflect.Struct { panic("invalid argument") }
	var struct_v = v.Elem()
	var struct_t = struct_v.Type()
	for i := 0; i < struct_t.NumField(); i += 1 {
		var field_name = struct_t.Field(i).Name
		var field_type = struct_t.Field(i).Type
		var t = ([] rune)(field_type.Name())
		t[0] = unicode.ToLower(t[0])
		var field_type_name = string(t)
		var child_name = (field_type_name + field_name)
		var err = obj.FindChild(child_name, struct_v.Field(i).Addr())
		if err != nil { return err }
	}
	return nil
}


