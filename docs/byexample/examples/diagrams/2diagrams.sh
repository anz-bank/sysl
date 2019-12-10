# First, make sure to set the environment variable SYSL_PLANTUML
export SYSL_PLANTUML=http://www.plantuml.com/plantuml

# Now run the sysl sd (sequence diagram) command
sysl sd -o "diagrams.png" -s "MobileApp <- Login" diagrams.sysl

# `-o` is the output file

# `-s` specifies a starting endpoint for the sequence diagram to initiate

#  `diagrams.sysl` is the input sysl file

# diagrams.png: