package record

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/impl/data/record/dumb"
	mongorec "github.com/ischenkx/kantoku/pkg/impl/data/record/mongo"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"math/rand"
	"sort"
	"testing"
)

func newMongoRecords(ctx context.Context) record.Storage {
	// Set connection configurations
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to the MongoDB server
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-ctx.Done()

		if err = client.Disconnect(ctx); err != nil {
			log.Println("failed to disconnect from mongodb:", err)
		}
	}()

	if err := client.Ping(ctx, readpref.Nearest()); err != nil {
		log.Fatal(err)
	}

	// Access the "main" database
	database := client.Database("testing")

	// Access the "testing" collection
	collection := database.Collection("records")

	return mongorec.New(collection)
}

var Keys = []string{
	"updated_at",
	"id",
	"name",
	"surname",
	"group",
}

var Groups = []string{
	"a",
	"b",
	"c",
	"d",
	"e",
}

const MinimumTestCases = 5
const MaximumTestCases = 15
const TotalRecords = 100

func HandTest(t *testing.T, impl record.Storage) {
	ctx := context.Background()
	t.Log("Dataset initialization...")
	t.Log("Generating records...")
	raw := []record.R{
		{
			"updated_at": 1,
			"id":         "1",
			"name":       "A",
			"surname":    "A",
			"group":      "a",
		},
		{
			"updated_at": 2,
			"id":         "2",
			"name":       "B",
			"surname":    "B",
			"group":      "a",
		},
		{
			"updated_at": 2,
			"id":         "3",
			"name":       "C",
			"surname":    "B",
			"group":      "b",
		},
		{
			"updated_at": 1,
			"id":         "4",
			"name":       "D",
			"surname":    "A",
			"group":      "c",
		},
		{
			"updated_at": 5,
			"id":         "5",
			"name":       "A",
			"surname":    "T",
			"group":      "c",
		},
		{
			"updated_at": 5,
			"id":         "6",
			"name":       "D",
			"surname":    "E",
			"group":      "b",
		},
		{
			"updated_at": 1,
			"id":         "7",
			"name":       "V",
			"surname":    "R",
			"group":      "a",
		},
	}

	for _, rec := range raw {
		_ = impl.Insert(ctx, rec)
	}

	t.Log("Generation successfully finished!")

	t.Run("Integrity", func(t *testing.T) {
		data, err := collect(ctx, impl.Cursor().Iter())
		assert.NoError(t, err)

		sort.Sort(recordSorter(data))
		sort.Sort(recordSorter(raw))

		assertDatasetsEqual(t, data, raw)
	})

	t.Run("Filters", func(t *testing.T) {
		t.Run("no filter", func(t *testing.T) {
			iter := impl.Filter().Filter().Filter().Cursor().Iter()
			data, err := collect(ctx, iter)
			assert.NoError(t, err)
			sorted(data, raw)

			assertDatasetsEqual(t, data, raw)
		})
		t.Run("group:a", func(t *testing.T) {
			iter := impl.Filter(record.E{"group", "a"}).Cursor().Iter()
			data, err := collect(ctx, iter)
			assert.NoError(t, err)

			filteredRaw := lo.Filter(raw, func(item record.R, _ int) bool {
				return item["group"] == "a"
			})
			sorted(data, filteredRaw)

			assertDatasetsEqual(t, data, filteredRaw)
		})
		t.Run("id:3", func(t *testing.T) {
			iter := impl.Filter(record.E{"id", "3"}).Cursor().Iter()
			data, err := collect(ctx, iter)
			assert.NoError(t, err)

			filteredRaw := lo.Filter(raw, func(item record.R, _ int) bool {
				return item["id"] == "3"
			})
			sorted(data, filteredRaw)

			assertDatasetsEqual(t, data, filteredRaw)
		})
		t.Run("id:100500 (not existing)", func(t *testing.T) {
			iter := impl.Filter(record.E{"id", "100500"}).Cursor().Iter()
			data, err := collect(ctx, iter)
			assert.NoError(t, err)

			filteredRaw := lo.Filter(raw, func(item record.R, _ int) bool {
				return item["id"] == "100500"
			})
			sorted(data, filteredRaw)

			assertDatasetsEqual(t, data, filteredRaw)
		})
		t.Run("updated_at:3", func(t *testing.T) {
			iter := impl.Filter(record.E{"updated_at", 3}).Cursor().Iter()
			data, err := collect(ctx, iter)
			assert.NoError(t, err)

			filteredRaw := lo.Filter(raw, func(item record.R, _ int) bool {
				return item["updated_at"] == 3
			})
			sorted(data, filteredRaw)

			assertDatasetsEqual(t, data, filteredRaw)
		})
	})

	t.Run("Inserts", func(t *testing.T) {
		rec := record.R{
			"updated_at": 100500,
			"id":         "42",
			"name":       "M",
			"surname":    "F",
			"group":      "y",
		}
		raw = append(raw, rec)
		err := impl.Insert(ctx, rec)
		assert.NoError(t, err)

		data, err := collect(ctx, impl.Cursor().Iter())
		assert.NoError(t, err)

		sorted(data, raw)
		assertDatasetsEqual(t, data, raw)
	})

	t.Run("Updates", func(t *testing.T) {
		rec := record.R{
			"updated_at": 100600,
			"id":         "24",
		}
		err := impl.Filter(record.E{"id", "42"}).Update(ctx, rec, nil)
		assert.NoError(t, err)

		data, err := collect(ctx, impl.Cursor().Iter())
		assert.NoError(t, err)

		for _, r := range raw {
			if r["id"] == "42" {
				for key, value := range rec {
					r[key] = value
				}
			}
		}

		sorted(data, raw)
		assertDatasetsEqual(t, data, raw)
	})
}

