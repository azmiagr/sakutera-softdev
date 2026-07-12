package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/azmiagr/sakutera-softdev/internal/handler/rest"
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/internal/service"
	"github.com/azmiagr/sakutera-softdev/pkg/bcrypt"
	"github.com/azmiagr/sakutera-softdev/pkg/config"
	"github.com/azmiagr/sakutera-softdev/pkg/database/mariadb"
	"github.com/azmiagr/sakutera-softdev/pkg/jwt"
	"github.com/azmiagr/sakutera-softdev/pkg/middleware"
	"github.com/azmiagr/sakutera-softdev/pkg/mlclient"
	"github.com/azmiagr/sakutera-softdev/pkg/supabase"
	"github.com/azmiagr/sakutera-softdev/pkg/whatsapp"
)

func main() {
	config.LoadEnvironment()

	db, err := mariadb.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = mariadb.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	err = mariadb.Seed(db)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	bcrypt := bcrypt.Init()
	jwt := jwt.Init()
	ml := mlclient.Init()
	wa := whatsapp.Init()
	storage := supabase.Init()
	svc := service.NewService(repo, bcrypt, jwt, ml, wa, storage)

	middleware := middleware.Init(svc, jwt)
	r := rest.NewRest(svc, middleware)
	r.MountEndpoint()

	go runRetentionJob(svc)

	r.Run()
}

func runRetentionJob(svc *service.Service) {
	days, _ := strconv.Atoi(os.Getenv("NOTIFICATION_RETENTION_DAYS"))
	if days <= 0 {
		days = 30
	}

	ticker := time.NewTicker(24 * time.Hour)
	for {
		_, _ = svc.CollectorService.RedactOldNotifications(time.Duration(days) * 24 * time.Hour)
		<-ticker.C
	}
}
