package exercise11

import (
	"reflect"
	"testing"

	"golang.org/x/tour/tree"
)

func TestWalk(t *testing.T) {
	testCases := []struct {
		name string
		tree *tree.Tree
		want []int
	}{
		{
			name: "Simple Tree",
			tree: &tree.Tree{
				Value: 1,
				Left: &tree.Tree{
					Value: 2,
				},
				Right: &tree.Tree{
					Value: 3,
				},
			},
			want: []int{2, 1, 3},
		},
		{
			name: "More Complex Tree",
			tree: &tree.Tree{
				Value: 4,
				Left: &tree.Tree{
					Value: 2,
					Left: &tree.Tree{
						Value: 1,
					},
					Right: &tree.Tree{
						Value: 3,
					},
				},
				Right: &tree.Tree{
					Value: 6,
					Left: &tree.Tree{
						Value: 5,
					},
					Right: &tree.Tree{
						Value: 7,
					},
				},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "Single Node Tree",
			tree: &tree.Tree{
				Value: 10,
			},
			want: []int{10},
		},
		{
			name: "Left Skewed Tree",
			tree: &tree.Tree{
				Value: 3,
				Left: &tree.Tree{
					Value: 2,
					Left: &tree.Tree{
						Value: 1,
					},
				},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Right Skewed Tree",
			tree: &tree.Tree{
				Value: 1,
				Right: &tree.Tree{
					Value: 2,
					Right: &tree.Tree{
						Value: 3,
					},
				},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Empty Tree (nil)",
			tree: nil,
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ch := make(chan int)
			go func() {
				Walk(tc.tree, ch)
				close(ch)
			}()

			var got []int
			for v := range ch {
				got = append(got, v)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Walk(%+v) = %v, want %v", tc.tree, got, tc.want)
			}
		})
	}
}

func TestSame(t *testing.T) {
	testCases := []struct {
		name string
		t1   *tree.Tree
		t2   *tree.Tree
		want bool
	}{
		{
			name: "Same Simple Trees",
			t1: &tree.Tree{
				Value: 1,
				Left: &tree.Tree{
					Value: 2,
				},
				Right: &tree.Tree{
					Value: 3,
				},
			},
			t2: &tree.Tree{
				Value: 1,
				Left: &tree.Tree{
					Value: 2,
				},
				Right: &tree.Tree{
					Value: 3,
				},
			},
			want: true,
		},
		{
			name: "Different Simple Trees",
			t1: &tree.Tree{
				Value: 1,
				Left: &tree.Tree{
					Value: 2,
				},
				Right: &tree.Tree{
					Value: 3,
				},
			},
			t2: &tree.Tree{
				Value: 1,
				Left: &tree.Tree{
					Value: 3,
				},
				Right: &tree.Tree{
					Value: 2,
				},
			},
			want: false,
		},
		{
			name: "Same Complex Trees",
			t1:   tree.New(4),
			t2:   tree.New(4),
			want: true,
		},
		{
			name: "Different Complex Trees",
			t1:   tree.New(4),
			t2:   tree.New(5),
			want: false,
		},
		{
			name: "One Empty Tree",
			t1: &tree.Tree{
				Value: 1,
			},
			t2:   nil,
			want: false,
		},
		{
			name: "Both Empty Trees",
			t1:   nil,
			t2:   nil,
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Same(tc.t1, tc.t2)
			if got != tc.want {
				t.Errorf("Same(%+v, %+v) = %v, want %v", tc.t1, tc.t2, got, tc.want)
			}
		})
	}
}
