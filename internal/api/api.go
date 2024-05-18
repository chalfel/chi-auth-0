package api

import (
	"encoding/json"
	"net/http"

	"github.com/chalfel/chi-auth-0/internal/data"
	"github.com/chalfel/chi-auth-0/pkg/db"
	"github.com/chalfel/chi-auth-0/pkg/exceptions"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type App struct {
	db     *db.Db
	router *chi.Mux
}

func NewApi(db *db.Db, router *chi.Mux) *App {
	app := &App{
		db,
		router,
	}

	return app
}

func (a *App) RegisterRoutes() {

	a.router.Get("/", a.Status)

	a.router.Route("/user", func(r chi.Router) {
		r.Post("/register", errorHandler(a.RegisterUser))
	})

}

type RegisterUserParams struct {
	ProviderID string `json:"provider_id"`
}

func (r *RegisterUserParams) Validate() error {
	if r.ProviderID == "" {
		return exceptions.New("provider id is required", "provider_id_is_required")
	}

	return nil
}

func (a *App) Status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ok"}`))
}

func (a *App) RegisterUser(w http.ResponseWriter, r *http.Request) (Response, error) {
	params := RegisterUserParams{}
	resp := Response{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return resp, err
	}

	if err := params.Validate(); err != nil {
		return resp, err
	}

	sql := `INSERT INTO users (id, provider_id) VALUES ($1, $2)`

	user := data.User{
		ID:         uuid.NewString(),
		ProviderID: params.ProviderID,
	}

	if err := a.db.Exec(r.Context(), sql, user.ID, user.ProviderID); err != nil {
		return resp, err
	}

	resp.StatusCode = http.StatusCreated
	resp.Data = user

	return resp, nil
}
