package abac

import (
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository"
	"github.com/permission-dev/internal/repository/dao"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenDomainPolicyRules(t *testing.T) {
	t.Parallel()
	now := time.Now().UnixMilli()
	tests := []struct {
		name   string
		rules  []dao.PolicyRule
		expect []domain.PolicyRule
	}{
		{
			name:   "empty rules",
			rules:  []dao.PolicyRule{},
			expect: []domain.PolicyRule{},
		},
		{
			name: "single rule",
			rules: []dao.PolicyRule{
				{
					ID:        1,
					AttrDefID: 100,
					Value:     "test",
					Operator:  "eq",
					Left:      0,
					Right:     0,
					Ctime:     now,
					Utime:     now,
				},
			},
			expect: []domain.PolicyRule{
				{
					ID:        1,
					AttrDef:   domain.AttributeDefinition{ID: 100},
					Value:     "test",
					LeftRule:  nil,
					RightRule: nil,
					Operator:  "eq",
					Ctime:     now,
					Utime:     now,
				},
			},
		},
		{
			name: "nested rules",
			rules: []dao.PolicyRule{
				{
					ID:        1,
					AttrDefID: 100,
					Value:     "root",
					Operator:  "and",
					Left:      2,
					Right:     3,
					Ctime:     now,
					Utime:     now,
				},
				{
					ID:        2,
					AttrDefID: 101,
					Value:     "left",
					Operator:  "neq",
					Left:      0,
					Right:     0,
					Ctime:     now,
					Utime:     now,
				},
				{
					ID:        3,
					AttrDefID: 102,
					Value:     "right",
					Operator:  "gt",
					Left:      0,
					Right:     0,
					Ctime:     now,
					Utime:     now,
				},
			},
			expect: []domain.PolicyRule{
				{
					ID: 1,
					AttrDef: domain.AttributeDefinition{
						ID: 100,
					},
					Value:    "root",
					Operator: "and",
					Ctime:    now,
					Utime:    now,
					LeftRule: &domain.PolicyRule{
						ID: 2,
						AttrDef: domain.AttributeDefinition{
							ID: 101,
						},
						Value:    "left",
						Operator: "neq",
						Ctime:    now,
						Utime:    now,
					},
					RightRule: &domain.PolicyRule{
						ID: 3,
						AttrDef: domain.AttributeDefinition{
							ID: 102,
						},
						Value:    "right",
						Operator: "gt",
						Ctime:    now,
						Utime:    now,
					},
				},
			},
		},
		{
			name: "multiple root rules",
			rules: []dao.PolicyRule{
				{
					ID:        1,
					AttrDefID: 100,
					Value:     "root1",
					Operator:  "eq",
					Left:      0,
					Right:     0,
					Ctime:     now,
					Utime:     now,
				},
				{
					ID:        2,
					AttrDefID: 101,
					Value:     "root2",
					Operator:  "neq",
					Left:      0,
					Right:     0,
					Ctime:     now,
					Utime:     now,
				},
			},
			expect: []domain.PolicyRule{
				{
					ID: 1,
					AttrDef: domain.AttributeDefinition{
						ID: 100,
					},
					Value:    "root1",
					Operator: "eq",
					Ctime:    now,
					Utime:    now,
				},
				{
					ID: 2,
					AttrDef: domain.AttributeDefinition{
						ID: 101,
					},
					Value:    "root2",
					Operator: "neq",
					Ctime:    now,
					Utime:    now,
				},
			},
		},
		{
			name: "complex nested rules with multiple levels",
			rules: []dao.PolicyRule{
				{
					ID:        1,
					AttrDefID: 100,
					Value:     "root",
					Operator:  "and",
					Left:      2,
					Right:     3,
					Ctime:     now,
					Utime:     now,
				},
				{
					ID:        2,
					AttrDefID: 101,
					Value:     "left",
					Operator:  "or",
					Left:      4,
					Right:     5,
					Ctime:     now,
					Utime:     now,
				},
				{
					ID:        3,
					AttrDefID: 102,
					Value:     "right",
					Operator:  "and",
					Left:      6,
					Right:     7,
					Ctime:     now,
					Utime:     now,
				},
				{
					ID:        4,
					AttrDefID: 103,
					Value:     "left-left",
					Operator:  "eq",
					Left:      0,
					Right:     0,
					Ctime:     now,
					Utime:     now,
				},
				{
					ID:        5,
					AttrDefID: 104,
					Value:     "left-right",
					Operator:  "neq",
					Left:      0,
					Right:     0,
					Ctime:     now,
					Utime:     now,
				},
				{
					ID:        6,
					AttrDefID: 105,
					Value:     "right-left",
					Operator:  "gt",
					Left:      0,
					Right:     0,
					Ctime:     now,
					Utime:     now,
				},
				{
					ID:        7,
					AttrDefID: 106,
					Value:     "right-right",
					Operator:  "lt",
					Left:      0,
					Right:     0,
					Ctime:     now,
					Utime:     now,
				},
			},
			expect: []domain.PolicyRule{
				{
					ID: 1,
					AttrDef: domain.AttributeDefinition{
						ID: 100,
					},
					Value:    "root",
					Operator: "and",
					Ctime:    now,
					Utime:    now,
					LeftRule: &domain.PolicyRule{
						ID: 2,
						AttrDef: domain.AttributeDefinition{
							ID: 101,
						},
						Value:    "left",
						Operator: "or",
						Ctime:    now,
						Utime:    now,
						LeftRule: &domain.PolicyRule{
							ID: 4,
							AttrDef: domain.AttributeDefinition{
								ID: 103,
							},
							Value:    "left-left",
							Operator: "eq",
							Ctime:    now,
							Utime:    now,
						},
						RightRule: &domain.PolicyRule{
							ID: 5,
							AttrDef: domain.AttributeDefinition{
								ID: 104,
							},
							Value:    "left-right",
							Operator: "neq",
							Ctime:    now,
							Utime:    now,
						},
					},
					RightRule: &domain.PolicyRule{
						ID: 3,
						AttrDef: domain.AttributeDefinition{
							ID: 102,
						},
						Value:    "right",
						Operator: "and",
						Ctime:    now,
						Utime:    now,
						LeftRule: &domain.PolicyRule{
							ID: 6,
							AttrDef: domain.AttributeDefinition{
								ID: 105,
							},
							Value:    "right-left",
							Operator: "gt",
							Ctime:    now,
							Utime:    now,
						},
						RightRule: &domain.PolicyRule{
							ID: 7,
							AttrDef: domain.AttributeDefinition{
								ID: 106,
							},
							Value:    "right-right",
							Operator: "lt",
							Ctime:    now,
							Utime:    now,
						},
					},
				},
			},
		},
	}
	for idx := range tests {
		tt := tests[idx]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := repository.GenDomainPolicyRules(tt.rules)
			assert.Equal(t, tt.expect, result)
		})
	}
}
