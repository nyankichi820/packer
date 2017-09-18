package vultr

import (
	"time"

	vultr "github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

type stepSnapshot struct {
	v *vultr.Client
}

func (s *stepSnapshot) Run(state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	server := state.Get("server").(vultr.Server)

	snapshot, err := s.v.CreateSnapshot(server.ID, c.Description)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Snapshot " + snapshot.ID + " created, waiting for it to complete...")
	for snapshot.Status != "complete" {
		time.Sleep(1 * time.Second)
		// Crude workaround for the lack of singular GetSnapshot() method
		if snapshots, err := s.v.GetSnapshots(); err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		} else {
			for _, s := range snapshots {
				if s.ID == snapshot.ID {
					snapshot = s
					break
				}
			}
		}
	}

	state.Put("snapshot", snapshot)
	return multistep.ActionContinue
}

func (s *stepSnapshot) Cleanup(state multistep.StateBag) {
}
