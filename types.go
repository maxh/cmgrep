package protocol

type Cursor struct {
	ByteOffset int64 `gob:"byte_offset"` 
	LineNumber int   `gob:"line_number"`
}

type GrepRequest struct {
	Pattern        string

	// CLI options
	ExtendedRegexp bool
	IgnoreCase     bool

	// todo: pagination
	// relay-inspired cursor-based pagination
	// First int
	// After *Cursor // nil for page 1

	// todo: option for summary-only
	// SummaryOnly bool
}

type MatchLine struct {
	LineNumber int
	Content    string
}

// todo: pagination
// type PageInfo struct {
// 	HasNextPage bool
// 	EndCursor   *Cursor
// }

type MatchInfo struct {
	HostNumber int // numeric per prompt
	TotalCount int
	Lines      []MatchLine

	// todo: pagination
	// PageInfo   PageInfo
}

type GrepResponse struct {
	MatchInfo MatchInfo
	Error   string
}
