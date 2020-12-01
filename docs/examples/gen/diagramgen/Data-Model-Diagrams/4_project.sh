

# If seperate data model diagrams are needed for every application, "%(epname).svg" can be used as the output file, and a data model will be rendered for every application
sysl data -o "%(epname).svg" -j Project 1_project.sysl


ls
# 1_project.sysl 4_project.sh  App.svg Server.svg

# App.svg and Server.svg separately:
