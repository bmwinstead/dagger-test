package utils

// import (
// 	"github.com/go-git/go-git/v5"
// 	"github.com/go-git/go-git/v5/plumbing"
// )
//
// func GetBuildTag(path string) (string, error) {
// 	r, err := git.PlainOpen(path)
// 	iter, err := r.Tags()
// 	if err != nil {
// 		return "", err
// 	}
// 	if err := iter.ForEach(func(ref *plumbing.Reference) error {
// 		obj, err := r.TagObject(ref.Hash())
// 		switch err {
// 		case nil:
// 			// Tag object present
// 		case plumbing.ErrObjectNotFound:
// 			// Not a tag object
// 		default:
// 			// Some other error
// 			return err
// 		}
// 	}); err != nil {
// 		// Handle outer iterator error
// 	}
// 	return "", nil
// }
