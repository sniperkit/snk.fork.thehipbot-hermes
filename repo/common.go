/*
Sniperkit-Bot
- Status: analyzed
*/

package repo

// Repository struct holds information for a repository
type Repository interface {
	Clone(path string) error
}
