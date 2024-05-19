package core

type FlagKind string
type FlagKindSet map[FlagKind]struct{}

const (
	FlagKindGenEntity  FlagKind = "gen_entity"
	FlagKindGenEnum    FlagKind = "gen_enum"
	FlagKindGenSpanner FlagKind = "gen_spanner"
	FlagKindGenRPC     FlagKind = "gen_rpc"
)

func (s FlagKindSet) Add(kind FlagKind) {
	s[kind] = struct{}{}
}

func (s FlagKindSet) Has(kind FlagKind) bool {
	_, ok := s[kind]
	return ok
}

func (s FlagKindSet) Size() int {
	return len(s)
}
