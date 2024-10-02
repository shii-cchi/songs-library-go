package service

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"songs-library-go/internal/delivery/dto"
	"songs-library-go/internal/domain"
	"strings"
	"time"
)

type SongsRepo interface {
	Create(groupName, songName string) (domain.Song, error)
	AddDetails(params domain.AddDetailsParams) error
	UpdateSong(params domain.UpdateParams) (domain.Song, error)
	Delete(songID int32) error
	GetSongs(page int, limit int, filtersMap map[string]string) ([]domain.Song, error)
	GetSong(songID int32) (string, error)
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

func (s SongsService) Create(createSongInput dto.CreateSongDto) (dto.SongDto, error) {
	song, err := s.repo.Create(createSongInput.Group, createSongInput.Song)
	if err != nil {
		if errors.Is(err, domain.ErrSongAlreadyExist) {
			return dto.SongDto{}, fmt.Errorf("%w (group name: %s, song name: %s)", err, createSongInput.Group, createSongInput.Song)
		}

		return dto.SongDto{}, fmt.Errorf("%w (group name: %s, song name: %s): %s", domain.ErrCreatingSong, createSongInput.Group, createSongInput.Song, err)
	}

	go s.getAndSaveDetails(song.ID, createSongInput.Group, createSongInput.Song)

	return dto.SongDto{
		ID:    song.ID,
		Group: song.Group,
		Song:  song.Song,
	}, nil
}

func (s SongsService) Update(updateSongInput dto.UpdateSongDto, songID int32) (dto.SongDto, error) {
	params, err := s.buildUpdateParams(songID, updateSongInput)
	if err != nil {
		return dto.SongDto{}, err
	}

	song, err := s.repo.UpdateSong(params)
	if err != nil {
		if errors.Is(err, domain.ErrSongAlreadyExist) || errors.Is(err, domain.ErrSongNotFound) {
			return dto.SongDto{}, fmt.Errorf("%w (id: %d)", err, songID)
		}

		return dto.SongDto{}, fmt.Errorf("%w (id: %d): %s", domain.ErrUpdatingSong, songID, err)
	}

	var releaseDate string
	if !song.ReleaseDate.IsZero() {
		releaseDate = song.ReleaseDate.Format(domain.DateFormat)
	}

	return dto.SongDto{
		ID:          songID,
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: releaseDate,
		Text:        song.Text,
		Link:        song.Link,
	}, nil
}

func (s SongsService) Delete(songID int32) error {
	if err := s.repo.Delete(songID); err != nil {
		if errors.Is(err, domain.ErrSongNotFound) {
			return fmt.Errorf("%w (id: %d)", domain.ErrSongNotFound, songID)
		}

		return fmt.Errorf("%w (id: %d): %s", domain.ErrDeletingSong, songID, err)
	}

	return nil
}

func (s SongsService) GetSongs(page int, limit int, filters map[string]string) ([]dto.SongDto, error) {
	songs, err := s.repo.GetSongs(page, limit, filters)
	if err != nil {
		return nil, err
	}

	return s.toSongDto(songs), nil
}

func (s SongsService) GetSong(songID int32, page int, limit int) (dto.VerseDto, error) {
	songText, err := s.repo.GetSong(songID)
	if err != nil {
		return dto.VerseDto{}, err
	}

	verses := strings.Split(songText, "\n\n")

	if page > len(verses)/limit+1 {
		return dto.VerseDto{}, fmt.Errorf("%w (page: %d, total page: %d)", domain.ErrPageDoesntExist, page, len(verses)/limit+1)
	}

	start := (page - 1) * limit
	end := start + limit
	if end > len(verses) {
		end = len(verses)
	}

	versesPage := verses[start:end]

	return dto.VerseDto{Verses: versesPage}, nil
}

//func (s SongsService) getAndSaveDetails(songID int32, groupName, songName string) {
//	params := domain.AddDetailsParams{ID: songID}
//
//	releaseDateTime, err := time.Parse(domain.DateFormat, "15.10.2010")
//	if err != nil {
//		log.WithError(err).Error(domain.ErrParsingReleaseDate)
//		return
//	}
//
//	params.ReleaseDate = &releaseDateTime
//
//	text := "some text"
//	link := "some link"
//
//	params.Text = &text
//	params.Link = &link
//
//	if err := s.repo.AddDetails(params); err != nil {
//		log.WithError(err).Error(domain.ErrAddingDetails)
//		return
//	}
//
//	log.Info(fmt.Sprintf("%s %d", domain.SuccessfulDetailAddition, songID))
//}

func (s SongsService) getAndSaveDetails(songID int32, groupName, songName string) {
	details, err := s.getDetails(songName, groupName)
	if err != nil {
		log.WithError(err).Error(domain.ErrGettingDetails)
		return
	}

	if !s.isAnyFieldProvided(details) {
		log.Errorf("%s (group name: %s, song name: %s)", domain.ErrDetailsNotFound, groupName, songName)
		return
	}

	params, err := s.buildDetailsParams(songID, details)
	if err != nil {
		log.Error(err)
		return
	}

	if err := s.repo.AddDetails(params); err != nil {
		log.WithError(err).Error(domain.ErrAddingDetails)
		return
	}

	log.Info(fmt.Sprintf("%s %d", domain.SuccessfulDetailAddition, songID))
}

func (s SongsService) getDetails(songName, groupName string) (domain.Details, error) {
	requestUrl := fmt.Sprintf("%s/info?group=%s&song=%s", s.musicInfoApiUrl, url.QueryEscape(groupName), url.QueryEscape(songName))

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return domain.Details{}, fmt.Errorf("%s: %s", domain.ErrCreatingRequest, err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return domain.Details{}, fmt.Errorf("%s: %s", domain.ErrSendingRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Details{}, fmt.Errorf("%s: %s", domain.ErrResponseError, resp.Status)
	}

	var details domain.Details
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return domain.Details{}, fmt.Errorf("%s: %s", domain.ErrDecodingResponse, err)
	}

	return details, nil
}

func (s SongsService) isAnyFieldProvided(details domain.Details) bool {
	return details.ReleaseDate != nil || details.Text != nil || details.Link != nil
}

func (s SongsService) buildDetailsParams(songID int32, details domain.Details) (domain.AddDetailsParams, error) {
	params := domain.AddDetailsParams{ID: songID}

	if details.ReleaseDate != nil {
		releaseDateTime, err := time.Parse(domain.DateFormat, *details.ReleaseDate)
		if err != nil {
			return domain.AddDetailsParams{}, fmt.Errorf("%w: %s", domain.ErrParsingReleaseDate, err)
		}
		params.ReleaseDate = &releaseDateTime
	}

	if details.Text != nil {
		params.Text = details.Text
	}

	if details.Link != nil {
		params.Link = details.Link
	}

	return params, nil
}

func (s SongsService) buildUpdateParams(songID int32, input dto.UpdateSongDto) (domain.UpdateParams, error) {
	params, err := s.buildDetailsParams(songID, domain.Details{
		ReleaseDate: input.ReleaseDate,
		Text:        input.Text,
		Link:        input.Link,
	})
	if err != nil {
		return domain.UpdateParams{}, err
	}

	var updateParams domain.UpdateParams

	if input.Group != nil {
		updateParams.Group = input.Group
	}

	if input.Song != nil {
		updateParams.Song = input.Song
	}

	updateParams.Details = params

	return updateParams, nil
}

func (s SongsService) toSongDto(songs []domain.Song) []dto.SongDto {
	var songsDto []dto.SongDto

	for _, song := range songs {
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

		songsDto = append(songsDto, songDto)
	}

	return songsDto
}
