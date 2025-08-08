package pcommon_test

import (
	"iter"
	"maps"
	"slices"
	"testing"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/xpdata"
)

type tb_t struct {
	t *testing.T
	b *testing.B
}

func (tb tb_t) Run(name string, run func(r tb_t)) {
	if tb.t != nil {
		tb.t.Run(name, func(t *testing.T) {
			run(tb_t{t: t})
		})
	}
	if tb.b != nil {
		tb.b.Run(name, func(b *testing.B) {
			run(tb_t{b: b})
		})
	}
}

type maplike interface {
	All() iter.Seq2[string, pcommon.Value]
}

func run[T maplike](tb tb_t, data [][]string, accumulate func(keys []string) T) {
	if tb.b != nil {
		i := 0
		tb.b.ResetTimer()
		for tb.b.Loop() {
			row := slices.Clone(data[i%len(data)])
			i++
			_ = accumulate(row)
		}
	}

	if tb.t != nil {
		for _, keys := range data {
			realKeyCount := map[string]int{}
			for _, k := range keys {
				realKeyCount[k]++
			}
			m := accumulate(keys)
			keyCount := map[string]int{}
			for k, v := range m.All() {
				count := 1
				if v.Type() == pcommon.ValueTypeSlice {
					count = v.Slice().Len()
				}
				if _, ok := keyCount[k]; ok {
					tb.t.Fatalf("not duplicated")
				}
				keyCount[k] = count
			}
			if !maps.Equal(realKeyCount, keyCount) {
				tb.t.Fatalf("wrong number of values for key")
			}
		}
	}
}

func run_b[T maplike](b *testing.B, data [][]string, accumulate func(keys []string) T) {
	run(tb_t{b: b}, data, accumulate)
}

