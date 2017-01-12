package psn

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProcCSV(t *testing.T) {
	dn, err := GetDevice("/boot")
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	nt, err := GetDefaultInterface()
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}

	fpath := filepath.Join(os.TempDir(), fmt.Sprintf("test-%010d.csv", time.Now().UnixNano()))
	defer os.RemoveAll(fpath)

	epath := filepath.Join(homeDir(), "etcd-client-num")
	defer os.RemoveAll(epath)

	if err := toFile([]byte("10"), epath); err != nil {
		t.Fatal(err)
	}
	c := NewCSV(fpath, 1, dn, nt, epath)
	for i := 0; i < 3; i++ {
		if i > 0 {
			if err := toFile([]byte(fmt.Sprintf("%d", 100*i)), epath); err != nil {
				t.Fatal(err)
			}
		}

		fmt.Printf("#%d: collecting data with %s and %s at %s\n", i, dn, nt, fpath)
		if err := c.Add(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second)
	}

	// test if Add fills in empty rows
	time.Sleep(2*time.Second + 50*time.Millisecond)

	if err := c.Add(); err != nil {
		t.Fatal(err)
	}
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	fmt.Println("CSV contents:", string(b))

	cv, err := ReadCSV(fpath)
	if err != nil {
		t.Fatal(err)
	}
	if c.PID != cv.PID {
		t.Fatalf("PID expected %d, got %d", c.PID, cv.PID)
	}
	if c.DiskDevice != cv.DiskDevice {
		t.Fatalf("DiskDevice expected %s, got %s", c.DiskDevice, cv.DiskDevice)
	}
	if c.NetworkInterface != cv.NetworkInterface {
		t.Fatalf("NetworkInterface expected %s, got %s", c.NetworkInterface, cv.NetworkInterface)
	}
	if c.MinUnixTS != cv.MinUnixTS {
		t.Fatalf("MinUnixTS expected %d, got %d", c.MinUnixTS, cv.MinUnixTS)
	}
	if c.MaxUnixTS != cv.MaxUnixTS {
		t.Fatalf("MaxUnixTS expected %d, got %d", c.MaxUnixTS, cv.MaxUnixTS)
	}
	if len(c.Rows) != len(cv.Rows) {
		t.Fatalf("len(Rows) expected %d, got %d", len(c.Rows), len(cv.Rows))
	}

	for i := range c.Rows {
		if i == 0 && string(c.Rows[i].Extra) != "10" {
			t.Fatalf("Rows[%d].Extra expected 10, got %s", i, c.Rows[i].Extra)
		} else if i > 0 && i < 3 && string(c.Rows[i].Extra) != fmt.Sprintf("%d", 100*i) {
			t.Fatalf("Rows[%d].Extra expected %d, got %s", i, 100*i, c.Rows[i].Extra)
		} else if i >= 3 && string(c.Rows[i].Extra) != "200" {
			t.Fatalf("Rows[%d].Extra expected 200, got %s", i, c.Rows[i].Extra)
		}

		if c.Rows[i].PSEntry.Program != cv.Rows[i].PSEntry.Program {
			t.Fatalf("Rows[%d].PSEntry.Program expected %s, got %s", i, c.Rows[i].PSEntry.Program, cv.Rows[i].PSEntry.Program)
		}
		if c.Rows[i].PSEntry.State != cv.Rows[i].PSEntry.State {
			t.Fatalf("Rows[%d].PSEntry.State expected %s, got %s", i, c.Rows[i].PSEntry.State, cv.Rows[i].PSEntry.State)
		}
		if c.Rows[i].PSEntry.PID != cv.Rows[i].PSEntry.PID {
			t.Fatalf("Rows[%d].PSEntry.PID expected %d, got %d", i, c.Rows[i].PSEntry.PID, cv.Rows[i].PSEntry.PID)
		}
		if c.Rows[i].PSEntry.PPID != cv.Rows[i].PSEntry.PPID {
			t.Fatalf("Rows[%d].PSEntry.PPID expected %d, got %d", i, c.Rows[i].PSEntry.PPID, cv.Rows[i].PSEntry.PPID)
		}
		if c.Rows[i].PSEntry.CPU != cv.Rows[i].PSEntry.CPU {
			t.Fatalf("Rows[%d].PSEntry.CPU expected %s, got %s", i, c.Rows[i].PSEntry.CPU, cv.Rows[i].PSEntry.CPU)
		}
		if (cv.Rows[i].PSEntry.CPUNum - c.Rows[i].PSEntry.CPUNum) > 0.01 {
			t.Fatalf("Rows[%d].PSEntry.CPUNum expected %f, got %f", i, c.Rows[i].PSEntry.CPUNum, cv.Rows[i].PSEntry.CPUNum)
		}
		if c.Rows[i].PSEntry.VMRSS != cv.Rows[i].PSEntry.VMRSS {
			t.Fatalf("Rows[%d].PSEntry.VMRSS expected %s, got %s", i, c.Rows[i].PSEntry.VMRSS, cv.Rows[i].PSEntry.VMRSS)
		}
		if c.Rows[i].PSEntry.VMRSSNum != cv.Rows[i].PSEntry.VMRSSNum {
			t.Fatalf("Rows[%d].PSEntry.VMRSSNum expected %d, got %d", i, c.Rows[i].PSEntry.VMRSSNum, cv.Rows[i].PSEntry.VMRSSNum)
		}
		if c.Rows[i].PSEntry.VMSize != cv.Rows[i].PSEntry.VMSize {
			t.Fatalf("Rows[%d].PSEntry.VMSize expected %s, got %s", i, c.Rows[i].PSEntry.VMSize, cv.Rows[i].PSEntry.VMSize)
		}
		if c.Rows[i].PSEntry.VMSizeNum != cv.Rows[i].PSEntry.VMSizeNum {
			t.Fatalf("Rows[%d].PSEntry.VMSizeNum expected %d, got %d", i, c.Rows[i].PSEntry.VMSizeNum, cv.Rows[i].PSEntry.VMSizeNum)
		}
		if c.Rows[i].PSEntry.FD != cv.Rows[i].PSEntry.FD {
			t.Fatalf("Rows[%d].PSEntry.FD expected %d, got %d", i, c.Rows[i].PSEntry.FD, cv.Rows[i].PSEntry.FD)
		}
		if c.Rows[i].PSEntry.Threads != cv.Rows[i].PSEntry.Threads {
			t.Fatalf("Rows[%d].PSEntry.Threads expected %d, got %d", i, c.Rows[i].PSEntry.Threads, cv.Rows[i].PSEntry.Threads)
		}

		if c.Rows[i].DSEntry.ReadsCompleted != cv.Rows[i].DSEntry.ReadsCompleted {
			t.Fatalf("Rows[%d].DSEntry.ReadsCompleted expected %d, got %d", i, c.Rows[i].DSEntry.ReadsCompleted, cv.Rows[i].DSEntry.ReadsCompleted)
		}
		if c.Rows[i].DSEntry.SectorsRead != cv.Rows[i].DSEntry.SectorsRead {
			t.Fatalf("Rows[%d].DSEntry.SectorsRead expected %d, got %d", i, c.Rows[i].DSEntry.SectorsRead, cv.Rows[i].DSEntry.SectorsRead)
		}
		if c.Rows[i].DSEntry.WritesCompleted != cv.Rows[i].DSEntry.WritesCompleted {
			t.Fatalf("Rows[%d].DSEntry.WritesCompleted expected %d, got %d", i, c.Rows[i].DSEntry.WritesCompleted, cv.Rows[i].DSEntry.WritesCompleted)
		}
		if c.Rows[i].DSEntry.SectorsWritten != cv.Rows[i].DSEntry.SectorsWritten {
			t.Fatalf("Rows[%d].DSEntry.SectorsWritten expected %d, got %d", i, c.Rows[i].DSEntry.SectorsWritten, cv.Rows[i].DSEntry.SectorsWritten)
		}
		if c.Rows[i].DSEntry.TimeSpentOnReading != cv.Rows[i].DSEntry.TimeSpentOnReading {
			t.Fatalf("Rows[%d].DSEntry.TimeSpentOnReading expected %s, got %s", i, c.Rows[i].DSEntry.TimeSpentOnReading, cv.Rows[i].DSEntry.TimeSpentOnReading)
		}
		if c.Rows[i].DSEntry.TimeSpentOnWriting != cv.Rows[i].DSEntry.TimeSpentOnWriting {
			t.Fatalf("Rows[%d].DSEntry.TimeSpentOnWriting expected %s, got %s", i, c.Rows[i].DSEntry.TimeSpentOnWriting, cv.Rows[i].DSEntry.TimeSpentOnWriting)
		}
		if c.Rows[i].DSEntry.TimeSpentOnReadingMs != cv.Rows[i].DSEntry.TimeSpentOnReadingMs {
			t.Fatalf("Rows[%d].DSEntry.TimeSpentOnReadingMs expected %d, got %d", i, c.Rows[i].DSEntry.TimeSpentOnReadingMs, cv.Rows[i].DSEntry.TimeSpentOnReadingMs)
		}
		if c.Rows[i].DSEntry.TimeSpentOnWritingMs != cv.Rows[i].DSEntry.TimeSpentOnWritingMs {
			t.Fatalf("Rows[%d].DSEntry.TimeSpentOnWritingMs expected %d, got %d", i, c.Rows[i].DSEntry.TimeSpentOnWritingMs, cv.Rows[i].DSEntry.TimeSpentOnWritingMs)
		}

		if c.Rows[i].ReadsCompletedDelta != cv.Rows[i].ReadsCompletedDelta {
			t.Fatalf("Rows[%d].ReadsCompletedDelta expected %d, got %d", i, c.Rows[i].ReadsCompletedDelta, cv.Rows[i].ReadsCompletedDelta)
		}
		if c.Rows[i].SectorsReadDelta != cv.Rows[i].SectorsReadDelta {
			t.Fatalf("Rows[%d].SectorsReadDelta expected %d, got %d", i, c.Rows[i].SectorsReadDelta, cv.Rows[i].SectorsReadDelta)
		}
		if c.Rows[i].WritesCompletedDelta != cv.Rows[i].WritesCompletedDelta {
			t.Fatalf("Rows[%d].WritesCompletedDelta expected %d, got %d", i, c.Rows[i].WritesCompletedDelta, cv.Rows[i].WritesCompletedDelta)
		}
		if c.Rows[i].SectorsWrittenDelta != cv.Rows[i].SectorsWrittenDelta {
			t.Fatalf("Rows[%d].SectorsWrittenDelta expected %d, got %d", i, c.Rows[i].SectorsWrittenDelta, cv.Rows[i].SectorsWrittenDelta)
		}

		if c.Rows[i].NSEntry.ReceiveBytes != cv.Rows[i].NSEntry.ReceiveBytes {
			t.Fatalf("Rows[%d].NSEntry.ReceiveBytes expected %s, got %s", i, c.Rows[i].NSEntry.ReceiveBytes, cv.Rows[i].NSEntry.ReceiveBytes)
		}
		if c.Rows[i].NSEntry.ReceiveBytesNum != cv.Rows[i].NSEntry.ReceiveBytesNum {
			t.Fatalf("Rows[%d].NSEntry.ReceiveBytesNum expected %d, got %d", i, c.Rows[i].NSEntry.ReceiveBytesNum, cv.Rows[i].NSEntry.ReceiveBytesNum)
		}
		if c.Rows[i].NSEntry.TransmitBytes != cv.Rows[i].NSEntry.TransmitBytes {
			t.Fatalf("Rows[%d].NSEntry.TransmitBytes expected %s, got %s", i, c.Rows[i].NSEntry.TransmitBytes, cv.Rows[i].NSEntry.TransmitBytes)
		}
		if c.Rows[i].NSEntry.TransmitBytesNum != cv.Rows[i].NSEntry.TransmitBytesNum {
			t.Fatalf("Rows[%d].NSEntry.TransmitBytesNum expected %d, got %d", i, c.Rows[i].NSEntry.TransmitBytesNum, cv.Rows[i].NSEntry.TransmitBytesNum)
		}
		if c.Rows[i].NSEntry.ReceivePackets != cv.Rows[i].NSEntry.ReceivePackets {
			t.Fatalf("Rows[%d].NSEntry.ReceivePackets expected %d, got %d", i, c.Rows[i].NSEntry.ReceivePackets, cv.Rows[i].NSEntry.ReceivePackets)
		}
		if c.Rows[i].NSEntry.TransmitPackets != cv.Rows[i].NSEntry.TransmitPackets {
			t.Fatalf("Rows[%d].NSEntry.TransmitPackets expected %d, got %d", i, c.Rows[i].NSEntry.TransmitPackets, cv.Rows[i].NSEntry.TransmitPackets)
		}

		if c.Rows[i].ReceiveBytesDelta != cv.Rows[i].ReceiveBytesDelta {
			t.Fatalf("Rows[%d].ReceiveBytesDelta expected %s, got %s", i, c.Rows[i].ReceiveBytesDelta, cv.Rows[i].ReceiveBytesDelta)
		}
		if c.Rows[i].ReceivePacketsDelta != cv.Rows[i].ReceivePacketsDelta {
			t.Fatalf("Rows[%d].ReceivePacketsDelta expected %d, got %d", i, c.Rows[i].ReceivePacketsDelta, cv.Rows[i].ReceivePacketsDelta)
		}
		if c.Rows[i].TransmitBytesDelta != cv.Rows[i].TransmitBytesDelta {
			t.Fatalf("Rows[%d].TransmitBytesDelta expected %s, got %s", i, c.Rows[i].TransmitBytesDelta, cv.Rows[i].TransmitBytesDelta)
		}
		if c.Rows[i].TransmitPacketsDelta != cv.Rows[i].TransmitPacketsDelta {
			t.Fatalf("Rows[%d].TransmitPacketsDelta expected %d, got %d", i, c.Rows[i].TransmitPacketsDelta, cv.Rows[i].TransmitPacketsDelta)
		}
		if c.Rows[i].ReceiveBytesNumDelta != cv.Rows[i].ReceiveBytesNumDelta {
			t.Fatalf("Rows[%d].ReceiveBytesNumDelta expected %d, got %d", i, c.Rows[i].ReceiveBytesNumDelta, cv.Rows[i].ReceiveBytesNumDelta)
		}
		if c.Rows[i].TransmitBytesNumDelta != cv.Rows[i].TransmitBytesNumDelta {
			t.Fatalf("Rows[%d].TransmitBytesNumDelta expected %d, got %d", i, c.Rows[i].TransmitBytesNumDelta, cv.Rows[i].TransmitBytesNumDelta)
		}

		if !bytes.Equal(c.Rows[i].Extra, cv.Rows[i].Extra) {
			t.Fatalf("Rows[%d].Extra expected %q, got %q", i, c.Rows[i].Extra, cv.Rows[i].Extra)
		}
	}
}