func TestImplementations(t *testing.T) {
	t.Run("Specimen hand-tests", func(t *testing.T) {
		HandTest(t, dumb.New())
	})

	rand.Seed(42)
	ctx := context.Background()
	implementations := map[string]func(ctx context.Context) record.Storage{
		"mongo": newMongoRecords,
	}

	for label, makeImpl := range implementations {
		t.Run(label, func(t *testing.T) {
			impl := makeImpl(ctx)
			assert.NoError(t, impl.Erase(ctx))

			RunImplementationTests(ctx, t, impl, dumb.New())
		})
	}
}

func RunImplementationTests(ctx context.Context, t *testing.T, impl record.Storage, specimen *dumb.Storage) {
	t.Run("Hand-tests", func(t *testing.T) {
		HandTest(t, impl)
	})

	t.Log("Dataset initialization...")
	t.Log("Erasing all data...")
	if err := impl.Erase(ctx); err != nil {
		t.Fatal("failed to erase:", err)
	}

	t.Log("Generating random records...")
	var raw []record.R
	for i := 0; i < TotalRecords; i++ {
		rec := randomRecord()
		_ = specimen.Insert(ctx, rec)
		raw = append(raw, rec)
		if err := impl.Insert(ctx, rec); err != nil {
			t.Fatal("failed to insert:", err)
		}
	}
	t.Log("Generation successfully finished!")

	t.Run("Automated", func(t *testing.T) {
		t.Run("Read operations", func(t *testing.T) {
			t.Run("Not filtered", func(t *testing.T) {
				RunReadOperationsTests(ctx, t, impl, specimen)
			})

			t.Run("Filtered", func(t *testing.T) {
				t.Run("Single filter", func(t *testing.T) {
					RunFilterTests(ctx, t, [][][]record.Entry{
						{
							{{"group", "a"}},
						},
					}, impl, specimen)
				})

				t.Run("Multiple filters (one term)", func(t *testing.T) {
					RunFilterTests(ctx, t, [][][]record.Entry{
						{
							lo.Map(Groups, func(group string, _ int) record.E {
								return record.E{"group", group}
							}),
						},
					}, impl, specimen)
				})

				t.Run("Multiple filters (multiple terms)", func(t *testing.T) {

					var groupFilters []record.E
					var nameFilters []record.E

					for i := 0; i < 5; i++ {
						groupFilters = append(groupFilters,
							record.E{"group", lo.Sample(Groups)})
					}

					for i := 0; i < 100; i++ {
						val := specimen.Sample()["name"]
						nameFilters = append(nameFilters, record.E{"name", val})
					}

					RunFilterTests(ctx, t, [][][]record.Entry{
						{
							groupFilters,
							nameFilters,
						},
					}, impl, specimen)
				})

				t.Run("Not existing property", func(t *testing.T) {
					RunFilterTests(ctx, t, [][][]record.Entry{
						{
							{{"x", "1"}, {"y", "2"}},
							{{"x", "4"}, {"y", "3"}},
						},
					}, impl, specimen)
				})

				t.Run("Random", func(t *testing.T) {
					var cases [][][]record.E

					testCases := randomTestCases()
					for i := 0; i < testCases; i++ {
						var workingSet []record.R
						for j := 0; j < TotalRecords/5; j++ {
							workingSet = append(workingSet, specimen.Sample())
						}

						var _case [][]record.E
						termsCount := 1 + rand.Intn(10)
						for j := 0; j < termsCount; j++ {
							var term []record.E
							for k := 0; k < TotalRecords/8; k++ {
								key := lo.Sample(Keys)
								rec := lo.Sample(workingSet)

								term = append(term, record.E{key, rec[key]})
							}
							_case = append(_case, term)
						}
						cases = append(cases, _case)
					}

					RunFilterTests(ctx, t, cases, impl, specimen)
				})
			})
		})

		t.Run("Write operations", func(t *testing.T) {
			t.Run("Insert", func(t *testing.T) {
				t.Run("Empty", func(t *testing.T) {
					_ = specimen.Insert(ctx, record.R{})
					if err := impl.Insert(ctx, record.R{}); err != nil {
						t.Fatal("failed to insert:", err)
					}
					RunReadOperationsTests(ctx, t, impl, specimen)
				})
				t.Run("Random", func(t *testing.T) {
					testCases := randomTestCases()
					for i := 0; i < testCases; i++ {
						rec := randomRecord()
						_ = specimen.Insert(ctx, rec)
						if err := impl.Insert(ctx, rec); err != nil {
							t.Fatal("failed to insert:", err)
						}
						RunReadOperationsTests(ctx, t, impl, specimen)
					}
				})

			})
			t.Run("Update", func(t *testing.T) {
				t.Run("By ID", func(t *testing.T) {
					testCases := randomTestCases()
					for i := 0; i < testCases; i++ {
						sample := specimen.Sample()
						key := lo.Sample(Keys)
						var newValue any
						switch key {
						case "group":
							newValue = lo.Sample(Groups)
						case "updated_at":
							newValue = rand.Intn(1 << 20)
						default:
							newValue = uuid.New().String()
						}
						_ = specimen.Filter(record.E{"id", sample["id"]}).Update(ctx, record.R{key: newValue}, nil)
						err := impl.Filter(record.E{"id", sample["id"]}).Update(ctx, record.R{key: newValue}, nil)
						assert.NoError(t, err)

						RunReadOperationsTests(ctx, t, impl, specimen)
					}
				})

				t.Run("By Group", func(t *testing.T) {
					testCases := randomTestCases()
					for i := 0; i < testCases; i++ {
						sample := specimen.Sample()
						key := lo.Sample(Keys)
						if key == "group" {
							key = "id"
						}

						var newValue any
						switch key {
						case "updated_at":
							newValue = rand.Intn(1 << 20)
						default:
							newValue = uuid.New().String()
						}
						_ = specimen.Filter(record.E{"group", sample["group"]}).Update(ctx, record.R{key: newValue}, nil)
						err := impl.Filter(record.E{"group", sample["group"]}).Update(ctx, record.R{key: newValue}, nil)
						assert.NoError(t, err)

						RunReadOperationsTests(ctx, t, impl, specimen)
					}
				})

				t.Run("Upsert", func(t *testing.T) {
					t.Run("Existing value", func(t *testing.T) {
						testCases := randomTestCases()
						for i := 0; i < testCases; i++ {
							rec := specimen.Sample()
							upsert := rec.Copy()
							upsert["surname"] = uuid.New().String()

							err := impl.
								Filter(record.E{"id", rec["id"]}).
								Update(ctx, rec, upsert)
							assert.NoError(t, err)

							_ = specimen.
								Filter(record.E{"id", rec["id"]}).
								Update(ctx, rec, upsert)

							RunReadOperationsTests(ctx, t, impl, specimen)
						}
					})

					t.Run("Not existing value", func(t *testing.T) {
						testCases := randomTestCases()
						for i := 0; i < testCases; i++ {
							rec := randomRecord()
							upsert := rec.Copy()
							upsert["surname"] = uuid.New().String()

							err := impl.
								Filter(record.E{"id", rec["id"]}).
								Update(ctx, rec, upsert)
							assert.NoError(t, err)

							_ = specimen.
								Filter(record.E{"id", rec["id"]}).
								Update(ctx, rec, upsert)

							RunReadOperationsTests(ctx, t, impl, specimen)
						}
					})
				})
				// TODO: implement random

			})
			t.Run("Erase", func(t *testing.T) {
				testCases := randomTestCases()
				for i := 0; i < testCases; i++ {
					sample := specimen.Sample()
					key := lo.Sample(Keys)
					value := sample[key]

					_ = specimen.Filter(record.E{key, value}).Erase(ctx)
					err := impl.Filter(record.E{key, value}).Erase(ctx)
					assert.NoError(t, err)

					RunReadOperationsTests(ctx, t, impl, specimen)
				}
			})
		})

		t.Run("Post Write Read operations", func(t *testing.T) {
			t.Run("Not filtered", func(t *testing.T) {
				RunReadOperationsTests(ctx, t, impl, specimen)
			})

			t.Run("Filtered", func(t *testing.T) {
				t.Run("Single filter", func(t *testing.T) {
					RunFilterTests(ctx, t, [][][]record.Entry{
						{
							{{"group", "a"}},
						},
					}, impl, specimen)
				})

				t.Run("Multiple filters (one term)", func(t *testing.T) {
					RunFilterTests(ctx, t, [][][]record.Entry{
						{
							lo.Map(Groups, func(group string, _ int) record.E {
								return record.E{"group", group}
							}),
						},
					}, impl, specimen)
				})

				t.Run("Multiple filters (multiple terms)", func(t *testing.T) {

					var groupFilters []record.E
					var nameFilters []record.E

					for i := 0; i < 5; i++ {
						groupFilters = append(groupFilters,
							record.E{"group", lo.Sample(Groups)})
					}

					for i := 0; i < 100; i++ {
						val := specimen.Sample()["name"]
						nameFilters = append(nameFilters, record.E{"name", val})
					}

					RunFilterTests(ctx, t, [][][]record.Entry{
						{
							groupFilters,
							nameFilters,
						},
					}, impl, specimen)
				})

				t.Run("Not existing property", func(t *testing.T) {
					RunFilterTests(ctx, t, [][][]record.Entry{
						{
							{{"x", "1"}, {"y", "2"}},
							{{"x", "4"}, {"y", "3"}},
						},
					}, impl, specimen)
				})

				t.Run("Random", func(t *testing.T) {
					var cases [][][]record.E

					testCases := randomTestCases()
					for i := 0; i < testCases; i++ {
						var workingSet []record.R
						for j := 0; j < TotalRecords/5; j++ {
							workingSet = append(workingSet, specimen.Sample())
						}

						var _case [][]record.E
						termsCount := 1 + rand.Intn(10)
						for j := 0; j < termsCount; j++ {
							var term []record.E
							for k := 0; k < TotalRecords/8; k++ {
								key := lo.Sample(Keys)
								rec := lo.Sample(workingSet)

								term = append(term, record.E{key, rec[key]})
							}
							_case = append(_case, term)
						}
						cases = append(cases, _case)
					}

					RunFilterTests(ctx, t, cases, impl, specimen)
				})
			})
		})
	})
}

