�'

Statements�'


Statements"[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag*�
OneOfStatements�
OneOfStatements"[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag:�� 
ts/test/all.sysl��&� 
ts/test/all.sysl��&2�
_
case1V� 
ts/test/all.sysl��#� 
ts/test/all.sysl��#B
ok <: string
d
case number 2S� 
ts/test/all.sysl�� � 
ts/test/all.sysl�� B
	ok <: int
f
"case 3"Z� 
ts/test/all.sysl��'� 
ts/test/all.sysl��'B
ok <: Types.Type
[Y� 
ts/test/all.sysl��&� 
ts/test/all.sysl��&B
error <: string� 
ts/test/all.sysl��&� 
ts/test/all.sysl��&*�
AnnotatedEndpoint�
AnnotatedEndpoint"}
annotation1n� 
ts/test/all.sysl��?� 
ts/test/all.sysl��?"&you can do string annotation like this"�
annotation2�� 
ts/test/all.sysl��2� 
ts/test/all.sysl��2:�
J� 
ts/test/all.sysl��� 
ts/test/all.sysl��"or
J� 
ts/test/all.sysl��"� 
ts/test/all.sysl��""in
J� 
ts/test/all.sysl�$�(� 
ts/test/all.sysl�$�("an
M� 
ts/test/all.sysl�*�1� 
ts/test/all.sysl�*�1"array"�
annotation3x� 
ts/test/all.sysl��� 
ts/test/all.sysl��"0you can also do
multiline annotations
like this
"[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag� 
ts/test/all.sysl��� 
ts/test/all.sysl��*�
Calls�
Calls"[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag:_� 
ts/test/all.sysl��� 
ts/test/all.sysl��


StatementsReturns:d� 
ts/test/all.sysl��"� 
ts/test/all.sysl��"

RestEndpoint
GET /param� 
ts/test/all.sysl��"� 
ts/test/all.sysl��"*�
Loops�
Loops"[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag:�� 
ts/test/all.sysl��	� 
ts/test/all.sysl��	"\	predicateM� 
ts/test/all.sysl��� 
ts/test/all.sysl��
...:�� 
ts/test/all.sysl��	� 
ts/test/all.sysl��	JZ
	predicateM� 
ts/test/all.sysl��� 
ts/test/all.sysl��
...:�� 
ts/test/all.sysl��� 
ts/test/all.sysl��"\	predicateM� 
ts/test/all.sysl��� 
ts/test/all.sysl��
...� 
ts/test/all.sysl��� 
ts/test/all.sysl��*�
Returns�
Returns"[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag:V� 
ts/test/all.sysl��� 
ts/test/all.sysl��B
ok <: string:Z� 
ts/test/all.sysl��� 
ts/test/all.sysl��B
ok <: Types.Type:]� 
ts/test/all.sysl��"� 
ts/test/all.sysl��"B
error <: Types.Type� 
ts/test/all.sysl��"� 
ts/test/all.sysl��"*�
GroupStatements�
GroupStatements"[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag:�� 
ts/test/all.sysl��� 
ts/test/all.sysl��:r
groupedg� 
ts/test/all.sysl��� 
ts/test/all.sysl��


StatementsGroupStatements� 
ts/test/all.sysl��� 
ts/test/all.sysl��*�
AnnotatedStatements�
AnnotatedStatements:e� 
ts/test/all.sysl��� 
ts/test/all.sysl��


StatementsMiscellaneous:�� 
ts/test/all.sysl��v� 
ts/test/all.sysl��vBi
gok <: string [annotation=["as", "an", "array"]] #Doesn't work, annos/tags/comments are part of the name:U� 
ts/test/all.sysl��� 
ts/test/all.sysl��
"statement"� 
ts/test/all.sysl��� 
ts/test/all.sysl��*�
Miscellaneous�
Miscellaneous:d� 
ts/test/all.sysl��"� 
ts/test/all.sysl��"
SimpleEndpoint -> SimpleEp� 
ts/test/all.sysl��"� 
ts/test/all.sysl��"*�
IfStmt�
IfStmt"[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag:�� 
ts/test/all.sysl��	� 
ts/test/all.sysl��	d

predicate1V� 
ts/test/all.sysl��� 
ts/test/all.sysl��B
ok <: string:�� 
ts/test/all.sysl��	� 
ts/test/all.sysl��	:t
else if predicate2^� 
ts/test/all.sysl��� 
ts/test/all.sysl��


StatementsIfStmt:�� 
ts/test/all.sysl��� 
ts/test/all.sysl��:U
elseM� 
ts/test/all.sysl��� 
ts/test/all.sysl��
...� 
ts/test/all.sysl��� 
ts/test/all.sysl���
ts/test/all.sysl��"�
ts/test/all.sysl��"l
ImportedApp]

ImportedApp*
...
...�
ts/test/imported.sysl �
ts/test/imported.sysl �
App�

App"\
patternsP:N
L�
ts/test/all.sysl�
ts/test/all.sysl"abstract*
...
...�
ts/test/all.sysl�
ts/test/all.sysl�
AppWithAnnotation�

AppWithAnnotation"W
patternsK:I
G�
ts/test/all.sysl�
ts/test/all.sysl"tag"\

annotationN�
ts/test/all.sysl�
ts/test/all.sysl"
annotation"{
annotation1l�
ts/test/all.sysl?�
ts/test/all.sysl?"(you can do "string" annotation like this"�
annotation2��
ts/test/all.sysl0�
ts/test/all.sysl0:�
F�
ts/test/all.sysl�
ts/test/all.sysl"or
F�
ts/test/all.sysl�
ts/test/all.sysl"in
F�
ts/test/all.sysl $�
ts/test/all.sysl $"an
��
ts/test/all.sysl&/�
ts/test/all.sysl&/:K
I�
ts/test/all.sysl'.�
ts/test/all.sysl'."array"�
annotation3u�
ts/test/all.sysl�
ts/test/all.sysl"1you can also do
multiline annotations

like this
�
ts/test/all.sysl�
ts/test/all.sysl�
SimpleEndpoint�

SimpleEndpoint"W
patternsK:I
G�
ts/test/all.sysl99�
ts/test/all.sysl99"tag*�
SimpleEp�
SimpleEp"{
annotation1l�
ts/test/all.sysl<<C�
ts/test/all.sysl<<C"(you can do "string" annotation like this"�
annotation2��
ts/test/all.sysl==4�
ts/test/all.sysl==4:�
F�
ts/test/all.sysl==�
ts/test/all.sysl=="or
F�
ts/test/all.sysl=="�
ts/test/all.sysl==""in
F�
ts/test/all.sysl=$=(�
ts/test/all.sysl=$=("an
��
ts/test/all.sysl=*=3�
ts/test/all.sysl=*=3:K
I�
ts/test/all.sysl=+=2�
ts/test/all.sysl=+=2"array"�
annotation3u�
ts/test/all.sysl>D�
ts/test/all.sysl>D"1you can also do
multiline annotations

like this
"_
patternsS:Q
O�
ts/test/all.sysl::�
ts/test/all.sysl::"SimpleEpTag"\

annotationN�
ts/test/all.sysl;;"�
ts/test/all.sysl;;""
annotation�
ts/test/all.sysl:D�
ts/test/all.sysl:D*�
SimpleEpWithParam�
SimpleEpWithParam"W
patternsK:I
G�
ts/test/all.syslD&D*�
ts/test/all.syslD&D*"tag:I�
ts/test/all.syslEE�
ts/test/all.syslEE
...J
untypedParamr �
ts/test/all.syslDE�
ts/test/all.syslDE*�
SimpleEpWithTypes�
SimpleEpWithTypes"W
patternsK:I
G�
ts/test/all.syslG*G.�
ts/test/all.syslG*G."tag:I�
ts/test/all.syslHH�
ts/test/all.syslHH
...J
native�
ts/test/all.syslGH�
ts/test/all.syslGH*�
SimpleEpWithArray�
SimpleEpWithArray"W
patternsK:I
G�
ts/test/all.syslJ[J_�
ts/test/all.syslJ[J_"tag:I�
ts/test/all.syslKK�
ts/test/all.syslKK
...J
	unlimitedRJ
limited
R
J
numR�
ts/test/all.syslJK�
ts/test/all.syslJK�
ts/test/all.sysl9K�
ts/test/all.sysl9K�6
Types�6

Types2�
Table�BW
patternsK:I
G�
ts/test/all.syslYY�
ts/test/all.syslYY"tag�
ts/test/all.syslYf&�
ts/test/all.syslYf&:�
x
int_with_bitwidthcR

���������
��������0@�
ts/test/all.syslee"�
ts/test/all.syslee"
_
float_with_bitwidthHR0@�
ts/test/all.syslff&�
ts/test/all.syslff&
{
	referencen�
ts/test/all.sysl\\&�
ts/test/all.sysl\\&J*


TypesTable

RestEndpointType
R
optionalF`�
ts/test/all.sysl]]�
ts/test/all.sysl]]
�
set��
ts/test/all.sysl^^�
ts/test/all.sysl^^jD�
ts/test/all.sysl^^�
ts/test/all.sysl^^
�
	with_anno�Bg

annotationY�
ts/test/all.syslaa1�
ts/test/all.syslaa1"this is an annotation�
ts/test/all.sysl`b	�
ts/test/all.sysl`b	
c
string_max_constraintJR�
ts/test/all.syslcc*�
ts/test/all.syslcc*
�

primaryKey�BV
patternsJ:H
F�
ts/test/all.syslZZ!�
ts/test/all.syslZZ!"pk�
ts/test/all.syslZZ"�
ts/test/all.syslZZ"
W
nativeTypeFieldD�
ts/test/all.sysl[[!�
ts/test/all.sysl[[!
�
sequence��
ts/test/all.sysl__&�
ts/test/all.sysl__&zD�
ts/test/all.sysl__&�
ts/test/all.sysl__&
h
decimal_with_precisionNR (�
ts/test/all.syslbb.�
ts/test/all.syslbb.
g
string_range_constraintLR
�
ts/test/all.sysldd0�
ts/test/all.sysldd0

primaryKey2�
Union�BW
patternsK:I
G�
ts/test/all.syslmm�
ts/test/all.syslmm"tag�
ts/test/all.syslms�
ts/test/all.syslms2�
D�
ts/test/all.syslnn�
ts/test/all.syslnn
D�
ts/test/all.sysloo�
ts/test/all.sysloo
��
ts/test/all.syslpp �
ts/test/all.syslpp zNR (�
ts/test/all.syslpp �
ts/test/all.syslpp 
n�
ts/test/all.syslqq�
ts/test/all.syslqqJ*


TypesUnion

RestEndpointType2�

EmptyUnion�BW
patternsK:I
G�
ts/test/all.syslss�
ts/test/all.syslss"tag�
ts/test/all.syslsv�
ts/test/all.syslsv2 2�
Alias�BW
patternsK:I
G�
ts/test/all.syslvv�
ts/test/all.syslvv"tagBy
annotation1j�
ts/test/all.syslww?�
ts/test/all.syslww?"&you can do string annotation like thisB�
annotation2��
ts/test/all.syslxx2�
ts/test/all.syslxx2:�
F�
ts/test/all.syslxx�
ts/test/all.syslxx"or
F�
ts/test/all.syslxx"�
ts/test/all.syslxx""in
F�
ts/test/all.syslx$x(�
ts/test/all.syslx$x("an
I�
ts/test/all.syslx*x1�
ts/test/all.syslx*x1"arrayB�
annotation3t�
ts/test/all.sysly}	�
ts/test/all.sysly}	"0you can also do
multiline annotations
like this
�
ts/test/all.syslv�
ts/test/all.syslv2�
AliasRef�B[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag� 
ts/test/all.sysl��� 
ts/test/all.sysl��J


TypesAliasRefType2�
Type�B\

annotationN�
ts/test/all.syslOO"�
ts/test/all.syslOO""
annotationBW
patternsK:I
G�
ts/test/all.syslNN�
ts/test/all.syslNN"tag�
ts/test/all.syslNW1�
ts/test/all.syslNW1�
�
nativeTypeField�BW
patternsK:I
G�
ts/test/all.syslP#P'�
ts/test/all.syslP#P'"tag�
ts/test/all.syslPP(�
ts/test/all.syslPP(
�
	reference�BW
patternsK:I
G�
ts/test/all.syslQ(Q,�
ts/test/all.syslQ(Q,"tag�
ts/test/all.syslQQ-�
ts/test/all.syslQQ-J)


TypesType

RestEndpointType
�
optional�BW
patternsK:I
G�
ts/test/all.syslRR!�
ts/test/all.syslRR!"tag`�
ts/test/all.syslRR"�
ts/test/all.syslRR"
�
set�BW
patternsK:I
G�
ts/test/all.syslSS"�
ts/test/all.syslSS""tag�
ts/test/all.syslSS#�
ts/test/all.syslSS#jD�
ts/test/all.syslSS#�
ts/test/all.syslSS#
�
sequence�BW
patternsK:I
G�
ts/test/all.syslT(T,�
ts/test/all.syslT(T,"tag�
ts/test/all.syslTT-�
ts/test/all.syslTT-zD�
ts/test/all.syslTT-�
ts/test/all.syslTT-
�
aliasSequence�BW
patternsK:I
G�
ts/test/all.syslU(U,�
ts/test/all.syslU(U,"tag�
ts/test/all.syslUU-�
ts/test/all.syslUU-J"


TypesTypeAliasSequence
�
	with_anno�BW
patternsK:I
G�
ts/test/all.syslVV!�
ts/test/all.syslVV!"tagBg

annotationY�
ts/test/all.syslWW1�
ts/test/all.syslWW1"this is an annotation�
ts/test/all.syslVY�
ts/test/all.syslVY2�
Enum�BW
patternsK:I
G�
ts/test/all.syslhh�
ts/test/all.syslhh"tag�
ts/test/all.syslhm�
ts/test/all.syslhm$


ENUM_3


ENUM_1


ENUM_22�
AliasSequence�BW
patternsK:I
G�
ts/test/all.sysl�
ts/test/all.sysl"tagB}
annotation1n� 
ts/test/all.sysl��?� 
ts/test/all.sysl��?"&you can do string annotation like thisB�
annotation2�� 
ts/test/all.sysl��2� 
ts/test/all.sysl��2:�
J� 
ts/test/all.sysl��� 
ts/test/all.sysl��"or
J� 
ts/test/all.sysl��"� 
ts/test/all.sysl��""in
J� 
ts/test/all.sysl�$�(� 
ts/test/all.sysl�$�("an
M� 
ts/test/all.sysl�*�1� 
ts/test/all.sysl�*�1"arrayB�
annotation3x� 
ts/test/all.sysl��	� 
ts/test/all.sysl��	"0you can also do
multiline annotations
like this
�
ts/test/all.sysl��
ts/test/all.sysl�zF�
ts/test/all.sysl��
ts/test/all.sysl�2�
AliasForeignRef�B[
patternsO:M
K� 
ts/test/all.sysl�� � 
ts/test/all.sysl�� "tag� 
ts/test/all.sysl��� 
ts/test/all.sysl��J4


TypesAliasForeignRef

RestEndpointType2�
AliasForeignRefSet�B[
patternsO:M
K� 
ts/test/all.sysl��#� 
ts/test/all.sysl��#"tag� 
ts/test/all.sysl��� 
ts/test/all.sysl��j� 
ts/test/all.sysl��� 
ts/test/all.sysl��J7


TypesAliasForeignRefSet

RestEndpointType�
ts/test/all.syslM��
ts/test/all.syslM��
Unsafe/Namespace :: Unsafe/App�

Unsafe/Namespace

Unsafe/App"[
patternsO:M
K� 
ts/test/all.sysl�$�(� 
ts/test/all.sysl�$�("tag2�
Unsafe.Type�B[
patternsO:M
K� 
ts/test/all.sysl��� 
ts/test/all.sysl��"tag� 
ts/test/all.sysl��5� 
ts/test/all.sysl��5�
�
Unsafe.Field�Bo
description`� 
ts/test/all.sysl��5� 
ts/test/all.sysl��5"Unsafe Field DescriptionB[
patternsO:M
K� 
ts/test/all.sysl��#� 
ts/test/all.sysl��#"tag� 
ts/test/all.sysl��� 
ts/test/all.sysl���
ts/test/all.sysl��5�
ts/test/all.sysl��5�
App :: with :: subpackages�

App
with
subpackages"W
patternsK:I
G�
ts/test/all.sysl �
ts/test/all.sysl "tag*
...
...�
ts/test/all.sysl!�
ts/test/all.sysl!�
RestEndpoint�

RestEndpoint"W
patternsK:I
G�
ts/test/all.sysl�
ts/test/all.sysl"tag*�
GET /�
GET /"
patterns
:
"rest:I�
ts/test/all.sysl�
ts/test/all.sysl
...B/�
ts/test/all.sysl!�
ts/test/all.sysl!*�
GET /pathwithtype/{native}�
GET /pathwithtype/{native}"
patterns
:
"rest:I�
ts/test/all.sysl##�
ts/test/all.sysl##
...Bj/pathwithtype/{native}"N
nativeD�
ts/test/all.sysl!!!�
ts/test/all.sysl!!!�
ts/test/all.sysl"%�
ts/test/all.sysl"%*�

GET /query�

GET /query"
patterns
:
"rest:I�
ts/test/all.sysl''�
ts/test/all.sysl''
...B�/queryN
nativeD�
ts/test/all.sysl&&�
ts/test/all.sysl&&R
optionalF`�
ts/test/all.sysl&&*�
ts/test/all.sysl&&*�
ts/test/all.sysl&)�
ts/test/all.sysl&)*�
PATCH /param�
PATCH /param"
patterns
:
"rest:I�
ts/test/all.sysl++�
ts/test/all.sysl++
...B
/paramJr
tmBX
patternsL:J
H�
ts/test/all.sysl* *%�
ts/test/all.sysl* *%"bodyJ

TypesType�
ts/test/all.sysl*-�
ts/test/all.sysl*-*�
POST /param�
POST /param"
patterns
:
"rest:I�
ts/test/all.sysl//�
ts/test/all.sysl//
...B
/paramJ
native�
ts/test/all.sysl.1�
ts/test/all.sysl.1*�

PUT /param�

PUT /param"
patterns
:
"rest:I�
ts/test/all.sysl33�
ts/test/all.sysl33
...B
/paramJ
	unlimitedRJ
limited
R
J
numR�
ts/test/all.sysl25�
ts/test/all.sysl25*�
GET /report.csv�
GET /report.csv"
patterns
:
"rest:I�
ts/test/all.sysl77�
ts/test/all.sysl77
...B/report.csv�
ts/test/all.sysl69�
ts/test/all.sysl69�
ts/test/all.sysl7�
ts/test/all.sysl7".
imported.sysl�
ts/test/all.sysl

