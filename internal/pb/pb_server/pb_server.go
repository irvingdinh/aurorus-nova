package pb_server

import (
	"context"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"go.uber.org/fx"

	"github.com/irvingdinh/aurorus-nova/internal/service/pb_helper_service"
)

type PbServer interface{}

func NewPbServer(
	lc fx.Lifecycle,
	pbHelperService pb_helper_service.PbHelperService,
) PbServer {
	i := &pbServerImpl{
		pbHelperService: pbHelperService,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return i.Start()
		},
		OnStop: func(ctx context.Context) error {
			return i.Stop()
		},
	})

	return i
}

type pbServerImpl struct {
	pbHelperService pb_helper_service.PbHelperService

	app *pocketbase.PocketBase
}

func (i *pbServerImpl) Start() error {
	app := pocketbase.New()

	app.OnRecordAfterCreateSuccess("digests").BindFunc(func(e *core.RecordEvent) error {
		err := i.pbHelperService.EnsureThumbnails(e, "featured_image")
		if err != nil {
			e.App.Logger().Error(err.Error())
		}

		return e.Next()
	})

	app.OnRecordAfterUpdateSuccess("digests").BindFunc(func(e *core.RecordEvent) error {
		err := i.pbHelperService.EnsureThumbnails(e, "featured_image")
		if err != nil {
			e.App.Logger().Error(err.Error())
		}

		return e.Next()
	})

	i.app = app

	go func() {
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (i *pbServerImpl) Stop() error {

	if i.app != nil {
		event := &core.TerminateEvent{App: i.app}
		return i.app.OnTerminate().Trigger(event, func(e *core.TerminateEvent) error {
			return e.App.ResetBootstrapState()
		})
	}

	return nil
}
