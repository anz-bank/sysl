package relmod

type App struct {
	AppName      []string
	AppLongName  string
	AppDocstring string
}

type Mixin struct {
	AppName   []string
	MixinName []string
}

type Endpoint struct {
	AppName     []string
	EpName      string
	EpLongName  string
	EpDocstring string
	EpEvent     EndpointEvent `arrai:",zeroempty"`
	Rest        Rest          `arrai:",zeroempty"`
}

type EndpointEvent struct {
	AppName   EndpointEventAppName
	EventName string
}
type EndpointEventAppName struct {
	Part []string
}

type QueryParam struct {
	Name string
	Type interface{} // Type*
}

// Backwards compatible.
// TODO: Design new structure.
type Rest struct {
	Method     string // GET, PUT, POST, DELETE, PATCH
	Path       string
	QueryParam []QueryParam `arrai:"query_param,omitempty"`
	URLParam   []QueryParam `arrai:"url_param,omitempty"`
}

type Param struct {
	AppName    []string
	EpName     string
	ParamName  string
	ParamLoc   string // enum
	ParamIndex int
	ParamType  interface{} // Type*
	ParamOpt   bool
}

type Statement struct {
	AppName     []string
	EpName      string
	StmtIndex   []int
	StmtParent  StatementParent        `arrai:",zeroempty"`
	StmtAction  string                 `arrai:",zeroempty"`
	StmtCall    map[string]interface{} `arrai:",zeroempty"`
	StmtCond    map[string]interface{} `arrai:",zeroempty"`
	StmtLoop    map[string]interface{} `arrai:",zeroempty"`
	StmtLoopN   map[string]interface{} `arrai:",zeroempty"`
	StmtForeach map[string]interface{} `arrai:",zeroempty"`
	StmtAlt     map[string]interface{} `arrai:",zeroempty"`
	StmtGroup   map[string]interface{} `arrai:",zeroempty"`
	StmtRet     StatementReturn        `arrai:",zeroempty"`
}

type StatementParent struct {
	AppName   []string
	EpName    string
	StmtIndex int
}

type StatementReturn struct {
	Status string
	Attr   StatementReturnAttrs
	Type   interface{} // Type*
}

type StatementReturnAttrs struct {
	Modifier []string
	Nvp      map[string]interface{}
}

type Type struct {
	AppName       []string
	TypeName      string
	TypeDocstring string
	TypeOpt       bool
}

type TypePrimitive struct {
	Primitive string
}

type TypeTuple struct {
	Tuple interface{} // TODO
}

type TypeRef struct {
	AppName  []string
	TypePath []string
	// TODO: Maybe add an optional fieldName field for field refs
}

type TypeSet struct {
	Set interface{} // Type*
}

type TypeSequence struct {
	Sequence interface{} // Type*
}

type Table struct {
	AppName  []string
	TypeName string
	Pk       []string
}

type Enum struct {
	AppName   []string
	EnumItems map[string]int64
	TypeName  string
}

type Alias struct {
	AppName   []string
	TypeName  string
	AliasType interface{} // Type*
}

type Event struct {
	AppName   []string
	EventName string
}

type Field struct {
	AppName         []string
	TypeName        string
	FieldName       string
	FieldOpt        bool
	FieldType       interface{}     // Type*
	FieldConstraint FieldConstraint `arrai:",zeroempty"`
}

type FieldConstraint struct {
	Length    FieldConstraintLength
	Precision int32
	Scale     int32
}

type FieldConstraintLength struct {
	Min int64
	Max int64
}

type View struct {
	AppName  []string
	ViewName string
	ViewType interface{} // Type*
}

type SourceContext struct {
	File  string
	Start Location
	End   Location
}

type Location struct {
	Line int
	Col  int
}

// ANNOTATIONS

type AppAnnotation struct {
	AppName      []string
	AppAnnoName  string
	AppAnnoValue interface{} // string, number, array of self
}

type MixinAnnotation struct {
	AppName        []string
	MixinName      []string
	MixinAnnoName  string
	MixinAnnoValue interface{} // string, number, array of self
}

type EndpointAnnotation struct {
	AppName     []string
	EpName      string
	EpAnnoName  string
	EpAnnoValue interface{} // string, number, array of self
}

type ParamAnnotation struct {
	AppName        []string
	EpName         string
	ParamName      string
	ParamLoc       string // enum
	ParamIndex     int
	ParamAnnoName  string
	ParamAnnoValue interface{} // string, number, array of self
}

type StatementAnnotation struct {
	AppName       []string
	EpName        string
	StmtAnnoName  string
	StmtAnnoValue interface{} // string, number, array of self
	StmtIndex     []int
}

type EventAnnotation struct {
	AppName        []string
	EventName      string
	EventAnnoName  string
	EventAnnoValue interface{} // string, number, array of self
}

type TypeAnnotation struct {
	AppName       []string
	TypeName      string
	TypeAnnoName  string
	TypeAnnoValue interface{} // string, number, array of self
}

