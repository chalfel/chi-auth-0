package api

import (
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/chalfel/chi-auth-0/internal/data"
	"github.com/chalfel/chi-auth-0/pkg/db"
	"github.com/chalfel/chi-auth-0/pkg/exceptions"
	"github.com/chalfel/chi-auth-0/pkg/jwt"
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
	authMiddleware := jwt.EnsureValidToken()
	a.router.Get("/", a.Status)
	a.router.Options("/", func(w http.ResponseWriter, r *http.Request) {})

	a.router.Route("/auth", func(r chi.Router) {
		r.With(authMiddleware).Post("/callback/", errorHandler(a.AuthCallback))
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

func (a *App) AuthCallback(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("#######################")
	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

	// sql := `INSERT INTO users (id, provider_id) VALUES ($1, $2)`

	user := data.User{
		ID:         uuid.NewString(),
		ProviderID: claims.RegisteredClaims.Subject,
	}

	// if err := a.db.Exec(r.Context(), sql, user.ID, user.ProviderID); err != nil {
	// 	return err
	// }
	fmt.Printf("%+v\n", user)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ok"}`))
	return nil
}
