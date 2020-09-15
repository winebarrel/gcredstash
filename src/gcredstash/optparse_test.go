package gcredstash

import (
	"reflect"
	"testing"
)

func TestParseOptionWithValue1(t *testing.T) {
	args := []string{"-a", "-b", "BBB", "-c", "CCC"}
	expectedArgs := []string{"-a", "-c", "CCC"}
	expectedValue := "BBB"

	newAags, value, err := ParseOptionWithValue(args, "-b")

	if !reflect.DeepEqual(expectedArgs, newAags) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedArgs, newAags)
	}

	if expectedValue != value {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedValue, value)
	}

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestParseOptionWithValue2(t *testing.T) {
	args := []string{"-a", "-c", "CCC"}
	expectedArgs := []string{"-a", "-c", "CCC"}
	expectedValue := ""

	newAags, value, err := ParseOptionWithValue(args, "-b")

	if !reflect.DeepEqual(expectedArgs, newAags) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedArgs, newAags)
	}

	if expectedValue != value {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedValue, value)
	}

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestErrParseOptionWithValue1(t *testing.T) {
	args := []string{"-a", "-b", "-c", "CCC"}
	expected := "option requires an argument: -b"

	_, _, err := ParseOptionWithValue(args, "-b")

	if err == nil || err.Error() != expected {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, err)
	}
}

func TestErrParseOptionWithValue2(t *testing.T) {
	args := []string{"-a", "-b", "-c"}
	expected := "option requires an argument: -c"

	_, _, err := ParseOptionWithValue(args, "-c")

	if err == nil || err.Error() != expected {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, err)
	}
}

func TestParseVersion1(t *testing.T) {
	args := []string{"-a", "-v", "1", "-c", "CCC"}
	expectedArgs := []string{"-a", "-c", "CCC"}
	expectedVersion := "0000000000000000001"

	newAags, version, err := ParseVersion(args)

	if !reflect.DeepEqual(expectedArgs, newAags) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedArgs, newAags)
	}

	if expectedVersion != version {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedVersion, version)
	}

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestParseVersion2(t *testing.T) {
	args := []string{"-a", "-c", "CCC"}
	expectedArgs := []string{"-a", "-c", "CCC"}
	expectedVersion := ""

	newAags, version, err := ParseVersion(args)

	if !reflect.DeepEqual(expectedArgs, newAags) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedArgs, newAags)
	}

	if expectedVersion != version {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedVersion, version)
	}

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestErrParseVersion1(t *testing.T) {
	args := []string{"-a", "-v", "-c", "CCC"}
	expected := "option requires an argument: -v"

	_, _, err := ParseVersion(args)

	if err == nil || err.Error() != expected {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, err)
	}
}

func TestErrParseVersion2(t *testing.T) {
	args := []string{"-a", "-v", "X", "-c", "CCC"}
	expected := `strconv.Atoi: parsing "X": invalid syntax`

	_, _, err := ParseVersion(args)

	if err == nil || err.Error() != expected {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, err)
	}
}

func TestParseContext(t *testing.T) {
	args := []string{"foo=100", "bar=ZOO"}
	expected := map[string]string{"foo": "100", "bar": "ZOO"}
	actual, err := ParseContext(args)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestErrParseContext1(t *testing.T) {
	args := []string{"foo=100", "bar"}
	expected := "invalid context: bar"
	_, err := ParseContext(args)

	if err == nil || err.Error() != expected {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, err)
	}
}

func TestErrParseContext2(t *testing.T) {
	args := []string{"foo=100", "bar="}
	expected := "invalid context: bar="
	_, err := ParseContext(args)

	if err == nil || err.Error() != expected {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, err)
	}
}

func TestHasOption1(t *testing.T) {
	args := []string{"-a", "-b", "BBB", "-c", "CCC"}
	expectedArgs := []string{"-b", "BBB", "-c", "CCC"}
	expectedValue := true

	newAags, value := HasOption(args, "-a")

	if !reflect.DeepEqual(expectedArgs, newAags) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedArgs, newAags)
	}

	if expectedValue != value {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedValue, value)
	}
}

func TestHasOption2(t *testing.T) {
	args := []string{"-b", "BBB", "-c", "CCC"}
	expectedArgs := []string{"-b", "BBB", "-c", "CCC"}
	expectedValue := false

	newAags, value := HasOption(args, "-a")

	if !reflect.DeepEqual(expectedArgs, newAags) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedArgs, newAags)
	}

	if expectedValue != value {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedValue, value)
	}
}