type FieldAnnotation struct {
	AppName        []string
	TypeName       string
	FieldName      string
	FieldAnnoName  string
	FieldAnnoValue interface{} // string, number, array of self
}

type ViewAnnotation struct {
	AppName       []string
	ViewName      string
	ViewAnnoName  string
	ViewAnnoValue interface{} // string, number, array of self
}

// TAGS
type AppTag struct {
	AppName []string
	AppTag  string
}

type MixinTag struct {
	AppName   []string
	MixinName []string
	MixinTag  string
}

type EndpointTag struct {
	AppName []string
	EpName  string
	EpTag   string
}

type ParamTag struct {
	AppName    []string
	EpName     string
	ParamName  string
	ParamLoc   string // enum
	ParamIndex int
	ParamTag   string
}

type StatementTag struct {
	AppName   []string
	EpName    string
	StmtIndex []int
	StmtTag   string
}

type EventTag struct {
	AppName   []string
	EventName string
	EventTag  string
}

type TypeTag struct {
	AppName  []string
	TypeName string
	TypeTag  string
}

type FieldTag struct {
	AppName   []string
	TypeName  string
	FieldName string
	FieldTag  string
}

type ViewTag struct {
	AppName  []string
	ViewName string
	ViewTag  string
}

// SOURCE CONTEXTS

type AppContext struct {
	AppName []string
	AppSrc  SourceContext
	AppSrcs []SourceContext
}

type MixinContext struct {
	AppName   []string
	MixinName []string
	MixinSrc  SourceContext
	MixinSrcs []SourceContext
}

type EndpointContext struct {
	AppName []string
	EpName  string
	EpSrc   SourceContext
	EpSrcs  []SourceContext
}

type ParamContext struct {
	AppName    []string
	EpName     string
	ParamName  string
	ParamLoc   string // enum
	ParamIndex int
	ParamSrc   SourceContext
	ParamSrcs  []SourceContext
}

type StatementContext struct {
	AppName   []string
	EpName    string
	StmtIndex []int
	StmtSrc   SourceContext
	StmtSrcs  []SourceContext
}

type EventContext struct {
	AppName   []string
	EventName string
	EventSrc  SourceContext
	EventSrcs []SourceContext
}

type TypeContext struct {
	AppName  []string
	TypeName string
	TypeSrc  SourceContext
	TypeSrcs []SourceContext
}

type FieldContext struct {
	AppName   []string
	TypeName  string
	FieldName string
	FieldSrc  SourceContext
	FieldSrcs []SourceContext
}

type ViewContext struct {
	AppName  []string
	ViewName string
	ViewSrc  SourceContext
	ViewSrcs []SourceContext
}

// SCHEMA

type Annotations struct {
	App   []AppAnnotation       `arrai:",unordered"`
	Mixin []MixinAnnotation     `arrai:",unordered"`
	Ep    []EndpointAnnotation  `arrai:",unordered"`
	Param []ParamAnnotation     `arrai:",unordered"`
	Stmt  []StatementAnnotation `arrai:",unordered"`
	Event []EventAnnotation     `arrai:",unordered"`
	Type  []TypeAnnotation      `arrai:",unordered"`
	Field []FieldAnnotation     `arrai:",unordered"`
	View  []ViewAnnotation      `arrai:",unordered"`
}

type Tags struct {
	App   []AppTag       `arrai:",unordered"`
	Mixin []MixinTag     `arrai:",unordered"`
	Ep    []EndpointTag  `arrai:",unordered"`
	Param []ParamTag     `arrai:",unordered"`
	Stmt  []StatementTag `arrai:",unordered"`
	Event []EventTag     `arrai:",unordered"`
	Type  []TypeTag      `arrai:",unordered"`
	Field []FieldTag     `arrai:",unordered"`
	View  []ViewTag      `arrai:",unordered"`
}

type SourceContexts struct {
	App   []AppContext       `arrai:",unordered"`
	Mixin []MixinContext     `arrai:",unordered"`
	Ep    []EndpointContext  `arrai:",unordered"`
	Param []ParamContext     `arrai:",unordered"`
	Stmt  []StatementContext `arrai:",unordered"`
	Event []EventContext     `arrai:",unordered"`
	Type  []TypeContext      `arrai:",unordered"`
	Field []FieldContext     `arrai:",unordered"`
	View  []ViewContext      `arrai:",unordered"`
}

type Schema struct {
	App   []App       `arrai:",unordered"`
	Mixin []Mixin     `arrai:",unordered"`
	Ep    []Endpoint  `arrai:",unordered"`
	Param []Param     `arrai:",unordered"`
	Stmt  []Statement `arrai:",unordered"`
	Event []Event     `arrai:",unordered"`
	Type  []Type      `arrai:",unordered"`
	Field []Field     `arrai:",unordered"`
	Table []Table     `arrai:",unordered"`
	Alias []Alias     `arrai:",unordered"`
	Enum  []Enum      `arrai:",unordered"`
	View  []View      `arrai:",unordered"`
	Anno  Annotations
	Tag   Tags
	Src   SourceContexts
}
