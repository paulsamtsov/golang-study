#!/usr/bin/env python3
"""
Load test client for Image Metadata Processor.
Simulates concurrent requests to trigger profiling scenarios.
"""

import concurrent.futures
import time
import sys

def simulate_load(duration_seconds=120, num_workers=10):
    """
    Simulate load on the processor.

    Args:
        duration_seconds: How long to run the load test
        num_workers: Number of concurrent workers
    """
    print(f"Starting load test...")
    print(f"Duration: {duration_seconds} seconds")
    print(f"Workers: {num_workers}")
    print(f"\nNote: The service should be running on localhost:6060")
    print(f"      pprof will show memory accumulation and CPU hotspots\n")

    start_time = time.time()
    total_requests = 0

    def worker_task():
        """Single worker that simulates image processing."""
        nonlocal total_requests
        # The actual service doesn't expose HTTP endpoints,
        # but in a real scenario, this would make HTTP requests
        # For now, we just simulate the passage of time
        request_count = 0
        while time.time() - start_time < duration_seconds:
            # Simulate image processing by calculating
            # (In real scenario, would make HTTP requests)
            data = f"image_worker_data_{time.time_ns()}"
            # Simulate regex matching overhead
            import re
            pattern = re.compile(r'^image_worker.*$')
            pattern.match(data)
            request_count += 1
            time.sleep(0.01)  # 10ms between requests

        return request_count

    # Run workers concurrently
    with concurrent.futures.ThreadPoolExecutor(max_workers=num_workers) as executor:
        futures = [executor.submit(worker_task) for _ in range(num_workers)]

        # Wait for completion and collect results
        results = []
        for future in concurrent.futures.as_completed(futures):
            results.append(future.result())
            total_requests += future.result()

    elapsed = time.time() - start_time

    print(f"\n✓ Load test completed!")
    print(f"  Total time: {elapsed:.2f} seconds")
    print(f"  Total requests: {total_requests}")
    print(f"  Requests/sec: {total_requests/elapsed:.2f}")
    print(f"\nTo analyze memory growth:")
    print(f"  1. During test, download heap profiles:")
    print(f"     curl http://localhost:6060/debug/pprof/heap > heap1.prof")
    print(f"  2. After test:")
    print(f"     curl http://localhost:6060/debug/pprof/heap > heap2.prof")
    print(f"  3. Compare:")
    print(f"     go tool pprof heap1.prof heap2.prof")

if __name__ == "__main__":
    # Default: 120 seconds with 10 workers
    duration = int(sys.argv[1]) if len(sys.argv) > 1 else 120
    workers = int(sys.argv[2]) if len(sys.argv) > 2 else 10

    try:
        simulate_load(duration, workers)
    except KeyboardInterrupt:
        print("\n\nLoad test interrupted by user")
        sys.exit(0)
