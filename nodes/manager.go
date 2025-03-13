package nodes

import (
	"sync"

	"github.com/onnasoft/ZenithSQL/statement"
	"github.com/sirupsen/logrus"
)

type NodeManager struct {
	nodes       map[string]*Node
	masters     map[string]*Node
	slaves      map[string]*Node
	taggedNodes map[string]map[string]*Node
	mu          sync.RWMutex
	logger      *logrus.Logger
}

func NewNodeManager(logger *logrus.Logger) *NodeManager {
	return &NodeManager{
		nodes:       make(map[string]*Node),
		masters:     make(map[string]*Node),
		slaves:      make(map[string]*Node),
		taggedNodes: make(map[string]map[string]*Node),
		logger:      logger,
	}
}

func (m *NodeManager) GetNode(id string) *Node {
	m.mu.RLock()
	defer m.mu.RUnlock()

	node, exists := m.nodes[id]
	if !exists {
		return nil
	}

	return node
}

func (m *NodeManager) AddNode(stmt *statement.LoginStatement, role NodeRole) *Node {
	m.mu.Lock()
	defer m.mu.Unlock()

	if node, exists := m.nodes[stmt.NodeID]; exists {
		return node
	}

	tags := make(map[string]struct{})
	for _, tag := range stmt.Tags {
		tags[tag] = struct{}{}
	}

	node := NewNode(stmt.NodeID, role, tags, m.logger)
	m.nodes[stmt.NodeID] = node

	for _, tag := range stmt.Tags {
		if _, exists := m.taggedNodes[tag]; !exists {
			m.taggedNodes[tag] = make(map[string]*Node)
		}
		m.taggedNodes[tag][stmt.NodeID] = node
	}

	if role == Master {
		m.masters[stmt.NodeID] = node
	} else {
		m.slaves[stmt.NodeID] = node
	}

	return node
}

func (m *NodeManager) RemoveNode(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	node, exists := m.nodes[id]
	if !exists {
		return
	}

	for conn := range node.Connections {
		conn.Close()
	}

	delete(m.nodes, id)

	if node.Role == Master {
		delete(m.masters, id)
	} else {
		delete(m.slaves, id)
	}
}

func (m *NodeManager) ClearAllNodes() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id := range m.nodes {
		m.RemoveNode(id)
	}
}
