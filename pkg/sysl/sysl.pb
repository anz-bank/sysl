
�V
pkg/sysl/sysl.protosysl"�
SourceContext
file (	Rfile2
start (2.sysl.SourceContext.LocationRstart.
end (2.sysl.SourceContext.LocationRend
text (	Rtext!
delta (2.sysl.DeltaRdelta
version (	Rversion0
Location
line (Rline
col (Rcol"�
Module*
apps (2.sysl.Module.AppsEntryRapps>
source_contextc (2.sysl.SourceContextBRsourceContext<
source_contextsd (2.sysl.SourceContextRsourceContextsJ
	AppsEntry
key (	Rkey'
value (2.sysl.ApplicationRvalue:8JJ"�
	Attribute
s (	H Rs
i (H Ri
n (H Rn%
a (2.sysl.Attribute.ArrayH Ra>
source_contextc (2.sysl.SourceContextBRsourceContext<
source_contextsd (2.sysl.SourceContextRsourceContexts*
Array!
elt (2.sysl.AttributeReltB
	attribute"
AppName
part (	Rpart"�
Application!
name (2.sysl.AppNameRname
	long_name (	RlongName
	docstring (	R	docstring2
attrs (2.sysl.Application.AttrsEntryRattrs>
	endpoints (2 .sysl.Application.EndpointsEntryR	endpoints2
types (2.sysl.Application.TypesEntryRtypes2
views
 (2.sysl.Application.ViewsEntryRviews)
mixin2 (2.sysl.ApplicationRmixin2+
wrapped	 (2.sysl.ApplicationRwrapped>
source_contextc (2.sysl.SourceContextBRsourceContext<
source_contextsd (2.sysl.SourceContextRsourceContexts4
DONOTUSE_mixin (2.sysl.AppNameRDONOTUSEMixinI

AttrsEntry
key (	Rkey%
value (2.sysl.AttributeRvalue:8L
EndpointsEntry
key (	Rkey$
value (2.sysl.EndpointRvalue:8D

TypesEntry
key (	Rkey 
value (2
.sysl.TypeRvalue:8D

ViewsEntry
key (	Rkey 
value (2
.sysl.ViewRvalue:8"�
Endpoint
name (	Rname
	long_name (	RlongName
	docstring (	R	docstring/
attrs (2.sysl.Endpoint.AttrsEntryRattrs
flag
 (	Rflag%
source (2.sysl.AppNameRsource
	is_pubsub (RisPubsub!
param	 (2.sysl.ParamRparam#
stmt (2.sysl.StatementRstmt:
rest_params (2.sysl.Endpoint.RestParamsR
restParams>
source_contextc (2.sysl.SourceContextBRsourceContext<
source_contextsd (2.sysl.SourceContextRsourceContextsI

AttrsEntry
key (	Rkey%
value (2.sysl.AttributeRvalue:8�

RestParams8
method (2 .sysl.Endpoint.RestParams.MethodRmethod
path (	RpathE
query_param (2$.sysl.Endpoint.RestParams.QueryParamR
queryParamA
	url_param (2$.sysl.Endpoint.RestParams.QueryParamRurlParamy

QueryParam
name (	Rname
type (2
.sysl.TypeRtype
loc (Rloc%
DONOTUSE_param (	RDONOTUSEParam"s
Method
	NO_Method 
GET
PUT
POST

DELETE	
PATCH
DONOTUSE_OPTIONS
DONOTUSE_HEAD";
Param
name (	Rname
type (2
.sysl.TypeRtype"�
	Statement&
action (2.sysl.ActionH Raction 
call (2
.sysl.CallH Rcall 
cond (2
.sysl.CondH Rcond 
loop (2
.sysl.LoopH Rloop$
loop_n (2.sysl.LoopNH RloopN)
foreach	 (2.sysl.ForeachH Rforeach
alt (2	.sysl.AltH Ralt#
group (2.sysl.GroupH Rgroup 
ret (2.sysl.ReturnH Rret0
attrs
 (2.sysl.Statement.AttrsEntryRattrs>
source_contextc (2.sysl.SourceContextBRsourceContext<
source_contextsd (2.sysl.SourceContextRsourceContextsI

AttrsEntry
key (	Rkey%
value (2.sysl.AttributeRvalue:8B
stmt" 
Action
action (	Raction"�
Call%
target (2.sysl.AppNameRtarget
endpoint (	Rendpoint 
arg (2.sysl.Call.ArgRargD
DONOTUSE_attrs (2.sysl.Call.DONOTUSEAttrsEntryRDONOTUSEAttrsm
Arg!
value (2.sysl.ValueRvalue
name (	Rname/
DONOTUSE_type (2
.sysl.TypeRDONOTUSETypeQ
DONOTUSEAttrsEntry
key (	Rkey%
value (2.sysl.AttributeRvalue:8"?
Cond
test (	Rtest#
stmt (2.sysl.StatementRstmt"�
Loop#
mode (2.sysl.Loop.ModeRmode
	criterion (	R	criterion#
stmt (2.sysl.StatementRstmt")
Mode
NO_Mode 	
WHILE	
UNTIL"B
LoopN
count (Rcount#
stmt (2.sysl.StatementRstmt"N
Foreach

collection (	R
collection#
stmt (2.sysl.StatementRstmt"r
Alt(
choice (2.sysl.Alt.ChoiceRchoiceA
Choice
cond (	Rcond#
stmt (2.sysl.StatementRstmt"B
Group
title (	Rtitle#
stmt (2.sysl.StatementRstmt""
Return
payload (	Rpayload"�
Type4
	primitive (2.sysl.Type.PrimitiveH R	primitive%
enum (2.sysl.Type.EnumH Renum(
tuple (2.sysl.Type.TupleH Rtuple%
list (2.sysl.Type.ListH Rlist"
map (2.sysl.Type.MapH Rmap)
one_of (2.sysl.Type.OneOfH RoneOf1
relation (2.sysl.Type.RelationH Rrelation,
type_ref	 (2.sysl.ScopedRefH RtypeRef
set (2
.sysl.TypeH Rset(
sequence (2
.sysl.TypeH Rsequence,
no_type (2.sysl.Type.NoTypeH RnoType+
attrs (2.sysl.Type.AttrsEntryRattrs5

constraint
 (2.sysl.Type.ConstraintR
constraint
	docstring (	R	docstring
opt (Ropt>
source_contextc (2.sysl.SourceContextBRsourceContext<
source_contextsd (2.sysl.SourceContextRsourceContextsI

AttrsEntry
key (	Rkey%
value (2.sysl.AttributeRvalue:8r
Enum0
items (2.sysl.Type.Enum.ItemsEntryRitems8

ItemsEntry
key (	Rkey
value (Rvalue:8�
Tuple;
	attr_defs (2.sysl.Type.Tuple.AttrDefsEntryRattrDefsG
FUTURE_fields (2".sysl.Type.Tuple.FUTUREFieldsEntryRFUTUREFieldsG
AttrDefsEntry
key (	Rkey 
value (2
.sysl.TypeRvalue:8W
FUTUREFieldsEntry
key (	Rkey,
value (2.sysl.Type.Tuple.FieldRvalue:8i
Field
type (2
.sysl.TypeRtype
min_repeats (R
minRepeats
max_repeats (R
maxRepeats&
List
type (2
.sysl.TypeRtypeE
Map
key (2
.sysl.TypeRkey 
value (2
.sysl.TypeRvalue'
OneOf
type (2
.sysl.TypeRtype�
Relation>
	attr_defs (2!.sysl.Type.Relation.AttrDefsEntryRattrDefs8
primary_key (2.sysl.Type.Relation.KeyR
primaryKey)
key (2.sysl.Type.Relation.KeyRkey
inject (	RinjectG
AttrDefsEntry
key (	Rkey 
value (2
.sysl.TypeRvalue:8"
Key
	attr_name (	RattrNamec
Foreign
app (2.sysl.AppNameRapp
relation (	Rrelation
	attr_name (	RattrName�

Constraint1
range (2.sysl.Type.Constraint.RangeRrange4
length (2.sysl.Type.Constraint.LengthRlength@

resolution (2 .sysl.Type.Constraint.ResolutionR
resolution
	precision (R	precision
scale (Rscale
	bit_width (RbitWidthE
Range
min (2.sysl.ValueRmin
max (2.sysl.ValueRmax,
Length
min (Rmin
max (Rmax6

Resolution
base (Rbase
index (Rindex
NoType"�
	Primitive
NO_Primitive 	
EMPTY
ANY
BOOL
INT	
FLOAT
DECIMAL

STRING	
BYTES
STRING_8
DATE	
DATETIME

XML
UUIDB
type"�
View!
param (2.sysl.ParamRparam%
ret_type (2
.sysl.TypeRretType
expr (2
.sysl.ExprRexpr+
views (2.sysl.View.ViewsEntryRviews+
attrs (2.sysl.View.AttrsEntryRattrs:
source_contextc (2.sysl.SourceContextRsourceContext<
source_contextsd (2.sysl.SourceContextRsourceContextsD

ViewsEntry
key (	Rkey 
value (2
.sysl.ViewRvalue:8I

AttrsEntry
key (	Rkey%
value (2.sysl.AttributeRvalue:8"�
Expr
name (	H Rname'
literal (2.sysl.ValueH Rliteral/
get_attr (2.sysl.Expr.GetAttrH RgetAttr4
	transform (2.sysl.Expr.TransformH R	transform+
ifelse (2.sysl.Expr.IfElseH Rifelse%
call (2.sysl.Expr.CallH Rcall+
unexpr (2.sysl.Expr.UnExprH Runexpr.
binexpr (2.sysl.Expr.BinExprH Rbinexpr.
relexpr (2.sysl.Expr.RelExprH Rrelexpr1
navigate	 (2.sysl.Expr.NavigateH Rnavigate%
list
 (2.sysl.Expr.ListH Rlist#
set (2.sysl.Expr.ListH Rset(
tuple (2.sysl.Expr.TupleH Rtuple
type (2
.sysl.TypeRtype>
source_contextc (2.sysl.SourceContextBRsourceContext<
source_contextsd (2.sysl.SourceContextRsourceContextsm
GetAttr
arg (2
.sysl.ExprRarg
attr (	Rattr
nullsafe (Rnullsafe
setof (Rsetof�
Navigate
arg (2
.sysl.ExprRarg
attr (	Rattr
nullsafe (Rnullsafe
setof (Rsetof
via (	Rvia&
List
expr (2
.sysl.ExprRexpr�
	Transform
arg (2
.sysl.ExprRarg
scopevar (	Rscopevar-
stmt (2.sysl.Expr.Transform.StmtRstmt
	all_attrs (RallAttrs!
except_attrs (	RexceptAttrs
nullsafe (Rnullsafe�
Stmt:
assign (2 .sysl.Expr.Transform.Stmt.AssignH Rassign4
let (2 .sysl.Expr.Transform.Stmt.AssignH Rlet$
inject (2
.sysl.ExprH RinjectR
Assign
name (	Rname
expr (2
.sysl.ExprRexpr
table (RtableB
stmt�
IfElse
cond (2
.sysl.ExprRcond#
if_true (2
.sysl.ExprRifTrue%
if_false (2
.sysl.ExprRifFalse
nullsafe (Rnullsafe8
Call
func (	Rfunc
arg (2
.sysl.ExprRarg�
UnExpr$
op (2.sysl.Expr.UnExpr.OpRop
arg (2
.sysl.ExprRarg"_
Op	
NO_Op 
NEG
POS
NOT
INV

SINGLE
SINGLE_OR_NULL

STRING�
BinExpr%
op (2.sysl.Expr.BinExpr.OpRop
lhs (2
.sysl.ExprRlhs
rhs (2
.sysl.ExprRrhs
scopevar (	Rscopevar
	attr_name (	RattrName"�
Op	
NO_Op 
EQ
NE
LT
LE
GT
GE
IN
CONTAINS

NOT_IN
NOT_CONTAINS
ADD
SUB
MUL	
DIV

MOD
POW
AND
OR

BUTNOT

BITAND	
BITOR

BITXOR
COALESCE	
WHERE
TO_MATCHING
TO_NOT_MATCHING
FLATTEN�
RelExpr%
op (2.sysl.Expr.RelExpr.OpRop"
target (2
.sysl.ExprRtarget
arg (2
.sysl.ExprRarg
scopevar (	Rscopevar

descending (R
descending
	attr_name (	RattrName"�
Op	
NO_Op 
MIN
MAX
SUM
AVERAGE
FUTURE_WHERE
FUTURE_FLATTEN
RANK
SNAPSHOT
FIRST_BY	�
Tuple1
attrs (2.sysl.Expr.Tuple.AttrsEntryRattrsD

AttrsEntry
key (	Rkey 
value (2
.sysl.ExprRvalue:8B
expr"�
Value
b (H Rb
i (H Ri
d (H Rd
s (	H Rs
decimal (	H Rdecimal
data (H Rdata
enum (H Renum&
list (2.sysl.Value.ListH Rlist#
map (2.sysl.Value.MapH Rmap$
set	 (2.sysl.Value.ListH Rset&
null
 (2.sysl.Value.NullH Rnull
uuid (H Ruuid)
List!
value (2.sysl.ValueRvalue~
Map0
items (2.sysl.Value.Map.ItemsEntryRitemsE

ItemsEntry
key (	Rkey!
value (2.sysl.ValueRvalue:8
NullB
value"Q
	ScopedRef%
context (2.sysl.ScopeRcontext
ref (2.sysl.ScopeRref"D
Scope'
appname (2.sysl.AppNameRappname
path (	Rpath*X
Delta
NO_Delta 

DELTA_SAME
DELTA_CHANGE
	DELTA_ADD
DELTA_REMOVEBZsyslbproto3