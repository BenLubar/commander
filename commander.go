package commander

import (
	"flag"
	"os"
	"runtime"
	"runtime/pprof"
)

// CPUProfileFlag is a filename to write a CPU profile.
var CPUProfileFlag string
var cpuProfileFile *os.File

// HeapProfileFlag is a filename to write a heap profile.
var HeapProfileFlag string
var heapProfileFile *os.File

// ThreadProfileFlag is a filename to write the stack traces that caused new OS
// threads to be created.
var ThreadProfileFlag string
var threadProfileFile *os.File

// BlockProfileFlag is a filename to write the stack traces that caused
// blocking on synchronization primitives.
var BlockProfileFlag string
var blockProfileFile *os.File

// RegisterFlags must be called before Init and before f.Parse. If you are not
// allocating your own FlagSet, pass flag.CommandLine as the argument. If you
// are using a different command line argument parsing package, you will need
// to assign the values of *ProfileFlag on your own.
func RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&CPUProfileFlag, "cpuprofile", "", "a filename to write a CPU profile.")
	f.StringVar(&HeapProfileFlag, "heapprofile", "", "a filename to write a heap profile.")
	f.StringVar(&ThreadProfileFlag, "threadprofile", "", "a filename to write the stack traces that caused new OS threads to be created.")
	f.StringVar(&BlockProfileFlag, "blockprofile", "", "a filename to write the stack traces that caused blocking on synchronization primitives.")
}

// Init creates the files named by any of the *ProfileFlag variables. If an
// error is returned, the program can still function, but the file that failed
// to open and any further files will not contain profiles. Close must still be
// called at the end of the program, however, as earlier profiles are not
// cancelled by later profiles failing.
func Init() error {
	if CPUProfileFlag != "" {
		f, err := os.Create(CPUProfileFlag)
		if err != nil {
			return err
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			_ = f.Close()
			return err
		}
		cpuProfileFile = f
	}
	if HeapProfileFlag != "" {
		f, err := os.Create(HeapProfileFlag)
		if err != nil {
			return err
		}
		heapProfileFile = f
	}
	if ThreadProfileFlag != "" {
		f, err := os.Create(ThreadProfileFlag)
		if err != nil {
			return err
		}
		threadProfileFile = f
	}
	if BlockProfileFlag != "" {
		f, err := os.Create(BlockProfileFlag)
		if err != nil {
			return err
		}
		blockProfileFile = f
	}
	return nil
}

// Close finishes the profiles and closes their files. Close should be run at
// the end of the program's execution, possibly through a defer statement in
// the main method.
func Close() {
	if cpuProfileFile != nil {
		pprof.StopCPUProfile()
		_ = cpuProfileFile.Close()
		cpuProfileFile = nil
	}
	if heapProfileFile != nil {
		runtime.GC() // make sure that the heap profile is complete.

		pprof.Lookup("heap").WriteTo(heapProfileFile, 0)
		_ = heapProfileFile.Close()
		heapProfileFile = nil
	}
	if threadProfileFile != nil {
		pprof.Lookup("threadcreate").WriteTo(threadProfileFile, 0)
		_ = threadProfileFile.Close()
		threadProfileFile = nil
	}
	if blockProfileFile != nil {
		pprof.Lookup("block").WriteTo(blockProfileFile, 0)
		_ = blockProfileFile.Close()
		blockProfileFile = nil
	}
}
