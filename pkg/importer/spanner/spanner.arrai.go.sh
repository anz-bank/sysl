#!/usr/bin/env sh

set -e

OUT=../spanner.arrai.go
arrai run concat_go.arrai import.arrai | arrai eval '$"package importer\n\nconst importSpannerScript = `\n${//os.stdin}`"' > $OUT
