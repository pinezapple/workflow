package services_test

import (
	"context"
	"testing"

	"workflow/heimdall/core"
	"workflow/heimdall/services"
)

func TestArborist(t *testing.T) {
	var (
		cfg = &core.ArboristConfig{
			Host: "https://genome.vinbigdata.org/authz",
		}
		ctx   = context.Background()
		token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImZlbmNlX2tleV9rZXlzIn0.eyJwdXIiOiJhY2Nlc3MiLCJhdWQiOlsiZGF0YSIsInVzZXIiLCJmZW5jZSIsIm9wZW5pZCJdLCJzdWIiOiIyOSIsImlzcyI6Imh0dHBzOi8vZ2Vub21lLnZpbmJpZ2RhdGEub3JnL3VzZXIiLCJpYXQiOjE2MjgzMzA5MTUsImV4cCI6MTYyODM2NjkxNSwianRpIjoiMzRlN2Q4OTgtOWI2Ny00YzgxLWExNzYtOWYwNjI2ZGE3ZmI0Iiwic2NvcGUiOlsiZGF0YSIsInVzZXIiLCJmZW5jZSIsIm9wZW5pZCJdLCJjb250ZXh0Ijp7InVzZXIiOnsibmFtZSI6InRoYW5ocGhhbnBodTE4QGdtYWlsLmNvbSIsImlzX2FkbWluIjpmYWxzZSwiZ29vZ2xlIjp7InByb3h5X2dyb3VwIjpudWxsfSwicHJvamVjdHMiOnsiVk5QR3giOlsidXBkYXRlIiwicmVhZC1zdG9yYWdlIiwicmVhZCIsInVwbG9hZCJdLCJWTjEwMDBHIjpbInJlYWQtc3RvcmFnZSIsInJlYWQiLCJ1cGxvYWQiXX19fSwiYXpwIjoiIn0.L3rfPdej4nzdSUVpjZ-mQLziCAUXh_Wbixg3XmfAS-QbbPEpV-gGRFs9Dkw-aIW5feV4eKU37wXwWMrmlsi5AWzNd8jt_TRmf-hqrKLtruoSxnt17NhAJXJKyPvR4xp7zEmpA3YaRZwF1g4Bw27lROubE7NXQOVAI-HbMOyNl5c1x9E2dchpNYXVe4s0_bu7nXJIHd9gbAbAUzqHilmSN4QhChguBD_47g1JQdNUyd860DN-Vw4ZCN-3y4QxxNqDAPUwy1OgfNh00ZsgGFDaSXqHOdHatuEjiWqTRC9xbYSFq-ZRz2F98faApbvAncst1NQ66Gm6thWSEEIVKbQuPg"
	)
	services.SetArborist(cfg)

	prjs, err := services.GetArboristService().GetShareProjects(ctx, token)
	if err != nil {
		t.Error(err)
		return
	}

	for i := range prjs {
		t.Logf("Project: %s", prjs[i])
	}
}
