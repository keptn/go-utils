package http

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddEvent(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t2", "e3")

	assert.True(t, cache.Contains("t1", "e1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.False(t, cache.Contains("t1", "e3"))
	assert.True(t, cache.Contains("t2", "e3"))
}

func TestAddEventTwice(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t1", "e2")
	assert.Equal(t, 2, cache.Length("t1"))
	assert.Equal(t, 2, len(cache.Get("t1")))
}

func TestAddRemoveEvent(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t1", "e3")

	assert.Equal(t, 3, cache.Length("t1"))

	cache.Remove("t1", "e1")
	assert.Equal(t, 2, cache.Length("t1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.True(t, cache.Contains("t1", "e3"))

	cache.Remove("t1", "e3")
	assert.Equal(t, 1, cache.Length("t1"))
	assert.True(t, cache.Contains("t1", "e2"))
}

func TestKeep_NonExistingEvent(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t1", "e3")

	require.Equal(t, 3, cache.Length("t1"))
	cache.Keep("t1", []string{"e0"})
	assert.Equal(t, 3, cache.Length("t1"))
}

func TestKeep_WithDuplicates(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")

	require.Equal(t, 2, cache.Length("t1"))
	cache.Keep("t1", []string{"e2", "e2"})
	assert.Equal(t, 1, cache.Length("t1"))
}

func TestKeep_WithEmptyEvents(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")

	require.Equal(t, 2, cache.Length("t1"))
	cache.Keep("t1", []string{})
	assert.Equal(t, 0, cache.Length("t1"))
}

func TestKeep(t *testing.T) {
	cache := NewCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t2", "e3")
	cache.Add("t2", "e4")
	cache.Add("t2", "e5")

	cache.Keep("t1", []string{"e2"})
	cache.Keep("t2", []string{"e3", "e5"})

	assert.Equal(t, 1, cache.Length("t1"))
	assert.Equal(t, 2, cache.Length("t2"))
	assert.False(t, cache.Contains("t1", "e1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.True(t, cache.Contains("t2", "e3"))
	assert.False(t, cache.Contains("t2", "e4"))
	assert.True(t, cache.Contains("t2", "e5"))
}

func Test_subscriptionDiffer(t *testing.T) {
	type args struct {
		new []models.EventSubscription
		old []models.EventSubscription
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equal empty",
			args: args{
				new: []models.EventSubscription{},
				old: []models.EventSubscription{},
			},
			want: false,
		},
		{
			name: "equal",
			args: args{
				new: []models.EventSubscription{{Event: "e1"}},
				old: []models.EventSubscription{{Event: "e1"}},
			},
			want: false,
		},
		{
			name: "differ1",
			args: args{
				new: []models.EventSubscription{{Event: "e1"}},
				old: []models.EventSubscription{{Event: "e2"}},
			},
			want: true,
		},
		{
			name: "differ2",
			args: args{
				new: []models.EventSubscription{{Event: "e2"}},
				old: []models.EventSubscription{{Event: "e1"}},
			},
			want: true,
		},
		{
			name: "differ3",
			args: args{
				new: []models.EventSubscription{{Event: "e2"}, {Event: "e1"}},
				old: []models.EventSubscription{{Event: "e1"}},
			},
			want: true,
		},
		{
			name: "differ4",
			args: args{
				new: []models.EventSubscription{},
				old: []models.EventSubscription{{Event: "e1"}},
			},
			want: true,
		},
		{
			name: "differ5",
			args: args{
				new: []models.EventSubscription{{Event: "e1"}},
				old: []models.EventSubscription{},
			},
			want: true,
		},
		{
			name: "differ6",
			args: args{
				new: []models.EventSubscription{{Event: "e2"}},
				old: []models.EventSubscription{{Event: "e1"}, {Event: "e1"}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, subscriptionDiffer(tt.args.new, tt.args.old), "subscriptionDiffer(%v, %v)", tt.args.new, tt.args.old)
		})
	}
}
