package main

import "fmt"

type Condition interface {
	Evaluate(data map[string]interface{}) bool
}

type Action interface {
	Execute(data map[string]interface{})
}

type ConditionalCondition struct {
	Key      string
	Val      interface{}
	Operator string
}

func (c *ConditionalCondition) Evaluate(data map[string]interface{}) bool {
	value, _ := data[c.Key]
	switch c.Operator {
	case "<":
		return value.(float64) < c.Val.(float64)
	case ">":
		return value.(float64) < c.Val.(float64)
	case "==":
		return value.(float64) == c.Val.(float64)
	default:
		return false
	}
}

type PrintAction struct {
	Message string
}

func (p *PrintAction) Execute(data map[string]interface{}) {
	fmt.Println(p.Message)
}

type Rule struct {
	ID         string
	Name       string
	Conditions []Condition
	Actions    []Action
}

func (r *Rule) Evaluate(data map[string]interface{}) bool {
	for _, condition := range r.Conditions {
		if !condition.Evaluate(data) {
			return false
		}
	}
	return true
}

func (r *Rule) Execute(data map[string]interface{}) {
	for _, action := range r.Actions {
		action.Execute(data)
	}
}

type RuleEngine struct {
	Rules map[string]*Rule
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		Rules: make(map[string]*Rule),
	}
}

func (r *RuleEngine) AddRule(rule *Rule) {
	r.Rules[rule.ID] = rule
}

func (r *RuleEngine) Remove(ID string) {
	delete(r.Rules, ID)
}

func (r *RuleEngine) EvaluateAndExecute(data map[string]interface{}) {
	for _, rule := range r.Rules {
		if rule.Evaluate(data) {
			rule.Execute(data)
		}
	}
}

func main() {
	engine := NewRuleEngine()
	rule := &Rule{
		ID:   "id1",
		Name: "name1",
		Conditions: []Condition{
			&ConditionalCondition{
				Key:      "key",
				Val:      18.0,
				Operator: "<",
			},
		},
		Actions: []Action{
			&PrintAction{
				Message: "adult",
			},
		},
	}
	engine.AddRule(rule)
	data := map[string]interface{}{
		"key": 2.0,
	}
	engine.EvaluateAndExecute(data)
}
