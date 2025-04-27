package interaction

import (
	"github.com/yxrxy/videoHub/app/interaction/controllers/rpc"
	"github.com/yxrxy/videoHub/app/interaction/domain/service"
	"github.com/yxrxy/videoHub/app/interaction/infrastructure/mysql"
	"github.com/yxrxy/videoHub/app/interaction/usecase"
	"github.com/yxrxy/videoHub/kitex_gen/interaction"
	"github.com/yxrxy/videoHub/pkg/base/client"
)

func InjectInteractionHandler() interaction.InteractionService {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}

	db := mysql.NewInteraction(gormDB)
	svc := service.NewInteractionService(db)
	uc := usecase.NewInteractionCase(db, svc)

	return rpc.NewInteractionHandler(uc)
}
