package resolver

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/gaia"
)

func Test_NewResolutionEngine(t *testing.T) {
	Convey("When I create a new resolution engine", t, func() {
		e := NewResolutionEngine()
		Convey("The data structures must be initialized", func() {
			So(e.allObjectTags, ShouldNotBeNil)
			So(e.allSubjectTags, ShouldNotBeNil)
			So(e.policies, ShouldNotBeNil)
		})
	})
}

func Test_Insert(t *testing.T) {
	Convey("Given a new resolution engine", t, func() {
		e := NewResolutionEngine()
		Convey("When I insert a new policy", func() {
			plc := gaia.NewPolicy()
			plc.ID = "1"
			plc.Subject = [][]string{
				[]string{"subject1-clause1", "subject1-clause2"},
				[]string{"subject2-clause1", "subject2-clause2"},
			}
			plc.Object = [][]string{
				[]string{"object1-clause1", "object1-clause2"},
				[]string{"object2-clause1", "object2-clause2"},
			}
			e.Insert(plc)
			Convey("All data must be correct", func() {
				So(len(e.policies), ShouldEqual, 1)
				So(e.policies["1"], ShouldResemble, plc)
				So(len(e.allSubjectTags), ShouldEqual, 4)
				So(e.allSubjectTags["subject1-clause1"], ShouldEqual, 1)
				So(e.allSubjectTags["subject1-clause2"], ShouldEqual, 1)
				So(e.allSubjectTags["subject2-clause1"], ShouldEqual, 1)
				So(e.allSubjectTags["subject2-clause2"], ShouldEqual, 1)
				So(len(e.allObjectTags), ShouldEqual, 4)
				So(e.allObjectTags["object1-clause1"], ShouldEqual, 1)
				So(e.allObjectTags["object1-clause2"], ShouldEqual, 1)
				So(e.allObjectTags["object2-clause1"], ShouldEqual, 1)
				So(e.allObjectTags["object2-clause2"], ShouldEqual, 1)
			})

			Convey("When I update the policy, the data must be correct", func() {
				plc = gaia.NewPolicy()
				plc.ID = "1"
				plc.Subject = [][]string{
					[]string{"app1-clause1", "app1-clause2"},
					[]string{"app2-clause1", "app2-clause2"},
				}
				plc.Object = [][]string{
					[]string{"appobject1-clause1", "appobject1-clause2"},
					[]string{"appobject2-clause1", "appobject2-clause2"},
				}
				e.Insert(plc)
				So(len(e.policies), ShouldEqual, 1)
				So(e.policies["1"], ShouldResemble, plc)
				So(len(e.allSubjectTags), ShouldEqual, 4)
				So(e.allSubjectTags["app1-clause1"], ShouldEqual, 1)
				So(e.allSubjectTags["app1-clause2"], ShouldEqual, 1)
				So(e.allSubjectTags["app2-clause1"], ShouldEqual, 1)
				So(e.allSubjectTags["app2-clause2"], ShouldEqual, 1)
				So(len(e.allObjectTags), ShouldEqual, 4)
				So(e.allObjectTags["appobject1-clause1"], ShouldEqual, 1)
				So(e.allObjectTags["appobject1-clause2"], ShouldEqual, 1)
				So(e.allObjectTags["appobject2-clause1"], ShouldEqual, 1)
				So(e.allObjectTags["appobject2-clause2"], ShouldEqual, 1)
			})
		})
	})
}

func Test_RemovePolicy(t *testing.T) {
	Convey("Given a resolver with two policies", t, func() {
		plc1 := gaia.NewPolicy()
		plc1.ID = "1"
		plc1.Subject = [][]string{
			[]string{"subject1-clause1", "subject1-clause2"},
		}
		plc1.Object = [][]string{
			[]string{"object1-clause1", "object1-clause2"},
		}

		plc2 := gaia.NewPolicy()
		plc2.ID = "2"
		plc2.Subject = [][]string{
			[]string{"subject2-clause1", "subject2-clause2"},
		}
		plc2.Object = [][]string{
			[]string{"object2-clause1", "object2-clause2"},
		}

		e := NewResolutionEngine()
		e.Insert(plc1)
		e.Insert(plc2)

		Convey("The data must be correct", func() {
			So(len(e.policies), ShouldEqual, 2)
			So(e.policies["1"], ShouldResemble, plc1)
			So(e.policies["2"], ShouldResemble, plc2)
			So(e.allSubjectTags["subject1-clause1"], ShouldEqual, 1)
			So(e.allSubjectTags["subject1-clause2"], ShouldEqual, 1)
			So(e.allSubjectTags["subject2-clause1"], ShouldEqual, 1)
			So(e.allSubjectTags["subject2-clause2"], ShouldEqual, 1)
			So(len(e.allObjectTags), ShouldEqual, 4)
			So(e.allObjectTags["object1-clause1"], ShouldEqual, 1)
			So(e.allObjectTags["object1-clause2"], ShouldEqual, 1)
			So(e.allObjectTags["object2-clause1"], ShouldEqual, 1)
			So(e.allObjectTags["object2-clause2"], ShouldEqual, 1)
		})

		Convey("When I remove the first policy, the data must be correct", func() {
			e.Remove("1")
			So(len(e.policies), ShouldEqual, 1)
			So(e.allSubjectTags["subject2-clause1"], ShouldEqual, 1)
			So(e.allSubjectTags["subject2-clause2"], ShouldEqual, 1)
			So(e.allObjectTags["object2-clause1"], ShouldEqual, 1)
			So(e.allObjectTags["object2-clause2"], ShouldEqual, 1)

			Convey("When I remove an non-existent policy, the data must be correct", func() {
				e.Remove("not-exist")
				So(len(e.policies), ShouldEqual, 1)
				So(e.allSubjectTags["subject2-clause1"], ShouldEqual, 1)
				So(e.allSubjectTags["subject2-clause2"], ShouldEqual, 1)
				So(e.allObjectTags["object2-clause1"], ShouldEqual, 1)
				So(e.allObjectTags["object2-clause2"], ShouldEqual, 1)
			})

			Convey("When I remove the second policy, the data store must be empty", func() {
				e.Remove("2")
				So(len(e.policies), ShouldEqual, 0)
				So(len(e.allObjectTags), ShouldEqual, 0)
				So(len(e.allSubjectTags), ShouldEqual, 0)
			})
		})
	})
}