func runWithDuplicates(tb tb_t, data [][]string) {
	// Using Map.Put to insert data
	tb.Run("Put", func(tb tb_t) {

		// Using Map.Get to detect duplicates
		tb.Run("Get", func(tb tb_t) {
			run(tb, data, func(keys []string) pcommon.Map {
				m := pcommon.NewMap()
				m.EnsureCapacity(len(keys))
				for _, key := range keys {
					if val, found := m.Get(key); found {
						var slice pcommon.Slice
						switch val.Type() {
						case pcommon.ValueTypeStr:
							str := val.Str()
							slice = val.SetEmptySlice()
							slice.AppendEmpty().SetStr(str)
						case pcommon.ValueTypeSlice:
							slice = val.Slice()
						default:
							panic("unreachable")
						}
						slice.AppendEmpty().SetStr(BENCH_VAL)
						continue
					}
					m.PutStr(key, BENCH_VAL)
				}
				return m
			})
		})

		// Using slices.Sort and neighbor comparison to detect duplicates
		tb.Run("Sort", func(tb tb_t) {
			run(tb, data, func(keys []string) pcommon.Map {
				m := pcommon.NewMap()
				m.EnsureCapacity(len(keys))

				slices.Sort(keys)

				sameKey := false
				var curSlice pcommon.Slice
				for i, key := range keys {
					nextSameKey := i+1 < len(keys) && key == keys[i+1]
					if sameKey { // continue slice
						curSlice.AppendEmpty().SetStr(BENCH_VAL)
					} else if nextSameKey { // start new slice
						curSlice = m.PutEmptySlice(key)
						curSlice.EnsureCapacity(3)
						curSlice.AppendEmpty().SetStr(BENCH_VAL)
					} else { // not a duplicate key
						m.PutStr(key, BENCH_VAL)
					}
					sameKey = nextSameKey
				}

				return m
			})
		})
	})

	// Using Map.FromRaw to insert data, and map accesses to detect duplicates
	tb.Run("FromRaw", func(tb tb_t) {

		// Making a new map and new slices every run
		tb.Run("make", func(tb tb_t) {
			run(tb, data, func(keys []string) pcommon.Map {
				m := make(map[string]any, len(keys))
				for _, key := range keys {
					if val, found := m[key]; found {
						var vals []any
						switch val := val.(type) {
						case string:
							vals = make([]any, 0, 2)
							vals = append(vals, val)
						case []any:
							vals = val
						default:
							panic("unreachable")
						}
						m[key] = append(vals, BENCH_VAL)
						continue
					}
					m[key] = BENCH_VAL
				}
				m2 := pcommon.NewMap()
				m2.FromRaw(m)
				return m2
			})
		})

		// Reusing map and slice allocations every run
		tb.Run("reuse", func(tb tb_t) {
			var reusedMap = make(map[string]any, 200)
			var reusedSlices = func() [][]any {
				slices := make([][]any, 0, 100)
				for range 100 {
					slices = append(slices, make([]any, 0, 100))
				}
				return slices
			}()

			run(tb, data, func(keys []string) pcommon.Map {
				clear(reusedMap)

				nextSlice := 0
				for _, key := range keys {
					if val, found := reusedMap[key]; found {
						var vals []any
						switch val := val.(type) {
						case string:
							vals = reusedSlices[nextSlice]
							nextSlice++
							vals = append(vals, val)
						case []any:
							vals = val
						default:
							panic("unreachable")
						}
						reusedMap[key] = append(vals, BENCH_VAL)
						continue
					}
					reusedMap[key] = BENCH_VAL
				}

				m2 := pcommon.NewMap()
				m2.FromRaw(reusedMap)

				return m2
			})
		})
	})

	// Using a new Map implementation which keeps its keys in sorted order
	// and uses binary search for Put and Get
	tb.Run("SortedMap", func(tb tb_t) {
		tb.Run("Put", func(tb tb_t) {
			tb.Run("Get", func(tb tb_t) {
				run(tb, data, func(keys []string) pcommon.SortedMap {
					m := pcommon.NewSortedMap()
					m.EnsureCapacity(len(keys))
					for _, key := range keys {
						if val, found := m.Get(key); found {
							var slice pcommon.Slice
							switch val.Type() {
							case pcommon.ValueTypeStr:
								str := val.Str()
								slice = val.SetEmptySlice()
								slice.AppendEmpty().SetStr(str)
							case pcommon.ValueTypeSlice:
								slice = val.Slice()
							default:
								panic("unreachable")
							}
							slice.AppendEmpty().SetStr(BENCH_VAL)
							continue
						}
						m.PutStr(key, BENCH_VAL)
					}
					return m
				})
			})
		})
	})

	// Using a new Map implementation which is backed by a Go map
	tb.Run("MapMap", func(tb tb_t) {
		tb.Run("Put", func(tb tb_t) {
			tb.Run("Get", func(tb tb_t) {
				run(tb, data, func(keys []string) pcommon.MapMap {
					m := pcommon.NewMapMap()
					m.EnsureCapacity(len(keys))
					for _, key := range keys {
						if val, found := m.Get(key); found {
							var slice pcommon.Slice
							switch val.Type() {
							case pcommon.ValueTypeStr:
								str := val.Str()
								slice = val.SetEmptySlice()
								slice.AppendEmpty().SetStr(str)
							case pcommon.ValueTypeSlice:
								slice = val.Slice()
							default:
								panic("unreachable")
							}
							slice.AppendEmpty().SetStr(BENCH_VAL)
							continue
						}
						m.PutStr(key, BENCH_VAL)
					}
					return m
				})
			})
		})
	})

	// Using a new Map.PutUnsafe method to insert data
	tb.Run("PutUnsafe", func(tb tb_t) {
		tb.Run("Sort", func(tb tb_t) {
			run(tb, data, func(keys []string) pcommon.Map {
				m := pcommon.NewMap()
				m.EnsureCapacity(len(keys))

				slices.Sort(keys)

				sameKey := false
				var curSlice pcommon.Slice
				for i, key := range keys {
					nextSameKey := i+1 < len(keys) && key == keys[i+1]
					if sameKey { // continue slice
						curSlice.AppendEmpty().SetStr(BENCH_VAL)
					} else if nextSameKey { // start new slice
						curSlice = m.PutEmptyUnsafe(key).SetEmptySlice()
						curSlice.EnsureCapacity(3)
						curSlice.AppendEmpty().SetStr(BENCH_VAL)
					} else { // not a duplicate key
						m.PutEmptyUnsafe(key).SetStr(BENCH_VAL)
					}
					sameKey = nextSameKey
				}

				return m
			})
		})
	})

	// Using a new xpdata.MapBuilder type to insert data
	tb.Run("MapBuilder", func(tb tb_t) {

		// Using slices.Sort and neighbor comparison to detect duplicates
		// and MapBuilder.UnsafeIntoMap to convert to a map with no checks
		tb.Run("Sort UnsafeIntoMap", func(tb tb_t) {
			run(tb, data, func(keys []string) pcommon.Map {
				var mb xpdata.MapBuilder
				mb.EnsureCapacity(len(keys))

				slices.Sort(keys)

				sameKey := false
				var curSlice pcommon.Slice
				for i, key := range keys {
					nextSameKey := i+1 < len(keys) && key == keys[i+1]
					if sameKey { // continue slice
						curSlice.AppendEmpty().SetStr(BENCH_VAL)
					} else if nextSameKey { // start new slice
						curSlice = mb.AppendEmpty(key).SetEmptySlice()
						curSlice.EnsureCapacity(3)
						curSlice.AppendEmpty().SetStr(BENCH_VAL)
					} else { // not a duplicate key
						mb.AppendEmpty(key).SetStr(BENCH_VAL)
					}
					sameKey = nextSameKey
				}

				m := pcommon.NewMap()
				mb.UnsafeIntoMap(m)
				return m
			})
		})

		// Using slices.Sort and neighbor comparison to detect duplicates
		// and MapBuilder.SortedIntoMap to check that the data is sorted
		tb.Run("Sort SortedIntoMap", func(tb tb_t) {
			run(tb, data, func(keys []string) pcommon.Map {
				var mb xpdata.MapBuilder
				mb.EnsureCapacity(len(keys))

				slices.Sort(keys)

				sameKey := false
				var curSlice pcommon.Slice
				for i, key := range keys {
					nextSameKey := i+1 < len(keys) && key == keys[i+1]
					if sameKey { // continue slice
						curSlice.AppendEmpty().SetStr(BENCH_VAL)
					} else if nextSameKey { // start new slice
						curSlice = mb.AppendEmpty(key).SetEmptySlice()
						curSlice.EnsureCapacity(3)
						curSlice.AppendEmpty().SetStr(BENCH_VAL)
					} else { // not a duplicate key
						mb.AppendEmpty(key).SetStr(BENCH_VAL)
					}
					sameKey = nextSameKey
				}

				m := pcommon.NewMap()
				mb.SortedIntoMap(m)
				return m
			})
		})

		// Using MapBuilder.MergeIntoMap to detect duplicates
		tb.Run("MergeIntoMap", func(tb tb_t) {
			run(tb, data, func(keys []string) pcommon.Map {
				var mb xpdata.MapBuilder
				mb.EnsureCapacity(len(keys))
				for _, key := range keys {
					mb.AppendEmpty(key).SetStr(BENCH_VAL)
				}
				m := pcommon.NewMap()
				mb.MergeIntoMap(m, func(vals []pcommon.Value) {
					first := vals[0].Str()
					slice := vals[0].SetEmptySlice()
					slice.EnsureCapacity(len(vals))
					slice.AppendEmpty().SetStr(first)
					for i := 1; i < len(vals); i++ {
						slice.AppendEmpty().SetStr(vals[i].Str())
					}
				})
				return m
			})
		})
	})
}

