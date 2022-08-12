package semaphore

import (
	"testing"
	"time"
)

func Test_Semaphore(t *testing.T) {
	var sem Semaphore
	if err := sem.Open("/testsem", 0644, 1); err != nil {
		t.Fatalf("Failed to open: %v", err)
	}

	if err := sem.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}

	if err := sem.Unlink(); err != nil {
		t.Fatalf("Failed to unlink: %v", err)
	}
}

func Test_SemWait(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem", 0644, 1)

	if err := sem.Wait(); err != nil {
		t.Fatalf("Failed to wait: %v", err)
	}

	if err := sem.Post(); err != nil {
		t.Fatalf("Failed to post: %v", err)
	}
	sem.Close()
}

func Test_SemTryWait(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem", 0644, 1)

	if err := sem.TryWait(); err != nil {
		t.Fatalf("Failed to wait: %v", err)
	}

	sem.Post()
	sem.Close()
}

func Test_SemTryWaitFail(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem", 0644, 1)

	sem.Wait()
	err := sem.TryWait()
	if err == nil {
		t.Fatal("TryWait should have failed")
	}
	sem.Post()
	sem.Close()
}

func Test_SemPostFail(t *testing.T) {
	var sem Semaphore
	err := sem.Post()
	if err == nil {
		t.Fatal("Post should have failed")
	}
}

func Test_SemTimedWait_nowait(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem", 0644, 1)
	if err := sem.TimedWait(1 * time.Second); err != nil {
		sem.Close()
		sem.Unlink()
		t.Fatalf("Should not have timedout: %v", err)
	}
	sem.Close()
	sem.Unlink()
}

func Test_SemTimedWait_wait(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem_wait", 0644, 1)
	sem.Wait()

	end := make(chan error, 1)
	go func() {
		var sem2 Semaphore
		sem2.Open("/testsem_wait", 0644, 1)
		semerr := sem2.TimedWait(2 * time.Second)
		sem2.Close()
		sem2.Unlink()
		end <- semerr
	}()

	time.Sleep(500 * time.Millisecond)
	sem.Post()
	sem.Close()

	err := <-end
	if err != nil {
		t.Fatalf("Should not have timedout: %v", err)
	}
}

func Test_SemTimedWait_timeout(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem_wait", 0644, 1)
	sem.Wait()

	end := make(chan error, 1)
	go func() {
		var sem2 Semaphore
		sem2.Open("/testsem_wait", 0644, 1)
		semerr := sem2.TimedWait(1 * time.Second)
		sem2.Close()
		end <- semerr
	}()

	time.Sleep(2 * time.Second)
	sem.Post()
	err := <-end
	sem.Close()
	sem.Unlink()
	if err == nil {
		t.Fatalf("Should have timedout: %v", err)
	}
}

func Test_SemDoubleClose(t *testing.T) {
	var sem Semaphore
	if err := sem.Open("/testsem", 0644, 1); err != nil {
		t.Fatalf("Failed to open: %v", err)
	}

	if err := sem.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}
	if err := sem.Close(); err == nil {
		t.Fatalf("Should have received error: %v", err)
	}
}

func Test_SemDoubleUnlink(t *testing.T) {
	var sem Semaphore
	if err := sem.Open("/testsem", 0644, 1); err != nil {
		t.Fatalf("Failed to open: %v", err)
	}

	if err := sem.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}

	if err := sem.Unlink(); err != nil {
		t.Fatalf("Failed to unlink: %v", err)
	}
	if err := sem.Unlink(); err == nil {
		t.Fatalf("Should have received error: %v", err)
	}
}

func Test_isSemaphoreInitialized(t *testing.T) {
	var sem Semaphore
	if err := sem.Close(); err == nil {
		t.Fatalf("Should have recived error: %v", err)
	}
	if err := sem.Post(); err == nil {
		t.Fatalf("Should have recived error: %v", err)
	}
	if err := sem.Wait(); err == nil {
		t.Fatalf("Should have recived error: %v", err)
	}
	if err := sem.TryWait(); err == nil {
		t.Fatalf("Should have recived error: %v", err)
	}
}
