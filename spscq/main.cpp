#include <iostream>
#include <chrono>
#include <cstdlib>
#include <thread>
#include <cassert>

#include "queue.h"

template <typename T>
void validate(T& queue, int n)
{
    auto producer = std::thread([&] {
        for (int i = 0; i < n; ++i) {
            while (!queue.push(i));
        }
    });
    auto consumer = std::thread([&] {
        for (int i = 0; i < n; ++i) {
            while (true) {
                auto v = queue.pop();
                if (v.has_value()) {
                    assert(*v == i);
                    break;
                }
            }
        }
    });

    producer.join();
    consumer.join();
}

template <typename T>
long bench_single(T& queue, int n)
{
    auto start = std::chrono::high_resolution_clock::now();
    validate(queue, n);
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::nanoseconds>(end - start).count();
    return duration;
}

template <typename T>
double bench_best(T& queue, int n)
{
    auto best = bench_single(queue, n);
    for (int i = 0; i < 10; ++i) {
        best = std::min(best, bench_single(queue, n));
    }
    return double(best);
}

void test_correctness()
{
    SPSCQueue<int> q;
    validate(q, 10000);
    std::cout << "test_correctness PASSED" << std::endl;
}

void test_performance(int n)
{
    SPSCQueue<int> q;
    double r = bench_best(q, n);
    double throughput = n*sizeof(int)/r * (1000000000./(1<<30));
    std::cout << "Pushed " << n << " values. Throughput: " << throughput << "GB/s" << std::endl;
    assert(throughput > 0.1);
}

void test_performance()
{
    for (int i = 17; i < 22; ++i) {
        test_performance(1<<i);
    }
    std::cout << "test_performance PASSED" << std::endl;
}

int main()
{
    test_correctness();
    test_performance();
}


