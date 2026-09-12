// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package mdns

import "syscall"

// sendSourcePort is the local UDP port outbound datagrams are sent from.
// Windows keeps the original ephemeral-port behavior: it will not reliably
// send from a socket bound to the contested mDNS port, and a Windows host is
// normally discovered directly rather than through a reflector. See
// openSendConn.
const sendSourcePort = 0

// sendFromRecvSocket is false on Windows: the platform refuses to send from
// the socket bound to the multicast group, so Run never shares its receive
// socket and sends always go through a fresh per-send socket.
const sendFromRecvSocket = false

// setReuseAddr is the Windows counterpart of the Unix build. The only
// difference is the socket handle type (syscall.Handle vs int). See the Unix
// file for the rationale; SO_REUSEADDR on Windows gives the shared-bind
// behavior we need and there is no SO_REUSEPORT to worry about.
func setReuseAddr(network, address string, c syscall.RawConn) error {
	var sockErr error
	if err := c.Control(func(fd uintptr) {
		sockErr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	}); err != nil {
		return err
	}
	return sockErr
}
