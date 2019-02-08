package resolver

import (
	"sync"

	mapset "github.com/deckarep/golang-set"
	"go.aporeto.io/gaia"
)

type policy struct {
}

// Resolver is a an in-memory policy resolver.
type Resolver struct {
	policies       map[string]*gaia.Policy
	allSubjectTags map[string]int
	allObjectTags  map[string]int
	sync.RWMutex
}

// NewResolutionEngine creates a new in-memory policy resolution
// engine.
func NewResolutionEngine() *Resolver {
	return &Resolver{
		policies:       map[string]*gaia.Policy{},
		allSubjectTags: map[string]int{},
		allObjectTags:  map[string]int{},
	}
}

// Insert will insert a new policy into the cache memory.
func (r *Resolver) Insert(p *gaia.Policy) {
	r.Lock()
	defer r.Unlock()

	if oldPolicy, ok := r.policies[p.ID]; ok {
		r.pruneTags(oldPolicy.Subject, true)
		r.pruneTags(oldPolicy.Object, false)
	}

	r.addTags(p.Subject, true)
	r.addTags(p.Object, false)
	r.policies[p.ID] = p
}

// Remove will remove a policy from the cache memory.
func (r *Resolver) Remove(id string) {
	r.Lock()
	defer r.Unlock()

	p, ok := r.policies[id]
	if !ok {
		return
	}

	delete(r.policies, id)
	r.pruneTags(p.Subject, true)
	r.pruneTags(p.Object, false)
}

// MatchingPolicies will return the list of matching policies give a set
// of incoming tags. It currently performs a linear search on policies
// using a set stucture for the tags.
func (r *Resolver) MatchingPolicies(tags []string, isSubject bool) gaia.PoliciesList {

	r.RLock()
	defer r.RUnlock()

	tagSet := r.craeteTagSet(tags, isSubject)

	result := gaia.PoliciesList{}

	for _, p := range r.policies {
		rules := p.Object
		if isSubject {
			rules = p.Subject
		}

		for _, clause := range rules {
			iSubset := make([]interface{}, len(clause))
			for i := range clause {
				iSubset[i] = clause[i]
			}
			if tagSet.Contains(iSubset...) {
				result = append(result, p)
				break
			}
		}
	}

	return result
}

func (r *Resolver) craeteTagSet(tags []string, isSubject bool) mapset.Set {
	allTags := r.allObjectTags
	if isSubject {
		allTags = r.allSubjectTags
	}

	set := mapset.NewThreadUnsafeSet()

	for _, t := range tags {
		if _, ok := allTags[t]; ok {
			set.Add(t)
		}
	}

	return set
}

func (r *Resolver) pruneTags(rules [][]string, isSubject bool) {
	allTags := r.allObjectTags
	if isSubject {
		allTags = r.allSubjectTags
	}

	for _, clause := range rules {
		for _, tag := range clause {
			if count, ok := allTags[tag]; ok {
				if count == 1 {
					delete(allTags, tag)
				} else {
					allTags[tag] = count - 1
				}
			}
		}
	}
}

func (r *Resolver) addTags(rules [][]string, isSubject bool) {
	allTags := r.allObjectTags
	if isSubject {
		allTags = r.allSubjectTags
	}

	for _, clause := range rules {
		for _, tag := range clause {
			if count, ok := allTags[tag]; ok {
				allTags[tag] = count + 1
			} else {
				allTags[tag] = 1
			}
		}
	}
}
