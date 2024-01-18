package peer

import "pan/memory"

type DialerManager = *memory.Map[NodeType, NodeDialer]
type HandshakeManager = *memory.Map[NodeType, NodeHandshake]
