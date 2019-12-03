.PHONY: diagrams
diagrams: seq-diagrams data-diagram

.PHONY: data-diagram
data-diagram: data-Relational-Model.png data-Object-Model.png

data-%.png: data.sysl
	sysl data -v -o 'data-%(epname).png' --project Project data

.PHONY: seq-diagrams
seq-diagrams: call-Seq.png if-else-Seq.png args-Seq.png for-loop-Seq.png app_icons-Seq.png events-Seq.png attribs-Seq.png

%-Seq.png: %.sysl
	sysl sd -v -o "$*-%(epname).png" -a Project $*

.PHONY: int-diagrams
int-diagrams: call-Integration-Diagram.png

%-Integration-Diagram.png:
	sysl ints -v -o "$*-%(epname).png" -a Project $*
