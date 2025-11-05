package app

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	authAPI "github.com/Daniil-Sakharov/RocketFactory/auth/internal/api/grpc/auth/v1"
	userAPI "github.com/Daniil-Sakharov/RocketFactory/auth/internal/api/grpc/user/v1"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/config"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
	sessionRepo "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/session"
	userRepo "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/user"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service"
	authService "github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/auth"
	userService "github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/user"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/cache"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/cache/redis"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/closer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/migrator"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/migrator/pg"
	authv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/auth/v1"
	userv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	userV1API userv1.UserServiceServer
	authV1API authv1.AuthServiceServer

	userService service.UserService
	authService service.AuthService

	userRepository         repository.UsersRepository
	sessionCacheRepository repository.SessionCacheRepository

	postgresDB *sqlx.DB
	migrator   migrator.Migrator

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) UserAPI(ctx context.Context) userv1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userAPI.NewAPI(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) AuthAPI(ctx context.Context) authv1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authAPI.NewAPI(d.AuthService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(
			d.UsersRepository(ctx),
			d.SessionRepository(),
			config.AppConfig().Redis.CacheTTL(),
		)
	}
	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(
			d.UsersRepository(ctx),
			d.SessionRepository(),
		)
	}
	return d.userService
}

func (d *diContainer) SessionRepository() repository.SessionCacheRepository {
	if d.sessionCacheRepository == nil {
		d.sessionCacheRepository = sessionRepo.NewRepository(d.RedisClient())
	}
	return d.sessionCacheRepository
}

func (d *diContainer) UsersRepository(ctx context.Context) repository.UsersRepository {
	if d.userRepository == nil {
		d.userRepository = userRepo.NewRepository(d.PostgresDB())
	}
	return d.userRepository
}

func (d *diContainer) Migrator(ctx context.Context) migrator.Migrator {
	if d.migrator == nil {
		db := d.PostgresDB()
		d.migrator = pg.NewMigrator(db.DB, config.AppConfig().Postgres.MigrationsDir())
	}
	return d.migrator
}

func (d *diContainer) PostgresDB() *sqlx.DB {
	if d.postgresDB == nil {
		db, err := sqlx.Connect("pgx", config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("Ошибка в подключении к PostgreSQL: %s\n", err.Error()))
		}

		err = db.Ping()
		if err != nil {
			panic(fmt.Sprintf("Ошибка в соединении с PostgreSQL: %s\n", err.Error()))
		}

		closer.AddNamed("PostgreSQL", func(ctx context.Context) error {
			return db.Close()
		})

		d.postgresDB = db
	}

	return d.postgresDB
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}
	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}
	return d.redisClient
}
