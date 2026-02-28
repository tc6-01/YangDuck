package recipe

import (
	"strings"
)

type Registry struct {
	recipes map[string]Recipe
	byType  map[RecipeType][]Recipe
}

func NewRegistry() *Registry {
	return &Registry{
		recipes: make(map[string]Recipe),
		byType:  make(map[RecipeType][]Recipe),
	}
}

func (r *Registry) Add(recipes ...Recipe) {
	for _, rec := range recipes {
		r.recipes[rec.ID] = rec
		r.byType[rec.Type] = append(r.byType[rec.Type], rec)
	}
}

func (r *Registry) Get(id string) (Recipe, bool) {
	rec, ok := r.recipes[id]
	return rec, ok
}

func (r *Registry) ListByType(t RecipeType) []Recipe {
	return r.byType[t]
}

func (r *Registry) All() []Recipe {
	out := make([]Recipe, 0, len(r.recipes))
	for _, rec := range r.recipes {
		out = append(out, rec)
	}
	return out
}

func (r *Registry) Search(keyword string) []Recipe {
	keyword = strings.ToLower(keyword)
	var results []Recipe
	for _, rec := range r.recipes {
		if matchesKeyword(rec, keyword) {
			results = append(results, rec)
		}
	}
	return results
}

func matchesKeyword(rec Recipe, kw string) bool {
	if strings.Contains(strings.ToLower(rec.Name), kw) {
		return true
	}
	if strings.Contains(strings.ToLower(rec.Description), kw) {
		return true
	}
	if strings.Contains(strings.ToLower(rec.ID), kw) {
		return true
	}
	for _, tag := range rec.Tags {
		if strings.Contains(strings.ToLower(tag), kw) {
			return true
		}
	}
	return false
}

func (r *Registry) Count() int {
	return len(r.recipes)
}

func (r *Registry) Bundles() []Recipe {
	return r.byType[TypeBundle]
}
