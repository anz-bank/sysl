Ù
AÓ

A*”
Fetch…
Fetch:

SourceRead:êBç
äok <: A.FetchResponse [dataflow={"A.FetchResponse.ax": "Source.Foo.x", "A.FetchResponse.ay": "Source.Foo.y"}, description="1:1 transform"]ö

model.sysl2x
FetchResponsegö

model.sysl!J
#
axö

model.sysl
#
ayö

model.syslö

model.syslÙ
BÓ

B*”
Fetch…
Fetch:

SourceRead:êBç
äok <: B.FetchResponse [dataflow={"B.FetchResponse.bx": "Source.Foo.x", "B.FetchResponse.by": "Source.Foo.y"}, description="1:1 transform"]ö

model.sysl"&2x
FetchResponsegö

model.sysl&*J
#
byö

model.sysl((
#
bxö

model.sysl''ö

model.sysl!!˛
C¯

C*ù
Fetchì
Fetch:

CFetchA:

CFetchB:ÃB…
∆ok <: C.FetchResponse [dataflow={"C.FetchResponse.cx": ["A.FetchResponse.ax", "B.FetchResponse.bx"], "C.FetchResponse.cy": ["A.FetchResponse.ay", "B.FetchResponse.by"]}, description="1:1 transform"]ö

model.sysl38*›
FetchA“
FetchA:

AFetch:úBô
ñok <: C.FetchResponse [dataflow={"C.FetchResponse.cx": "A.FetchResponse.ax", "C.FetchResponse.cy": "A.FetchResponse.ay"}, description="1:1 transform"]ö

model.sysl+/*›
FetchB“
FetchB:

BFetch:úBô
ñok <: C.FetchResponse [dataflow={"C.FetchResponse.cx": "B.FetchResponse.bx", "C.FetchResponse.cy": "B.FetchResponse.by"}, description="1:1 transform"]ö

model.sysl/32x
FetchResponsegö

model.sysl8<J
#
cxö

model.sysl99
#
cyö

model.sysl::ö

model.sysl**º
D∂

D*õ
Fetchë
Fetch:

AFetch:

CFetch:ÃB…
∆ok <: D.FetchResponse [dataflow={"D.FetchResponse.dx": ["A.FetchResponse.ax", "C.FetchResponse.cx"], "D.FetchResponse.dy": ["A.FetchResponse.ay", "C.FetchResponse.cy"]}, description="1:1 transform"]ö

model.sysl=B2x
FetchResponsegö

model.syslBFJ
#
dxö

model.syslCC
#
dyö

model.syslDDö

model.sysl<<Û
ClientË

Client*œ
Do»
Do:

DFetch:ñBì
êok <: Client.Screen [dataflow={"Client.Screen.xx": "D.FetchResponse.dx", "Client.Screen.yy": "D.FetchResponse.dy"}, description="1:1 transform"]ö

model.syslGK2q
Screengö

model.syslKOJ
#
xxö

model.syslLL
#
yyö

model.syslMMö

model.syslFFÇ
all{

all*Y
allR
all:

Source:

Client:
A:
B:
C:
Dö

model.syslPWö

model.syslOOı
SourceÍ

Source"
patterns:
"db"+
description"A database.
 Stores data.
*?
Read7
Read:B
ok <: Source.Fooö

model.sysl2Ò
FooÈB0
description!"A Foo.
 Represents foo things.
ö

model.sysl:ô
A
x<B
description"The x value.ö

model.sysl

T
yOB0
description!"A Foo.
 Represents foo things.
ö

model.sysl2 
Bar¬B
description"A bar table.ö

model.sysl:Ö
W
aRB
patterns:
"pkB
description"A bar table.ö

model.sysl
F
bAB 
description"An optional int`ö

model.sysl
]
xXB
description"A foreign keyö

model.syslJ


SourceBarFoox
aö

model.sysl
