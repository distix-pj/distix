package model

import (
)

type System struct {
	HostName string
	Packages []Package
}

// func isEdgeExists(nl *sbom.NodeList, fromID, toID string, edgeType sbom.Edge_Type) bool {
// 	for _, edge := range nl.GetEdges() {
// 		if edge.GetFrom() == fromID && edge.GetType() == edgeType {
// 			for _, to := range edge.GetTo() {
// 				if to == toID {
// 					return true
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

// func isSelfDepend(fromID, toID string) bool {
// 	if fromID == toID {
// 		return true
// 	}
// 	return false
// }

