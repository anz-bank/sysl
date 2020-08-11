.PHONY: all diagrams init-diagram plantuml compare clean

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

# Platform

# Compiles Sysl specs for arr.ai import.
.PRECIOUS: $(OUT)/%.pb
$(OUT)/%.pb: $(OUT) %.sysl
	sysl protobuf --mode=pb $*.sysl > $(OUT)/$*.pb

# 
blueprint:
	arrai run blueprint.arrai | arrai run blueprint_make.arrai > out/Makefile

# Initializes the environment to support diagram generation targets.
init-diagram: plantuml
	which plantuml-encode || go install github.com/orlade/plantuml-encode

plantuml:
	# Version pinned to reduce variability of generated outputs.
	docker start plantuml || docker run -d -p 8080:8080 --name=plantuml plantuml/plantuml-server:jetty-v1.2020.14

# Testing

TEST=test
SYSL_DIAGRAMS=$(TEST)/data_diagram.svg \
	  		  $(TEST)/integration_diagram.svg \
	  		  $(TEST)/sequence_diagram.svg
COMPARISONS=$(TEST)/data_diagram.md \
	  		$(TEST)/integration_diagram.md \
	  		$(TEST)/sequence_diagram.md

$(TEST):
	mkdir -p $(TEST)

$(TEST)/data_diagram.svg: $(TEST)
	sysl data -d -o $@ model.sysl

$(TEST)/integration_diagram.svg: $(TEST)
	sysl ints -j all -o $@ model.sysl

$(TEST)/sequence_diagram.svg: $(TEST)
	sysl sd -s "Client <- Do" -o $@ model.sysl

$(TEST)/%.md: $(TEST)
	echo "![](../out/$*.svg)\n![]($*.svg)" > $(TEST)/$*.md

# Generates all diagrams with arr.ai and Sysl, and side-by-side Markdown comparisons.
compare: $(DIAGRAMS) $(SYSL_DIAGRAMS) $(COMPARISONS)

# Deletes all generated outputs.
clean:
	$(RM) -r $(OUT) $(TEST)
