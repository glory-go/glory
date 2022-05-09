package interceptor

import (
	"github.com/glory-go/glory/debug/api/glory/boot"
)

type MetadataSorter []*boot.ServiceMetadata

func (m MetadataSorter) Len() int {
	return len(m)
}

func (m MetadataSorter) Less(i, j int) bool {
	return m[i].InterfaceName+m[i].ImplementationName < m[j].InterfaceName+m[j].ImplementationName
}

func (m MetadataSorter) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type MethodSorter []string

func (m MethodSorter) Len() int {
	return len(m)
}

func (m MethodSorter) Less(i, j int) bool {
	return m[i] < m[j]
}

func (m MethodSorter) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
