package service

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"songs-library-go/internal/delivery/dto"
	"songs-library-go/internal/domain"
	"strings"
)

type SongsRepo interface {
	GetSongs(page int, limit int, filtersMap map[string]string) ([]domain.Song, int, error)
	GetSongText(songID int32) (string, error)
	Delete(songID int32) error
	UpdateSong(songID int32, paramsMap map[string]string) (domain.Song, error)
	Create(groupName, songName string) (domain.Song, error)
	AddDetails(songID int32, paramsMap map[string]string) error
}

type SongsService struct {
	repo            SongsRepo
	musicInfoApiUrl string
}

func NewSongsService(repo SongsRepo, musicInfoApiUrl string) *SongsService {
	return &SongsService{
		repo:            repo,
		musicInfoApiUrl: musicInfoApiUrl,
	}
}

func (s SongsService) GetSongs(params dto.GetSongsDto) ([]domain.Song, int, error) {
	filtersMap := s.makeSongParamsMap(params.Filters)

	songs, totalPages, err := s.repo.GetSongs(params.PaginationParams.Page, params.PaginationParams.Limit, filtersMap)
	if err != nil {
		return nil, 0, err
	}

	return songs, totalPages, nil
}

func (s SongsService) GetSongText(songID int32, params dto.PaginationParamsDto) ([]string, int, error) {
	songText, err := s.repo.GetSongText(songID)
	if err != nil {
		return nil, 0, err
	}

	verses := strings.Split(songText, "\n\n")

	start := (params.Page - 1) * params.Limit
	end := start + params.Limit
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end], len(verses)/params.Limit + 1, nil
}

func (s SongsService) Delete(songID int32) error {
	return s.repo.Delete(songID)
}

func (s SongsService) Update(songID int32, updateSongInput dto.SongParamsDto) (domain.Song, error) {
	paramsMap := s.makeSongParamsMap(updateSongInput)

	song, err := s.repo.UpdateSong(songID, paramsMap)
	if err != nil {
		return domain.Song{}, err
	}

	return song, nil
}

func (s SongsService) Create(createSongInput dto.CreateSongDto) (domain.Song, error) {
	song, err := s.repo.Create(createSongInput.Group, createSongInput.Song)
	if err != nil {
		return domain.Song{}, err
	}

	go s.getAndSaveDetails(song.ID, createSongInput.Group, createSongInput.Song)

	return song, nil
}

//	func (s SongsService) getAndSaveDetails(songID int32, groupName, songName string) {
//		params := domain.AddDetailsParams{ID: songID}
//
//		releaseDateTime, err := time.Parse(domain.DateFormat, "15.10.2010")
//		if err != nil {
//			log.WithError(err).Error(domain.ErrParsingReleaseDate)
//			return
//		}
//
//		params.ReleaseDate = &releaseDateTime
//
//		text := "some text"
//		link := "some link"
//
//		params.Text = &text
//		params.Link = &link
//
//		if err := s.repo.AddDetails(params); err != nil {
//			log.WithError(err).Error(domain.ErrAddingDetails)
//			return
//		}
//
//		log.Info(fmt.Sprintf("%s %d", domain.SuccessfulDetailAddition, songID))
//	}

func (s SongsService) makeSongParamsMap(params dto.SongParamsDto) map[string]string {
	paramsMap := make(map[string]string)

	if params.Group != nil {
		paramsMap["group"] = *params.Group
	}

	if params.Song != nil {
		paramsMap["song"] = *params.Song
	}

	if params.ReleaseDate != nil {
		paramsMap["release_date"] = *params.ReleaseDate
	}

	if params.Text != nil {
		paramsMap["text"] = *params.Text
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
	requestUrl := fmt.Sprintf("%s/info?group=%s&song=%s", s.musicInfoApiUrl, url.QueryEscape(groupName), url.QueryEscape(songName))

	req, err := http.NewRequest("GET", requestUrl, nil)
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
