package version

import (
	"fmt"
)

var v *Version

type Version struct {
	Number     string
	CommitHash string
	TimeStamp  string
}

func Set(number, commitHash, timeStamp string) {
	if v == nil {
		v = &Version{
			Number:     number,
			CommitHash: commitHash,
			TimeStamp:  timeStamp,
		}
	}
}

func Get(verbose bool) string {
	if verbose {
		return fmt.Sprintf("%s-%s (%s)", v.Number, v.CommitHash, v.TimeStamp)
	}
	return v.Number
}
