package kraken

import "testing"

func TestFormatLimitPrice(t *testing.T) {
	tests := []struct {
		askPrice        float64
		limitAdjustment float64
		want            string
	}{
		{69946.5, 0.1, "69946.4"},
		{69946.3, 0.1, "69946.2"},
		{69946.3, 0.30000000000000004, "69946.0"}, // 0.1 + 0.1 + 0.1
		{100000, 0.1, "99999.9"},
		{69946.25, 0.1, "69946.2"},
	}

	for _, tt := range tests {
		if got := formatLimitPrice(tt.askPrice, tt.limitAdjustment); got != tt.want {
			t.Errorf("formatLimitPrice(%v, %v) = %s, want %s", tt.askPrice, tt.limitAdjustment, got, tt.want)
		}
	}
}

func TestLimitOrderArgsArePostOnly(t *testing.T) {
	args := limitOrderArgs(69946.5, 0.1)

	if args["oflags"] != "post" {
		t.Errorf("oflags = %v, want post", args["oflags"])
	}
	if args["validate"] != false {
		t.Errorf("validate = %v, want false", args["validate"])
	}
	if _, ok := args["price"].(string); !ok {
		t.Errorf("price should be a string, got %T", args["price"])
	}
}