func BenchmarkMapExpRealistic(b *testing.B) {
	data := generateRealisticBenchData()
	runWithDuplicates(tb_t{b: b}, data)
}

func TestMapExpRealistic(t *testing.T) {
	data := generateRealisticBenchData()
	runWithDuplicates(tb_t{t: t}, data)
}

func BenchmarkMapExpWorstCase(b *testing.B) {
	data := generateWorstCaseBenchData()
	runWithDuplicates(tb_t{b: b}, data)
}

func BenchmarkMapExpNoDuplicates(b *testing.B) {
	data := generateDeduplicatedBenchData()

	b.Run("Put", func(b *testing.B) {
		run_b(b, data, func(keys []string) pcommon.Map {
			m := pcommon.NewMap()
			m.EnsureCapacity(len(keys))
			for _, key := range keys {
				m.PutStr(key, BENCH_VAL)
			}
			return m
		})
	})

	b.Run("FromRaw", func(b *testing.B) {

		b.Run("make", func(b *testing.B) {
			run_b(b, data, func(keys []string) pcommon.Map {
				m := make(map[string]any, len(keys))
				for _, key := range keys {
					m[key] = BENCH_VAL
				}
				m2 := pcommon.NewMap()
				m2.FromRaw(m)
				return m2
			})
		})

		b.Run("reuse", func(b *testing.B) {
			var reusedMap = make(map[string]any, 200)

			run_b(b, data, func(keys []string) pcommon.Map {
				clear(reusedMap)
				for _, key := range keys {
					reusedMap[key] = BENCH_VAL
				}
				m := pcommon.NewMap()
				m.FromRaw(reusedMap)
				return m
			})
		})
	})

	b.Run("PutUnsafe", func(b *testing.B) {
		run_b(b, data, func(keys []string) pcommon.Map {
			m := pcommon.NewMap()
			m.EnsureCapacity(len(keys))
			for _, key := range keys {
				m.PutEmptyUnsafe(key).SetStr(BENCH_VAL)
			}
			return m
		})
	})

	b.Run("MapBuilder", func(b *testing.B) {

		b.Run("UnsafeIntoMap", func(b *testing.B) {
			run_b(b, data, func(keys []string) pcommon.Map {
				m := pcommon.NewMap()
				var mb xpdata.MapBuilder
				mb.EnsureCapacity(len(keys))
				for _, key := range keys {
					mb.AppendEmpty(key).SetStr(BENCH_VAL)
				}
				mb.UnsafeIntoMap(m)
				return m
			})
		})

		b.Run("Sort SortedIntoMap", func(b *testing.B) {
			run_b(b, data, func(keys []string) pcommon.Map {
				m := pcommon.NewMap()
				var mb xpdata.MapBuilder
				mb.EnsureCapacity(len(keys))
				slices.Sort(keys)
				for _, key := range keys {
					mb.AppendEmpty(key).SetStr(BENCH_VAL)
				}
				mb.SortedIntoMap(m)
				return m
			})
		})

		b.Run("DistinctIntoMap", func(b *testing.B) {
			run_b(b, data, func(keys []string) pcommon.Map {
				m := pcommon.NewMap()
				var mb xpdata.MapBuilder
				mb.EnsureCapacity(len(keys))
				for _, key := range keys {
					mb.AppendEmpty(key).SetStr(BENCH_VAL)
				}
				mb.DistinctIntoMap(m)
				return m
			})
		})
	})
}
