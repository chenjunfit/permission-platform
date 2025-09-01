package bitring

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type step struct {
	event bool
	want  bool
}

func TestBitRing_isConditions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		size        int
		threshold   float64
		consecutive int
		steps       []step
	}{
		{
			name:        "永不触发",
			size:        8,
			threshold:   0.4,
			consecutive: 3,
			steps: []step{
				{false, false},
				{false, false},
				{false, false},
			},
		},
		{
			name:        "连续3次事件触发条件",
			size:        16,
			threshold:   1.0,
			consecutive: 3,
			steps: []step{
				{true, false},
				{true, false},
				{true, true}, // 连续3次事件
				{false, false},
			},
		},
		{
			name:        "连续5次事件触发条件（自定义阈值）",
			size:        10,
			threshold:   1.0,
			consecutive: 5,
			steps: []step{
				{true, false},
				{true, false},
				{true, false},
				{true, false},
				{true, true},
			},
		},
		{
			name:        "事件率超过阈值触发条件",
			size:        4,
			threshold:   0.5,
			consecutive: 3,
			steps: []step{
				{false, false}, // 0%
				{true, false},  // 50% (=阈值) 不触发
				{true, true},   // 66% (>阈值) 触发
			},
		},
		{
			name:        "环形覆盖后仍正确计数",
			size:        5,
			threshold:   0.4,
			consecutive: 3,
			steps: []step{
				{true, false},  // 1/1 =100% 触发
				{true, false},  // 2/2 =100% 触发
				{false, true},  // 2/3 ≈67%  触发
				{false, true},  // 2/4 =50%  触发
				{false, false}, // 2/5 =40%  不触发
				{false, false}, // 覆盖首位true，变成1/5=20% 不触发
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			br := NewBitRing(tc.size, tc.consecutive, tc.threshold)
			for i, st := range tc.steps {
				br.Add(st.event)
				got := br.isConditionMet()
				assert.Equalf(t, st.want, got, "步骤 %d 期望 %v 得到 %v", i, st.want, got)
			}
		})
	}
}

func TestBitRing_InternalCounter(t *testing.T) {
	t.Parallel()
	br := NewBitRing(3, 3, 0.6)

	assert.False(t, br.isConditionMet())

	br.Add(true)
	br.Add(false)
	br.Add(true) // filled
	assert.Equal(t, 2, br.eventCount)

	br.Add(false) // 覆盖 idx0 的 true
	assert.Equal(t, 1, br.eventCount)
	assert.False(t, br.isConditionMet())
}
func TestBitRing_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	br := NewBitRing(100, 3, 0.5)
	var wg sync.WaitGroup

	// 并发添加数据
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				br.Add(j%2 == 0)
			}
		}()
	}

	// 并发读取
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 30; j++ {
				_ = br.isConditionMet()
			}
		}()
	}

	wg.Wait()
	// 无需断言，只要不发生race条件或panic即为通过
}
func TestBitRing_EdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("无效参数处理", func(t *testing.T) {
		t.Parallel()
		br := NewBitRing(-1, -2, -0.5)
		assert.Equal(t, defaultSize, br.size, "应使用默认size")
		assert.Equal(t, defaultConsecutive, br.consecutive, "应使用默认consecutive")
		assert.Equal(t, 0.0, br.threshold, "阈值应限制为非负")

		br = NewBitRing(10, 20, 2.0)
		assert.Equal(t, 10, br.size, "应保留有效size")
		assert.Equal(t, 10, br.consecutive, "consecutive不应大于size")
		assert.Equal(t, 1.0, br.threshold, "阈值应限制为不超过1.0")
	})

	t.Run("空缓冲区处理", func(t *testing.T) {
		t.Parallel()
		br := NewBitRing(10, 3, 0.5)
		assert.False(t, br.isConditionMet(), "空缓冲区不应触发条件")
	})
}
