.PHONY: all diagrams clean init-diagram plantuml

SYSL_PLANTUML=http://localhost:8080
OUT=out

# Raw targets that are generated directly by arr.ai.
RAW=$(OUT)/data_dictionary.csv \
	$(OUT)/terraform_gdc.tf \
	$(OUT)/sysl_demo.arrai

# Diagram targets that require rendering.
DIAGRAMS=$(OUT)/data_diagram.svg \
		 $(OUT)/integration_diagram.svg \

all: raw diagrams

raw: $(OUT) $(RAW)
diagrams: $(OUT) init-diagram $(DIAGRAMS)

# Deletes all generated outputs.
clean:
	$(RM) -r $(OUT)

# Creates the output directory.
$(OUT):
	mkdir -p $(OUT)

# Generates outputs from the model by running named arr.ai scripts.
.PRECIOUS: $(OUT)/%.puml
$(OUT)/%.arrai $(OUT)/%.csv $(OUT)/%.puml $(OUT)/%.tf: $(OUT)/model.pb
	arrai run $$(echo $* | cut -d'.' -f1).arrai > $@

# Generates diagram outputs from the model by calling PlantUML with the generated .puml output.
$(OUT)/%.png $(OUT)/%.svg: $(OUT)/%.puml
	@$(eval FORMAT=`echo $@ | rev | cut -d'.' -f1 | rev`)
	curl $(SYSL_PLANTUML)/$(FORMAT)/~1$(shell cat $< | plantuml-encode) > $@

# Platform

# Compiles Sysl specs for arr.ai import.
.PRECIOUS: $(OUT)/%.pb
$(OUT)/%.pb: %.sysl
	sysl protobuf --mode=pb $*.sysl > $(OUT)/$*.pb

# Initializes the environment to support diagram generation targets.
init-diagram: plantuml
	which plantuml-encode || go install github.com/orlade/plantuml-encode

plantuml:
	# Version pinned to reduce variability of generated outputs.
	docker start plantuml || docker run -d -p 8080:8080 --name=plantuml plantuml/plantuml-server:jetty-v1.2020.14
