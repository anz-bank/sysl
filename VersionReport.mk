######## Build information report
originalVersion=$(shell git describe --tags)

ifeq (-, $(findstring -, $(originalVersion))) #it is in branch
tagName= $(firstword $(subst -,  ,$(originalVersion)))
diffLogs = $(foreach item, $(shell git log --pretty=format:"%h" $(tagName)..HEAD),$(item))

ifeq (true, $(words $(diffLogs)) = 0 || $(shell git status -s) = ) # no changes
Version=$(tagName)
else

ifeq (, $(shell git status -s)) # it doesn't have uncommitted or untracked changes
Version=$(FullCommit) (from $(tagName))
else
Version=DIRTY-$(FullCommit) (from $(tagName))
endif

endif

else

# it is in tag
ifeq (,$(shell git status -s)) # no changes
Version=$(originalVersion)
else
Version=DIRTY-$(originalVersion)
endif

endif

FullCommit=$(shell git log --pretty=format:"%H" -1)
GoVersion=$(strip $(subst  go version, ,$(shell go version)))
BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS='-X "main.Version=$(Version)" -X "main.GitFullCommit=$(FullCommit)" -X "main.BuildDate=$(BuildDate)" -X "main.GoVersion=$(GoVersion)"'
