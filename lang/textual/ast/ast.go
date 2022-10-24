package ast

import "rxgui/lang/source"


type Code = source.Code

type Node struct {
	Location  source.Location
}
func NodeFromLocation(loc source.Location) Node {
	return Node { Location: loc }
}

func IterateNodeRegistry(f func(interface{})) {
	for _, node := range nodeRegistry {
		f(node)
	}
}
var nodeRegistry = [...] interface{} {
	Root {},
	Alias {},
	VariousAliasTarget {},
	AliasToNamespace {},
	AliasToRefBase {},
	VariousStatement {},
	DeclEntry {},
	DeclType {},
	DeclFunction {},
	FunctionSignature {},
	Inputs {},
	DeclMethod {},
	DeclConst {},
	VariousBody {},
	NativeBody {},
	Doc {},
	Ref {},
	RefBase {},
	Identifier {},
	Type {},
	VariousTypeDef {},
	NativeTypeDef {},
	Interface {},
	Method {},
	Record {},
	RecordDef {},
	Field {},
	Union {},
	Enum {},
	EnumItem {},
	Expr {},
	Cast {},
	VariousTerm {},
	VariousPipe {},
	PipeCast {},
	PipeGet {},
	PipeInterior {},
	InfixTerm {},
	PipeInfix {},
	VariousPipeCall {},
	CallOrdered {},
	CallUnordered {},
	ArgumentMapping {},
	Lambda {},
	Block {},
	VariousBinding {},
	BindingPlain {},
	BindingCps {},
	VariousPattern {},
	PatternSingle {},
	PatternMultiple {},
	RefTerm {},
	New {},
	If {},
	ElIf {},
	Cond {},
	When {},
	Case {},
	Each {},
	String {},
	StringPart {},
	VariousStringPartContent {},
	Text {},
	Char {},
	Bytes {},
	Byte {},
	Int {},
	Float {},
	VariousReplCmd {},
	ReplAssign {},
	ReplRun {},
	ReplEval {},
}

type Root struct {
	Node                              `part:"root"`
	Namespace   MaybeIdentifier       `part?:"ns.name?.name"`
	Aliases     [] Alias              `list:"alias*"`
	Statements  [] VariousStatement   `list:"stmt*"`
}
type Alias struct {
	Node                         `part:"alias"`
	Off     bool                 `option:"off.#"`
	Name    MaybeIdentifier      `part?:"alias_name.name"`
	Target  VariousAliasTarget   `part:"alias_target"`
}
type AliasTarget interface { impl(AliasTarget) }
type VariousAliasTarget struct {
	Node                       `part:"alias_target"`
	AliasTarget  AliasTarget   `use:"first"`
}
func (AliasToNamespace) impl(AliasTarget) {}
type AliasToNamespace struct {
	Node                    `part:"alias_to_ns"`
	Namespace  Identifier   `part:"name"`
}
func (AliasToRefBase) impl(AliasTarget) {}
type AliasToRefBase struct {
	Node               `part:"alias_to_ref_base"`
	RefBase  RefBase   `part:"ref_base"`
}
type Statement interface { impl(Statement) }
type VariousStatement struct {
	Node                   `part:"stmt"`
	Statement  Statement   `use:"first"`
}
// func (DeclAsset) impl(Statement) {}
// type DeclAsset struct {
//     Node
//     Name     string
//     Content  LoadedAsset
// }
// func (LoadedAsset) implBody() {}
// type LoadedAsset struct {
//     Path  string
//     Data  [] byte
// }
func (DeclEntry) impl(Statement) {}
type DeclEntry struct {
	Node              `part:"decl_entry"`
	Docs     [] Doc   `list:"docs.doc+"`
	Off      bool     `option:"off.#"`
	Content  Block    `part:"block"`
}
func (DeclType) impl(Statement) {}
type DeclType struct {
	Node                         `part:"decl_type"`
	Docs        [] Doc           `list:"docs.doc+"`
	Off         bool             `option:"off.#"`
	Name        Identifier       `part:"name"`
	TypeParams  [] Identifier    `list:"type_params.name*,"`
	Implements  [] RefBase       `list:"impl.ref_base*,"`
	TypeDef     VariousTypeDef   `part:"type_def"`
}
func (DeclFunction) impl(Statement) {}
type DeclFunction struct {
	Node                           `part:"decl_func"`
	Docs       [] Doc              `list:"docs.doc+"`
	Off        bool                `option:"off.#"`
	Operator   bool                `option:"function.@operator"`
	Variadic   bool                `option:"variadic.@variadic"`
	Name       Identifier          `part:"name"`
	Signature  FunctionSignature   `part:"sig"`
	Body       VariousBody         `part:"body"`
}
type FunctionSignature struct {
	Node                        `part:"sig"`
	TypeParams  [] Identifier   `list:"type_params.name*,"`
	Inputs      Inputs          `part:"inputs"`
	Implicit    MaybeInputs     `part?:"implicit.inputs"`
	Output      Type            `part:"output.type"`
}
type MaybeInputs interface { impl(MaybeInputs) }
func (Inputs) impl(MaybeInputs) {}
type Inputs struct {
	Node                 `part:"inputs"`
	Content  RecordDef   `part:"record_def"`
}
func (DeclMethod) impl(Statement) {}
type DeclMethod struct {
	Node                    `part:"decl_method"`
	Docs      [] Doc        `list:"docs.doc+"`
	Off       bool          `option:"off.#"`
	Receiver  Identifier    `part:"receiver.name"`
	Name      Identifier    `part:"name"`
	Type      Type          `part:"type"`
	Body      VariousBody   `part:"body"`
}
func (DeclConst) impl(Statement) {}
type DeclConst struct {
	Node                `part:"decl_const"`
	Docs  [] Doc        `list:"docs.doc+"`
	Off   bool          `option:"off.#"`
	Name  Identifier    `part:"name"`
	Type  Type          `part:"type"`
	Body  VariousBody   `part:"body"`
}
type Body interface { implBody() }
type VariousBody struct {
	Node         `part:"body"`
	Body  Body   `use:"first"`
}
func (NativeBody) implBody() {}
type NativeBody struct {
	Node         `part:"native_body"`
	Id    Text   `part:"text"`
}
type Doc struct {
	Node               `part:"doc"`
	RawContent  Code   `content:"Doc"`
}

type Ref struct {
	Node                `part:"ref"`
	Base      RefBase   `part:"ref_base"`
	TypeArgs  [] Type   `list:"type_args.type*,"`
}
type RefBase struct {
	Node                    `part:"ref_base"`
	NS    MaybeIdentifier   `part?:"ns_prefix.name"`
	Item  Identifier        `part:"name"`
}
type MaybeIdentifier interface { impl(MaybeIdentifier) }
func (Identifier) impl(MaybeIdentifier) {}
type Identifier struct {
	Node         `part:"name"`
	Name  Code   `content:"Name"`
}

type Type struct {
	Node        `part:"type"`
	Ref   Ref   `part:"ref"`
}
type TypeDef interface { impl(TypeDef) }
type VariousTypeDef struct {
	Node               `part:"type_def"`
	TypeDef  TypeDef   `use:"first"`
}
func (NativeTypeDef) impl(TypeDef) {}
type NativeTypeDef struct {
	Node `part:"native_type_def"`
}
func (Interface) impl(TypeDef) {}
type Interface struct {
	Node                 `part:"interface"`
	Methods  [] Method   `list:"method*,"`
}
type Method struct {
	Node               `part:"method"`
	Docs  [] Doc       `list:"docs.doc+"`
	Name  Identifier   `part:"name"`
	Type  Type         `part:"type"`
}
func (Record) impl(TypeDef) {}
type Record struct {
	Node              `part:"record"`
	Def   RecordDef   `part:"record_def"`
}
type RecordDef struct {
	Node               `part:"record_def"`
	Fields  [] Field   `list:"field*,"`
}
type Field struct {
	Node                  `part:"field"`
	Docs     [] Doc       `list:"docs.doc+"`
	Name     Identifier   `part:"name"`
	Type     Type         `part:"type"`
	Default  MaybeExpr    `part?:"field_default.expr"`
}
func (Union) impl(TypeDef) {}
type Union struct {
	Node             `part:"union"`
	Items  [] Type   `list:"type+,"`
}
func (Enum) impl(TypeDef) {}
type Enum struct {
	Node                 `part:"enum"`
	Items  [] EnumItem   `list:"enum_item+,"`
}
type EnumItem struct {
	Node               `part:"enum_item"`
	Docs  [] Doc       `list:"docs.doc+"`
	Name  Identifier   `part:"name"`
}

type MaybeExpr interface { impl(MaybeExpr) }
func (Expr) impl(MaybeExpr)  {}
type Expr struct {
	Node                       `part:"expr"`
	Casts     [] Cast          `list:"cast*"`
	Term      VariousTerm      `part:"term"`
	Pipeline  [] VariousPipe   `list:"pipe*"`
}
type Cast struct {
	Node           `part:"cast"`
	Target  Type   `part:"type"`
}
type Term interface { impl(Term) }
type VariousTerm struct {
	Node         `part:"term"`
	Term  Term   `use:"first"`
}
type Pipe interface { impl(Pipe) }
type VariousPipe struct {
	Node         `part:"pipe"`
	Pipe  Pipe   `use:"first"`
}

func (PipeCast) impl(Pipe) {}
type PipeCast struct {
	Node         `part:"pipe_cast"`
	Cast  Cast   `part:"cast"`
}
func (PipeGet) impl(Pipe) {}
type PipeGet struct {
	Node               `part:"pipe_get"`
	Key   Identifier   `part:"name"`
}
func (PipeInterior) impl(Pipe) {}
type PipeInterior struct {
	Node               `part:"pipe_interior"`
	RefBase  RefBase   `part:"ref_base"`
}

func (InfixTerm) impl(Term) {}
type InfixTerm struct {
	Node             `part:"infix_term"`
	Operator  Ref    `part:"operator.ref"`
	Left      Expr   `part:"infix_left.expr"`
	Right     Expr   `part:"infix_right.expr"`
}
func (PipeInfix) impl(Pipe) {}
type PipeInfix struct {
	Node                        `part:"pipe_infix"`
	Off       bool              `option:"off.#"`
	Callee    Ref               `part:"ref"`
	PipeCall  VariousPipeCall   `part:"pipe_call"`
}
type PipeCall interface { impl(PipeCall) }
func (VariousPipeCall) impl(Pipe) {}
type VariousPipeCall struct {
	Node                 `part:"pipe_call"`
	PipeCall  PipeCall   `use:"first"`
}
func (CallOrdered) impl(PipeCall) {}
type CallOrdered struct {
	Node                 `part:"call_ordered"`
	Arguments  [] Expr   `list:"expr*,"`
}
func (CallUnordered) impl(PipeCall) {}
type CallUnordered struct {
	Node                           `part:"call_unordered"`
	Mappings  [] ArgumentMapping   `list:"arg_mapping*,"`
}
type ArgumentMapping struct {
	Node                `part:"arg_mapping"`
	Name   Identifier   `part:"name"`
	Value  MaybeExpr    `part?:"arg_mapping_to.expr"`
}

func (Lambda) impl(Term) {}
type Lambda struct {
	Node                            `part:"lambda"`
	InputPattern  MaybePattern      `part?:"pattern?.pattern"`
	OutputExpr    Expr              `part:"expr"`
	SelfRefName   MaybeIdentifier   `part?:"lambda_self.name"`
}
func (Block) impl(Term) {}
func (Block) implBody() {}
type Block struct {
	Node                          `part:"block"`
	Bindings  [] VariousBinding   `list:"binding*"`
	Return    Expr                `part:"expr"`
}
type Binding interface { impl(Binding) }
type VariousBinding struct {
	Node               `part:"binding"`
	Binding  Binding   `use:"first"`
}
func (BindingPlain) impl(Binding) {}
type BindingPlain struct {
	Node                      `part:"binding_plain"`
	Off      bool             `option:"off.#"`
	Const    bool             `option:"let.Const"`
	Pattern  VariousPattern   `part:"pattern"`
	Value    Expr             `part:"expr"`
}
func (BindingCps) impl(Binding) {}
type BindingCps struct {
	Node                    `part:"binding_cps"`
	Off      bool           `option:"off.#"`
	Callee   Ref            `part:"ref"`
	Pattern  MaybePattern   `part?:"cps_pattern.pattern"`
	Value    Expr           `part:"expr"`
}
type Pattern interface { impl(Pattern) }
type MaybePattern interface { impl(MaybePattern) }
func (VariousPattern) impl(MaybePattern) {}
type VariousPattern struct {
	Node               `part:"pattern"`
	Pattern  Pattern   `use:"first"`
}
func (PatternSingle) impl(Pattern) {}
type PatternSingle struct {
	Node                `part:"pattern_single"`
	Name   Identifier   `part:"name"`
}
func (PatternMultiple) impl(Pattern) {}
type PatternMultiple struct {
	Node                   `part:"pattern_multiple"`
	Names  [] Identifier   `list:"name+,"`
}

func (RefTerm) impl(Term) {}
type RefTerm struct {
	Node            `part:"ref_term"`
	New  MaybeNew   `part?:"new"`
	Ref  Ref        `part:"ref"`
}
type MaybeNew interface { impl(MaybeNew) }
func (New) impl(MaybeNew) {}
type New struct {
	Node                    `part:"new"`
	Tag   MaybeIdentifier   `part?:"new_tag.name"`
}
func (ImplicitRefTerm) impl(Term) {}
type ImplicitRefTerm struct {
	Node
	Name  Identifier
}

func (If) impl(Term) {}
type If struct {
	Node             `part:"if"`
	Conds  [] Cond   `list:"cond+,"`
	Yes    Block     `part:"if_yes.block"`
	No     Block     `part:"if_no.block"`
	ElIfs  [] ElIf   `list:"elif*"`
}
type ElIf struct {
	Node             `part:"elif"`
	Conds  [] Cond   `list:"cond+,"`
	Yes    Block     `part:"block"`
}
type Cond struct {
	Node                    `part:"cond"`
	Expr     Expr           `part:"expr"`
	Pattern  MaybePattern   `part?:"cond_pattern.pattern"`
}
func (When) impl(Term) {}
type When struct {
	Node               `part:"when"`
	Operand  Expr      `part:"expr"`
	Cases    [] Case   `list:"case+,"`
}
type Case struct {
	Node                          `part:"case"`
	Off           bool            `option:"off.#"`
	Names         [] Identifier   `list:"name+_bar"`
	InputPattern  MaybePattern    `part?:"pattern?.pattern"`
	OutputExpr    Expr            `part:"expr"`
}
func (Each) impl(Term) {}
type Each struct {
	Node               `part:"each"`
	Operand  Type      `part:"type"`
	Cases    [] Case   `list:"case+,"`
}

func (String) impl(Term) {}
type String struct {
	Node                   `part:"string"`
	First  Text            `part:"text"`
	Parts  [] StringPart   `list:"string_part*"`
}
type StringPart struct {
	Node                                `part:"string_part"`
	Content  VariousStringPartContent   `part:"string_part_content"`
}
type StringPartContent interface { implStringPartContent() }
type VariousStringPartContent struct {
	Node                                   `part:"string_part_content"`
	StringPartContent  StringPartContent   `use:"first"`
}
func (Text) implStringPartContent() {}
type Text struct {
	Node          `part:"text"`
	Value  Code   `content:"Text"`
}
func (Char) impl(Term) {}
func (Char) implStringPartContent() {}
type Char struct {
	Node          `part:"char"`
	Value  Code   `content:"Char"`
}
func (Bytes) impl(Term) {}
type Bytes struct {
	Node             `part:"bytes"`
	Bytes  [] Byte   `list:"byte+"`
}
type Byte struct {
	Node          `part:"byte"`
	Value  Code   `content:"Byte"`
}
func (Int) impl(Term) {}
type Int struct {
	Node          `part:"int"`
	Value  Code   `content:"Int"`
}
func (Float) impl(Term) {}
type Float struct {
	Node          `part:"float"`
	Value  Code   `content:"Float"`
}

type VariousReplCmd struct {
	Node               `part:"repl_cmd"`
	ReplCmd  ReplCmd   `use:"first"`
}
type ReplCmd interface { impl(ReplCmd) }
func (ReplAssign) impl(ReplCmd) {}
type ReplAssign struct {
	Node               `part:"repl_assign"`
	Name  Identifier   `part:"name"`
	Expr  Expr         `part:"expr"`
}
func (ReplRun) impl(ReplCmd) {}
type ReplRun struct {
	Node         `part:"repl_run"`
	Expr  Expr   `part:"expr"`
}
func (ReplEval) impl(ReplCmd) {}
type ReplEval struct {
	Node         `part:"repl_eval"`
	Expr  Expr   `part:"expr"`
}