func RunReadOperationsTests(ctx context.Context, t *testing.T, impl, specimen record.Set) {
	t.Run("Read Operations Suite", func(t *testing.T) {
		t.Run("Common", func(t *testing.T) {
			RunCursorTests(ctx, t, impl.Cursor(), specimen.Cursor())
		})
		t.Run("Distinct", func(t *testing.T) {
			RunDistinctTests(ctx, t, impl, specimen)
		})
	})
}

func RunDistinctTests(ctx context.Context, t *testing.T, impl, specimen record.Set) {
	t.Run("Group", func(t *testing.T) {
		RunDistinctKeyTests(ctx, t, [][]string{{"group"}}, impl, specimen)
	})

	t.Run("Not existing keys", func(t *testing.T) {
		RunDistinctKeyTests(ctx, t,
			[][]string{
				{"NE1"},
				{"NE2", "NE3"},
				{"NE5", "NE6", "NE7"},
			},
			impl, specimen)
	})

	t.Run("Random", func(t *testing.T) {
		var keySets [][]string

		testCases := randomTestCases()
		for i := 0; i < testCases; i++ {
			keys := randomKeySubset(0.3)
			keySets = append(keySets, keys)
		}

		RunDistinctKeyTests(ctx, t, keySets, impl, specimen)
	})
}

