package gohavoq

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func splitLine(line string) []string {
	s := regexp.MustCompile(`[\,\s]+`).Split(line, -1)
	return s
}

func PartitionFromEdgeList(inFn string, outFn string, partFn func(uint32, uint32, int) (uint64, uint64, int), nPartitions int, offset uint32) error {

	// inFn is a text file of edges, space or comma delimeted.
	inF, err := os.OpenFile(inFn, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("openfile inF: %v", err)
	}
	defer inF.Close()

	type partition struct {
		outFn   string
		outF    *os.File
		offset  int64
		entries uint64
	}
	partitions := make([]partition, nPartitions)
	// Set up the partitions.
	for i := 0; i < nPartitions; i++ {
		pFn := fmt.Sprintf("%s-%d", outFn, i)
		// Reserve some space for the number of entries, written at the end.
		partitions[i] = partition{outFn: pFn, offset: 8, entries: 0}
	}

	scanner := bufio.NewScanner(inF)
	var l string
	for scanner.Scan() {
		l = scanner.Text()

		if strings.HasPrefix(l, "#") {
			continue

		}
		// Split the text into two uint32s.
		pieces := splitLine(l)
		if len(pieces) != 2 {
			return fmt.Errorf("Parsing error: got %s", l)
		}
		u64, err := strconv.ParseUint(pieces[0], 10, 32)
		if err != nil {
			return fmt.Errorf("Parsing error: got %s", l)
		}
		v64, err := strconv.ParseUint(pieces[1], 10, 32)
		if err != nil {
			return fmt.Errorf("Parsing error: got %s", l)
		}

		// if we have an offset (think 1-based), apply it here.
		u := uint32(u64) - offset
		v := uint32(v64) - offset

		// get the hashes and the partition number
		uHash, vHash, p := partFn(u, v, nPartitions)

		part := &partitions[p]
		if part.outF == nil { // we haven't seen this partition yet
			pF, err := os.OpenFile(part.outFn, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return fmt.Errorf("openfile outFn %s: %v", part.outFn, err)
			}
			part.outF = pF
		}
		// fmt.Println("u = ", u, "v = ", v, "part = ", part)
		uBytes := make([]byte, 8)
		vBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(uBytes, uHash)
		binary.LittleEndian.PutUint64(vBytes, vHash)

		// fmt.Println("  writing u:", uHash, " to offset ", part.offset)
		_, err = part.outF.WriteAt(uBytes, part.offset)
		if err != nil {
			return fmt.Errorf("Writing error for entry %d, offset %d: %v", part.entries, part.offset, err)
		}
		part.offset += 8
		_, err = part.outF.WriteAt(vBytes, part.offset)
		// fmt.Println("  writing v: ", vHash, " to offset ", part.offset)
		if err != nil {
			return fmt.Errorf("Writing error for entry %d, offset %d: %v", part.entries, part.offset, err)
		}
		part.offset += 8
		part.entries++
		// fmt.Printf("entries for part %d = %d, offset = %d\n", p, part.entries, part.offset)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Other error: %v", err)
	}
	for i := 0; i < nPartitions; i++ {
		part := partitions[i]
		if part.entries != 0 { // if we have entries for this partition, write the length at the beginning
			lenB := make([]byte, 8)
			// fmt.Printf("writing nEntries = %d for partition %d\n", part.entries, i)
			binary.LittleEndian.PutUint64(lenB, part.entries)
			n, err := part.outF.WriteAt(lenB, 0)
			if err != nil {
				return fmt.Errorf("WriteAt failed: %v", err)
			}

			if n != 8 {
				return fmt.Errorf("Only wrote %d of 8 bytes", n)
			}
			part.outF.Close()
		}
	}

	return nil
}