func Test_MatchingPolicies(t *testing.T) {
	Convey("Given a policy resolver with valid policies", t, func() {
		plc1 := gaia.NewPolicy()
		plc1.ID = "1"
		plc1.Subject = [][]string{
			[]string{"subject1-clause1", "subject1-clause2"},
		}
		plc1.Object = [][]string{
			[]string{"object1-clause1", "object1-clause2"},
		}

		plc2 := gaia.NewPolicy()
		plc2.ID = "2"
		plc2.Subject = [][]string{
			[]string{"subject2-clause1", "subject2-clause2"},
		}
		plc2.Object = [][]string{
			[]string{"object2-clause1", "object2-clause2"},
		}

		plc3 := gaia.NewPolicy()
		plc3.ID = "3"
		plc3.Subject = [][]string{
			[]string{"subject1-clause1", "subject2-clause2"},
		}
		plc3.Object = [][]string{
			[]string{"object1-clause1", "object2-clause2"},
		}

		e := NewResolutionEngine()
		e.Insert(plc1)
		e.Insert(plc2)
		e.Insert(plc3)

		Convey("When I search for a set of tags that match policy1", func() {
			tags := []string{"subject1-clause1", "subject1-clause2", "bad", "ignore", "ignore2"}
			p := e.MatchingPolicies(tags, true)
			So(len(p), ShouldEqual, 1)
			So(p[0], ShouldResemble, plc1)
		})

		Convey("When I search for a set of tags that match policy2", func() {
			tags := []string{"subject2-clause1", "subject2-clause2", "bad", "ignore", "ignore2"}
			p := e.MatchingPolicies(tags, true)
			So(len(p), ShouldEqual, 1)
			So(p[0], ShouldResemble, plc2)
		})

		Convey("When I search for a set of tags that match no policy", func() {
			tags := []string{"bad", "ignore", "ignore2"}
			p := e.MatchingPolicies(tags, true)
			So(len(p), ShouldEqual, 0)
		})

		Convey("When I search for a set of tags with partial match", func() {
			tags := []string{"subject2-clause1", "bad", "ignore", "ignore2"}
			p := e.MatchingPolicies(tags, true)
			So(len(p), ShouldEqual, 0)
		})

		Convey("When I search for a set of tags that match all policies", func() {
			tags := []string{"subject1-clause1", "subject1-clause2", "subject2-clause1", "subject2-clause2", "bad", "ignore", "ignore2"}
			p := e.MatchingPolicies(tags, true)
			So(len(p), ShouldEqual, 3)
		})
	})
}

func Test_Benchmark(t *testing.T) {
	Convey("When I insert 1000 policies", t, func() {
		e := NewResolutionEngine()
		for i := 0; i < 10000; i++ {
			p := gaia.NewPolicy()
			p.ID = strconv.Itoa(i)
			p.Subject = [][]string{
				[]string{
					fmt.Sprintf("c1=%d", i),
					fmt.Sprintf("c2=%d", i),
					fmt.Sprintf("c3=%d", i),
					fmt.Sprintf("c4=%d", i),
				},
			}
			p.Object = [][]string{
				[]string{
					fmt.Sprintf("c1=%d", i),
					fmt.Sprintf("c2=%d", i),
					fmt.Sprintf("c3=%d", i),
					fmt.Sprintf("c4=%d", i),
				},
			}
			e.Insert(p)
		}

		tagBase := []string{}
		for i := 0; i < 50; i++ {
			tagBase = append(tagBase, fmt.Sprintf("sometag=%d", i))
		}

		Convey("I should be able to search for any of them", func() {
			start := time.Now()
			for i := 0; i < 1000; i++ {
				tags := append(tagBase, []string{
					fmt.Sprintf("c1=%d", i),
					fmt.Sprintf("c2=%d", i),
					fmt.Sprintf("c3=%d", i),
					fmt.Sprintf("c4=%d", i),
				}...)
				plc := e.MatchingPolicies(tags, true)
				So(len(plc), ShouldEqual, 1)
				So(plc[0], ShouldResemble, e.policies[strconv.Itoa(i)])
			}
			message := fmt.Sprintf("%s", time.Since(start))
			Convey("Time taken is "+message, func() {})
		})

		Convey("I should be able to search one policy very fast", func() {
			i := 1
			start := time.Now()
			tags := append(tagBase, []string{
				fmt.Sprintf("c1=%d", i),
				fmt.Sprintf("c2=%d", i),
				fmt.Sprintf("c3=%d", i),
				fmt.Sprintf("c4=%d", i),
			}...)
			plc := e.MatchingPolicies(tags, true)
			So(len(plc), ShouldEqual, 1)
			So(plc[0], ShouldResemble, e.policies[strconv.Itoa(i)])
			message := fmt.Sprintf("%s", time.Since(start))
			Convey("Time taken is "+message, func() {})
		})
	})
}
