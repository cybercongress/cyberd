package types

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
	Jobs Jobs 	  `json:"routes" yaml:"routes"`
}

func NewGenesisState(params Params, jobs []Job) GenesisState {
	return GenesisState{
		Params: 	 params,
		Jobs: 	 	 jobs,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
