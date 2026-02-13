package data

import (
	"reflect"
	"testing"
)

// TestTrimPrefixMapsString tests the TrimPrefixMapsString function
func TestTrimPrefixMapsString(t *testing.T) {
	tests := []struct {
		name   string
		m      map[string]string
		prefix map[string]string
		want   map[string]string
	}{
		{
			name:   "Normal: Match all",
			m:      map[string]string{"a": "prefix_value_a", "b": "prefix_value_b"},
			prefix: map[string]string{"a": "prefix_", "b": "prefix_"},
			want:   map[string]string{"a": "value_a", "b": "value_b"},
		},
		{
			name:   "Part: Match part",
			m:      map[string]string{"a": "prefix_value_a", "b": "no_prefix_b"},
			prefix: map[string]string{"a": "prefix_"},
			want:   map[string]string{"a": "value_a", "b": "no_prefix_b"},
		},
		{
			name:   "No match: No match",
			m:      map[string]string{"a": "value_a", "b": "value_b"},
			prefix: map[string]string{"c": "prefix_"},
			want:   map[string]string{"a": "value_a", "b": "value_b"},
		},
		{
			name:   "Null map: Main map is null",
			m:      map[string]string{},
			prefix: map[string]string{"a": "prefix_"},
			want:   map[string]string{},
		},
		{
			name:   "Null map: Prefix map is null",
			m:      map[string]string{"a": "prefix_value_a"},
			prefix: map[string]string{},
			want:   map[string]string{"a": "prefix_value_a"},
		},
		{
			name:   "Edge: The Same",
			m:      map[string]string{"a": "prefix_"},
			prefix: map[string]string{"a": "prefix_"},
			want:   map[string]string{"a": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimPrefixMapsString(tt.m, tt.prefix)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrimPrefixMapsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTrimSuffixMapsString tests the TrimSuffixMapsString function
func TestTrimSuffixMapsString(t *testing.T) {
	tests := []struct {
		name   string
		m      map[string]string
		suffix map[string]string
		want   map[string]string
	}{
		{
			name:   "Normal: Match all",
			m:      map[string]string{"a": "value_a_suffix", "b": "value_b_suffix"},
			suffix: map[string]string{"a": "_suffix", "b": "_suffix"},
			want:   map[string]string{"a": "value_a", "b": "value_b"},
		},
		{
			name:   "Part: Match part",
			m:      map[string]string{"a": "value_a_suffix", "b": "value_b_no_suffix"},
			suffix: map[string]string{"a": "_suffix"},
			want:   map[string]string{"a": "value_a", "b": "value_b_no_suffix"},
		},
		{
			name:   "No match: No match",
			m:      map[string]string{"a": "value_a", "b": "value_b"},
			suffix: map[string]string{"c": "_suffix"},
			want:   map[string]string{"a": "value_a", "b": "value_b"},
		},
		{
			name:   "Null map: Main map is null",
			m:      map[string]string{},
			suffix: map[string]string{"a": "_suffix"},
			want:   map[string]string{},
		},
		{
			name:   "Null map: Prefix map is null",
			m:      map[string]string{"a": "value_a_suffix"},
			suffix: map[string]string{},
			want:   map[string]string{"a": "value_a_suffix"},
		},
		{
			name:   "Edge: The Same",
			m:      map[string]string{"a": "_suffix"},
			suffix: map[string]string{"a": "_suffix"},
			want:   map[string]string{"a": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimSuffixMapsString(tt.m, tt.suffix)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrimSuffixMapsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
