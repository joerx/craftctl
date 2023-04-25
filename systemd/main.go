package systemd

import (
	"context"
	"fmt"
	"joerx/minecraft-cli/handler"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/coreos/go-systemd/v22/unit"
)

type unitController struct {
	name string
	conn *dbus.Conn
}

func NewUnitController(ctx context.Context, name string) (*unitController, error) {
	conn, err := dbus.NewUserConnectionContext(ctx)
	if err != nil {
		return nil, err
	}
	return &unitController{name, conn}, nil
}

func (uc *unitController) Start(ctx context.Context) (string, error) {
	en := unit.UnitNameEscape(uc.name)
	ch := make(chan string)

	// TODO: add timeout
	if _, err := uc.conn.StartUnitContext(ctx, en, "replace", ch); err != nil {
		return "", err
	}
	state := <-ch

	return state, nil
}

func (uc *unitController) Stop(ctx context.Context) (string, error) {
	en := unit.UnitNameEscape(uc.name)
	ch := make(chan string)

	if _, err := uc.conn.StopUnitContext(ctx, en, "replace", ch); err != nil {
		return "", err
	}
	state := <-ch

	return state, nil
}

func (uc *unitController) GetState(ctx context.Context) (handler.State, error) {
	state := handler.State{}

	en := unit.UnitNameEscape(uc.name)
	us, err := uc.conn.ListUnitsByNamesContext(ctx, []string{en})
	if err != nil {
		return state, err
	}

	if len(us) == 0 {
		return state, fmt.Errorf("no units found for name '%s'", uc.name)
	}

	switch us[0].ActiveState {
	case "inactive":
		state.State = handler.StateInactive
	case "active":
		state.State = handler.StateActive
	default:
		state.State = handler.StateUnknown
	}

	state.StateDetail = fmt.Sprintf("%s (%s)", us[0].ActiveState, us[0].SubState)
	state.Name = unit.UnitNameUnescape(us[0].Name)

	return state, nil
}
