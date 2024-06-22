#include <iostream>
#include <random>
#include <algorithm>
#include <unordered_set>
#include <cassert>
#include <chrono>
#include <cstdlib>

#include "search.h"

bool search_stl(const std::vector<int>& data, int value)
{
    return std::binary_search(data.begin(), data.end(), value);
}

template <typename F>
long bench_single(const F& binsearch, const std::vector<int>& sorted, const std::vector<int>& unsorted, size_t nlookups)
{
    auto start = std::chrono::high_resolution_clock::now();
    auto lookup = [&] (size_t i) {
        if (!binsearch(sorted, unsorted[i])) {
            std::cout << "Fail binsearch" << std::endl;
        }
    };
    for (size_t j = 0; j < nlookups / unsorted.size(); ++j) {
        for (size_t i = 0; i < unsorted.size(); ++i) {
            lookup(i);
        }
    }
    for (size_t i = 0; i < nlookups % unsorted.size(); ++i) {
        lookup(i);
    }
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::nanoseconds>(end - start).count();
    return duration;
}

template <typename F>
double bench_best(const F& binsearch, const std::vector<int>& sorted, const std::vector<int>& unsorted)
{
    size_t nlookups = std::max(1000UL, unsorted.size());
    nlookups = std::min(nlookups, 1UL<<16);
    auto best = bench_single(binsearch, sorted, unsorted, nlookups);
    for (int i = 0; i < 1000; ++i) {
        asm volatile ("");
        best = std::min(best, bench_single(binsearch, sorted, unsorted, nlookups));
    }
    return double(best) / nlookups;
}

void bench(const std::vector<int>& sorted, const std::vector<int>& unsorted)
{
    double basic = bench_best(search_stl, sorted, unsorted);
    double good = bench_best(search, sorted, unsorted);
    double ratio = basic / good;
    std::cout << "Search in " << sorted.size() << " elements.";
    std::cout << " std::binary_search " << basic << "ns per search.";
    std::cout << " Your " << good << "ns per search. Speedup: " << ratio << std::endl;
    assert(ratio > 1.5);
}

void test_performance(int n)
{
    std::vector<int> data(n);
    std::vector<int> sdata(n);
    std::random_device device;
    std::mt19937 mt(device());
    for (int i = 0; i < n; ++i) {
        data[i] = mt();
        sdata[i] = data[i];
    }
    std::sort(sdata.begin(), sdata.end());
    bench(sdata, data);
}

void validate(int n)
{
    std::vector<int> data(n);
    std::vector<int> sdata(n);
    std::random_device device;
    std::mt19937 mt(device());
    std::unordered_set<int> map;
    for (int i = 0; i < n; ++i) {
        data[i] = mt();
        sdata[i] = data[i];
        map.insert(data[i]);
    }
    std::sort(sdata.begin(), sdata.end());
    for (auto v : data) {
        assert(search(sdata, v));
    }
    for (int i = 0; i < n; ++i) {
        int v = mt();
        assert(search(sdata, v) == (map.find(v) != map.end()));
    }
}

void test_correctness()
{
    for (int i = 0; i < 100; ++i) {
        validate(i);
    }
    std::cout << "test_correctness PASSED" << std::endl; 
}

void test_performance()
{
    for (int i = 6; i < 50; ++i) {
        test_performance(i);
    }
    std::cout << "test_performance PASSED" << std::endl; 
}

int main()
{
    test_correctness();
    test_performance();
}
