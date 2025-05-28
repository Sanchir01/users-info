package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"

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
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				return
			}
		}
	}()
	genderuser, err := s.GetGenderUser(ctx, name)
	if err != nil {
		slog.Error("ошибка получения пола пользователя:", err.Error())
		return err
	}
	ageuser, err := s.GetAgeUser(ctx, name)
	if err != nil {
		slog.Error("ошибка получения пола пользователя:", err.Error())
		return err
	}
	nationalityuser, err := s.GetNationalityUser(ctx, name)
	if err != nil {
		slog.Error("ошибка получения пола пользователя:", err.Error())
		return err
	}
	slog.Info("пол пользователя", genderuser, "user age", ageuser, "user nationality", nationalityuser)
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
	slog.Warn("gender url", urlgender)
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
	fmt.Println("data gender", data)
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
	slog.Warn("gender url", urlgender)
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
	slog.Info("user age", data)
	if err := validator.New().Struct(&data); err != nil {
		slog.Error("ошибка валидации данных:", err.Error())
		return 0, fmt.Errorf("validation error: %w", err)
	}
	fmt.Println("data gender", data)
	return data.Age, nil

}

func (s *Service) GetNationalityUser(ctx context.Context, name string) (string, error) {
	urlgender := fmt.Sprintf(nationalityUrl, name)
	slog.Warn("gender url", urlgender)
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
	slog.Info("user nationality", data)
	if err := validator.New().Struct(&data); err != nil {
		slog.Error("ошибка валидации данных:", err.Error())
		return "", fmt.Errorf("validation error: %w", err)
	}
	fmt.Println("data gender", data)
	return data.Country[0].CountryID, nil

}
