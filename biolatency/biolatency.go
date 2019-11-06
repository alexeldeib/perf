package biolatency

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	"go.opencensus.io/stats"
)

const (
	biolatency = "biolatency-bpfcc"
	// pattern    = `\[([0-9]+[KM]{1}), ([0-9]+[KM]{1})\)[ ]+([0-9]+) \|[ @]*\|`
	kilo    = 1024
	mega    = kilo ^ 2
	pattern = `^[ ]*([0-9]+) -> ([0-9]+)[ ]*: ([0-9]+).*$`
)

func New(mBiolatency *stats.Int64Measure) (map[string]int, error) {
	var regex, err = regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	if _, err = exec.LookPath(biolatency); err != nil {
		return nil, err
	}

	var cmd = exec.Command(biolatency)
	var buf bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		// forcibly kill, we really don't care.
		time.Sleep(10 * time.Second)
		_ = cmd.Process.Signal(os.Interrupt)
		time.Sleep(10 * time.Second)
		_ = cmd.Process.Signal(os.Kill)
	}()

	if err := cmd.Wait(); err != nil {
		return nil, errors.Wrapf(err, "faled to wait for cmd: %s", stderr.String())
	}

	var out = map[string]int{}

	var scanner = bufio.NewScanner(&buf)
	for scanner.Scan() {
		var line = scanner.Text()
		if strings.Contains(line, "Attaching") || strings.Contains(line, "usecs") {
			continue
		}
		if !regex.MatchString(line) {
			continue
		}
		matches := regex.FindAllStringSubmatch(line, -1)
		if matches == nil {
			return nil, errors.New("failed to read biolatency after passing validation")
		}
		litter.Dump(line)
		litter.Dump(matches)
		if len(matches) != 1 || len(matches[0]) != 4 {
			return nil, errors.New("expected to parse 1 match and 3 submatches")
		}

		count, err := strconv.Atoi(matches[0][3])
		if err != nil {
			return nil, err
		}

		bucket, err := strconv.Atoi(matches[0][2])
		if err != nil {
			return nil, err
		}

		// if unit == "K" {
		// 	bucket = bucket * int(kilo)
		// }
		// if unit == "M" {
		// 	bucket = bucket * int(mega)
		// }
		litter.Dump("writing bucket  ", bucket, "  count", count)

		// count represents a bucket in the histogram. We map it to the oc histogram.
		for n := 0; n < count; n++ {
			stats.Record(context.Background(), mBiolatency.M(int64(bucket)))
		}
		out[matches[0][2]] = count
	}
	return out, nil
}