func RunDistinctKeyTests(ctx context.Context, t *testing.T, keySets [][]string, impl, specimen record.Set) {
	for _, keys := range keySets {
		t.Logf("Keys = %s", keys)
		t.Run("CursorTests", func(t *testing.T) {
			RunCursorTests(ctx, t, impl.Distinct(keys...), specimen.Distinct(keys...))
		})
	}
}

func RunCursorTests(ctx context.Context, t *testing.T, impl, specimen record.Cursor[record.Record]) {
	t.Run("Count", func(t *testing.T) {
		retrieved, err := impl.Count(ctx)
		assert.NoError(t, err)

		expected := lo.Must(specimen.Count(ctx))

		assert.Equal(t, expected, retrieved)
	})

	t.Run("Integrity", func(t *testing.T) {
		retrievedDataset, err := collect(ctx, impl.Iter())
		if err != nil {
			t.Fatal("failed to retrieve data:", err)
		}
		specimenDataset := lo.Must(collect(ctx, specimen.Iter()))

		sort.Sort(recordSorter(retrievedDataset))
		sort.Sort(recordSorter(specimenDataset))

		assertDatasetsEqual(t, retrievedDataset, specimenDataset)
	})

	t.Run("Sort", func(t *testing.T) {
		t.Run("Random", func(t *testing.T) {
			t.Log("Generating random sorters...")
			testCases := randomTestCases()

			for i := 0; i < testCases; i++ {
				sorters := randomSorters()
				t.Logf("Sorters = %s", sorters)
				retrievedDataset, err := collect(ctx, impl.Sort(sorters...).Iter())
				if err != nil {
					t.Fatal("failed to retrieve data:", err)
				}
				specimenDataset := lo.Must(collect(ctx, specimen.Sort(sorters...).Iter()))

				assertDatasetsEqual(t, retrievedDataset, specimenDataset)
			}

		})

		t.Run("Empty", func(t *testing.T) {
			retrievedDataset, err := collect(ctx, impl.Sort().Iter())
			if err != nil {
				t.Fatal("failed to retrieve data:", err)
			}
			specimenDataset := lo.Must(collect(ctx, specimen.Sort().Iter()))

			sort.Sort(recordSorter(retrievedDataset))
			sort.Sort(recordSorter(specimenDataset))

			assertDatasetsEqual(t, retrievedDataset, specimenDataset)
		})
	})

	t.Run("Skips and Limits", func(t *testing.T) {
		t.Run("Random", func(t *testing.T) {
			t.Log("Generating random offsets and limits")
			testCases := randomTestCases()

			count := lo.Must(specimen.Count(ctx))

			for i := 0; i < testCases; i++ {
				offset := rand.Intn(count*1/2 + 1)
				limit := rand.Intn(count*3/4 + 1)
				sorters := randomSorters()

				t.Logf("Offset = %d, Limit = %d", offset, limit)
				t.Logf("Sorters = %s", sorters)

				retrievedDataset, err := collect(ctx, impl.
					Sort(sorters...).
					Skip(offset).
					Limit(limit).
					Iter())
				if err != nil {
					t.Fatal("failed to retrieve data:", err)
				}

				specimenDataset := lo.Must(collect(ctx, specimen.
					Sort(sorters...).
					Skip(offset).
					Limit(limit).
					Iter()))

				assertDatasetsEqual(t, specimenDataset, retrievedDataset)
			}
		})
	})

	t.Run("Masks", func(t *testing.T) {
		t.Run("Common", func(t *testing.T) {
			RunMaskTests(ctx, t, [][][]record.Mask{
				{
					{record.Include("name")},
				},
				{
					{record.Include("name"), record.Include("name"), record.Include("surname")},
					{record.Include("name")},
					{record.Include("name")},
				},
				{
					{record.Include("name")},
					{record.Include("surname")},
				},
				{
					{record.Include("DOES NOT EXIST1")},
					{record.Include("DOES NOT EXIST2")},
				},
				{
					{record.Include("DOES NOT EXIST")},
					{record.Include("name")},
				},
				{
					{record.Exclude("name")},
					{record.Exclude("surname")},
				},
			}, impl, specimen)
		})
	})
}

