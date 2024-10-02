package domain

import "errors"

var (
	ErrSongAlreadyExist   = errors.New("song with this name by this group already exist")
	ErrCreatingSong       = errors.New("error creating song in db")
	ErrUpdatingSong       = errors.New("error updating song")
	ErrDeletingSong       = errors.New("error deleting song")
	ErrParsingReleaseDate = errors.New("error parsing release date")
	ErrSongNotFound       = errors.New("song with this id not found")
	ErrCreatingRequest    = errors.New("error creating request to another server for getting song details")
	ErrSendingRequest     = errors.New("error sending request to another server")
	ErrResponseError      = errors.New("response error with status code")
	ErrDecodingResponse   = errors.New("error decoding response from another server")
	ErrGettingDetails     = errors.New("error getting song details from another server")
	ErrDetailsNotFound    = errors.New("details for song not found")
	ErrAddingDetails      = errors.New("error adding song details in db")
	ErrPageDoesntExist    = errors.New("this page doesn't exist")
	ErrSongTextNotFound   = errors.New("text for this song not found")
)
