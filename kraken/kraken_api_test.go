package kraken

import "testing"

func TestParseSpreadAskPrice(t *testing.T) {
	const pair = "XXBTZUSD"

	valid := `{"error":[],"result":{"XXBTZUSD":[[1700000000,"69940.1","69946.0"],[1700000001,"69945.0","69946.5"]],"last":1700000001}}`
	got, err := parseSpreadAskPrice([]byte(valid), pair)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 69946.5 {
		t.Errorf("got %v, want 69946.5", got)
	}

	bad := map[string]string{
		"kraken error":     `{"error":["EQuery:Unknown asset pair"]}`,
		"missing pair":     `{"error":[],"result":{"last":1}}`,
		"empty list":       `{"error":[],"result":{"XXBTZUSD":[]}}`,
		"entry not a list": `{"error":[],"result":{"XXBTZUSD":["x"]}}`,
		"empty entry":      `{"error":[],"result":{"XXBTZUSD":[[]]}}`,
		"price not string": `{"error":[],"result":{"XXBTZUSD":[[1,"1",2]]}}`,
		"price not number": `{"error":[],"result":{"XXBTZUSD":[[1,"1","abc"]]}}`,
		"zero price":       `{"error":[],"result":{"XXBTZUSD":[[1,"1","0"]]}}`,
		"invalid json":     `<html>`,
		"result not map":   `{"error":[],"result":[]}`,
	}

	for name, body := range bad {
		if _, err := parseSpreadAskPrice([]byte(body), pair); err == nil {
			t.Errorf("%s: expected an error, got nil", name)
		}
	}
}
