package explorer

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespaceLabel = "Namespace"
)

type ClusterExplorer struct {
	Items []*MenuItem
}

func (n *ClusterExplorer) List() error {
	n.Items = []*MenuItem{}

	ns, err := k8sclient.CoreV1().Namespaces().List(v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, item := range ns.Items {
		m := &MenuItem{}
		m.SetKind(namespaceLabel)
		m.SetName(item.Name)
		n.Items = append(n.Items, m)
	}
	n.Items = AddExit(n.Items)
	// TODO add CRDs
	return nil
}

func (n *ClusterExplorer) RunPrompt() (string, error) {
	var strs []string
	for _, item := range n.Items {
		strs = append(strs, item.GetReadable())
	}
	selection := transXY.Prompt("Select cluster resource", strs)
	checkExitItem(selection)
	return selection, nil
}

func (n *ClusterExplorer) Execute(selection string) error {
	item := NewMenuItemFromReadable(selection)
	switch item.GetKind() {
	case namespaceLabel:
		namespaceExplorer := &NamespaceExplorer{
			PreviousItem:       item,
			NamespaceToExplore: item.name,
			PreviousExplorer:   n,
		}
		return Explore(namespaceExplorer)
	case exitLabel:
		return Exit()
	default:
		return fmt.Errorf("unable to parse selection: %s", selection)
	}
}

func (n *ClusterExplorer) Kind() string {
	return "cluster"
}
