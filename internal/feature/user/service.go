package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	genderctx "github.com/Sanchir01/users-info/internal/gender"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	genderUrl      = "https://api.genderize.io/?name=%s"
	nationalityUrl = "https://api.nationalize.io/?name=%s"
	ageUrl         = "https://api.agify.io/?name=%s"
)

type Service struct {
	repo       *Repository
	primaryDB  *pgxpool.Pool
	httpClient *http.Client
}

func NewService(repo *Repository, primaryDB *pgxpool.Pool) *Service {
	return &Service{repo: repo, primaryDB: primaryDB, httpClient: &http.Client{Timeout: 5 * time.Second}}
}

func (s *Service) CreateUserService(
	name, surname, patronymic string,
	ctx context.Context,
) error {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	type result struct {
		value interface{}
		err   error
	}

	genderCh := make(chan result, 1)
	ageCh := make(chan result, 1)
	nationalityCh := make(chan result, 1)
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				return
			}
		}
	}()

	go func() {
		gender, err := s.GetGenderUser(ctx, name)
		genderCh <- result{value: gender, err: err}
	}()
	go func() {
		age, err := s.GetAgeUser(ctx, name)
		ageCh <- result{value: age, err: err}
	}()
	go func() {
		nationality, err := s.GetNationalityUser(ctx, name)
		nationalityCh <- result{value: nationality, err: err}
	}()

	genderRes := <-genderCh
	ageRes := <-ageCh
	nationalityRes := <-nationalityCh

	if genderRes.err != nil {
		slog.Error("ошибка получения пола пользователя:", genderRes.err.Error())
		return genderRes.err
	}
	if ageRes.err != nil {
		slog.Error("ошибка получения возраста пользователя:", ageRes.err.Error())
		return ageRes.err
	}
	if nationalityRes.err != nil {
		slog.Error("ошибка получения национальности пользователя:", nationalityRes.err.Error())
		return nationalityRes.err
	}

	genderuser, ok := genderRes.value.(genderctx.Gender)
	if !ok {
		return fmt.Errorf("не удалось привести genderuser к string")
	}
	ageuser, ok := ageRes.value.(int)
	if !ok {
		return fmt.Errorf("не удалось привести ageuser к int")
	}
	nationalityuser, ok := nationalityRes.value.(string)
	if !ok {
		return fmt.Errorf("не удалось привести nationalityuser к string")
	}
	if err := s.repo.CreateUserRepository(name, surname, patronymic, nationalityuser, ageuser, genderuser, tx, ctx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (s *Service) GetGenderUser(ctx context.Context, name string) (genderctx.Gender, error) {
	urlgender := fmt.Sprintf(genderUrl, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlgender, nil)
	if err != nil {
		return genderctx.Unknown, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return genderctx.Unknown, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("ошибка получения данных с genderize:", resp.Status)
		return genderctx.Unknown, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	var data GenderizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.Error("ошибка парсинга:", err.Error())
		return genderctx.Unknown, err
	}
	if err := validator.New().Struct(&data); err != nil {
		slog.Error("ошибка валидации данных:", err.Error())
		return genderctx.Unknown, fmt.Errorf("validation error: %w", err)
	}
	switch data.Gender {
	case genderctx.GenderMale:
		return genderctx.GenderMale, nil
	case genderctx.GenderFemale:
		return genderctx.GenderFemale, nil
	default:
		return genderctx.Unknown, nil
	}

}

func (s *Service) GetAgeUser(ctx context.Context, name string) (int, error) {
	urlgender := fmt.Sprintf(ageUrl, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlgender, nil)
	if err != nil {
		return 0, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("ошибка получения данных с genderize:", resp.Status)
		return 0, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	var data UserAgeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.Error("ошибка парсинга:", err.Error())
		return 0, err
	}
	if err := validator.New().Struct(&data); err != nil {
		slog.Error("ошибка валидации данных:", err.Error())
		return 0, fmt.Errorf("validation error: %w", err)
	}
	return data.Age, nil

}

func (s *Service) GetNationalityUser(ctx context.Context, name string) (string, error) {
	urlgender := fmt.Sprintf(nationalityUrl, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlgender, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("ошибка получения данных с genderize:", resp.Status)
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	var data NationalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.Error("ошибка парсинга:", err.Error())
		return "", err
	}
	if err := validator.New().Struct(&data); err != nil {
		slog.Error("ошибка валидации данных:", err.Error())
		return "", fmt.Errorf("validation error: %w", err)
	}
	return data.Country[0].CountryID, nil

}
