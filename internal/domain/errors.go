package domain

import "errors"

// Error variables for various song-related operations.
var (
	ErrSongNotFound     = errors.New("song with this id not found")
	ErrSongAlreadyExist = errors.New("song with this name by this group already exist")
	ErrCreatingRequest  = errors.New("error creating request to another server for getting song details")
	ErrSendingRequest   = errors.New("error sending request to another server")
	ErrResponseError    = errors.New("response error with status code")
	ErrDecodingResponse = errors.New("error decoding response from another server")
	ErrGettingDetails   = errors.New("error getting song details from another server")
	ErrDetailsNotFound  = errors.New("details for song not found")
	ErrAddingDetails    = errors.New("error adding song details in db")
)
