# First, make sure to set the environment variable SYSL_PLANTUML
export SYSL_PLANTUML=http://www.plantuml.com/plantuml

# Now run the sysl data model command
sysl data -o "3_project.svg" -j Project 1_project.sysl


# "-o" is the output file
# "%(epname).svg" is a special "hack" and will generate a seperate data model diagram for all the applications defined within the project:
ls
# 1_project.sysl 2_project.sh  App.svg Server.svg


# "-j" specifies the project to render

#  "1_project.sysl" is the input sysl file

# NOTE: there is currently a bug where data-types defined in different applications don't render correctly

# See https://github.com/anz-bank/sysl/issues/474 for progress updates

# 3_project.svg:

