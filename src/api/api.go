package api

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	// TODO: check if we can remove those imports
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
	"github.com/treenq/treenq/pkg/crypto"
	"github.com/treenq/treenq/pkg/vel"
	"github.com/treenq/treenq/pkg/vel/auth"
	"github.com/treenq/treenq/src/domain"
	"github.com/treenq/treenq/src/repo"
	"github.com/treenq/treenq/src/repo/artifacts"
	"github.com/treenq/treenq/src/repo/extract"
	"github.com/treenq/treenq/src/resources"

	authService "github.com/treenq/treenq/src/services/auth"
	"github.com/treenq/treenq/src/services/cdk"
)

func New(conf Config) (http.Handler, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	l := vel.NewLogger(os.Stdout, slog.LevelDebug)

	db, err := resources.OpenDB(conf.DbDsn, conf.MigrationsDir)
	if err != nil {
		return nil, err
	}
	store, err := repo.NewStore(db)
	if err != nil {
		return nil, err
	}

	githubJwtIssuer := auth.NewJwtIssuer(conf.GithubClientID, []byte(conf.GithubPrivateKey), nil, conf.JwtTtl)
	authJwtIssuer := auth.NewJwtIssuer("treenq-api", []byte(conf.AuthPrivateKey), []byte(conf.AuthPublicKey), conf.AuthTtl)
	githubClient := repo.NewGithubClient(githubJwtIssuer, http.DefaultClient)
	gitDir := filepath.Join(wd, "gits")
	gitClient := repo.NewGit(gitDir)
	docker, err := artifacts.NewDockerArtifactory(conf.DockerRegistry)
	if err != nil {
		return nil, err
	}
	extractor := extract.NewExtractor()

	authMiddleware := auth.NewJwtMiddleware(authJwtIssuer, l)
	githubAuthMiddleware := vel.NoopMiddleware
	if conf.GithubWebhookSecretEnable {
		sha256Verifier := crypto.NewSha256SignatureVerifier(conf.GithubWebhookSecret, "sha256=")
		githubAuthMiddleware = crypto.NewSha256SignatureVerifierMiddleware(sha256Verifier, l)
	}

	oauthProvider := authService.New(conf.GithubClientID, conf.GithubSecret, conf.GithubRedirectURL)
	kube := cdk.NewKube(conf.Host)
	handlers := domain.NewHandler(
		store,
		githubClient,
		gitClient,
		extractor,
		docker,
		kube,
		string(conf.KubeConfig),
		oauthProvider,
		authJwtIssuer,
		conf.GithubWebhookURL,
		l,
	)
	return resources.NewRouter(handlers, authMiddleware, githubAuthMiddleware, vel.NewLoggingMiddleware(l)).Mux(), nil
}
