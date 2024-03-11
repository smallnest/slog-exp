package rotate_test

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/smallnest/slog-exp/rotate"
	"github.com/stretchr/testify/assert"
)

func TestSatisfiesIOWriter(t *testing.T) {
	var w io.Writer
	w, _ = rotate.New("/foo/bar")
	_ = w
}

func TestSatisfiesIOCloser(t *testing.T) {
	var c io.Closer
	c, _ = rotate.New("/foo/bar")
	_ = c
}

func TestLogRotate(t *testing.T) {
	dir, err := ioutil.TempDir("", "file-rotatelogs-test")
	if !assert.NoError(t, err, "creating temporary directory should succeed") {
		return
	}
	defer os.RemoveAll(dir)

	// Change current time, so we can safely purge old logs
	dummyTime := time.Now().Add(-7 * 24 * time.Hour)
	dummyTime = dummyTime.Add(time.Duration(-1 * dummyTime.Nanosecond()))
	linkName := filepath.Join(dir, "log")
	rl, err := rotate.New(
		filepath.Join(dir, "log%Y%m%d%H%M%S"),
		rotate.WithMaxAge(24*time.Hour),
		rotate.WithLinkName(linkName),
	)
	if !assert.NoError(t, err, `rotate.New should succeed`) {
		return
	}
	defer rl.Close()

	str := "Hello, World"
	n, err := rl.Write([]byte(str))
	if !assert.NoError(t, err, "rl.Write should succeed") {
		return
	}

	if !assert.Len(t, str, n, "rl.Write should succeed") {
		return
	}

	fn := rl.CurrentFileName()
	if fn == "" {
		t.Errorf("Could not get filename %s", fn)
	}

	content, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Errorf("Failed to read file %s: %s", fn, err)
	}

	if string(content) != str {
		t.Errorf(`File content does not match (was "%s")`, content)
	}

	err = os.Chtimes(fn, dummyTime, dummyTime)
	if err != nil {
		t.Errorf("Failed to change access/modification times for %s: %s", fn, err)
	}

	fi, err := os.Stat(fn)
	if err != nil {
		t.Errorf("Failed to stat %s: %s", fn, err)
	}

	if !fi.ModTime().Equal(dummyTime) {
		t.Errorf("Failed to chtime for %s (expected %s, got %s)", fn, fi.ModTime(), dummyTime)
	}
}

func CreateRotationTestFile(dir string, base time.Time, d time.Duration, n int) {
	timestamp := base
	for i := 0; i < n; i++ {
		// %Y%m%d%H%M%S
		suffix := timestamp.Format("20060102150405")
		path := filepath.Join(dir, "log"+suffix)
		os.WriteFile(path, []byte("rotation test file\n"), os.ModePerm)
		os.Chtimes(path, timestamp, timestamp)
		timestamp = timestamp.Add(d)
	}
}

func TestLogRotationCount(t *testing.T) {
	dir, err := ioutil.TempDir("", "file-rotatelogs-rotationcount-test")
	if !assert.NoError(t, err, "creating temporary directory should succeed") {
		return
	}
	defer os.RemoveAll(dir)

	dummyTime := time.Now().Add(-7 * 24 * time.Hour)
	dummyTime = dummyTime.Add(time.Duration(-1 * dummyTime.Nanosecond()))

	t.Run("Either maxAge or rotationCount should be set", func(t *testing.T) {
		rl, err := rotate.New(
			filepath.Join(dir, "log%Y%m%d%H%M%S"),
			rotate.WithMaxAge(-1),
			rotate.WithRotationCount(-1),
		)
		if !assert.NoError(t, err, `Both of maxAge and rotationCount is disabled`) {
			return
		}
		defer rl.Close()
	})

	t.Run("Either maxAge or rotationCount should be set", func(t *testing.T) {
		rl, err := rotate.New(
			filepath.Join(dir, "log%Y%m%d%H%M%S"),
			rotate.WithMaxAge(1),
			rotate.WithRotationCount(1),
		)
		if !assert.NoError(t, err, `Both of maxAge and rotationCount is enabled`) {
			return
		}
		defer rl.Close()
	})

	t.Run("Only latest log file is kept", func(t *testing.T) {
		rl, err := rotate.New(
			filepath.Join(dir, "log%Y%m%d%H%M%S"),
			rotate.WithMaxAge(-1),
			rotate.WithRotationCount(1),
		)
		if !assert.NoError(t, err, `rotate.New should succeed`) {
			return
		}
		defer rl.Close()

		n, err := rl.Write([]byte("dummy"))
		if !assert.NoError(t, err, "rl.Write should succeed") {
			return
		}
		if !assert.Len(t, "dummy", n, "rl.Write should succeed") {
			return
		}
		time.Sleep(time.Second)
		files, err := filepath.Glob(filepath.Join(dir, "log*"))
		if !assert.Equal(t, 1, len(files), "Only latest log is kept") {
			return
		}
	})

	t.Run("Old log files are purged except 2 log files", func(t *testing.T) {
		CreateRotationTestFile(dir, dummyTime, time.Duration(time.Hour), 5)
		rl, err := rotate.New(
			filepath.Join(dir, "log%Y%m%d%H%M%S"),
			rotate.WithMaxAge(-1),
			rotate.WithRotationCount(2),
		)
		if !assert.NoError(t, err, `rotate.New should succeed`) {
			return
		}
		defer rl.Close()

		n, err := rl.Write([]byte("dummy"))
		if !assert.NoError(t, err, "rl.Write should succeed") {
			return
		}
		if !assert.Len(t, "dummy", n, "rl.Write should succeed") {
			return
		}
		time.Sleep(time.Second)
		files, err := filepath.Glob(filepath.Join(dir, "log*"))
		if !assert.Equal(t, 2, len(files), "One file is kept") {
			return
		}
	})

}

func TestLogSetOutput(t *testing.T) {
	dir, err := ioutil.TempDir("", "file-rotatelogs-test")
	if err != nil {
		t.Errorf("Failed to create temporary directory: %s", err)
	}
	defer os.RemoveAll(dir)

	rl, err := rotate.New(filepath.Join(dir, "log%Y%m%d%H%M%S"))
	if !assert.NoError(t, err, `rotate.New should succeed`) {
		return
	}
	defer rl.Close()

	log.SetOutput(rl)
	defer log.SetOutput(os.Stderr)

	str := "Hello, World"
	log.Print(str)

	fn := rl.CurrentFileName()
	if fn == "" {
		t.Errorf("Could not get filename %s", fn)
	}

	content, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Errorf("Failed to read file %s: %s", fn, err)
	}

	if !strings.Contains(string(content), str) {
		t.Errorf(`File content does not contain "%s" (was "%s")`, str, content)
	}
}
