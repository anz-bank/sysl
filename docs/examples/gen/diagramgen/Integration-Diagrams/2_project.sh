# First, make sure to set the environment variable SYSL_PLANTUML
export SYSL_PLANTUML=http://www.plantuml.com/plantuml

# Now run the sysl sd (sequence diagram) command
sysl ints -o 3_project.svg --project Project 1_project.sysl

# `-o` is the output file

# `-s` specifies a starting endpoint for the sequence diagram to initiate

#  `project.sysl` is the input sysl file

# project.svg:
