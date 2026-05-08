package model

type PackageDependency struct {
	From *Package
	To *Package
}
type PackageDependencies struct {
	From *Package
	Tos []*Package
}

type DependencyGraph struct {
	deps []PackageDependency
}


func (g *DependencyGraph) Dependencies() []PackageDependency {
	return g.deps
}

func (g *DependencyGraph) GroupedDependencies() []PackageDependencies {
	seen := map[*Package]int{}
	result := []PackageDependencies{}
	for _, dep := range g.deps {
		if idx, exists := seen[dep.From]; exists {
			result[idx].Tos = append(result[idx].Tos, dep.To)
		} else {
			seen[dep.From] = len(result)
			result = append(result, PackageDependencies{
				From: dep.From,
				Tos: []*Package{dep.To},
			})
		}
	}
	return result
}


func (sys *System) BuildDependencyGraph() *DependencyGraph {
	providerMap := map[string]*Package{}
	// for _, pkg := range sys.Packages{
	for i := range sys.Packages{
		pkg := &sys.Packages[i]
		for _, prov := range pkg.Provides {
			providerMap[prov.Name] = pkg
		}
	}

	seen := map[[2]*Package]bool{}
	graph := DependencyGraph{}
	// for _, pkg := range sys.Packages{
	for i := range sys.Packages{
		pkg := &sys.Packages[i]
		for _, req := range pkg.Requires {
			if provider, ok := providerMap[req.Name]; ok {
				key := [2]*Package{pkg, provider}
				if !seen[key] {
					seen[key] = true
					graph.deps = append(graph.deps, PackageDependency{
						From: pkg,
						To: provider,
					})
				}
			}
			// else {} // error?
		}
	}
	return &graph
}

