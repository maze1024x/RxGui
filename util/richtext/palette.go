package richtext

import (
    "os"
    "runtime"
)


type AnsiPalette (func(string) string)
var defaultAnsiPalette = (func() AnsiPalette {
    if runtime.GOOS == "windows" {
        var _, is_mintty = os.LookupEnv("TERM")
        var _, is_wt = os.LookupEnv("WT_SESSION")
        if is_mintty || is_wt {
            return lightAnsiPalette
        } else {
            return nil
        }
    } else {
        return lightAnsiPalette
    }
})()

type CssPalette (func([] string) ([] func()(string,string)))
var defaultCssPalette = CssPalette(func(tags ([] string)) ([] func()(string,string)) {
    var rules = make([] func()(string,string), 0)
    var add = func(key string, value string) {
        rules = append(rules, func() (string, string) {
            return key, value
        })
    }
    for _, tag := range tags {
        switch tag {
        case TAG_B:          add("font-weight", "bold")
        case TAG_I:          add("font-style", "italic")
        case TAG_DEL:        add("text-decoration", "line-through")
        case TAG_HIGHLIGHT:  add("font-weight", "bold"); add("color", "#F33"); add("background-color", "#FF3")
        case TAG_ERR:        add("font-weight", "bold"); add("color", "#F33")
        case TAG_ERR_INLINE: add("font-weight", "bold"); add("color", "#A33")
        case TAG_ERR_NOTE:   add("font-weight", "bold"); add("color", "#A1A")
        case TAG_INPUT:      add("color", "#33D")
        case TAG_OUTPUT:     add("color", "#777")
        case TAG_SUCCESS:    add("color", "#3C3")
        case TAG_FAILURE:    add("color", "#D33")
        // colors from Atom One Light
        case TAG_DBG_TYPE:     add("color", "#C18401")
        case TAG_DBG_FIELD:    add("color", "#E45649")
        case TAG_DBG_STRING:   add("color", "#50A14F")
        case TAG_DBG_NUMBER:   add("color", "#4078F2")
        case TAG_DBG_CONSTANT: add("color", "#986801")
        }
    }
    return rules
})

const bold = "\033[1m"
const red = "\033[31m"
const green = "\033[32m"
const orange = "\033[33m"
const blue = "\033[34m"
const magenta = "\033[35m"
const cyan = "\033[36m"
const reset = "\033[0m"
func LightAnsiPalette() AnsiPalette { return lightAnsiPalette }
var lightAnsiPalette AnsiPalette = func(tag string) string {
    switch tag {
    case TAG_B:          return bold
    case TAG_HIGHLIGHT:  return (bold + red)
    case TAG_ERR:        return bold
    case TAG_ERR_INLINE: return (bold + red)
    case TAG_ERR_NOTE:   return (bold + magenta)
    default: return ""
    }
}
func DarkAnsiPalette() AnsiPalette { return darkAnsiPalette }
var darkAnsiPalette AnsiPalette = func(tag string) string {
    switch tag {
    case TAG_B:          return bold
    case TAG_HIGHLIGHT:  return (bold + orange)
    case TAG_ERR:        return bold
    case TAG_ERR_INLINE: return (bold + orange)
    case TAG_ERR_NOTE:   return (bold + cyan)
    default: return ""
    }
}


