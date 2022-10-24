package atom

import (
	"path/filepath"
	"rxgui/lang/source"
	"rxgui/lang/textual/cst"
	"rxgui/lang/textual/parser"
	"rxgui/lang/textual/scanner"
	"rxgui/lang/textual/transformer"
	"rxgui/interpreter"
)


type LintRequest struct {
	Path            string      `json:"path"`
	VisitedModules  [] string   `json:"visitedModules"`
}

type LintResponse struct {
	Module   string         `json:"module"`
	Errors   [] LintError   `json:"errors"`
}

type LintError struct {
	Severity     string         `json:"severity"`
	Location     LintLocation   `json:"location"`
	Excerpt      string         `json:"excerpt"`
	Description  string         `json:"description"`
}

type LintLocation struct {
	File      string   `json:"file"`
	Position  Range    `json:"position"`
}

type Range struct {
	Start  Point   `json:"start"`
	End    Point   `json:"end"`
}

type Point struct {
	Row  int   `json:"row"`
	Col  int   `json:"column"`
}

func GetPoint(point scanner.Point) Point {
	var row int
	var col int
	row = (point.Row - 1)
	if point.Col == 0 {
		col = 1
	} else {
		col = (point.Col - 1)
	}
	return Point { row, col }
}

func GetLocation(file string, info scanner.RowColInfo, span scanner.Span) LintLocation {
	if span == (scanner.Span {}) {
		// empty span
		return LintLocation {
			File:     file,
			Position: Range {
				Start: Point { 0, 0 },
				End:   Point { 0, 1 },
			},
		}
	} else {
		var start = info[span.Start]
		var end = info[span.End]
		return LintLocation {
			File:     file,
			Position: Range {
				Start: GetPoint(start),
				End:   GetPoint(end),
			},
		}
	}
}

func AdaptLocation(loc source.Location) LintLocation {
	var key = loc.File.GetKey()
	var f, is_textual = loc.File.(*transformer.File)
	if is_textual {
		var span = loc.Pos.Span
		return GetLocation(key.Path, f.Info, span)
	} else {
		return LintLocation {
			File: key.Path,
		}
	}
}

func GetError(e *source.Error, tip string) LintError {
	var loc = e.Location
	var desc = e.Description()
	return LintError {
		Severity:    "error",
		Location:    AdaptLocation(loc),
		Excerpt:     desc.RenderPlainText(),
		Description: tip,
	}
}

func Lint(req LintRequest, ctx LangServerContext) LintResponse {
	if filepath.Base(filepath.Dir(req.Path)) == "builtin" {
		return LintResponse {
			Module: req.Path,
		}
	}
	var errs1, errs2 = interpreter.Lint(req.Path, nil)
	// ctx.WriteDebugLog("Lint Path: " + req.Path)
	if errs1 != nil {
		var err = errs1[0]
		if e, ok := err.(*parser.Error);
		(ok && !(e.IsScannerError || e.IsEmptyTree())) {
			var token = cst.GetNodeFirstToken(e.Tree, e.NodeIndex)
			var loc = GetLocation(req.Path, e.Tree.Info, token.Span)
			return LintResponse {
				Module: req.Path,
				Errors: [] LintError {{
					Severity:    "error",
					Location:    loc,
					Excerpt:     e.Desc().RenderPlainText(),
					Description: "",
				}},
			}
		} else {
			// var msg = err.Message().RenderPlainText()
			// ctx.WriteDebugLog(req.Path + " unable to lint: " + msg)
			return LintResponse {
				Module: req.Path,
			}
		}
	}
	if errs2 != nil {
		var all = make([] LintError, 0)
		for _, e := range errs2 {
			all = append(all, GetError(e, ""))
		}
		return LintResponse {
			Module: req.Path,
			Errors: all,
		}
	}
	return LintResponse {
		Module: req.Path,
	}
}


