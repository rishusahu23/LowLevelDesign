package main

type Processor struct{}

func (p *Processor) Process(input string) string {
	return "processed string " + input
}

type DecisionMaking struct {
}

func (d *DecisionMaking) MakeDecision(input string) string {
	return "make decision " + input
}

type BrainManager struct {
	memory         *Memory
	Processor      *Processor
	DecisionMaking *DecisionMaking
}

func NewBrainManager() *BrainManager {
	return &BrainManager{
		memory:         NewMemory(),
		Processor:      &Processor{},
		DecisionMaking: &DecisionMaking{},
	}
}

func (b *BrainManager) Think() string {
	b.memory.Store("input", "output")
	output, _ := b.memory.Recall("input")
	output = b.Processor.Process(output)
	output = b.DecisionMaking.MakeDecision(output)
	return output
}