func RunMaskTests(ctx context.Context, t *testing.T, maskingOperations [][][]record.Mask, impl, specimen record.Cursor[record.Record]) {
	for _, op := range maskingOperations {
		impl1, specimen1 := impl, specimen

		for _, masks := range op {
			impl1 = impl1.Mask(masks...)
			specimen1 = specimen1.Mask(masks...)
		}

		t.Logf("Masks = %s", op)

		retrievedDataset, err := collect(ctx, impl1.Iter())
		if err != nil {
			t.Fatal("failed to retrieve data:", err)
		}

		specimenDataset := lo.Must(collect(ctx, specimen1.Iter()))

		sort.Sort(recordSorter(retrievedDataset))
		sort.Sort(recordSorter(specimenDataset))

		assertDatasetsEqual(t, specimenDataset, retrievedDataset)
	}
}

func RunFilterTests(ctx context.Context, t *testing.T, filteringOperations [][][]record.Entry, impl, specimen record.Set) {
	for _, op := range filteringOperations {
		impl1, specimen1 := impl, specimen

		for _, filters := range op {
			impl1 = impl1.Filter(filters...)
			specimen1 = specimen1.Filter(filters...)
		}

		t.Logf("Filters = %s", op)
		RunReadOperationsTests(ctx, t, impl1, specimen1)
	}
}

