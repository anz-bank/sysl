# First, make sure to set the environment variable SYSL_PLANTUML
export SYSL_PLANTUML=http://www.plantuml.com/plantuml

# Now run the sysl sd (sequence diagram) command
sysl data -o project.png -j Project project.sysl

# "-o" is the output file

# "-j" specifies the project to render

#  "project.sysl" is the input sysl file

# project.png:
