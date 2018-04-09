package cmds

import (
	"github.com/ipfs/go-ipfs-cmdkit"
)

// Flag names
const (
	EncShort     = "enc"
	EncLong      = "encoding"
	RecShort     = "r"
	RecLong      = "recursive"
	ChanOpt      = "stream-channels"
	TimeoutOpt   = "timeout"
	OptShortHelp = "h"
	OptLongHelp  = "help"
	DerefLong    = "dereference-command-line"
	DerefAllLong = "dereference-links"
)

// options that are used by this package
var OptionEncodingType = cmdkit.StringOption(EncLong, EncShort, "The encoding type the output should be encoded with (json, xml, or text)").WithDefault("text")
var OptionRecursivePath = cmdkit.BoolOption(RecLong, RecShort, "Add directory paths recursively").WithDefault(false)
var OptionStreamChannels = cmdkit.BoolOption(ChanOpt, "Stream channel output")
var OptionTimeout = cmdkit.StringOption(TimeoutOpt, "Set a global timeout on the command")
var OptionDerefArgs = cmdkit.BoolOption(DerefLong, "Resolve link arguments instead of adding them as links").WithDefault(false)
var OptionDerefAll = cmdkit.BoolOption(DerefAllLong, "Resolve links instead of adding them as links").WithDefault(false)