//func randomTestCases() int {
//	return MinimumTestCases + rand.Intn(MaximumTestCases-MinimumTestCases+1)
//}

func randomTestCases() int {
	return (MinimumTestCases + MaximumTestCases) / 2
}

func randomRecord() record.R {
	return record.R{
		"updated_at": rand.Intn(1e9),
		"id":         uuid.New().String(),
		"name":       uuid.New().String(),
		"surname":    uuid.New().String(),
		"group":      lo.Sample(Groups),
	}
}

func randomSorters() []record.Sorter {
	var sorters []record.Sorter

	orderings := []record.Ordering{
		record.Asc,
		record.Desc,
		//record.NoOrdering,
	}

	lo.ForEach(lo.Shuffle(Keys), func(key string, _ int) {
		ordering := lo.Sample(orderings)
		sorters = append(sorters, record.Sorter{
			Key:      key,
			Ordering: ordering,
		})
	})

	return sorters
}

func randomKeySubset(dropRate float64) []string {
	return lo.Shuffle(lo.Filter(Keys, func(string, int) bool { return rand.Float64() >= dropRate }))
}

func collect[Item any](ctx context.Context, iter record.Iter[Item]) ([]Item, error) {
	var collected []Item
	defer iter.Close(ctx)
	for {
		rec, err := iter.Next(ctx)
		if errors.Is(err, record.ErrIterEmpty) {
			break
		}
		if err != nil {
			return nil, err
		}
		collected = append(collected, rec)
	}
	return collected, nil
}

func sorted(datasets ...[]record.R) {
	for _, ds := range datasets {
		sort.Sort(recordSorter(ds))
	}
}

func assertDatasetsEqual(t *testing.T, d1, d2 []record.R) {
	assert.True(t, datasetsEqual(d1, d2))
}

func datasetsEqual(d1, d2 []record.R) bool {
	return len(datasetsElementWiseDifference(d1, d2)) == 0
}

func datasetsElementWiseDifference(d1, d2 []record.R) []lo.Tuple2[record.R, record.R] {
	return lo.Filter(lo.Zip2(d1, d2), func(item lo.Tuple2[record.R, record.R], _ int) bool {
		return !recordsEqual(item.A, item.B)
	})
}

func recordsEqual(r1, r2 record.R) bool {
	if len(r1) != len(r2) {
		return false
	}
	for key, v1 := range r1 {
		v2, ok := r2[key]
		if !ok {
			return false
		}
		if !assert.ObjectsAreEqualValues(v1, v2) {
			return false
		}
	}
	return true
}

type recordSorter []record.R

func (sorter recordSorter) Len() int {
	return len(sorter)
}

func (sorter recordSorter) Less(i, j int) bool {
	r1, r2 := sorter[i], sorter[j]

	for _, key := range Keys {
		v1, v2 := r1[key], r2[key]
		if v1 == nil || v2 == nil {
			continue
		}
		if x, ok := v1.(string); ok {
			y := v2.(string)
			if x < y {
				return true
			} else if x > y {
				return false
			}
		}
		if x, ok := v1.(int); ok {
			y := v2.(int)
			if x < y {
				return true
			} else if x > y {
				return false
			}
		}
		if x, ok := v1.(int32); ok {
			y := v2.(int32)
			if x < y {
				return true
			} else if x > y {
				return false
			}
		}
		if x, ok := v1.(int64); ok {
			y := v2.(int64)
			if x < y {
				return true
			} else if x > y {
				return false
			}
		}
		if x, ok := v1.(float64); ok {
			y := v2.(float64)
			if x < y {
				return true
			} else if x > y {
				return false
			}
		}
	}

	return false
}

func (sorter recordSorter) Swap(i, j int) {
	sorter[i], sorter[j] = sorter[j], sorter[i]
}
