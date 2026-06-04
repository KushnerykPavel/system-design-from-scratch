package main

import "testing"

func TestRecommendedResponse(t *testing.T) {
	if got := RecommendedResponse(MovePushTradeOff); got != "name_two_options_and_explain_what_you_gain_and_give_up" {
		t.Fatalf("unexpected response: %s", got)
	}
	if got := RecommendedResponse("unknown"); got != "clarify_intent_and_keep_the_answer_structured" {
		t.Fatalf("unexpected fallback response: %s", got)
	}
}
