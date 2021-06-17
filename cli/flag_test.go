package cli

import (
	"reflect"
	"testing"
)

func TestFlag_Bool(t *testing.T) {
	type args struct {
		nonInteractive   bool
		forceInteractive bool
		isTerminal       bool
		includeDefault   bool
	}
	tests := []struct {
		name    string
		flag    Flag
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.flag.Bool(tt.args.nonInteractive, tt.args.forceInteractive, tt.args.isTerminal, tt.args.includeDefault)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Bool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlag_IsPresent(t *testing.T) {
	tests := []struct {
		name string
		flag Flag
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flag.IsPresent(); got != tt.want {
				t.Errorf("IsPresent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlag_IsPresentInArgument(t *testing.T) {
	tests := []struct {
		name string
		flag Flag
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flag.IsPresentInArgument(); got != tt.want {
				t.Errorf("IsPresentInArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlag_Set(t *testing.T) {
	type args struct {
		newValue interface{}
	}
	tests := []struct {
		name    string
		flag    Flag
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.flag.Set(tt.args.newValue); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFlag_String(t *testing.T) {
	type args struct {
		nonInteractive   bool
		forceInteractive bool
		isTerminal       bool
		includeDefault   bool
	}
	tests := []struct {
		name    string
		flag    Flag
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.flag.String(tt.args.nonInteractive, tt.args.forceInteractive, tt.args.isTerminal, tt.args.includeDefault)
			if (err != nil) != tt.wantErr {
				t.Errorf("String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("String() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlag_StringSlice(t *testing.T) {
	type args struct {
		nonInteractive   bool
		forceInteractive bool
		isTerminal       bool
		includeDefault   bool
	}
	tests := []struct {
		name    string
		flag    Flag
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.flag.StringSlice(tt.args.nonInteractive, tt.args.forceInteractive, tt.args.isTerminal, tt.args.includeDefault)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
