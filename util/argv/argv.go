package argv

import (
    "fmt"
    "errors"
    "reflect"
    "strings"
    "strconv"
)

func ParseArgs[T any] (args ([] string)) (T, string, error) {
    var struct_ T
    var help string
    var err = parseArgs(args, &help, reflect.ValueOf(&struct_))
    return struct_, help, err
}
func parseArgs(args ([] string), help *string, ptr reflect.Value) error {
    if ptr.Kind() != reflect.Ptr {
        panic("invalid argument: expect a pointer to write data")
    }
    if ptr.Elem().Kind() != reflect.Struct {
        panic("invalid argument: expect a struct pointer to write fields")
    }
    const key_prefix = "--"
    const key_val_sep = '='
    var positional = make([] string, 0)
    var named = make(map[string] string)
    var named_used = make(map[string] bool)
    var no_more_named = false
    var arg0 = args[0]
    for _, arg := range args[1:] {
        if arg == key_prefix {
            no_more_named = true
        } else if !(no_more_named) && strings.HasPrefix(arg, key_prefix) {
            var arg = strings.TrimPrefix(arg, key_prefix)
            var n = strings.IndexRune(arg, key_val_sep)
            var key, value = (func() (string, string) {
                if n == -1 {
                    return arg, ""
                } else {
                    return arg[:n], arg[(n + 1):]
                }
            })()
            named[key] = value
        } else {
            positional = append(positional, arg)
        }
    }
    var positional_hint string
    var named_hint = make(map[string] string)
    var named_desc = make(map[string] string)
    var commands ([] string)
    var arity = make(map[string] int)
    var default_command string
    var options = make([] string, 0)
    var obj = ptr.Elem()
    var t = obj.Type()
    var first_err = error(nil)
    for i := 0; i < t.NumField(); i += 1 {
        var field = t.Field(i)
        var kind = field.Tag.Get("arg")
        var value, err = (func() (interface{}, error) {
            switch kind {
            case "positional":
                positional_hint = field.Tag.Get("hint")
                return ([] string)(positional), nil
            case "command":
                const sep = "; "
                var key = strings.Split(field.Tag.Get("key"), sep)
                var desc = strings.Split(field.Tag.Get("desc"), sep)
                for i := range key {
                    var parts = strings.Split(key[i], "-")
                    if len(parts) >= 2 {
                        key[i] = parts[0]
                        var n, _ = strconv.Atoi(parts[1])
                        arity[key[i]] = n
                    } else {
                        arity[key[i]] = -1
                    }
                }
                commands = key
                default_command = field.Tag.Get("default")
                for i := range key {
                    if i < len(desc) {
                        named_desc[key[i]] = desc[i]
                    }
                }
                var command = ""
                for _, item := range key {
                    var _, has_item = named[item]
                    if has_item {
                        if command == "" {
                            command = item
                            named_used[item] = true
                        } else {
                            return nil, errors.New(fmt.Sprintf(
                                "ambiguous command: '%s' or '%s' ?",
                                command, item,
                            ))
                        }
                    }
                }
                if command == "" {
                    command = default_command
                }
                var command_arity = arity[command]
                if command_arity == 0 {
                    if len(positional) > 0 {
                        return nil, errors.New(fmt.Sprintf(
                            "redundant argument(s) for '%s' command",
                            command,
                        ))
                    }
                } else if command_arity > 0 {
                    if len(positional) < command_arity {
                        return nil, errors.New(fmt.Sprintf(
                            "not enough arguments for '%s' command",
                            command,
                        ))
                    }
                }
                return string(command), nil
            case "flag-enable":
                var key = field.Tag.Get("key")
                options = append(options, key)
                named_desc[key] = field.Tag.Get("desc")
                var _, flag_is_set = named[key]
                var enabled = flag_is_set
                named_used[key] = true
                return bool(enabled), nil
            case "flag-disable":
                var key = field.Tag.Get("key")
                options = append(options, key)
                named_desc[key] = field.Tag.Get("desc")
                var _, flag_is_set = named[key]
                var enabled = !(flag_is_set)
                named_used[key] = true
                return bool(enabled), nil
            case "value-string":
                var key = field.Tag.Get("key")
                options = append(options, key)
                named_hint[key] = field.Tag.Get("hint")
                named_desc[key] = field.Tag.Get("desc")
                named_used[key] = true
                return string(named[key]), nil
            default:
                return nil, errors.New("invalid argument kind: " + kind)
            }
        })()
        if err == nil {
            obj.Field(i).Set(reflect.ValueOf(value))
        } else {
            if first_err == nil {
                first_err = err
            }
        }
    }
    if first_err == nil {
        for key, _ := range named {
            if !(named_used[key]) {
                first_err = errors.New("unknown option: " + key)
                break
            }
        }
    }
    var buf strings.Builder
    var cmd_opt_hint = "[COMMAND|OPTION]..."
    if positional_hint == "" {
        positional_hint = "[ARGUMENT]..."
    }
    fmt.Fprintf(&buf, "Usage: %s %s %s\n", arg0, cmd_opt_hint, positional_hint)
    const pad1 = "  "
    const pad2 = "  \t"
    if len(commands) > 0 {
        buf.WriteRune('\n')
        buf.WriteString("Commands:")
        buf.WriteRune('\n')
        for _, item := range commands {
            var desc = named_desc[item]
            buf.WriteString(pad1 + key_prefix + item)
            if item == default_command {
                buf.WriteString(pad2 + "(default) " + desc)
            } else {
                buf.WriteString(pad2 + desc)
            }
            buf.WriteRune('\n')
        }
    }
    if len(options) > 0 {
        buf.WriteRune('\n')
        buf.WriteString("Options:")
        buf.WriteRune('\n')
        for _, item := range options {
            var hint = named_hint[item]
            var desc = named_desc[item]
            if hint != "" {
                var key_val_sep = string([] rune { key_val_sep })
                buf.WriteString(pad1 + key_prefix + item + key_val_sep + hint)
            } else {
                buf.WriteString(pad1 + key_prefix + item)
            }
            buf.WriteString(pad2 + desc)
            buf.WriteRune('\n')
        }
    }
    *help = buf.String()
    return first_err
}


