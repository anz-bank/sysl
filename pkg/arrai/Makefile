.PHONY: all diagrams init-diagram plantuml test compare clean

SYSL_PLANTUML=http://localhost:8080
OUT=out

# Raw targets that are generated directly by arr.ai.
RAW=$(OUT)/data_dictionary.csv \
	$(OUT)/terraform_gdc.tf \
	$(OUT)/sysl_polish.arrai

# Diagram targets that require rendering.
DIAGRAMS=$(OUT)/data_diagram.svg \
		 $(OUT)/integration_diagram.svg \
		 $(OUT)/sequence_diagram.svg

all: raw diagrams compare

raw: $(OUT) $(RAW)
diagrams: $(OUT) init-diagram $(DIAGRAMS)

# Creates the output directory.
$(OUT):
	mkdir -p $(OUT)

# Generates outputs from the model by running named arr.ai scripts.
.PRECIOUS: $(OUT)/%.puml
$(OUT)/%.arrai $(OUT)/%.csv $(OUT)/%.puml $(OUT)/%.tf: %.arrai $(OUT)/model.pb
	arrai run $$(echo $* | cut -d'.' -f1).arrai > $@

# Generates diagram outputs from the model by calling PlantUML with the generated .puml output.
$(OUT)/%.png $(OUT)/%.svg: $(OUT)/%.puml
	@$(eval FORMAT=`echo $@ | rev | cut -d'.' -f1 | rev`)
	curl $(SYSL_PLANTUML)/$(FORMAT)/~1$(shell cat $< | plantuml-encode) > $@

# Generates a catalog view of the model.
catalog: $(OUT)/model.pb
	arrai run --out=dir:out/catalog catalog.arrai

# Platform

# Compiles Sysl specs for arr.ai import.
.PRECIOUS: $(OUT)/%.pb
$(OUT)/%.pb: $(OUT) %.sysl
	sysl protobuf --mode=pb $*.sysl > $(OUT)/$*.pb

# Initializes the environment to support diagram generation targets.
init-diagram: plantuml
	which plantuml-encode || go install github.com/orlade/plantuml-encode

plantuml:
	# Version pinned to reduce variability of generated outputs.
	docker start plantuml || docker run -d -p 8080:8080 --name=plantuml plantuml/plantuml-server:jetty-v1.2020.14

# Testing

# Run arr.ai unit tests.
test:
	arrai test

# Directory for Sysl/arr.ai comparisons.
CMP_DIR=compare
SYSL_DIAGRAMS=$(CMP_DIR)/data_diagram.svg \
			  $(CMP_DIR)/integration_diagram.svg \
			  $(CMP_DIR)/sequence_diagram.svg
COMPARISONS=$(CMP_DIR)/data_diagram.md \
			$(CMP_DIR)/integration_diagram.md \
			$(CMP_DIR)/sequence_diagram.md

$(CMP_DIR):
	mkdir -p $(CMP_DIR)

$(CMP_DIR)/data_diagram.svg: $(CMP_DIR)
	sysl data -d -o $@ model.sysl

$(CMP_DIR)/integration_diagram.svg: $(CMP_DIR)
	sysl ints -j all -o $@ model.sysl

$(CMP_DIR)/sequence_diagram.svg: $(CMP_DIR)
	sysl sd -s "Client <- Do" -o $@ model.sysl

$(CMP_DIR)/%.md: $(CMP_DIR)
	echo "![](../out/$*.svg)\n![]($*.svg)" > $(CMP_DIR)/$*.md

# Generates all diagrams with arr.ai and Sysl, and side-by-side Markdown comparisons.
compare: $(DIAGRAMS) $(SYSL_DIAGRAMS) $(COMPARISONS)

# Deletes all generated outputs.
clean:
	$(RM) -r $(OUT) $(CMP_DIR)
