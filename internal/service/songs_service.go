package service

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"net/http"
	"net/url"
	"songs-library-go/internal/delivery/dto"
	"songs-library-go/internal/domain"
	"strings"
	"time"
)

// SongsRepo defines methods for interacting with the song data store, including retrieval, creation, updating, and deletion of songs.
type SongsRepo interface {
	GetSongs(page int, limit int, filtersMap map[string]interface{}) ([]domain.Song, int, error)
	GetSongText(songID int32) (string, error)
	Delete(songID int32) error
	UpdateSong(songID int32, paramsMap map[string]interface{}) (domain.Song, error)
	Create(groupName, songName string) (domain.Song, error)
	AddDetails(songID int32, paramsMap map[string]interface{}) error
}

// SongsService manages song operations and interacts with the repository and external music information API.
type SongsService struct {
	repo            SongsRepo
	musicInfoAPIURL string
}

// NewSongsService initializes and returns a new instance of SongsService with the provided repository and music info API URL.
func NewSongsService(repo SongsRepo, musicInfoAPIURL string) *SongsService {
	return &SongsService{
		repo:            repo,
		musicInfoAPIURL: musicInfoAPIURL,
	}
}

// GetSongs retrieves songs from the repository based on the provided filtering and pagination parameters.
func (s SongsService) GetSongs(params dto.GetSongsDto) ([]domain.Song, int, error) {
	filtersMap := s.makeSongParamsMap(params.Filters)

	songs, totalPages, err := s.repo.GetSongs(params.PaginationParams.Page, params.PaginationParams.Limit, filtersMap)
	if err != nil {
		return nil, 0, err
	}

	return songs, totalPages, nil
}

// GetSongText retrieves the text of a song by its ID and paginates the verses based on the provided parameters.
func (s SongsService) GetSongText(songID int32, params dto.PaginationParamsDto) ([]string, int, error) {
	songText, err := s.repo.GetSongText(songID)
	if err != nil {
		return nil, 0, err
	}

	if songText == "" {
		return make([]string, 0), 0, nil
	}

	verses := strings.Split(songText, "\n\n")

	totalPages := int(math.Ceil(float64(len(verses)) / float64(params.Limit)))

	if params.Page > totalPages {
		return make([]string, 0), 0, nil
	}

	start := (params.Page - 1) * params.Limit
	end := start + params.Limit
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end], totalPages, nil
}

// Delete removes a song by its ID from the repository.
func (s SongsService) Delete(songID int32) error {
	return s.repo.Delete(songID)
}

// Update modifies an existing song's details based on the provided parameters.
func (s SongsService) Update(songID int32, updateSongInput dto.SongParamsDto) (domain.Song, error) {
	paramsMap := s.makeSongParamsMap(updateSongInput)

	song, err := s.repo.UpdateSong(songID, paramsMap)
	if err != nil {
		return domain.Song{}, err
	}

	return song, nil
}

// Create adds a new song to the repository and initiates the process to fetch and save its details.
func (s SongsService) Create(createSongInput dto.CreateSongDto) (domain.Song, error) {
	song, err := s.repo.Create(createSongInput.Group, createSongInput.Song)
	if err != nil {
		return domain.Song{}, err
	}

	go s.getAndSaveDetails(song.ID, createSongInput.Group, createSongInput.Song)

	return song, nil
}

func (s SongsService) makeSongParamsMap(params dto.SongParamsDto) map[string]interface{} {
	paramsMap := make(map[string]interface{})

	if params.Group != nil {
		paramsMap["group"] = *params.Group
	}

	if params.Song != nil {
		paramsMap["song"] = *params.Song
	}

	if params.ReleaseDate != nil {
		date, _ := time.Parse(domain.DateFormat, *params.ReleaseDate)
		paramsMap["release_date"] = date
	}

	if params.Text != nil {
		paramsMap["text"] = strings.TrimSpace(*params.Text)
	}

	if params.Link != nil {
		paramsMap["link"] = *params.Link
	}

	return paramsMap
}

func (s SongsService) toSongDto(song domain.Song) dto.SongDto {
	var releaseDate string
	if !song.ReleaseDate.IsZero() {
		releaseDate = song.ReleaseDate.Format(domain.DateFormat)
	}

	songDto := dto.SongDto{
		ID:          song.ID,
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: releaseDate,
		Text:        song.Text,
		Link:        song.Link,
	}

	return songDto
}

func (s SongsService) getAndSaveDetails(songID int32, groupName, songName string) {
	details, err := s.getDetails(songName, groupName)
	if err != nil {
		log.WithError(err).Error(domain.ErrGettingDetails)
		return
	}

	paramsMap := s.makeSongParamsMap(details)
	if len(paramsMap) == 0 {
		log.Errorf("%s (group name: %s, song name: %s)", domain.ErrDetailsNotFound, groupName, songName)
		return
	}

	if err := s.repo.AddDetails(songID, paramsMap); err != nil {
		log.WithError(err).Error(domain.ErrAddingDetails)
		return
	}

	log.Info(fmt.Sprintf("%s %d", domain.SuccessfulDetailAddition, songID))
}

func (s SongsService) getDetails(songName, groupName string) (dto.SongParamsDto, error) {
	requestURL := fmt.Sprintf("%s/info?group=%s&song=%s", s.musicInfoAPIURL, url.QueryEscape(groupName), url.QueryEscape(songName))

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return dto.SongParamsDto{}, fmt.Errorf("%s: %s", domain.ErrCreatingRequest, err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return dto.SongParamsDto{}, fmt.Errorf("%s: %s", domain.ErrSendingRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.SongParamsDto{}, fmt.Errorf("%s: %s", domain.ErrResponseError, resp.Status)
	}

	var details dto.SongParamsDto
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return dto.SongParamsDto{}, fmt.Errorf("%s: %s", domain.ErrDecodingResponse, err)
	}

	return details, nil
}
