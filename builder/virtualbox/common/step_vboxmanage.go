package common

import (
	"context"
	"fmt"

	"github.com/txgruppi/parseargs-go"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type commandTemplate struct {
	Name string
}

// This step executes additional VBoxManage commands as specified by the
// template.
//
// Uses:
//   driver Driver
//   ui packer.Ui
//   vmName string
//
// Produces:
type StepVBoxManage struct {
	Commands []string
	Ctx      interpolate.Context
}

func (s *StepVBoxManage) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	if len(s.Commands) > 0 {
		ui.Say("Executing custom VBoxManage commands...")
	}

	s.Ctx.Data = &commandTemplate{
		Name: vmName,
	}

	for _, originalCommand := range s.Commands {
		command := originalCommand
		var err error
		command, err = interpolate.Render(command, &s.Ctx)
		if err != nil {
			err := fmt.Errorf("Error preparing vboxmanage command: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		commandWords, err := parseargs.Parse(command)
		if err != nil {
			err := fmt.Errorf("Error preparing vboxmanage command: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		ui.Message(fmt.Sprintf("Executing: %s", command))
		if err := driver.VBoxManage(commandWords...); err != nil {
			err := fmt.Errorf("Error executing command: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (s *StepVBoxManage) Cleanup(state multistep.StateBag) {}
