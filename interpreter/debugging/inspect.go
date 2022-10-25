package debugging

import (
	"fmt"
	"strconv"
	"strings"
	"rxgui/util/richtext"
	"rxgui/interpreter/lang/source"
	"rxgui/interpreter/lang/typsys"
	"rxgui/interpreter/core"
	"rxgui/interpreter/compiler"
)


type InspectContext struct {
	context  *compiler.NsHeaderMap
}
func MakeInspectContext(ctx *compiler.NsHeaderMap) InspectContext {
	return InspectContext { ctx }
}
func (ctx InspectContext) FindType(ref source.Ref) (*typsys.TypeDef, bool) {
	return ctx.context.FindType(ref)
}

func Inspect(v core.Object, t typsys.Type, ctx InspectContext) richtext.Block {
	var opaque = func(t_desc string) richtext.Block {
		var b richtext.Block
		b.WriteRawSpan(t_desc, richtext.TAG_DBG_TYPE)
		return b
	}
	var primitive = func(v_desc string, v_tag string) richtext.Block {
		var b richtext.Block
		b.WriteSpan(v_desc, v_tag)
		return b
	}
	var write_type_line = func(b *richtext.Block) {
		var t_desc = typsys.Describe(t)
		b.WriteLine(t_desc, richtext.TAG_DBG_TYPE)
	}
	var collection = func(inner_t typsys.Type, size int, forEach func(func(core.Object))) richtext.Block {
		var b richtext.Block
		write_type_line(&b)
		var size_desc = fmt.Sprintf("<%d>", size)
		b.WriteLine(size_desc, richtext.TAG_DBG_NUMBER)
		forEach(func(item core.Object) {
			b.Append(Inspect(item, inner_t, ctx))
		})
		return b
	}
	if t == nil {
		return opaque("(?)")
	}
	switch T := t.(type) {
	case typsys.InferringType:
		return opaque(T.Id)
	case typsys.ParameterType:
		return opaque(T.Name)
	case typsys.RefType:
		var t_ref = T.Def
		if t_ref.Namespace == "" {
			switch t_ref.ItemName {
			case core.T_Null:
				var b richtext.Block
				b.WriteLine("Null", richtext.TAG_DBG_CONSTANT)
				return b
			case core.T_Bool:
				if core.GetBool(v) {
					return primitive("Yes", richtext.TAG_DBG_CONSTANT)
				} else {
					return primitive("No", richtext.TAG_DBG_CONSTANT)
				}
			case core.T_Char:
				var char = core.GetChar(v)
				var char_str = string([] rune { rune(char) })
				var char_str_desc = strconv.Quote(char_str)
				var char_desc = fmt.Sprintf("\\u%X %s", char, char_str_desc)
				return primitive(char_desc, richtext.TAG_DBG_NUMBER)
			case core.T_String:
				var str_desc = strconv.Quote(core.GetString(v))
				return primitive(str_desc, richtext.TAG_DBG_STRING)
			case core.T_Bytes:
				var bin = core.GetBytes(v)
				var buf strings.Builder
				for _, n := range bin {
					var hi = ((n & 0xF0) >> 4)
					var lo = (n & 0x0F)
					fmt.Fprintf(&buf, `\x%X%X`, hi, lo)
				}
				var bin_desc = buf.String()
				return primitive(bin_desc, richtext.TAG_DBG_NUMBER)
			case core.T_Int:
				var n = core.GetIntAsRawBigInt(v)
				var n_desc = n.String()
				return primitive(n_desc, richtext.TAG_DBG_NUMBER)
			case core.T_Float:
				var x = core.GetFloat(v)
				var x_desc = fmt.Sprint(x)
				if !(strings.Contains(x_desc, ".") ||
				strings.Contains(x_desc, "e") ||
				strings.Contains(x_desc, "E") ||
				strings.Contains(x_desc, "NaN") ||
				strings.Contains(x_desc, "Inf")) {
					x_desc = (x_desc + ".0")
				}
				return primitive(x_desc, richtext.TAG_DBG_NUMBER)
			case core.T_Time:
				var tau = core.GetTime(v)
				var tau_desc = tau.String()
				return primitive(tau_desc, richtext.TAG_DBG_NUMBER)
			case core.T_File:
				var f = core.GetFile(v)
				var f_desc = f.Path
				return primitive(f_desc, richtext.TAG_DBG_STRING)
			case core.T_Error:
				var b richtext.Block
				b.WriteLine(typsys.Describe(t), richtext.TAG_DBG_TYPE)
				var child richtext.Block
				var err = core.GetError(v)
				var err_desc = strconv.Quote(err.Error())
				child.WriteLine(err_desc, richtext.TAG_DBG_STRING)
				b.Append(child)
				return b
			case core.T_Lambda:
				return opaque(typsys.Describe(t))
			case core.T_List:
				var l = core.GetList(v)
				return collection(T.Args[0], l.Length(), l.ForEach)
			case core.T_Seq:
				var s = core.GetSeq(v)
				return collection(T.Args[0], s.Length(), s.ToList().ForEach)
			case core.T_Queue:
				var q = core.GetQueue(v)
				return collection(T.Args[0], q.Size(), q.ForEach)
			case core.T_Heap:
				var h = core.GetHeap(v)
				return collection(T.Args[0], h.Size(), h.ForEach)
			case core.T_Set:
				var s = core.GetSet(v)
				return collection(T.Args[0], s.Size(), s.ForEach)
			case core.T_Map:
				var b richtext.Block
				b.WriteLine(typsys.Describe(t), richtext.TAG_DBG_TYPE)
				var m = core.GetMap(v)
				var size_desc = fmt.Sprintf("<%d>", m.Size())
				b.WriteLine(size_desc, richtext.TAG_DBG_NUMBER)
				m.ForEach(func(key core.Object, val core.Object) {
					var item richtext.Block
					item.WriteLine("*")
					item.Append(Inspect(key, T.Args[0], ctx))
					item.Append(Inspect(val, T.Args[1], ctx))
					b.Append(item)
				})
				return b
			}
		}
		var def, ok = ctx.FindType(t_ref)
		if ok {
			switch content := def.Content.(type) {
			case typsys.Record:
				var V = (*v).(core.Record)
				var b richtext.Block
				write_type_line(&b)
				for i, field := range content.FieldList {
					var field_t = typsys.Inflate(
						field.Type, def.Parameters, T.Args,
					)
					var child richtext.Block
					child.WriteSpan(field.Name, richtext.TAG_DBG_FIELD)
					child.Append(Inspect(V.Objects[i], field_t, ctx))
					b.Append(child)
				}
				return b
			case typsys.Union:
				var u = (*v).(core.Union)
				var b richtext.Block
				write_type_line(&b)
				var index_desc = fmt.Sprintf("(%d)", u.Index)
				b.WriteLine(index_desc, richtext.TAG_DBG_NUMBER)
				var item = content.FieldList[u.Index]
				var case_t = typsys.Inflate(item.Type, def.Parameters, T.Args)
				var inner = Inspect(u.Object, case_t, ctx)
				b.Append(inner)
				return b
			case typsys.Enum:
				var index = int((*v).(core.Enum))
				var item = content.FieldList[index]
				var b richtext.Block
				write_type_line(&b)
				var item_desc = fmt.Sprintf("%s(%d)", item.Name, index)
				b.WriteLine(item_desc, richtext.TAG_DBG_NUMBER)
				return b
			}
		}
		var b richtext.Block
		write_type_line(&b)
		return b
	default:
		panic("impossible branch")
	}
}


