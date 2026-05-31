package action

import (
	"flclash-headless/model"
)

func BuildGlobalExitOptions(groups []model.GroupSummary, detail *model.GroupDetail) []model.GroupNode {
	if detail == nil {
		return nil
	}

	groupNames := make(map[string]bool)
	for _, g := range groups {
		if g.Name == "GLOBAL" {
			continue
		}
		groupNames[g.Name] = true
	}

	options := make([]model.GroupNode, 0, len(detail.Nodes))

	for _, node := range detail.Nodes {
		if node.Name == "DIRECT" || node.Name == "REJECT" {
			continue
		}
		if groupNames[node.Name] {
			options = append(options, node)
		}
	}

	for _, node := range detail.Nodes {
		if node.Name == "DIRECT" || node.Name == "REJECT" {
			continue
		}
		if groupNames[node.Name] {
			continue
		}
		options = append(options, node)
	}

	for _, node := range detail.Nodes {
		if node.Name == "DIRECT" || node.Name == "REJECT" {
			options = append(options, node)
		}
	}

	return options
}
