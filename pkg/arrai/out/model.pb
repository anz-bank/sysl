þ
Cø

C*Ý
FetchAÒ
FetchA:

AFetch:œB™
–ok <: C.FetchResponse [dataflow={"C.FetchResponse.cx": "A.FetchResponse.ax", "C.FetchResponse.cy": "A.FetchResponse.ay"}, description="1:1 transform"]š

model.sysl+/*Ý
FetchBÒ
FetchB:

BFetch:œB™
–ok <: C.FetchResponse [dataflow={"C.FetchResponse.cx": "B.FetchResponse.bx", "C.FetchResponse.cy": "B.FetchResponse.by"}, description="1:1 transform"]š

model.sysl/3*
Fetch“
Fetch:

CFetchA:

CFetchB:ÌBÉ
Æok <: C.FetchResponse [dataflow={"C.FetchResponse.cx": ["A.FetchResponse.ax", "B.FetchResponse.bx"], "C.FetchResponse.cy": ["A.FetchResponse.ay", "B.FetchResponse.by"]}, description="1:1 transform"]š

model.sysl382x
FetchResponsegš

model.sysl8<J
#
cxš

model.sysl99
#
cyš

model.sysl::š

model.sysl**¼
D¶

D*›
Fetch‘
Fetch:

AFetch:

CFetch:ÌBÉ
Æok <: D.FetchResponse [dataflow={"D.FetchResponse.dx": ["A.FetchResponse.ax", "C.FetchResponse.cx"], "D.FetchResponse.dy": ["A.FetchResponse.ay", "C.FetchResponse.cy"]}, description="1:1 transform"]š

model.sysl=B2x
FetchResponsegš

model.syslBFJ
#
dxš

model.syslCC
#
dyš

model.syslDDš

model.sysl<<ó
Clientè

Client*Ï
DoÈ
Do:

DFetch:–B“
ok <: Client.Screen [dataflow={"Client.Screen.xx": "D.FetchResponse.dx", "Client.Screen.yy": "D.FetchResponse.dy"}, description="1:1 transform"]š

model.syslGK2q
Screengš

model.syslKOJ
#
xxš

model.syslLL
#
yyš

model.syslMMš

model.syslFF‚
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
Dš

model.syslPWš

model.syslOO¡
Source–

Source"
patterns:
"db"+
description"A database.
 Stores data.
*k
Readc
Read:@B>
<ok <: Source.Foo [dataflow=["Source.Foo.x", "Source.Bar.a"]]š

model.sysl2ñ
FooéB0
description!"A Foo.
 Represents foo things.
š

model.sysl:™
A
x<B
description"The x value.š

model.sysl

T
yOB0
description!"A Foo.
 Represents foo things.
š

model.sysl2Ê
BarÂB
description"A bar table.š

model.sysl:…
W
aRB
patterns:
"pkB
description"A bar table.š

model.sysl
F
bAB 
description"An optional int`š

model.sysl
]
xXB
description"A foreign keyš

model.syslJ


SourceBarFoox
aš

model.sysl
ô
Aî

A*Ó
FetchÉ
Fetch:

SourceRead:B
Šok <: A.FetchResponse [dataflow={"A.FetchResponse.ax": "Source.Foo.x", "A.FetchResponse.ay": "Source.Foo.y"}, description="1:1 transform"]š

model.sysl2x
FetchResponsegš

model.sysl!J
#
ayš

model.sysl
#
axš

model.syslš

model.syslô
Bî

B*Ó
FetchÉ
Fetch:

SourceRead:B
Šok <: B.FetchResponse [dataflow={"B.FetchResponse.bx": "Source.Foo.x", "B.FetchResponse.by": "Source.Foo.y"}, description="1:1 transform"]š

model.sysl"&2x
FetchResponsegš

model.sysl&*J
#
bxš

model.sysl''
#
byš

model.sysl((š

model.sysl!!