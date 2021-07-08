package bstfsuring

import (
	"bytes"
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/dshulyak/uring"
	qt "github.com/frankban/quicktest"
)

func TestWhee(t *testing.T) {
	// Have a temp file to play with.
	// (This is a dummy for now -- we need to test file creation too.)
	f, err := ioutil.TempFile("", "writev-tests-")
	qt.Assert(t, err, qt.IsNil)
	defer f.Close()

	// Make the ring.
	ring, err := uring.Setup(4, nil) // TODO add the poll mode flags and see if we can survive.
	qt.Assert(t, err, qt.IsNil)
	defer ring.Close()

	// Okay, some state and memory we'll use for writes.
	// (Is there not a way where we get some memmap'd memory to work with, rather than having to copy this into it later?)
	// (Need to play with this to see how which things heap escape, and if we can minimize that.)
	var offset uint64
	buf := bytes.Repeat([]byte{'a'}, 8)
	vecs := []syscall.Iovec{{ // If you received several different buffers (say from some other IO source like network) you can batch them here, and cram them into one write call.
		Base: &buf[0],
		Len:  uint64(len(buf)),
	}}

	// Push info for this write into an entry on the ring.
	sqe := ring.GetSQEntry()
	uring.Writev(sqe, f.Fd(), vecs, offset, 0)
	offset += uint64(len(buf))

	// Tell the ring to roll that stuff over.
	// (There appear to be several methods for this and idk how they differ yet.)
	_, err = ring.Submit(1)
	qt.Assert(t, err, qt.IsNil)

	// ... kernel be workin' now ...

	// Read completion queue.
	// (I think this isn't prepared for nonblocking mode; making it so will probably require some nontrivial code (e.g. EAGAIN checking).)
	cqe, err := ring.GetCQEntry(0)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, cqe.Result() >= 0, qt.IsTrue, qt.Commentf("failed with %v", syscall.Errno(-cqe.Result())))

	// You should now be able to read the file back normally to check success.
}
