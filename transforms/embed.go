package transforms

import (
	"embed"
)

//go:embed **
// EmbedFs contains a filesystem representation of all the embedded transforms and other auxiliary files
var EmbedFs embed.FS
