package main

type Move string

const (
	MoveInterruptScale   Move = "interrupt_scale"
	MovePushTradeOff     Move = "push_tradeoff"
	MoveForceConstraint  Move = "force_constraint"
	MoveChallengeAssume  Move = "challenge_assumption"
	MoveAskImplementation Move = "ask_implementation"
)

func RecommendedResponse(move Move) string {
	switch move {
	case MoveInterruptScale:
		return "restate_changed_assumption_then_identify_first_design_impacts"
	case MovePushTradeOff:
		return "name_two_options_and_explain_what_you_gain_and_give_up"
	case MoveForceConstraint:
		return "re-prioritize_requirements_before_changing_components"
	case MoveChallengeAssume:
		return "surface_assumption_and_offer_range_based_design"
	case MoveAskImplementation:
		return "answer_briefly_then_return_to_architecture_signal"
	default:
		return "clarify_intent_and_keep_the_answer_structured"
	}
}

func main() {}
