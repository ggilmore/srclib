package graph

import (
	"strconv"

	"sourcegraph.com/sourcegraph/srclib/repo"
)

type RefSymbolKey struct {
	DefRepo     repo.URI `db:"symbol_repo" json:",omitempty"`
	DefUnitType string   `db:"symbol_unit_type" json:",omitempty"`
	DefUnit     string   `db:"symbol_unit" json:",omitempty"`
	DefPath     DefPath  `db:"symbol_path" json:",omitempty"`
}

type RefKey struct {
	DefRepo     repo.URI `db:"symbol_repo" json:",omitempty"`
	DefUnitType string   `db:"symbol_unit_type" json:",omitempty"`
	DefUnit     string   `db:"symbol_unit" json:",omitempty"`
	DefPath     DefPath  `db:"symbol_path" json:",omitempty"`
	Def         bool     `json:",omitempty"`
	Repo        repo.URI `json:",omitempty"`
	UnitType    string   `db:"unit_type" json:",omitempty"`
	Unit        string   `json:",omitempty"`
	File        string   `json:",omitempty"`
	CommitID    string   `db:"commit_id" json:",omitempty"`
	Start       int      `json:",omitempty"`
	End         int      `json:",omitempty"`
}

func (r *RefKey) RefSymbolKey() RefSymbolKey {
	return RefSymbolKey{
		DefRepo:     r.DefRepo,
		DefUnitType: r.DefUnitType,
		DefUnit:     r.DefUnit,
		DefPath:     r.DefPath,
	}
}

// START Ref OMIT
// Ref represents a reference from source code to a symbol.
type Ref struct {
	// The definition that this reference points to
	DefRepo     repo.URI `db:"symbol_repo"`
	DefUnitType string   `db:"symbol_unit_type"`
	DefUnit     string   `db:"symbol_unit"`
	DefPath     DefPath  `db:"symbol_path"`

	// Def is true if this ref is the original definition or a redefinition
	Def bool

	Repo repo.URI

	// CommitID is the immutable commit ID (not the branch name) of the VCS
	// revision that this ref was found in.
	CommitID string `db:"commit_id" json:",omitempty"`

	UnitType string `db:"unit_type" json:",omitempty"`
	Unit     string `json:",omitempty"`

	File  string
	Start int
	End   int
}

// END Ref OMIT

func (r *Ref) RefKey() RefKey {
	return RefKey{
		DefRepo:     r.DefRepo,
		DefUnitType: r.DefUnitType,
		DefUnit:     r.DefUnit,
		DefPath:     r.DefPath,
		Def:         r.Def,
		Repo:        r.Repo,
		UnitType:    r.UnitType,
		Unit:        r.Unit,
		File:        r.File,
		Start:       r.Start,
		End:         r.End,
	}
}

func (r *Ref) RefSymbolKey() RefSymbolKey {
	return RefSymbolKey{
		DefRepo:     r.DefRepo,
		DefUnitType: r.DefUnitType,
		DefUnit:     r.DefUnit,
		DefPath:     r.DefPath,
	}
}

func (r *Ref) DefKey() DefKey {
	return DefKey{
		Repo:     r.DefRepo,
		UnitType: r.DefUnitType,
		Unit:     r.DefUnit,
		Path:     r.DefPath,
	}
}

func (r *Ref) SetFromSymbolKey(k DefKey) {
	r.DefPath = k.Path
	r.DefUnitType = k.UnitType
	r.DefUnit = k.Unit
	r.DefRepo = k.Repo
}

// Sorting

type Refs []*Ref

func (r *Ref) sortKey() string {
	return string(r.DefPath) + string(r.DefRepo) + r.DefUnitType + r.DefUnit + string(r.Repo) + r.UnitType + r.Unit + r.File + strconv.Itoa(r.Start) + strconv.Itoa(r.End)
}
func (vs Refs) Len() int           { return len(vs) }
func (vs Refs) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs Refs) Less(i, j int) bool { return vs[i].sortKey() < vs[j].sortKey() }

// RefSet is a set of Refs. It can used to determine whether a grapher emits
// duplicate refs.
type RefSet struct {
	refs map[Ref]struct{}
}

func NewRefSet() *RefSet {
	return &RefSet{make(map[Ref]struct{})}
}

// AddAndCheckUnique adds ref to the set of seen refs, and returns whether the
// ref already existed in the set.
func (c *RefSet) AddAndCheckUnique(ref Ref) (duplicate bool) {
	_, present := c.refs[ref]
	if present {
		return true
	}
	c.refs[ref] = struct{}{}
	return false
}
