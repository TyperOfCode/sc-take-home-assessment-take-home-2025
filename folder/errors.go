package folder

import "errors"

var ErrUnexpectedError = errors.New("unexpected error")

// get_folder errors
var ErrFolderDoesNotExist = errors.New("folder doesn't exist")
var ErrFolderDoesNotExistInOrg = errors.New("folder doesn't exist in the specified organization")

// move_folder errors
var ErrInvalidArguments = errors.New("empty folder name in source or destination")
var ErrMoveToSource = errors.New("cannot move a folder to itself")
var ErrMoveToDescendant = errors.New("cannot move a folder to its descendant")
var ErrMoveToDifferentOrg = errors.New("cannot move a folder to a different organization")

var ErrSourceDoesNotExist = errors.New("source folder doesn't exist")
var ErrDestDoesNotExist = errors.New("destination folder doesn't exist")
